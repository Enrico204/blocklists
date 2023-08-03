package azure

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/file"
	"github.com/coredns/coredns/plugin/pkg/fall"
	"github.com/coredns/coredns/plugin/pkg/upstream"
	"github.com/coredns/coredns/request"

	publicdns "github.com/Azure/azure-sdk-for-go/profiles/latest/dns/mgmt/dns"
	privatedns "github.com/Azure/azure-sdk-for-go/profiles/latest/privatedns/mgmt/privatedns"
	"github.com/miekg/dns"
)

type zone struct {
	id      string
	z       *file.Zone
	zone    string
	private bool
}

type zones map[string][]*zone

// Azure is the core struct of the azure plugin.
type Azure struct {
	zoneNames     []string
	publicClient  publicdns.RecordSetsClient
	privateClient privatedns.RecordSetsClient
	upstream      *upstream.Upstream
	zMu           sync.RWMutex
	zones         zones

	Next plugin.Handler
	Fall fall.F
}

// New validates the input DNS zones and initializes the Azure struct.
func New(ctx context.Context, publicClient publicdns.RecordSetsClient, privateClient privatedns.RecordSetsClient, keys map[string][]string, accessMap map[string]string) (*Azure, error) {
	zones := make(map[string][]*zone, len(keys))
	names := make([]string, len(keys))
	var private bool

	for resourceGroup, znames := range keys {
		for _, name := range znames {
			switch accessMap[resourceGroup+name] {
			case "public":
				if _, err := publicClient.ListAllByDNSZone(context.Background(), resourceGroup, name, nil, ""); err != nil {
					return nil, err
				}
				private = false
			case "private":
				if _, err := privateClient.ListComplete(context.Background(), resourceGroup, name, nil, ""); err != nil {
					return nil, err
				}
				private = true
			}

			fqdn := dns.Fqdn(name)
			if _, ok := zones[fqdn]; !ok {
				names = append(names, fqdn)
			}
			zones[fqdn] = append(zones[fqdn], &zone{id: resourceGroup, zone: name, private: private, z: file.NewZone(fqdn, "")})
		}
	}

	return &Azure{
		publicClient:  publicClient,
		privateClient: privateClient,
		zones:         zones,
		zoneNames:     names,
		upstream:      upstream.New(),
	}, nil
}

// Run updates the zone from azure.
func (h *Azure) Run(ctx context.Context) error {
	if err := h.updateZones(ctx); err != nil {
		return err
	}
	go func() {
		delay := 1 * time.Minute
		timer := time.NewTimer(delay)
		defer timer.Stop()
		for {
			timer.Reset(delay)
			select {
			case <-ctx.Done():
				log.Debugf("Breaking out of Azure update loop for %v: %v", h.zoneNames, ctx.Err())
				return
			case <-timer.C:
				if err := h.updateZones(ctx); err != nil && ctx.Err() == nil {
					log.Errorf("Failed to update zones %v: %v", h.zoneNames, err)
				}
			}
		}
	}()
	return nil
}

func (h *Azure) updateZones(ctx context.Context) error {
	var err error
	var publicSet publicdns.RecordSetListResultPage
	var privateSet privatedns.RecordSetListResultPage
	errs := make([]string, 0)
	for zName, z := range h.zones {
		for i, hostedZone := range z {
			newZ := file.NewZone(zName, "")
			if hostedZone.private {
				for privateSet, err = h.privateClient.List(ctx, hostedZone.id, hostedZone.zone, nil, ""); privateSet.NotDone(); err = privateSet.NextWithContext(ctx) {
					updateZoneFromPrivateResourceSet(privateSet, newZ)
				}
			} else {
				for publicSet, err = h.publicClient.ListByDNSZone(ctx, hostedZone.id, hostedZone.zone, nil, ""); publicSet.NotDone(); err = publicSet.NextWithContext(ctx) {
					updateZoneFromPublicResourceSet(publicSet, newZ)
				}
			}
			if err != nil {
				errs = append(errs, fmt.Sprintf("failed to list resource records for %v from azure: %v", hostedZone.zone, err))
			}
			newZ.Upstream = h.upstream
			h.zMu.Lock()
			(*z[i]).z = newZ
			h.zMu.Unlock()
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("errors updating zones: %v", errs)
	}
	return nil
}

func updateZoneFromPublicResourceSet(recordSet publicdns.RecordSetListResultPage, newZ *file.Zone) {
	for _, result := range *(recordSet.Response().Value) {
		resultFqdn := *(result.RecordSetProperties.Fqdn)
		resultTTL := uint32(*(result.RecordSetProperties.TTL))
		if result.RecordSetProperties.ARecords != nil {
			for _, A := range *(result.RecordSetProperties.ARecords) {
				a := &dns.A{Hdr: dns.RR_Header{Name: resultFqdn, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: resultTTL},
					A: net.ParseIP(*(A.Ipv4Address))}
				newZ.Insert(a)
			}
		}

		if result.RecordSetProperties.AaaaRecords != nil {
			for _, AAAA := range *(result.RecordSetProperties.AaaaRecords) {
				aaaa := &dns.AAAA{Hdr: dns.RR_Header{Name: resultFqdn, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: resultTTL},
					AAAA: net.ParseIP(*(AAAA.Ipv6Address))}
				newZ.Insert(aaaa)
			}
		}

		if result.RecordSetProperties.MxRecords != nil {
			for _, MX := range *(result.RecordSetProperties.MxRecords) {
				mx := &dns.MX{Hdr: dns.RR_Header{Name: resultFqdn, Rrtype: dns.TypeMX, Class: dns.ClassINET, Ttl: resultTTL},
					Preference: uint16(*(MX.Preference)),
					Mx:         dns.Fqdn(*(MX.Exchange))}
				newZ.Insert(mx)
			}
		}

		if result.RecordSetProperties.PtrRecords != nil {
			for _, PTR := range *(result.RecordSetProperties.PtrRecords) {
				ptr := &dns.PTR{Hdr: dns.RR_Header{Name: resultFqdn, Rrtype: dns.TypePTR, Class: dns.ClassINET, Ttl: resultTTL},
					Ptr: dns.Fqdn(*(PTR.Ptrdname))}
				newZ.Insert(ptr)
			}
		}

		if result.RecordSetProperties.SrvRecords != nil {
			for _, SRV := range *(result.RecordSetProperties.SrvRecords) {
				srv := &dns.SRV{Hdr: dns.RR_Header{Name: resultFqdn, Rrtype: dns.TypeSRV, Class: dns.ClassINET, Ttl: resultTTL},
					Priority: uint16(*(SRV.Priority)),
					Weight:   uint16(*(SRV.Weight)),
					Port:     uint16(*(SRV.Port)),
					Target:   dns.Fqdn(*(SRV.Target))}
				newZ.Insert(srv)
			}
		}

		if result.RecordSetProperties.TxtRecords != nil {
			for _, TXT := range *(result.RecordSetProperties.TxtRecords) {
				txt := &dns.TXT{Hdr: dns.RR_Header{Name: resultFqdn, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: resultTTL},
					Txt: *(TXT.Value)}
				newZ.Insert(txt)
			}
		}

		if result.RecordSetProperties.NsRecords != nil {
			for _, NS := range *(result.RecordSetProperties.NsRecords) {
				ns := &dns.NS{Hdr: dns.RR_Header{Name: resultFqdn, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: resultTTL},
					Ns: *(NS.Nsdname)}
				newZ.Insert(ns)
			}
		}

		if result.RecordSetProperties.SoaRecord != nil {
			SOA := result.RecordSetProperties.SoaRecord
			soa := &dns.SOA{Hdr: dns.RR_Header{Name: resultFqdn, Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: resultTTL},
				Minttl:  uint32(*(SOA.MinimumTTL)),
				Expire:  uint32(*(SOA.ExpireTime)),
				Retry:   uint32(*(SOA.RetryTime)),
				Refresh: uint32(*(SOA.RefreshTime)),
				Serial:  uint32(*(SOA.SerialNumber)),
				Mbox:    dns.Fqdn(*(SOA.Email)),
				Ns:      *(SOA.Host)}
			newZ.Insert(soa)
		}

		if result.RecordSetProperties.CnameRecord != nil {
			CNAME := result.RecordSetProperties.CnameRecord.Cname
			cname := &dns.CNAME{Hdr: dns.RR_Header{Name: resultFqdn, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: resultTTL},
				Target: dns.Fqdn(*CNAME)}
			newZ.Insert(cname)
		}
	}
}

func updateZoneFromPrivateResourceSet(recordSet privatedns.RecordSetListResultPage, newZ *file.Zone) {
	for _, result := range *(recordSet.Response().Value) {
		resultFqdn := *(result.RecordSetProperties.Fqdn)
		resultTTL := uint32(*(result.RecordSetProperties.TTL))
		if result.RecordSetProperties.ARecords != nil {
			for _, A := range *(result.RecordSetProperties.ARecords) {
				a := &dns.A{Hdr: dns.RR_Header{Name: resultFqdn, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: resultTTL},
					A: net.ParseIP(*(A.Ipv4Address))}
				newZ.Insert(a)
			}
		}
		if result.RecordSetProperties.AaaaRecords != nil {
			for _, AAAA := range *(result.RecordSetProperties.AaaaRecords) {
				aaaa := &dns.AAAA{Hdr: dns.RR_Header{Name: resultFqdn, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: resultTTL},
					AAAA: net.ParseIP(*(AAAA.Ipv6Address))}
				newZ.Insert(aaaa)
			}
		}

		if result.RecordSetProperties.MxRecords != nil {
			for _, MX := range *(result.RecordSetProperties.MxRecords) {
				mx := &dns.MX{Hdr: dns.RR_Header{Name: resultFqdn, Rrtype: dns.TypeMX, Class: dns.ClassINET, Ttl: resultTTL},
					Preference: uint16(*(MX.Preference)),
					Mx:         dns.Fqdn(*(MX.Exchange))}
				newZ.Insert(mx)
			}
		}

		if result.RecordSetProperties.PtrRecords != nil {
			for _, PTR := range *(result.RecordSetProperties.PtrRecords) {
				ptr := &dns.PTR{Hdr: dns.RR_Header{Name: resultFqdn, Rrtype: dns.TypePTR, Class: dns.ClassINET, Ttl: resultTTL},
					Ptr: dns.Fqdn(*(PTR.Ptrdname))}
				newZ.Insert(ptr)
			}
		}

		if result.RecordSetProperties.SrvRecords != nil {
			for _, SRV := range *(result.RecordSetProperties.SrvRecords) {
				srv := &dns.SRV{Hdr: dns.RR_Header{Name: resultFqdn, Rrtype: dns.TypeSRV, Class: dns.ClassINET, Ttl: resultTTL},
					Priority: uint16(*(SRV.Priority)),
					Weight:   uint16(*(SRV.Weight)),
					Port:     uint16(*(SRV.Port)),
					Target:   dns.Fqdn(*(SRV.Target))}
				newZ.Insert(srv)
			}
		}

		if result.RecordSetProperties.TxtRecords != nil {
			for _, TXT := range *(result.RecordSetProperties.TxtRecords) {
				txt := &dns.TXT{Hdr: dns.RR_Header{Name: resultFqdn, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: resultTTL},
					Txt: *(TXT.Value)}
				newZ.Insert(txt)
			}
		}

		if result.RecordSetProperties.SoaRecord != nil {
			SOA := result.RecordSetProperties.SoaRecord
			soa := &dns.SOA{Hdr: dns.RR_Header{Name: resultFqdn, Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: resultTTL},
				Minttl:  uint32(*(SOA.MinimumTTL)),
				Expire:  uint32(*(SOA.ExpireTime)),
				Retry:   uint32(*(SOA.RetryTime)),
				Refresh: uint32(*(SOA.RefreshTime)),
				Serial:  uint32(*(SOA.SerialNumber)),
				Mbox:    dns.Fqdn(*(SOA.Email)),
				Ns:      dns.Fqdn(*(SOA.Host))}
			newZ.Insert(soa)
		}

		if result.RecordSetProperties.CnameRecord != nil {
			CNAME := result.RecordSetProperties.CnameRecord.Cname
			cname := &dns.CNAME{Hdr: dns.RR_Header{Name: resultFqdn, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: resultTTL},
				Target: dns.Fqdn(*CNAME)}
			newZ.Insert(cname)
		}
	}
}

// ServeDNS implements the plugin.Handler interface.
func (h *Azure) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	qname := state.Name()

	zone := plugin.Zones(h.zoneNames).Matches(qname)
	if zone == "" {
		return plugin.NextOrFailure(h.Name(), h.Next, ctx, w, r)
	}

	zones, ok := h.zones[zone] // ok true if we are authoritative for the zone.
	if !ok || zones == nil {
		return dns.RcodeServerFailure, nil
	}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	var result file.Result
	for _, z := range zones {
		h.zMu.RLock()
		m.Answer, m.Ns, m.Extra, result = z.z.Lookup(ctx, state, qname)
		h.zMu.RUnlock()

		// record type exists for this name (NODATA).
		if len(m.Answer) != 0 || result == file.NoData {
			break
		}
	}

	if len(m.Answer) == 0 && result != file.NoData && h.Fall.Through(qname) {
		return plugin.NextOrFailure(h.Name(), h.Next, ctx, w, r)
	}

	switch result {
	case file.Success:
	case file.NoData:
	case file.NameError:
		m.Rcode = dns.RcodeNameError
	case file.Delegation:
		m.Authoritative = false
	case file.ServerFailure:
		return dns.RcodeServerFailure, nil
	}

	w.WriteMsg(m)
	return dns.RcodeSuccess, nil
}

// Name implements plugin.Handler.Name.
func (h *Azure) Name() string { return "azure" }
