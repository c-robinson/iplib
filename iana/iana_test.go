package iana

import (
	"github.com/c-robinson/iplib"
	"net"
	"testing"
)

var IPTests = []struct {
	name     string
	address  string
	resCount int
}{
	{
		"NotReservedv4",
		"144.21.1.19",
		0,
	},
	{
		"Reservedv4",
		"192.168.123.49",
		1,
	},
	{
		"NotReservedv6",
		"25:100:200::195:16",
		0,
	},
	{
		"Reservedv6",
		"2001:db8:1::250:3",
		1,
	},
}

func TestGetReservationsForIP(t *testing.T) {
	for _, tt := range IPTests {
		ip := net.ParseIP(tt.address)
		r := GetReservationsForIP(ip)
		if len(r) != tt.resCount {
			t.Errorf("'%s' want %d reservations, got %d", tt.name, tt.resCount, len(r))
			for _, v := range r {
				t.Logf("%s, %s", v.Title, v.Network.String())
			}
		}
	}
}

var NetTests = []struct {
	name           string
	resCount       int
	network        string
	rfcList        []string
	valForwardable bool
	valGlobal      bool
	valReserved    bool
}{
	{
		"NotReservedv4",
		0,
		"1.0.0.0/8",
		[]string{},
		true,
		true,
		false,
	},
	{
		"Reservedv4",
		1,
		"10.0.0.0/8",
		[]string{"RFC1918"},
		true,
		false,
		false,
	},
	{
		"ContainedInReservedv4",
		1,
		"10.44.1.0/24",
		[]string{"RFC1918"},
		true,
		false,
		false,
	},
	{
		"ContainsReservedv4",
		1,
		"192.168.0.0/16",
		[]string{"RFC1918"},
		true,
		false,
		false,
	},
	{
		"MultipleReservationsv4",
		8,
		"192.0.0.0/12",
		[]string{"RFC5737", "RFC6890", "RFC7050", "RFC7335", "RFC7600", "RFC7723", "RFC8155"},
		false,
		false,
		true,
	},

	{
		"NotReservedv6",
		0,
		"2001:db7::/32",
		[]string{},
		true,
		true,
		false,
	},
	{
		"Reservedv6",
		1,
		"2001:db8::/32",
		[]string{"RFC3849"},
		false,
		false,
		false,
	},
	{
		"ContainedInReservedv6",
		1,
		"2001:db8:100::/40",
		[]string{"RFC3849"},
		false,
		false,
		false,
	},
	{
		"ContainsReservedv6",
		1,
		"2001:d0::/28",
		[]string{"RFC2928"},
		false,
		false,
		false,
	},
	{
		"MultipleReservationsv6",
		10,
		"2001::/16",
		[]string{"RFC1752", "RFC2928", "RFC3849", "RFC4380", "RFC5180", "RFC7343", "RFC7450", "RFC7535", "RFC7723", "RFC7954", "RFC8155", "RFC8190"},
		false,
		false,
		true,
	},
	{
		"KnowsAbout4in6",
		1,
		"::ffff:c0a9:0101/16",
		[]string{"RFC4291"},
		false,
		false,
		true,
	},
}

func TestGetReservationsForNetwork(t *testing.T) {
	for _, tt := range NetTests {
		_, n, _ := iplib.ParseCIDR(tt.network)
		r := GetReservationsForNetwork(n)
		if len(r) != tt.resCount {
			t.Errorf("'%s' want %d reservations, got %d", tt.name, tt.resCount, len(r))
			for _, v := range r {
				t.Logf("%s, %s", v.Title, v.Network.String())
			}
		}
	}
}

func TestGetRFCsForNetwork(t *testing.T) {
	for _, tt := range NetTests {
		_, n, _ := iplib.ParseCIDR(tt.network)
		rfclist := GetRFCsForNetwork(n)
		if v := equalList(rfclist, tt.rfcList); v != true {
			t.Errorf("'%s' (%s) want %v, got %v", tt.name, tt.network, tt.rfcList, rfclist)
		}
	}
}

func TestIsForwardable(t *testing.T) {
	for _, tt := range NetTests {
		_, n, _ := iplib.ParseCIDR(tt.network)
		if tt.valForwardable != IsForwardable(n) {
			t.Errorf("'%s' (%s) want %t, got %t", tt.name, tt.network, tt.valForwardable, IsForwardable(n))
		}
	}
}

func TestIsGlobal(t *testing.T) {
	for _, tt := range NetTests {
		_, n, _ := iplib.ParseCIDR(tt.network)
		if tt.valGlobal != IsGlobal(n) {
			t.Errorf("'%s' (%s) want %t, got %t", tt.name, tt.network, tt.valGlobal, IsGlobal(n))
		}
	}
}

func TestIsReserved(t *testing.T) {
	for _, tt := range NetTests {
		_, n, _ := iplib.ParseCIDR(tt.network)
		if tt.valReserved != IsReserved(n) {
			t.Errorf("'%s' (%s) want %t, got %t", tt.name, tt.network, tt.valReserved, IsReserved(n))
		}
	}
}

func equalList(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
