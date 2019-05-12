/*
Package iana imports the Internet Assigned Numbers Authority (IANA) IP Special
Registries as a data structure and implements functions to compare the reserved
networks against iplib.Net objects. The IANA registry is used to reserve
portions of IP network space for special use, examples being the IPv4 Private
Use blocks (10.0.0/8, 172.16.0.0/12 and 192.168.0.0/16) and the IPv6 netblock
set aside for documentation purposes (2001:db8::/32).

Note that this package does not contain historical reservations. So IPv6
2001:10::/28 (ORCHIDv1) is listed in the document as "deprecated" and not
present in this library.

The data-set for the IANA registries is available from:
- https://www.iana.org/assignments/iana-ipv4-special-registry/iana-ipv4-special-registry.xhtml
- https://www.iana.org/assignments/iana-ipv6-special-registry/iana-ipv6-special-registry.xhtml
*/
package iana

import (
	"github.com/c-robinson/iplib"
	"net"
	"sort"
)

// Registry holds the aggregated network list from IANA's v4 and v6 registries.
// Only the following fields were imported: Address Block, Name, RFC,
// Forwardable, Globally Reachable and Reserved-by-Protocol
var Registry []*Reservation

// Reservation describes an entry in the IANA IP Special Registry
type Reservation struct {

	// Network is the reserved network
	Network iplib.Net

	// Title is a name given to the reservation
	Title string

	// RFC is the list of relevant RFCs
	RFC []string

	// true if a router may forward packets bound for this network between
	// external interfaces
	Forwardable bool

	// true if a router may pass packets bound for this network outside of
	// a private network
	Global bool

	// true if an IP implementation must implement this policy in order to
	// be compliant
	Reserved bool
}

func init() {
	Registry = []*Reservation{
		{getFromCIDR("0.0.0.0/8"), "This host on this network", []string{"RFC1122"}, false, false, true},
		{getFromCIDR("10.0.0.0/8"), "Private-Use", []string{"RFC1918"}, true, false, false},
		{getFromCIDR("100.64.0.0/10"), "Shared Address Space", []string{"RFC6598"}, false, false, true},
		{getFromCIDR("127.0.0.0/8"), "Loopback", []string{"RFC1122"}, false, false, true},
		{getFromCIDR("169.254.0.0/16"), "Link Local", []string{"RFC3927"}, false, false, true},
		{getFromCIDR("172.16.0.0/12"), "Private-Use", []string{"RFC1918"}, true, false, false},
		{getFromCIDR("192.0.0.0/24"), "IETF Protocol Assignments", []string{"RFC6890"}, false, false, false},
		{getFromCIDR("192.0.0.0/29"), "IPv4 Service Continuity Prefix", []string{"RFC7335"}, true, false, false},
		{getFromCIDR("192.0.0.8/32"), "IPv4 dummy address", []string{"RFC7600"}, false, false, false},
		{getFromCIDR("192.0.0.9/32"), "Port Control Protocol Anycast", []string{"RFC7723"}, true, true, true},
		{getFromCIDR("192.0.0.10/32"), "Traversal Using Relays around NAT Anycast", []string{"RFC8155"}, true, true, false},
		{getFromCIDR("192.0.0.170/32"), "NAT64/DNS64 Discovery", []string{"RFC7050"}, false, false, true},
		{getFromCIDR("192.0.0.171/32"), "NAT64/DNS64 Discovery", []string{"RFC7050"}, false, false, true},
		{getFromCIDR("192.0.2.0/24"), "Documentation (TEST-NET-1)", []string{"RFC5737"}, false, false, false},
		{getFromCIDR("192.31.196.0/24"), "AS112-v4", []string{"RFC7535"}, true, true, false},
		{getFromCIDR("192.52.193.0/24"), "AMT", []string{"RFC7450"}, true, true, false},
		{getFromCIDR("192.168.0.0/16"), "Private-Use", []string{"RFC1918"}, true, false, false},
		{getFromCIDR("192.175.48.0/24"), "Direct Delegation AS112 Service", []string{"RFC7534"}, true, true, false},
		{getFromCIDR("198.18.0.0/15"), "Benchmarking", []string{"RFC2544"}, true, false, false},
		{getFromCIDR("198.51.100.0/24"), "Documentation (TEST-NET-2)", []string{"RFC5737"}, false, false, false},
		{getFromCIDR("203.0.113.0/24"), "Documentation (TEST-NET-3)", []string{"RFC5737"}, false, false, false},
		{getFromCIDR("240.0.0.0/4"), "Reserved", []string{"RFC1112"}, false, false, true},
		{getFromCIDR("255.255.255.255/32"), "Limited Broadcast", []string{"RFC8190", "RFC919"}, false, false, true},
		{getFromCIDR("::1/128"), "Loopback Address", []string{"RFC4291"}, false, false, true},
		{getFromCIDR("::/128"), "Unspecified Address", []string{"RFC4291"}, false, false, true},
		{getFromCIDR("::ffff:0:0/96"), "IPv4-mapped Address", []string{"RFC4291"}, false, false, true},
		{getFromCIDR("64:ff9b::/96"), "IPv4-IPv6 Translation", []string{"RFC6052"}, true, true, false},
		{getFromCIDR("64:ff9b:1::/48"), "IPv4-IPv6 Translation", []string{"RFC8215"}, true, false, false},
		{getFromCIDR("100::/64"), "Discard-Only Address Block", []string{"RFC6666"}, true, false, false},
		{getFromCIDR("2001::/23"), "IETF Protocol Assignments", []string{"RFC2928"}, false, false, false},
		{getFromCIDR("2001::/32"), "TEREDO", []string{"RFC4380", "RFC8190"}, true, true, false},
		{getFromCIDR("2001:1::1/128"), "Port Control Protocol Anycast", []string{"RFC7723"}, true, true, false},
		{getFromCIDR("2001:1::2/128"), "Traversal Using Relays around NAT Anycast", []string{"RFC8155"}, true, true, false},
		{getFromCIDR("2001:2::/48"), "Benchmarking", []string{"RFC5180", "RFC1752"}, true, false, false},
		{getFromCIDR("2001:3::/32"), "AMT", []string{"RFC7450"}, true, true, false},
		{getFromCIDR("2001:4:112::/48"), "AS112-v6", []string{"RFC7535"}, true, true, false},
		{getFromCIDR("2001:5::/32"), "EID Space for LISP (Managed by RIPE NCC)", []string{"RFC7954"}, true, true, true},
		{getFromCIDR("2001:20::/28"), "ORCHIDv2", []string{"RFC7343"}, true, true, false},
		{getFromCIDR("2001:db8::/32"), "Documentation", []string{"RFC3849"}, false, false, false},
		{getFromCIDR("2002::/16"), "6to4", []string{"RFC3056"}, true, true, false},
		{getFromCIDR("2620:4f:8000::/48"), "Direct Delegation AS112 Service", []string{"RFC7534"}, true, true, false},
		{getFromCIDR("fc00::/7"), "Unique-Local", []string{"RFC4193", "RFC8190"}, true, false, false},
		{getFromCIDR("fe80::/10"), "Link-Local Unicast", []string{"RFC4291"}, false, false, true},
	}
}

// GetReservationsForNetwork returns a list of any IANA reserved networks
// that are either part of the supplied network or that the supplied network
// is part of
func GetReservationsForNetwork(n iplib.Net) []*Reservation {
	reservations := []*Reservation{}
	for _, r := range Registry {
		if r.Network.ContainsNet(n) || n.ContainsNet(r.Network) {
			reservations = append(reservations, r)
		}
	}
	return reservations
}

// GetReservationsForIP returns a list of any IANA reserved networks that
// the supplied IP is part of
func GetReservationsForIP(ip net.IP) []*Reservation {
	reservations := []*Reservation{}
	for _, r := range Registry {
		if r.Network.Contains(ip) {
			if iplib.EffectiveVersion(ip) == 4 && r.Title == "IPv4-mapped Address" {
				continue
			}
			reservations = append(reservations, r)
		}
	}
	return reservations
}

// GetRFCsForNetwork returns a list of all RFCs that apply to the given
// network
func GetRFCsForNetwork(n iplib.Net) []string {
	rfclist := []string{}
	reservations := GetReservationsForNetwork(n)
	if len(reservations) > 0 {
		for _, r := range reservations {
		LOOP:
			for _, rfc := range r.RFC {
				for _, xrfc := range rfclist {
					if xrfc == rfc {
						continue LOOP
					}
				}
				rfclist = append(rfclist, rfc)
			}
		}
		sort.Strings(rfclist)
	}
	return rfclist
}

// IsForwardable will return false if the given iplib.Net contains or is
// contained in a network that is marked not-forwardable in the IANA registry
func IsForwardable(n iplib.Net) bool {
	reservations := GetReservationsForNetwork(n)
	for _, r := range reservations {
		if r.Forwardable == false {
			return false
		}
	}
	return true
}

// IsGlobal will return false if the given iplib.Net contains or is contained
// in a network that is marked not-global in the IANA registry
func IsGlobal(n iplib.Net) bool {
	reservations := GetReservationsForNetwork(n)
	for _, r := range reservations {
		if r.Global == false {
			return false
		}
	}
	return true
}

// IsReserved  will return true if the given iplib.Net contains or is
// contained in a network that is marked reserved-by-protocol in the IANA
// registry
func IsReserved(n iplib.Net) bool {
	reservations := GetReservationsForNetwork(n)
	for _, r := range reservations {
		if r.Reserved == true {
			return true
		}
	}
	return false
}

func getFromCIDR(s string) iplib.Net {
	_, n, _ := iplib.ParseCIDR(s)
	return n
}
