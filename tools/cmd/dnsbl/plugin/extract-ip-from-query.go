package dnsbl

import (
	"fmt"
	"net"
	"strings"
)

func extractIPFromQuery(name string) net.IP {
	if strings.Count(name, ".") == 4 {
		// IPv4
		var ipAddrParts = strings.Split(name, ".")
		var ipaddr = net.ParseIP(fmt.Sprintf("%s.%s.%s.%s", ipAddrParts[3], ipAddrParts[2], ipAddrParts[1], ipAddrParts[0]))
		return ipaddr
	} else if strings.Count(name, ".") == 32 {
		// IPv6
		var ipAddrParts = strings.Split(name[:len(name)-1], ".")
		var ipAddrBuf strings.Builder

		for i := len(ipAddrParts) - 1; i >= 0; i-- {
			ipAddrBuf.WriteString(ipAddrParts[i])
			if i > 0 && i%4 == 0 {
				ipAddrBuf.WriteRune(':')
			}
		}
		return net.ParseIP(ipAddrBuf.String())
	}
	log.Warningf("the client sent a query that I don't understand, received %s", name)
	return nil
}
