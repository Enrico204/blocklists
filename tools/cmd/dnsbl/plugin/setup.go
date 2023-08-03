package dnsbl

import (
	"fmt"
	"git.netsplit.it/enrico204/blocklists/tools/pkg/blocklists"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
	"github.com/yl2chen/cidranger"
	"strconv"
	"time"
)

// init registers this plugin.
func init() {
	plugin.Register("dnsbl", setup)
}

// parseConfig parses the CoreDNS configuration and loads IP lists in memory for matching.
func parseConfig(c *caddy.Controller) (Zones, plugin.Zones, error) {
	var configuredZones = make(map[string]*Zone)
	var configuredZonesNames []string

	for c.Next() {
		// DNSBL zone block.
		args := c.RemainingArgs()
		if len(args) != 1 {
			return nil, nil, fmt.Errorf("expected 1 args, got %d", len(args))
		}
		zone := &Zone{
			name:    dns.Fqdn(args[0]),
			matcher: cidranger.NewPCTrieRanger(),
		}

		// Exit with error if the config adds a duplicated zone
		_, ok := configuredZones[zone.name]
		if ok {
			return nil, nil, fmt.Errorf("duplicate zone name %s", zone.name)
		}
		configuredZonesNames = append(configuredZonesNames, zone.name)
		configuredZones[zone.name] = zone

		// Configure default SOA record
		serial, _ := strconv.ParseUint(time.Now().Format("2006010215"), 10, 32)
		zone.soa = &dns.SOA{
			Hdr: dns.RR_Header{
				Name:   zone.name,
				Rrtype: dns.TypeSOA,
				Class:  dns.ClassINET,
				Ttl:    60,
			},
			Ns:      fmt.Sprintf("ns1.%s", zone.name),
			Mbox:    fmt.Sprintf("hostmaster.%s", zone.name),
			Serial:  uint32(serial),
			Refresh: 3600,
			Retry:   180,
			Expire:  3600000,
			Minttl:  10,
		}

		// Enter the inner block {}
		var listsToLoad []string
		for c.NextBlock() {
			switch c.Val() {
			case "lists":
				// lists <name> [path ...]
				args = c.RemainingArgs()
				if len(args) < 1 {
					return nil, nil, fmt.Errorf("no lists specified for zone %s", zone.name)
				}
				listsToLoad = args
			case "update_every":
				// update_every <name>
				args = c.RemainingArgs()
				if len(args) != 1 {
					return nil, nil, fmt.Errorf("error in update_every for %s", zone.name)
				}
				updateEvery, err := time.ParseDuration(args[1])
				if err != nil {
					return nil, nil, fmt.Errorf("error parsing update_every: %w", err)
				}
				zone.soa.Refresh = uint32(updateEvery / time.Second)
			case "soa_nameserver":
				// soa_hostname <name>
				args = c.RemainingArgs()
				if len(args) != 1 {
					return nil, nil, fmt.Errorf("error in soa_nameserver for %s", zone.name)
				}
				zone.soa.Ns = args[1]
			case "soa_email":
				// soa_email <name>
				args = c.RemainingArgs()
				if len(args) != 1 {
					return nil, nil, fmt.Errorf("error in soa_email for %s", zone.name)
				}
				zone.soa.Mbox = args[1]

			default:
				return nil, nil, c.ArgErr()
			}
		}

		// If no list has been specified, we have a configuration error.
		if len(listsToLoad) == 0 {
			return nil, nil, fmt.Errorf("no lists specified for zone %s", zone.name)
		}

		// Load all IP lists specified in the config.
		for _, path := range listsToLoad {
			addresses, err := blocklists.ReadBlocklistFile(zapToClog(), path)
			if err != nil {
				return nil, nil, fmt.Errorf("error loading list %s for zone %s: %w", path, zone.name, err)
			} else if len(addresses) == 0 {
				return nil, nil, fmt.Errorf("no address loaded for zone %s, check the list validity", zone.name)
			}

			for _, addr := range addresses {
				err := zone.matcher.Insert(cidranger.NewBasicRangerEntry(*addr))
				if err != nil {
					return nil, nil, fmt.Errorf("error adding address to matcher %s: %w", zone.name, err)
				}
			}
		}
	}

	return configuredZones, configuredZonesNames, nil
}

// setup is the function that gets called when the config parser see the token "dnsbl".
func setup(c *caddy.Controller) error {
	zones, names, err := parseConfig(c)
	if err != nil {
		return plugin.Error("dnsbl", err)
	}

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return &DNSBL{
			Next:  next,
			Zones: zones,
			Names: names,
		}
	})
	return nil
}
