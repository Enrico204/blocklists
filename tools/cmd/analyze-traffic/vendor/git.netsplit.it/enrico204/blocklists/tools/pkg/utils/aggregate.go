package utils

import (
	agg "github.com/ldkingvivi/go-aggregate"
	"net"
)

// Aggregate returns the aggregated list of networks.
func Aggregate(addresses []*net.IPNet) []*net.IPNet {
	var aggregate = make([]agg.CidrEntry, 0, len(addresses))
	for _, addr := range addresses {
		aggregate = append(aggregate, agg.NewBasicCidrEntry(addr))
	}
	aggregate = agg.Aggregate(aggregate, func(_, _ agg.CidrEntry) {})

	var ret = make([]*net.IPNet, 0, len(aggregate))
	for _, addr := range aggregate {
		ret = append(ret, addr.GetNetwork())
	}
	return ret
}
