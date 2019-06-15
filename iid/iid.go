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
	"net"

	"github.com/c-robinson/iplib"
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

// Registry holds the aggregated network list from IANA's "Reserved IPv6
// Interface Identifiers" as specified in RFC5453. In order to be compliant
// with RFC 7217's algorithm for "Semantically Opaque Interface Identifiers"
// addresses should be checked against this registry to make sure there are
// no conflicts
var Registry []*Reservation

// Reservation describes an entry in the IANA IP Special Registry
type Reservation struct {
	// FirstAddr is the first address in the reservation
	FirstAddr net.IP

	// LastAddr is the last address  in the reservation
	LastAddr net.IP

	// Title is a name given to the reservation
	Title string

	// RFC is the list of relevant RFCs
	RFC string
}

func init() {
	Registry = []*Reservation{
		{
			getFromString("::"),
			getFromString("0000:0000:0000:0000:FFFF:FFFF:FFFF:FFFF"),
			"Subnet-Router Anycast",
			"RFC4291",
		},
		{
			getFromString("0200:5EFF:FE00:0000::"),
			getFromString("0200:5EFF:FE00:5212:FFFF:FFFF:FFFF:FFFF"),
			"Reserved IPv6 Interface Identifiers corresponding to the IANA Ethernet Block",
			"RFC4291",
		},
		{
			getFromString("0200:5EFF:FE00:5213::"),
			getFromString("0200:5EFF:FE00:5213:FFFF:FFFF:FFFF:FFFF"),
			"Proxy Mobile IPv6",
			"RFC6543",
		},
		{
			getFromString("0200:5EFF:FE00:5214::"),
			getFromString("0200:5EFF:FEFF:FFFF:FFFF:FFFF:FFFF:FFFF"),
			"Reserved IPv6 Interface Identifiers corresponding to the IANA Ethernet Block",
			"RFC4291",
		},
		{
			getFromString("FDFF:FFFF:FFFF:FF80::"),
			getFromString("FDFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF"),
			"Reserved Subnet Anycast Addresses",
			"RFC2526",
		},
	}
}

// GetReservationsForIP returns a list of any IANA reserved networks that
// the supplied IP is part of
func GetReservationsForIP(ip net.IP) []*Reservation {
	reservations := []*Reservation{}
	for _, r := range Registry {
		f := iplib.CompareIPs(r.FirstAddr, ip)
		l := iplib.CompareIPs(r.LastAddr, ip)

		if f >= 0 && l <= 0 {
			reservations = append(reservations, r)
		}
	}
	return reservations
}

func getFromString(s string) net.IP {
	return net.ParseIP(s)
}
