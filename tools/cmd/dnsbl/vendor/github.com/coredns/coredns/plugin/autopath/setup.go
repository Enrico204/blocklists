package autopath

import (
	"fmt"
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/miekg/dns"
)

func init() { plugin.Register("autopath", setup) }

func setup(c *caddy.Controller) error {
	ap, mw, err := autoPathParse(c)
	if err != nil {
		return plugin.Error("autopath", err)
	}

	// Do this in OnStartup, so all plugin has been initialized.
	c.OnStartup(func() error {
		m := dnsserver.GetConfig(c).Handler(mw)
		if m == nil {
			return nil
		}
		if x, ok := m.(AutoPather); ok {
			ap.searchFunc = x.AutoPath
		} else {
			return plugin.Error("autopath", fmt.Errorf("%s does not implement the AutoPather interface", mw))
		}
		return nil
	})

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		ap.Next = next
		return ap
	})

	return nil
}

func autoPathParse(c *caddy.Controller) (*AutoPath, string, error) {
	ap := &AutoPath{}
	mw := ""

	for c.Next() {
		zoneAndresolv := c.RemainingArgs()
		if len(zoneAndresolv) < 1 {
			return ap, "", fmt.Errorf("no resolv-conf specified")
		}
		resolv := zoneAndresolv[len(zoneAndresolv)-1]
		if strings.HasPrefix(resolv, "@") {
			mw = resolv[1:]
		} else {
			// assume file on disk
			rc, err := dns.ClientConfigFromFile(resolv)
			if err != nil {
				return ap, "", fmt.Errorf("failed to parse %q: %v", resolv, err)
			}
			ap.search = rc.Search
			plugin.Zones(ap.search).Normalize()
			ap.search = append(ap.search, "") // sentinel value as demanded.
		}
		zones := zoneAndresolv[:len(zoneAndresolv)-1]
		ap.Zones = plugin.OriginsFromArgsOrServerBlock(zones, c.ServerBlockKeys)
	}
	return ap, mw, nil
}
