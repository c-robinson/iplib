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
	start net.IP
	end   net.IP
	xnet  string
	exact bool
	err   error
}{
	{
		net.ParseIP("192.168.0.255"), net.ParseIP("192.168.2.0"),
		"192.168.1.0/24", false, nil,
	},
	{
		net.ParseIP("192.168.0.255"), net.ParseIP("10.0.0.0"),
		"", false, ErrNoValidRange,
	},
	{
		net.ParseIP("192.168.0.255"), net.ParseIP("2001:db8:0:1::"),
		"", false, ErrNoValidRange,
	},
	{
		net.ParseIP("2001:db8:0:1::"), net.ParseIP("192.168.0.255"),
		"", false, ErrNoValidRange,
	},
	{
		net.ParseIP("192.168.0.255"), net.ParseIP("192.168.0.255"),
		"", false, ErrNoValidRange,
	},
	{
		net.ParseIP("192.168.0.255"), net.ParseIP("192.168.1.1"),
		"192.168.1.0/32", true, nil,
	},
	{
		net.ParseIP("192.168.1.0"), net.ParseIP("192.168.1.2"),
		"192.168.1.1/32", true, nil,
	},
	{
		net.ParseIP("192.168.0.255"), net.ParseIP("192.168.1.2"),
		"192.168.1.0/31", true, nil,
	},
	{
		net.ParseIP("192.168.0.255"), net.ParseIP("192.168.1.3"),
		"192.168.1.0/31", false, nil,
	},
	{
		net.ParseIP("192.168.1.0"), net.ParseIP("192.168.1.3"),
		"192.168.1.0/30", true, nil,
	},
	{
		net.ParseIP("192.168.0.255"), net.ParseIP("192.168.1.4"),
		"192.168.1.0/30", false, nil,
	},
	{
		net.ParseIP("192.168.0.255"), net.ParseIP("192.168.1.5"),
		"192.168.1.0/30", false, nil,
	},
	{
		net.ParseIP("192.168.0.254"), net.ParseIP("192.168.2.0"),
		"192.168.0.255/32", false, nil,
	},
	{
		net.ParseIP("192.168.0.255"), net.ParseIP("192.168.2.0"),
		"192.168.1.0/24", false, nil,
	},
	{
		net.ParseIP("12.168.0.254"), net.ParseIP("12.168.0.255"),
		"12.168.0.254/32", true, nil,
	},
	{
		net.ParseIP("2001:db7:ffff:ffff:ffff:ffff:ffff:ffff"), net.ParseIP("2001:db8:0:1::"),
		"2001:db8::/64", true, nil,
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
