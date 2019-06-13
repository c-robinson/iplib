package iid

import "net"

// SetEUI64Addr takes an IPv6 address, a hardware MAC address and a scope as
// input and uses them to generate an Interface Identifier suitable for use
// in link local, global unicast and Stateless Address Autoconfiguration
// (SLAAC) addresses (but see RFC4941 for problems with this last case).
//
// The IP is assumed to be a /64, and the hardware address must be either 48
// or 64 bits. The hardware address will be appended to the IP address as per
// RFC4291 section 2.5.1 and be modified as follows:
//
// * the 7th bit of the first octet (the 'X' bit in the EUI-64 format) is set
//   to 1 if the address is globally scoped, or 0 if it is locally scoped
//
// * if the address is 48 bits, the octets 0xfffe are inserted in the middle
//   of the address to pad it to 64 bits
func SetEUI64Addr(ip net.IP, hw net.HardwareAddr, global bool) net.IP {
	tag := []byte{0xff, 0xfe}
	var x uint = 0
	if global == true {
		x = 1
	}

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

	if global {
		hw[0] &= 1 << 1
	} else {
		hw[0] |= 1 << 1
	}

	hw[0] ^= 1 << x

	copy(eui64[8:], hw)

	return eui64
}