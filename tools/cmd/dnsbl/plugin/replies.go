package dnsbl

import (
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"net"
)

// replyWithNxDomain returns a NXDOMAIN for the query.
func replyWithNxDomain(state request.Request, soa *dns.SOA) (int, error) {
	m := new(dns.Msg)
	m.SetReply(state.Req)
	m.Authoritative = true
	m.Rcode = dns.RcodeNameError
	m.Ns = []dns.RR{soa}
	err := state.W.WriteMsg(m)
	if err != nil {
		return dns.RcodeServerFailure, plugin.Error("dnsbl", err)
	}
	return dns.RcodeSuccess, nil
}

// replyWithSOA returns a SOA for the query.
func replyWithSOA(state request.Request, soa *dns.SOA) (int, error) {
	m := new(dns.Msg)
	m.SetReply(state.Req)
	m.Authoritative = true
	m.Answer = append(m.Answer, soa)
	err := state.W.WriteMsg(m)
	if err != nil {
		return dns.RcodeServerFailure, plugin.Error("dnsbl", err)
	}
	return dns.RcodeSuccess, nil
}

// replyWithFound returns an A record for the query. The record point to 127.0.0.2.
func replyWithFound(state request.Request) (int, error) {
	m := new(dns.Msg)
	m.SetReply(state.Req)
	m.Authoritative = true
	m.Answer = append(m.Answer, &dns.A{
		Hdr: dns.RR_Header{
			Name:   state.Name(),
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    0,
		},
		A: net.IPv4(127, 0, 0, 2),
	})
	err := state.W.WriteMsg(m)
	if err != nil {
		return dns.RcodeServerFailure, plugin.Error("dnsbl", err)
	}
	return dns.RcodeSuccess, nil
}

// replyWithTXT returns an TXT record for the query.
func replyWithTXT(state request.Request) (int, error) {
	m := new(dns.Msg)
	m.SetReply(state.Req)
	m.Authoritative = true
	m.Answer = append(m.Answer, &dns.TXT{
		Hdr: dns.RR_Header{
			Name:   state.Name(),
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    0,
		},
		Txt: []string{"CoreDNS + DNSBL plugin"},
	})
	err := state.W.WriteMsg(m)
	if err != nil {
		return dns.RcodeServerFailure, plugin.Error("dnsbl", err)
	}
	return dns.RcodeSuccess, nil
}
