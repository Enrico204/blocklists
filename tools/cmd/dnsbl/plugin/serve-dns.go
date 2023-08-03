package dnsbl

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"strings"
)

// ServeDNS implements the plugin.Handler interface. This method gets called when example is used
// in a Server.
func (d DNSBL) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	// Check if the request is for a zone we are serving. If it doesn't match we pass the request on to the next plugin.
	zoneName := d.Names.Matches(state.Name())
	if zoneName == "" {
		return plugin.NextOrFailure(d.Name(), d.Next, ctx, w, r)
	}
	state.Zone = zoneName

	// Export metric with the server label set to the current server handling the request.
	requestCount.WithLabelValues(metrics.WithServer(ctx)).Inc()

	// Get the zone matcher.
	zone, ok := d.Zones[zoneName]
	if !ok {
		return dns.RcodeServerFailure, nil
	}

	// Strip zone from name.
	name := strings.TrimSuffix(state.Name(), zoneName)
	queryType := state.QType()

	log.Debugf("received query for: %s type: %s", name, dns.TypeToString[queryType])

	// DNSBL protocol uses only A or TXT queries.
	switch queryType {
	case dns.TypeSOA:
		return replyWithSOA(state, zone.soa)
	case dns.TypeA:
		return handleLookupQuery(ctx, name, zone, state)
	case dns.TypeTXT:
		return replyWithTXT(state)
	default:
		log.Warningf("client asked for wrong qtype: %s", dns.TypeToString[queryType])
		return dns.RcodeNotImplemented, nil
	}
}

// handleLookupQuery handles the DNSBL lookup query by reconstructing the IP and doing the lookup in IP lists.
func handleLookupQuery(ctx context.Context, name string, zone *Zone, state request.Request) (int, error) {
	// Extract and build the IP address from the DNS query.
	ipaddr := extractIPFromQuery(name)
	if ipaddr == nil {
		log.Warningf("invalid ip address format: %s", name)
		return dns.RcodeFormatError, nil
	}

	// Look for the IP address in lists we loaded.
	if found, err := zone.matcher.Contains(ipaddr); err != nil {
		log.Warningf("error in matcher for %s: %s", name, err.Error())
		return dns.RcodeServerFailure, err
	} else if found {
		// The IP address is in the lists, reply with a record.
		matchesCount.WithLabelValues(metrics.WithServer(ctx)).Inc()
		return replyWithFound(state)
	}

	// No match is found, reply with NXDOMAIN.
	return replyWithNxDomain(state, zone.soa)
}
