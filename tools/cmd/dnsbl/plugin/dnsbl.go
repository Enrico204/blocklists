// Package dnsbl is a CoreDNS plugin that matches incoming DNSBL queries against IP lists.
//
// A DNSBL query is a special DNS query for A record. The FQDN is composed by a zone (e.g., bl.example.com) and the IPv4
// address written in reverse octet form.
// E.g., to query the address 1.2.3.4 in bl.example.com, the query is towards 4.3.2.1.bl.example.com.
// For IPv6, the DNSBL is queries using each hexadecimal character as label. E.g.,
// Zone: bl.example.com
// IPv6 to query: 2001:0db8:0000:0000:0000:0000:0000:0001
// Query: 1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.bl.example.com
package dnsbl

import (
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"
	"github.com/yl2chen/cidranger"
)

// Define log to be a logger with the plugin name in it. This way we can just use log.Info and
// friends to log.
var log = clog.NewWithPlugin("dnsbl")

// DNSBL is an example plugin to show how to write a plugin.
type DNSBL struct {
	Next  plugin.Handler
	Zones Zones
	Names plugin.Zones
}

func (d DNSBL) Name() string { return "dnsbl" }

type Zones map[string]*Zone

type Zone struct {
	// name is the name of the zone we are authoritative for
	name string

	// matcher is the matcher
	matcher cidranger.Ranger

	soa *dns.SOA
}
