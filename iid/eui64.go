package iid

import (
	"net"
)

// MakeEUI64Addr takes an IPv6 address, a hardware MAC address and a scope as
// input and uses them to generate an Interface Identifier suitable for use
// in link local, global unicast and Stateless Address Autoconfiguration
// (SLAAC) addresses (but see RFC4941 for problems with this last case).
//
// The IP is assumed to be a /64, and the hardware address must be either 48
// or 64 bits. The hardware address will be appended to the IP address as per
// RFC4291 section 2.5.1 and be modified as follows:
//
// * the 7th bit of the first octet (the 'X' bit in the EUI-64 format) may be
//   modified. If ScopeGlobal is passed, the bit will be set to 1, it will be
//   set to 0 for ScopeLocal, and ScopeInvert will cause 0 to become 1 or 1 to
//   become 0. If ScopeNone is passed the bit is left alone. See 'NOTE' below
//   for the rationale here
//
// * if the address is 48 bits, the octets 0xFFFE are inserted in the middle
//   of the address to pad it to 64 bits
//
// NOTE: there is some ambiguity to the RFC here. Most discussions I've seen
// on the topic say that the 7th bit should _always_ be inverted, but the RFC
// reads like the IPv6 EUI64 format uses the _inverse signal_ from the IEEE
// EUI64 format; so where the IEEE uses 0 for global scoping, the IPv6 IID
// should use 1. This function punts on the question and provides for all
// interpretations via the Scope parameter
func MakeEUI64Addr(ip net.IP, hw net.HardwareAddr, scope Scope) net.IP {
	tag := []byte{0xff, 0xfe}

	if len(ip) < 16 {
		return nil
	}

	if len(hw) < 6 || len(hw) > 8 {
		return nil
	}

	eui64 := make([]byte, 16)
	copy(eui64, ip)

	if len(hw) == 6 {
		hw = append(hw[:3], append(tag, hw[3:]...)...)
	}

	copy(eui64[8:], hw)

	switch scope {
	case ScopeGlobal:
		eui64[8] |= 1 << 1  // set 0 or 1 -> 1

	case ScopeLocal:
		eui64[8] &^= 1 << 1 // set 0 or 1 -> 0

	case ScopeInvert:
		eui64[8] ^= 1 << 1  // set 0 -> 1 or 1 -> 0
	default:
	}

	return eui64
}