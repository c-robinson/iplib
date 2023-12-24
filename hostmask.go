package iplib

import (
	"encoding/hex"
	"net"

	"lukechampine.com/uint128"
)

// HostMask is a mask that can be applied to IPv6 addresses to mask out bits
// from the right side of the address instead of the left (which is the
// purview of a netmask), the intended use-case is for situations where there
// is a desire to reserve a portion of the address for some other purpose and
// only allow iplib to manage the remainder. A concrete example would be
// IPv6 Interface Identifiers as described in RFC4291, RFC4941 or RFC7217 in
// which the final 64bits of the address are used to construct a unique host
// identifier and the allocator only has control of the first 64bits. So the
// next IP from 2001:db8:1234:5678:: would be 2001:db8:1234:5679 instead of
// 2001:db8:1234:5678::1. Here is a Net6 object eing initialized without a
// hostmask:
//
//	n := NewNet6(2001:db8::, 56, 0)
//	Address            2001:db8::
//	Netmask            ffff:ffff:ffff:ff00:0000:0000:0000:0000
//	Hostmask           0000:0000:0000:0000:0000:0000:0000:0000
//	First              2001:0db8:0000:0000:0000:0000:0000:0000
//	Last               2001:0db8:0000:00ff:ffff:ffff:ffff:ffff
//	Count              4722366482869645213696
//
// This creates a block with 4.7 sextillion usable addresses. Below is he same
// block with a hostmask of /60. The mask is applied from the rightmost byte,
// leaving 12 unmasked bits for a total of 4096 allocatable addresses:
//
//	n:= NewNet6(2001:db8::, 56, 60)
//	Address            2001:db8::
//	Netmask            ffff:ffff:ffff:ff00:0000:0000:0000:0000
//	Hostmask           0000:0000:0000:0000:0fff:ffff:ffff:ffff
//	First              2001:0db8:0000:0000:0000:0000:0000:0000
//	Last               2001:0db8:0000:00ff:f000:0000:0000:0000
//	Count              4096
//
// In the first example the second IP address of the netblock is 2001:db8::1,
// in the second example it is 2001:db8:0:1::
//
// One important note: even though bytes are filled in from the right the bits
// within those bytes are still blocked out left-to-right, so that address
// incrementing/decrementing makes sense to the end user, as shown here:
//
//	BINARY      Base16  Base10  Example Max16  Max10
//	0000 0000     0x00       0      /56  0xFF    255
//	1000 0000     0x80     128      /57  0x7F    127
//	1100 0000     0xC0     192      /58  0x3F     63
//	1110 0000     0xE0     224      /59  0x1F     31
//	1111 0000     0xF0     240      /60  0x0F     15
//	1111 1000     0xF8     248      /61  0x07      7
//	1111 1100     0xFC     252      /62  0x03      3
//	1111 1110     0xFE     254      /63  0x01      1
//
// A hostmask of /1 will block out the left-most bit of the 16th byte
// while a /8 will block the entire 16th byte.
//
// To initialize a hostmask you must give it an integer value between 1 and
// 128, which represent the number of bits in the mask.
type HostMask []byte

// NewHostMask returns a HostMask initialized to masklen
func NewHostMask(masklen int) HostMask {
	mask := make([]byte, 16)
	if masklen == 0 {
		return mask
	}
	for i := 15; i >= 0; i-- {
		if masklen < 8 {
			mask[i] = ^byte(0xff >> uint(masklen))
			break
		}
		mask[i] = 0xff
		masklen -= 8
	}
	return mask
}

// BoundaryByte returns the rightmost byte in the mask in which any bits fall
// inside the hostmask, as well as the position of that byte. For example a
// masklength of 58 would return "0xc0, 8" while 32 would return "0xff, 12".
// If the hostmask is unset "0x00, -1" will be returned
func (m HostMask) BoundaryByte() (byte, int) {
	hmlen, _ := m.Size() // will be between 0 and 128, where 0 is "no mask"

	if hmlen == 0 {
		return 0x00, -1
	}

	quo, mod := hmlen/8, hmlen%8
	if mod == 0 {
		quo--
	}
	pos := 15 - quo

	return m[pos], pos
}

// Size returns the number of ones and total bits in the mask
func (m HostMask) Size() (int, int) {
	ones := 0
	bits := 128
	for i := len(m) - 1; i >= 0; i-- {
		b := m[i]
		if b == 0xff {
			ones += 8
			continue
		}
		for b&0x80 != 0 {
			ones++
			b <<= 1
		}
		break
	}
	return ones, bits
}

// String returns the hexadecimal form of m, with no punctuation
func (m HostMask) String() string {
	return hex.EncodeToString(m)
}

// DecrementIP6WithinHostmask returns a net.IP that is less than the unmasked
// portion of the supplied net.IP by the supplied integer value. If the
// input or output value fall outside the boundaries of the hostmask a
// ErrAddressOutOfRange will be returned
func DecrementIP6WithinHostmask(ip net.IP, hm HostMask, count uint128.Uint128) (net.IP, error) {
	bb, bbpos := hm.BoundaryByte()
	if bbpos == 0 {
		return net.IP{}, ErrBadMaskLength
	}

	if bbpos == -1 {
		return DecrementIP6By(ip, count), nil
	}

	// check if ip is outside of hostmask already
	if bbpos < 15 {
		for _, b := range ip[bbpos+1:] {
			if b > 0 {
				return net.IP{}, ErrAddressOutOfRange
			}
		}
		if ip[bbpos]+bb < bb {
			return net.IP{}, ErrAddressOutOfRange
		}
	}

	count, bb = decrementBoundaryByte(bb, ip[bbpos], count)
	xip := decrementUnmaskedBytes(ip[:bbpos], count)
	if len(xip) == 0 {
		return xip, ErrAddressOutOfRange
	}

	if len(xip) > bbpos {
		return net.IP{}, ErrAddressOutOfRange
	}

	xip = append(xip, bb)

	xip = append(xip, make([]byte, 15-bbpos)...)

	return xip, nil
}

// IncrementIP6WithinHostmask returns a net.IP that is greater than the
// unmasked portion of the supplied net.IP by the supplied integer value. If
// the input or output value fall outside the boundaries of the hostmask a
// ErrAddressOutOfRange will be returned
func IncrementIP6WithinHostmask(ip net.IP, hm HostMask, count uint128.Uint128) (net.IP, error) {
	bb, bbpos := hm.BoundaryByte()
	if bbpos == 0 {
		return net.IP{}, ErrBadMaskLength
	}

	if bbpos == -1 {
		return IncrementIP6By(ip, count), nil
	}

	// check if ip is outside of hostmask already
	for _, b := range ip[bbpos+1:] {
		if b > 0 {
			return net.IP{}, ErrAddressOutOfRange
		}
	}

	count, bb = incrementBoundaryByte(bb, ip[bbpos], count)
	xip := incrementUnmaskedBytes(ip[:bbpos], count)

	if len(xip) > bbpos {
		return net.IP{}, ErrAddressOutOfRange
	}

	xip = append(xip, bb)

	xip = append(xip, make([]byte, 15-bbpos)...)
	if CompareIPs(xip, ip) < 0 {
		return net.IP{}, ErrAddressOutOfRange
	}
	return xip, nil
}

// NextIP6WithinHostmask takes a net.IP and Hostmask as arguments and attempts
// to increment the IP by one, within the boundary of the hostmask. If bits
// inside the hostmask are set, an empty net.IP{} and an ErrAddressOutOfRange
// will be returned
func NextIP6WithinHostmask(ip net.IP, hm HostMask) (net.IP, error) {
	xip := CopyIP(ip)

	for i := len(xip) - 1; i >= 0; i-- {
		if hm[i] == 0xff {
			if xip[i] > 0 {
				return net.IP{}, ErrAddressOutOfRange
			}
			continue
		}
		if (xip[i] | hm[i]) == 0xff {
			// xip[i] is the boundary byte, and the accessible bits are at max
			xip[i] = 0
			continue
		}
		xip[i]++

		if xip[i] > 0 {
			return xip, nil
		}
	}
	return net.IP{}, ErrAddressOutOfRange
}

// PreviousIP6WithinHostmask takes a net.IP and Hostmask as arguments and
// attempts to decrement the IP by one, within the boundary of the hostmask.
// If bits inside the hostmask are set, an empty net.IP{} and an
// ErrAddressOutOfRange will be returned
func PreviousIP6WithinHostmask(ip net.IP, hm HostMask) (net.IP, error) {
	xip := CopyIP(ip)
	bb, bbpos := hm.BoundaryByte()
	bbmax := 0xff - bb

	for i := len(xip) - 1; i >= 0; i-- {
		if hm[i] == 0xff {
			if xip[i] > 0 {
				return net.IP{}, ErrAddressOutOfRange
			}
			continue
		}

		xip[i]--

		if xip[i] != 255 {
			if i == bbpos-1 {
				// if we underflowed into the boundary byte we need to adjust
				// it to it's actual max, not 0xff
				xip[bbpos] = bbmax
			}
			return xip, nil
		}
	}
	return net.IP{}, ErrAddressOutOfRange
}

// decrementBoundaryByte takes a boundary-byte, a boundary-value and a count
// as input and returns a modified boundary byte and count for further
// processing. bb is used to calculate the maximum value for bv and then the
// count + bv is divided by that max. The function returns a new count and
// boundary-byte
func decrementBoundaryByte(bb, bv byte, count uint128.Uint128) (uint128.Uint128, byte) {
	if count.IsZero() {
		return count, bv
	}

	byteMax := uint128.From64(256 - uint64(bb)) // max value of unmasked bits in the byte
	byteVal := uint128.From64(uint64(bv))       // cur value of unmasked bits in the byte

	mod := uint128.New(0, 0)

	count, mod = count.QuoRem(byteMax)

	// extract the actual modulus into bmod
	rb := make([]byte, 16)
	mod.PutBytesBE(rb)
	bmod := rb[15]

	if bmod > bv {
		count = count.Add64(1)

		byteVal = byteVal.Add(byteMax)
		byteVal = byteVal.Sub(mod)

		// convert to byte
		byteVal.PutBytesBE(rb)

		return count, rb[15]
	}
	return count, bv - bmod
}

// decrementUnmaskedBytes decrements an arbitrary []byte by count and returns
// a []byte of the same length
func decrementUnmaskedBytes(nb []byte, count uint128.Uint128) []byte {
	if count.IsZero() {
		return nb

	}

	// convert the []byte to a uint128, which requires a [16]byte
	pnb := append(make([]byte, 16-len(nb)), nb...)
	n := uint128.FromBytesBE(pnb)

	if count.Cmp(n) > 0 {
		return []byte{}
	}

	n = n.Sub(count)

	// convert the uint128 back to a []byte
	xb := make([]byte, 16)
	n.PutBytesBE(xb)

	// return only as many elements as were passed in
	return xb[16-len(nb):]
}

// incrementBoundaryByte takes a boundary-byte, a boundary-value and a count
// as input and returns a modified boundary byte and count for further
// processing. bb is used to calculate the maximum value for bv and then the
// count + bv is divided by that max. The function will return the modulus
// as a byte value, and the pointer to count will have the quotient
func incrementBoundaryByte(bb, bv byte, count uint128.Uint128) (uint128.Uint128, byte) {
	if count.IsZero() {
		return count, bv
	}
	byteMax := uint128.From64(256 - uint64(bb)) // max value of unmasked bits in the byte
	byteVal := uint128.From64(uint64(bv))       // cur value of unmasked bits in the byte

	count = count.Add(byteVal)
	if count.Cmp(byteMax) < 0 {
		return uint128.Uint128{}, byte(count.Lo)
	}

	mod := uint128.New(0, 0)
	count, mod = count.QuoRem(byteMax)
	rb := make([]byte, 16)
	mod.PutBytesBE(rb)
	return count, rb[15]
}

// incrementUnmaskedBytes increments an arbitrary []byte by count and returns
// a []byte of the same length
func incrementUnmaskedBytes(nb []byte, count uint128.Uint128) []byte {
	if count.IsZero() {
		return nb
	}

	// convert the []byte to a uint128, which requires a [16]byte
	pnb := append(make([]byte, 16-len(nb)), nb...)
	n := uint128.FromBytesBE(pnb)

	n = n.Add(count)

	// convert the uint128 back to a []byte
	xb := make([]byte, 16)
	n.PutBytesBE(xb)

	// return only as many elements as were passed in
	return xb[16-len(nb):]
}
