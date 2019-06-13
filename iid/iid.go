package iid

import (
	"net"

	"github.com/c-robinson/iplib"
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
