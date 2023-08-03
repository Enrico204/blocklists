package acl

import (
	"net"
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/infobloxopen/go-trees/iptree"
	"github.com/miekg/dns"
)

const pluginName = "acl"

func init() { plugin.Register(pluginName, setup) }

func newDefaultFilter() *iptree.Tree {
	defaultFilter := iptree.NewTree()
	_, IPv4All, _ := net.ParseCIDR("0.0.0.0/0")
	_, IPv6All, _ := net.ParseCIDR("::/0")
	defaultFilter.InplaceInsertNet(IPv4All, struct{}{})
	defaultFilter.InplaceInsertNet(IPv6All, struct{}{})
	return defaultFilter
}

func setup(c *caddy.Controller) error {
	a, err := parse(c)
	if err != nil {
		return plugin.Error(pluginName, err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		a.Next = next
		return a
	})

	return nil
}

func parse(c *caddy.Controller) (ACL, error) {
	a := ACL{}
	for c.Next() {
		r := rule{}
		args := c.RemainingArgs()
		r.zones = plugin.OriginsFromArgsOrServerBlock(args, c.ServerBlockKeys)

		for c.NextBlock() {
			p := policy{}

			action := strings.ToLower(c.Val())
			if action == "allow" {
				p.action = actionAllow
			} else if action == "block" {
				p.action = actionBlock
			} else if action == "filter" {
				p.action = actionFilter
			} else if action == "drop" {
				p.action = actionDrop
			} else {
				return a, c.Errf("unexpected token %q; expect 'allow', 'block', 'filter' or 'drop'", c.Val())
			}

			p.qtypes = make(map[uint16]struct{})
			p.filter = iptree.NewTree()

			hasTypeSection := false
			hasNetSection := false

			remainingTokens := c.RemainingArgs()
			for len(remainingTokens) > 0 {
				if !isPreservedIdentifier(remainingTokens[0]) {
					return a, c.Errf("unexpected token %q; expect 'type | net'", remainingTokens[0])
				}
				section := strings.ToLower(remainingTokens[0])

				i := 1
				var tokens []string
				for ; i < len(remainingTokens) && !isPreservedIdentifier(remainingTokens[i]); i++ {
					tokens = append(tokens, remainingTokens[i])
				}
				remainingTokens = remainingTokens[i:]

				if len(tokens) == 0 {
					return a, c.Errf("no token specified in %q section", section)
				}

				switch section {
				case "type":
					hasTypeSection = true
					for _, token := range tokens {
						if token == "*" {
							p.qtypes[dns.TypeNone] = struct{}{}
							break
						}
						qtype, ok := dns.StringToType[token]
						if !ok {
							return a, c.Errf("unexpected token %q; expect legal QTYPE", token)
						}
						p.qtypes[qtype] = struct{}{}
					}
				case "net":
					hasNetSection = true
					for _, token := range tokens {
						if token == "*" {
							p.filter = newDefaultFilter()
							break
						}
						token = normalize(token)
						_, source, err := net.ParseCIDR(token)
						if err != nil {
							return a, c.Errf("illegal CIDR notation %q", token)
						}
						p.filter.InplaceInsertNet(source, struct{}{})
					}
				default:
					return a, c.Errf("unexpected token %q; expect 'type | net'", section)
				}
			}

			// optional `type` section means all record types.
			if !hasTypeSection {
				p.qtypes[dns.TypeNone] = struct{}{}
			}

			// optional `net` means all ip addresses.
			if !hasNetSection {
				p.filter = newDefaultFilter()
			}

			r.policies = append(r.policies, p)
		}
		a.Rules = append(a.Rules, r)
	}
	return a, nil
}

func isPreservedIdentifier(token string) bool {
	identifier := strings.ToLower(token)
	return identifier == "type" || identifier == "net"
}

// normalize appends '/32' for any single IPv4 address and '/128' for IPv6.
func normalize(rawNet string) string {
	if idx := strings.IndexAny(rawNet, "/"); idx >= 0 {
		return rawNet
	}

	if idx := strings.IndexAny(rawNet, ":"); idx >= 0 {
		return rawNet + "/128"
	}
	return rawNet + "/32"
}
