package iplib

import (
	"fmt"
	"net"
	"testing"
)

var NewNetTests = []struct {
	ip      net.IP
	masklen int
	out     string
}{
	{
		net.ParseIP("192.168.0.0"), 32, "192.168.0.0/32",
	},
	{
		net.ParseIP("192.168.0.0"), 24, "192.168.0.0/24",
	},
	{
		net.ParseIP("192.168.0.7"), 32, "192.168.0.7/32",
	},
	{
		net.ParseIP("192.168.0.7"), 24, "192.168.0.0/24",
	},
	{
		net.ParseIP("2001:db8::"), 64, "2001:db8::/64",
	},
	{
		net.ParseIP("::ffff:c0a8:0101"), 16, "192.168.0.0/16",
	},
}

func TestNewNet(t *testing.T) {
	for i, tt := range NewNetTests {
		xnet := NewNet(tt.ip, tt.masklen)
		_, pnet, _ := net.ParseCIDR(tt.out)
		if xnet.String() != pnet.String() {
			t.Errorf("[%d] NewNet(%s, %d) expected %s got %s", i, tt.ip.String(), tt.masklen, pnet.String(), xnet.String())
		}
	}
}

var NewNetBetweenTests = []struct {
	start   net.IP
	end     net.IP
	xnet    string
	exact   bool
	err     error
	netslen int
}{
	{ // 0
		net.ParseIP("192.168.0.255"), net.ParseIP("10.0.0.0"),
		"", false, ErrNoValidRange, 0,
	},
	{
		net.ParseIP("192.168.0.255"), net.ParseIP("2001:db8:0:1::"),
		"", false, ErrNoValidRange, 0,
	},
	{
		net.ParseIP("2001:db8:0:1::"), net.ParseIP("192.168.0.255"),
		"", false, ErrNoValidRange, 0,
	},
	{
		net.ParseIP("2001:db8:0:1::"), net.ParseIP("2001:db8::"),
		"", false, ErrNoValidRange, 0,
	},
	{
		net.ParseIP("192.168.0.255"), net.ParseIP("192.168.0.255"),
		"192.168.0.255/32", true, nil, 1,
	},
	{ // 5
		net.ParseIP("2001:db8:0:1::"), net.ParseIP("2001:db8:0:1::"),
		"2001:db8:0:1::/128", true, nil, 1,
	},
	{
		net.ParseIP("192.168.1.0"), net.ParseIP("192.168.2.0"),
		"192.168.1.0/24", false, nil, 2,
	},
	{
		net.ParseIP("2001:db8:1::"), net.ParseIP("2001:db8:2::"),
		"2001:db8:1::/48", false, nil, 2,
	},
	{
		net.ParseIP("192.168.1.0"), net.ParseIP("192.168.1.255"),
		"192.168.1.0/24", true, nil, 1,
	},
	{
		net.ParseIP("2001:db8:1::"), net.ParseIP("2001:db8:1:ffff:ffff:ffff:ffff:ffff"),
		"2001:db8:1::/48", true, nil, 1,
	},
	{ // 10
		net.ParseIP("192.168.1.0"), net.ParseIP("192.168.1.1"),
		"192.168.1.0/31", true, nil, 1,
	},
	{
		net.ParseIP("2001:db8:1::"), net.ParseIP("2001:db8:1::1"),
		"2001:db8:1::/127", true, nil, 1,
	},
	{
		net.ParseIP("192.168.0.255"), net.ParseIP("192.168.1.2"),
		"192.168.0.255/32", false, nil, 3,
	},
	{
		net.ParseIP("2001:db8:0:ffff:ffff:ffff:ffff:ffff"), net.ParseIP("2001:db8:1::1"),
		"2001:db8:0:ffff:ffff:ffff:ffff:ffff/128", false, nil, 2,
	},
	{
		net.ParseIP("10.0.0.0"), net.ParseIP("255.0.0.0"),
		"10.0.0.0/7", false, nil, 13,
	},
	{ // 15
		net.ParseIP("2001:db8::"), net.ParseIP("2001:db8:ffff:ffff:ffff:ffff:ffff::"),
		"2001:db8::/33", false, nil, 81,
	},
}

func TestNewNetBetween(t *testing.T) {
	for i, tt := range NewNetBetweenTests {
		xnet, exact, err := NewNetBetween(tt.start, tt.end)
		if e := compareErrors(err, tt.err); len(e) > 0 {
			t.Errorf("[%d] NewNetBetween(%s, %s) expected error '%v', got '%v'", i, tt.start, tt.end, tt.err, err)
		} else if tt.err == nil {
			if xnet == nil && len(tt.xnet) != 0 {
				t.Fatalf("[%d] should not be nil!", i)
			}
			if xnet.String() != tt.xnet {
				t.Errorf("[%d] NewNetBetween(%s, %s) expected '%s', got '%s'", i, tt.start, tt.end, tt.xnet, xnet.String())
			}
			if exact != tt.exact {
				t.Errorf("[%d] NewNetBetween(%s, %s) expected '%t', got '%t'", i, tt.start, tt.end, tt.exact, exact)
			}
		}
	}
}

func TestAllNetsBetween(t *testing.T) {
	for i, tt := range NewNetBetweenTests {
		//t.Logf("[%d] nets between %s and %s", i, tt.start, tt.end)
		xnets, err := AllNetsBetween(tt.start, tt.end)
		//t.Logf("[%d] got %+v, '%v'", i, xnets, err)
		if e := compareErrors(err, tt.err); len(e) > 0 {
			t.Errorf("[%d] expected error '%v', got '%v'", i, tt.err, err)
		}
		if tt.err == nil {
			if len(xnets) != tt.netslen {
				t.Logf("[%d] AllNetsBetween(%s, %s) [%+v]", i, tt.start, tt.end, xnets)
				t.Errorf("[%d] expected %d networks, got %d", i, tt.netslen, len(xnets))
			}
		}
	}
}

var ParseCIDRTests = []struct {
	s    string
	xnet string
	err  error
	ver  int
}{
	{"not.legit/22", "", fmt.Errorf("invalid CIDR address: not.legit/22"), 0},
	{"192.168.1.1", "", fmt.Errorf("invalid CIDR address: 192.168.1.1"), 0},
	{"192.168.1.0/24", "192.168.1.0/24", nil, 4},
	{"2001:db8::/64", "2001:db8::/64", nil, 6},
	{"::ffff:c0a8:0101/32", "192.168.1.1/32", nil, 4},
	{"::ffff:c0a9:0101/16", "192.169.0.0/16", nil, 4},
	{"::ffff:c0a8:0101/64", "::/64", nil, 6},
}

func TestParseCIDR(t *testing.T) {
	for i, tt := range ParseCIDRTests {
		_, n, err := ParseCIDR(tt.s)
		if e := compareErrors(err, tt.err); len(e) > 0 {
			t.Errorf("[%d] ParseCIDR(%s) expected error '%v', got '%v'", i, tt.s, tt.err, err)
		} else if tt.err == nil {
			if n.Version() != tt.ver {
				t.Errorf("[%d] expected IPNet version '%d' got '%d'", i, tt.ver, n.Version())
			}
			if n.String() != tt.xnet {
				fmt.Println(n)
				t.Errorf("[%d] expected '%s' for '%s'", i, tt.xnet, n.String())
			}
		}
	}
}
