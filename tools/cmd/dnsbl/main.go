// This main package compiles a CoreDNS server embedding the DNSBL plugin.
//
// If you are interested in the DNSBL plugin itself, see the plugin/ directory.
package main

import (
	_ "git.netsplit.it/enrico204/blocklists/tools/cmd/dnsbl-coredns/plugin"
	"github.com/coredns/coredns/core/dnsserver"
	_ "github.com/coredns/coredns/core/plugin"
	"github.com/coredns/coredns/coremain"
)

func init() {
	for i, name := range dnsserver.Directives {
		if name == "file" {
			dnsserver.Directives = append(dnsserver.Directives[:i],
				append([]string{"dnsbl"}, dnsserver.Directives[i:]...)...)
			return
		}
	}
	panic("file plugin not found in dnsserver.Directives")
}

func main() {
	coremain.Run()
}
