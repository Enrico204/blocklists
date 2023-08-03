package utils

import (
	"math/big"
	"net"
	"net/netip"
)

// IPrangeToCIDR returns the list of CIDR that describes the range indicated by startIP and endIP (included).
func IPrangeToCIDR(startIP, endIP net.IP) []*net.IPNet {
	var ret []*net.IPNet

	start, ok := netip.AddrFromSlice(startIP)
	if !ok {
		panic(startIP.String())
	}

	end, ok := netip.AddrFromSlice(endIP)
	if !ok {
		panic(endIP.String())
	}

	if start.Is4() != end.Is4() {
		panic("start and end should be of the same type")
	}

	// If the range was specified with start and end swapped, restore the correct order.
	if start.Compare(end) > 0 {
		start, end = end, start
	}

	// Start, end, "next" and "mask" addresses as big integers to support IPv6.
	var (
		startInt = new(big.Int).SetBytes(start.AsSlice())
		endInt   = new(big.Int).SetBytes(end.AsSlice())
		nextIP   = new(big.Int)
		mask     = new(big.Int)
	)

	// Here we don't need big integers strictly, however it's useful to have them to avoid unnecessary type conversions
	// in the loop when doing the math.
	var (
		maxBit = new(big.Int)
		cmpSh  = new(big.Int)
		bits   = new(big.Int)
		one    = big.NewInt(1)
	)

	// Buffer for bigint to address conversion.
	var addrbuf []byte

	// Determine maximum mask bits and initialize buffer for IPv4 or IPv6.
	if start.Is4() {
		maxBit.SetUint64(32)
		addrbuf = make([]byte, 4)
	} else {
		maxBit.SetUint64(128)
		addrbuf = make([]byte, 16)
	}

	// Main loop to generate CIDR blocks
	for {
		// Initialize the bits counter and initial mask to 1.
		bits.SetUint64(1)
		mask.SetUint64(1)

		// Loop until bits reach the maximum subnet mask.
		for bits.Cmp(maxBit) < 0 {
			// Calculate the next IP address by applying the mask to the current IP.
			nextIP.Or(startInt, mask)

			// Convert the current bits to an unsigned integer for shifting.
			var bitShiftBuffer = uint(bits.Uint64())

			// Zero the least significant bits.
			cmpSh.Rsh(startInt, bitShiftBuffer)
			cmpSh.Lsh(cmpSh, bitShiftBuffer)

			// If the next IP is greater than the end IP or the shifted value is not equal to the current IP,
			// then it means the current mask is the largest possible, so we break the loop.
			if (nextIP.Cmp(endInt) > 0) || (cmpSh.Cmp(startInt) != 0) {
				bits.Sub(bits, one)
				mask.Rsh(mask, 1)
				break
			}

			// Increment the bits counter and shift the mask one bit to the left and add one to it.
			bits.Add(bits, one)
			mask.Add(mask.Lsh(mask, 1), one)
		}

		// Convert the current IP and calculated mask to an IPv4 or IPv6 address and CIDR block,
		// then add it to the cidr array.
		addr, _ := netip.AddrFromSlice(startInt.FillBytes(addrbuf))
		ret = append(ret, &net.IPNet{
			IP:   addr.AsSlice(),
			Mask: net.CIDRMask(int(bits.Sub(maxBit, bits).Uint64()), int(maxBit.Uint64())),
		})

		// Update the next IP address
		if nextIP.Or(startInt, mask); nextIP.Cmp(endInt) >= 0 {
			// If the next IP is greater than or equal to the end IP, we have covered the entire range,
			// so we break out of the loop.
			break
		}

		// Increment the current IP to the next IP.
		startInt.Add(nextIP, one)
	}
	return ret
}
