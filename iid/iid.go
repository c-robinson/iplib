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
	_ "crypto/sha256"
	"encoding/binary"
	"errors"
	"net"

	"github.com/c-robinson/iplib"
)

// Scope describes the availability of an IPv6 IID and determines how IID-
// generating functions treat the 7th bit in the 9th octet of the address
// (the 'X' bit in the EUI-64 format, or the 'u' bit in RFC4291)
//
// NOTE: there is some ambiguity to the RFC here. Most discussions I've seen
// on the topic say that the bit should _always_ be inverted, but the RFC
// reads like the IPv6 EUI64 format uses the _inverse signal_ from the IEEE
// EUI64 format; so where the IEEE uses 0 for global scoping, the IPv6 IID
// should use 1. This module punts on the question and provides for all
// interpretations via the scope parameter but recommends passing an explicit
// ScopeGlobal or ScopeLocal
type Scope int

const (
	// ScopeNone is an undefined IID scope, the X bit will not be modified
	ScopeNone   Scope = iota

	// ScopeInvert will cause the X bit to be inverted, setting 0 to 1 and 1
	// to 0. This behavior is widely interpreted as the correct behavior
	ScopeInvert

	// ScopeGlobal will cause the X bit to be set to 1, indicating that the
	// IID should be globally scoped
	ScopeGlobal

	// ScopeLocal will cause the X bit to be set to 0, indicating that the IID
	// should only be locally scoped
	ScopeLocal
)

var (
	ErrIIDAddressCollision = errors.New("proposed IID collides with IANA reserved IID list")
)

// Registry holds the aggregated network list from IANA's "Reserved IPv6
// Interface Identifiers" as specified in RFC5453. In order to be compliant
// with RFC 7217's algorithm for "Semantically Opaque Interface Identifiers"
// addresses should be checked against this registry to make sure there are
// no conflicts
var Registry []*Reservation

// Reservation describes an entry in the IANA IP Special Registry
type Reservation struct {
	// FirstRes is the first address in the reservation
	FirstRes []byte

	// LastRes is the last address in the reservation
	LastRes []byte

	// Title is a name given to the reservation
	Title string

	// RFC is the list of relevant RFCs
	RFC string
}

func init() {
	Registry = []*Reservation{
		{
			[]byte{ 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00 },
			[]byte{ 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00 },
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
			"Reserved IPv6 Interface Identifiers corresponding to the IANA Ethernet Block (2)",
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
// ip: v6 net.IP, only the first 64bits will be used
//
// hw: 48- or 64-bit net.HardwareAddr, typically of the interface that
// this address will be assigned to
//
// counter: a monotonically incrementing number read from some non-volatile
// local storage. This variable provides the velocity to the entire algorithm
// and should be incremented after each use. There is no guarantee that a
// generated address wont accidentally fall within the range of reserved IPv6
// IIDs and, should this happen, an ErrIIDAddressCollision will be returned.
// This is harmless and if it happens counter should be incremented and the
// function called again
//
// netid: some piece of information identifying the local subnet, such as an
// 802.11 SSID. RFC6059 lists other interesting options. This field may be
// left blank ([]byte{})
//
// secret: a local, closely held, secret key. This is the sauce that makes the
// address opaque
//
// htype: a crypto.Hash function to use when generating the IID.
//
// scope: the scope of the IID
//
// NOTE that MD5 is specifically prohibited for being too easily guessable.
//
// NOTE that unless you use sha256 you will need to import the hash function
// you intend to use, (e.g. import _ "crypto/sha512")
func GenerateRFC7217Addr(ip net.IP, hw net.HardwareAddr, counter int64, netid, secret []byte, htype crypto.Hash, scope Scope) (net.IP, error) {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(counter))

	bs = append(hw, bs...)
	bs = append(bs, netid...)
	bs = append(bs, secret...)

	f := htype.New()

	iid := make([]byte, 16)
	copy(iid, ip)

	f.Write(bs)
	rid := f.Sum(nil)
	rid = setScopeBit(rid, scope)

	copy(iid[8:], rid[0:8])

	if r := GetReservationsForIP(iid); r != nil {
		return nil, ErrIIDAddressCollision
	}

	return iid, nil
}

// GetReservationsForIP returns a list of any IANA reserved networks that
// the supplied IP is part of
func GetReservationsForIP(ip net.IP) *Reservation {
	if iplib.EffectiveVersion(ip) != 6 {
		return nil
	}
	for _, r := range Registry {
		f := bytes.Compare(ip[8:], r.FirstRes)
		l := bytes.Compare(ip[8:], r.LastRes)

		if f >= 0 && l <= 0 {
			return r
		}
	}
	return nil
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
// modified per the definition for each constant.
//
// * if the address is 48 bits, the octets 0xFFFE are inserted in the middle
// of the address to pad it to 64 bits
func MakeEUI64Addr(ip net.IP, hw net.HardwareAddr, scope Scope) net.IP {
	tag := []byte{0xff, 0xfe}

	if iplib.EffectiveVersion(ip) != 6 {
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


// MakeOpaqueAddr offers one implementation of RFC7217's algorithm for
// generating a "semantically opaque interface identifier". The caller must
// supply a counter and secret and MAY supply an additional "netid".
// Ultimately this function calls GenerateRFC7217Addr() with scope set to
// "global" and an htype of SHA256, but please see the documentation in that
// function for an explanation of all the input fields
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