package utils

import (
	"net"
	"strings"
)

func ParseCIDR(addr string) (*net.IPNet, error) {
	addr = strings.TrimSpace(addr)

	if !strings.ContainsRune(addr, '/') {
		if !strings.ContainsRune(addr, ':') {
			addr += "/32"
		} else {
			addr += "/128"
		}
	}
	_, addrnet, err := net.ParseCIDR(addr)
	return addrnet, err
}
