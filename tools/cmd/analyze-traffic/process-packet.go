package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
)

// processPacket matches src and dst of IPv4/IPv6 packets against matchers.
func processPacket(packet gopacket.Packet) error {
	var src, dst net.IP
	if packet.NetworkLayer().LayerType() == layers.LayerTypeIPv4 {
		src = make(net.IP, 4)
		dst = make(net.IP, 4)
	} else if packet.NetworkLayer().LayerType() == layers.LayerTypeIPv6 {
		src = make(net.IP, 16)
		dst = make(net.IP, 16)
	} else {
		// Skip non IP packets
		return nil
	}
	copy(src, packet.NetworkLayer().NetworkFlow().Src().Raw()[:len(src)])
	copy(dst, packet.NetworkLayer().NetworkFlow().Dst().Raw()[:len(dst)])

	if *aggregateLists {
		if srcMatch, err := aggregatedRanger.Contains(src); err != nil {
			logger.Errorw("error looking for match on SRC address", "err", err, "src", src)
			return err
		} else if srcMatch {
			aggregatedSrcMatches++
		}

		if dstMatch, err := aggregatedRanger.Contains(dst); err != nil {
			logger.Errorw("error looking for match on DST address", "err", err, "dst", src)
			return err
		} else if dstMatch {
			aggregatedDstMatches++
		}
	} else {
		for tag, matcher := range perBlocklistRanger {
			if srcMatch, err := matcher.Contains(src); err != nil {
				logger.Errorw("error looking for match on SRC address", "err", err, "src", src)
				return err
			} else if srcMatch {
				perBlocklistSrcMatches[tag]++
			}

			if dstMatch, err := matcher.Contains(dst); err != nil {
				logger.Errorw("error looking for match on DST address", "err", err, "dst", src)
				return err
			} else if dstMatch {
				perBlocklistDstMatches[tag]++
			}
		}
	}
	return nil
}
