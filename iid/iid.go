/*
Package iid provides functions for generating and validating IPv6 Interface
Identifiers (IID's). For the purposes of this module an IID is an IPv6 address
constructed, somehow, from information which uniquely identifies a given
interface on a network, and is unique within that network.

As part of validation this package imports the Internet Assigned Numbers
Authority (IANA) Reserved IPv6 Interface Identifiers as a data structure and
implements functions to compare the reserved networks against IID's to avoid
conflicts. The data set for the IANA registry is available from:

- https://www.iana.org/assignments/ipv6-interface-ids/ipv6-interface-ids.xhtml
*/
package iid

import (
	"bytes"
	"crypto"
	"encoding/binary"
	"errors"
	"net"
)

// Scope describes the availability of an IPv6 IID
type Scope int

const (
	// ScopeNone is an undefined IPv6 IID scope
	ScopeNone   Scope = iota

	// ScopeInvert will invert the scope of an IID
	ScopeInvert

	// ScopeGlobal is a global IPv6 IID scope
	ScopeGlobal

	// ScopeLocal is a local IPv6 IID scope
	ScopeLocal
)

var (
	ErrIIDAddressCollision = errors.New("proposed IID collides with IANA reserved IID list")
	ErrInsufficientHashLength = errors.New("hash function must return a digest of 64bits or more")
)

// Registry holds the aggregated network list from IANA's "Reserved IPv6
// Interface Identifiers" as specified in RFC5453. In order to be compliant
// with RFC 7217's algorithm for "Semantically Opaque Interface Identifiers"
// addresses should be checked against this registry to make sure there are
// no conflicts
var Registry []*Reservation

// Reservation describes an entry in the IANA IP Special Registry
type Reservation struct {
	// FirstAddr is the first address in the reservation
	FirstAddr []byte

	// LastAddr is the last address  in the reservation
	LastAddr []byte

	// Title is a name given to the reservation
	Title string

	// RFC is the list of relevant RFCs
	RFC string
}

func init() {
	Registry = []*Reservation{
		{
			[]byte{  0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00 },
			[]byte{  0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00 },
			"Subnet-Router Anycast",
			"RFC4291",
		},
		{
			[]byte{ 0x02, 0x00, 0x5e, 0xff, 0xfe, 0x00, 0x00, 0x00 },
			[]byte{ 0x02, 0x00, 0x5e, 0xff, 0xfe, 0x00, 0x52, 0x12 },
			"Reserved IPv6 Interface Identifiers corresponding to the IANA Ethernet Block",
			"RFC4291",
		},
		{
			[]byte{ 0x02, 0x00, 0x5e, 0xff, 0xfe, 0x00, 0x52, 0x13 },
			[]byte{ 0x02, 0x00, 0x5e, 0xff, 0xfe, 0x00, 0x52, 0x13 },
			"Proxy Mobile IPv6",
			"RFC6543",
		},
		{
			[]byte{ 0x02, 0x00, 0x5e, 0xff, 0xfe, 0x00, 0x52, 0x14 },
			[]byte{ 0x02, 0x00, 0x5e, 0xff, 0xfe, 0xff, 0xff, 0xff },
			"Reserved IPv6 Interface Identifiers corresponding to the IANA Ethernet Block",
			"RFC4291",
		},
		{
			[]byte{ 0xfd, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x80 },
			[]byte{ 0xfd, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff },
			"Reserved Subnet Anycast Addresses",
			"RFC2526",
		},
	}
}

// GenerateRFC7217Addr generates a pseudo-random IID from supplied input
// parameters, in compliance with RFC7217. The signature of this function
// deviates from the one specified in that RFC only insomuch as is necessary
// to conform to the implementing language. The input fields are:
//
// ip      - v6 net.IP, only the first 64bits will be used
// hw      - 48- or 64-bit net.HardwareAddr, typically of the interface that
//           this address will be assigned to
// counter - a monotonically incrementing number read from some non-volatile
//           local storage. If the function returns ErrIIDAddressCollision
//           the counter should be incremented and the function called again
// netid   - some piece of information identifying the local subnet, such as
//           an 802.11 SSID. RFC6059 lists other interesting options. This
//           field may be left blank ([]byte{})
// secret  - a local, closely held, secret key. This is the sauce that makes
//           the address opaque
// htype   - a crypto.Hash function to use when generating the IID. Note that
//           MD5 is specifically prohibited for being too easily guessable
// scope   - the scope of the IID
func GenerateRFC7217Addr(ip net.IP, hw net.HardwareAddr, counter int64, netid, secret []byte, htype crypto.Hash, scope Scope) (net.IP, error) {
	var bs []byte
	binary.LittleEndian.PutUint64(bs, uint64(counter))

	bs = append(hw, bs...)
	bs = append(bs, netid...)
	bs = append(bs, secret...)

	f := htype.New()
	if f.Size() < 8 {
		return nil, ErrInsufficientHashLength
	}
	iid := make([]byte, 16)
	copy(iid, ip)

	rid := f.Sum(bs)
	rid = setScopeBit(rid, scope)

	copy(iid[8:], rid[0:8])

	if r := GetReservationsForIP(iid); len(r) > 0 {
		return nil, ErrIIDAddressCollision
	}

	return iid, nil
}

// GetReservationsForIP returns a list of any IANA reserved networks that
// the supplied IP is part of
func GetReservationsForIP(ip net.IP) []*Reservation {
	reservations := []*Reservation{}
	for _, r := range Registry {
		f := bytes.Compare(r.FirstAddr, ip[8:])
		l := bytes.Compare(r.LastAddr, ip[8:])

		if f >= 0 && l <= 0 {
			reservations = append(reservations, r)
		}
	}
	return reservations
}

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
// interpretations via the Scope parameter but recommends passing an explicit
// ScopeGlobal or ScopeLocal
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
	return setScopeBit(eui64, scope)
}


// MakeOpaqueAddr offers one implemention of RFC7217s algorithm for generating
// a "semantically opaque interface identifier". The caller must supply a
// counter and secret and MAY supply an additional "netid". Ultimately this
// function calls GenerateRFC7217Addr() with scope set to "global" and an
// htype of SHA256, but please see the documentation in that function for an
// explanation of all the input fields
func MakeOpaqueAddr(ip net.IP, hw net.HardwareAddr, counter int64, netid, secret []byte) (net.IP, error) {
	return GenerateRFC7217Addr(ip, hw, counter, netid, secret, crypto.SHA256, ScopeGlobal)
}

func setScopeBit(ip net.IP, scope Scope) net.IP {
	switch scope {
	case ScopeGlobal:
		ip[8] |= 1 << 1  // set 0 or 1 -> 1

	case ScopeLocal:
		ip[8] &^= 1 << 1 // set 0 or 1 -> 0

	case ScopeInvert:
		ip[8] ^= 1 << 1  // set 0 -> 1 or 1 -> 0
	default:
	}

	return ip
}