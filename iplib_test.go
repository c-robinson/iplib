package iplib

import (
	"math/big"
	"net"
	"sort"
	"strconv"
	"testing"
)

var IPTests = []struct {
	ipaddr net.IP
	next   net.IP
	prev   net.IP
	intval uint32
	hexval string
}{
	{
		net.IP{10, 1, 2, 3},
		net.IP{10, 1, 2, 4},
		net.IP{10, 1, 2, 2},
		167838211,
		"0a010203",
	},
	{
		net.IP{10, 1, 2, 255},
		net.IP{10, 1, 3, 0},
		net.IP{10, 1, 2, 254},
		167838463,
		"0a0102ff",
	},
	{
		net.IP{10, 1, 2, 0},
		net.IP{10, 1, 2, 1},
		net.IP{10, 1, 1, 255},
		167838208,
		"0a010200",
	},
	{
		net.IP{255, 255, 255, 255},
		net.IP{255, 255, 255, 255},
		net.IP{255, 255, 255, 254},
		4294967295,
		"ffffffff",
	},
	{
		net.IP{0, 0, 0, 0},
		net.IP{0, 0, 0, 1},
		net.IP{0, 0, 0, 0},
		0,
		"00000000",
	},
}

func TestNextIP(t *testing.T) {
	for _, tt := range IPTests {
		x := CompareIPs(tt.next, NextIP(tt.ipaddr))
		if x != 0 {
			t.Errorf("On NextIP(%+v) expected %+v, got %+v", tt.ipaddr, tt.next, NextIP(tt.ipaddr))
		}
	}
}

func TestPrevIP(t *testing.T) {
	for _, tt := range IPTests {
		x := CompareIPs(tt.prev, PrevIP(tt.ipaddr))
		if x != 0 {
			t.Errorf("On PrevIP(%+v) expected %+v, got %+v", tt.ipaddr, tt.prev, PrevIP(tt.ipaddr))
		}
	}
}

func TestIP4ToUint32(t *testing.T) {
	for _, tt := range IPTests {
		i := IP4ToUint32(tt.ipaddr)
		if i != tt.intval {
			t.Errorf("On IP4ToUint32(%+v) expected %d, got %d", tt.ipaddr, tt.intval, i)
		}
	}
}

func TestIPToHexString(t *testing.T) {
	for _, tt := range IPTests {
		s := IPToHexString(tt.ipaddr)
		if s != tt.hexval {
			t.Errorf("On IPToHexString(%+v) expected %s, got %s", tt.ipaddr, tt.hexval, s)
		}
	}
}

func TestHexStringToIP(t *testing.T) {
	for _, tt := range IPTests {
		ip := HexStringToIP(tt.hexval)
		x := CompareIPs(tt.ipaddr, ip)
		if x != 0 {
			t.Errorf("On HexStringToIP(%s) expected %s, got %s", tt.hexval, tt.ipaddr, ip)
		}
	}
}

func TestUint32ToIP4(t *testing.T) {
	for _, tt := range IPTests {
		ip := Uint32ToIP4(tt.intval)
		x := CompareIPs(ip, tt.ipaddr)
		if x != 0 {
			t.Errorf("On Uint32ToIP4(%d) expected %+v, got %+v", tt.intval, tt.ipaddr, ip)
		}
	}
}

var IP6Tests = []struct {
	ipaddr net.IP
	next   net.IP
	prev   net.IP
	intval string
	hexval string
	expand string
}{
	{
		net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52},
		net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 53},
		net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 51},
		"42540766452641154071740215577757643572",
		"2001:db8:85a3::8a2e:370:7334",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
	},
	{
		net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"0",
		"::",
		"0000:0000:0000:0000:0000:0000:0000:0000",
	},
	{
		net.IP{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		net.IP{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		net.IP{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 254},
		"340282366920938463463374607431768211455",
		"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
		"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
	},
}

func TestNextIP6(t *testing.T) {
	for _, tt := range IP6Tests {
		x := CompareIPs(tt.next, NextIP(tt.ipaddr))
		if x != 0 {
			t.Errorf("On IPv6 NextIP(%+v) expected %+v, got %+v", tt.ipaddr, tt.next, NextIP(tt.ipaddr))
		}
	}
}

func TestPrevIP6(t *testing.T) {
	for _, tt := range IP6Tests {
		x := CompareIPs(tt.prev, PrevIP(tt.ipaddr))
		if x != 0 {
			t.Errorf("On IPv6 PrevIP(%+v) expected %+v, got %+v", tt.ipaddr, tt.prev, PrevIP(tt.ipaddr))
		}
	}
}

func TestIP6ToBigint(t *testing.T) {
	for _, tt := range IP6Tests {
		i := IPToBigint(tt.ipaddr)
		if i.String() != tt.intval {
			t.Errorf("On IPToBigint(%+v) expected %s, got %v", tt.ipaddr, tt.intval, i)
		}
	}
}

func TestIP6ToHexString(t *testing.T) {
	for _, tt := range IP6Tests {
		s := IPToHexString(tt.ipaddr)
		if s != tt.hexval {
			t.Errorf("On IPv6 IPToHexString(%+v) expected %s, got %s", tt.ipaddr, tt.hexval, s)
		}
	}
}

func TestBigintToIP6(t *testing.T) {
	for _, tt := range IP6Tests {
		z := big.Int{}
		z.SetString(tt.intval, 10)
		ip := BigintToIP6(&z)
		x := CompareIPs(ip, tt.ipaddr)
		if x != 0 {
			t.Errorf("On BigintToIP6(%s) expected %+v, got %+v", tt.intval, tt.ipaddr, ip)
		}
	}
}

func TestExpandIP6(t *testing.T) {
	for _, tt := range IP6Tests {
		s := ExpandIP6(tt.ipaddr)
		if s != tt.expand {
			t.Errorf("On ExpandIP6(%s) expected '%s', got '%s'", tt.ipaddr, tt.expand, s)
		}
	}
}

var IPDeltaTests = []struct {
	ipaddr net.IP
	decr   net.IP
	incr   net.IP
	intval uint32
	incres uint32
	decres uint32
}{
	{
		net.IP{192, 168, 2, 2},
		net.IP{192, 168, 1, 1},
		net.IP{192, 168, 3, 3},
		257,
		257,
		257,
	},
	{
		net.IP{10, 0, 0, 0},
		net.IP{9, 0, 0, 0},
		net.IP{11, 0, 0, 0},
		16777216,
		16777216,
		16777216,
	},
	{
		net.IP{255, 255, 255, 0},
		net.IP{255, 255, 252, 0},
		net.IP{255, 255, 255, 255},
		768,
		255,
		768,
	},
	{
		net.IP{0, 0, 0, 255},
		net.IP{0, 0, 0, 0},
		net.IP{0, 0, 3, 255},
		768,
		768,
		255,
	},
	{
		net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52},
		net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 28},
		net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 76},
		24,
		24,
		24,
	},
}

func TestDeltaIP(t *testing.T) {
	for _, tt := range IPDeltaTests {
		i := DeltaIP(tt.ipaddr, tt.incr)
		if i != tt.incres {
			t.Errorf("On DeltaIP(%s, %s) expected %d got %d", tt.ipaddr, tt.incr, tt.incres, i)
		}

		i = DeltaIP(tt.ipaddr, tt.decr)
		if i != tt.decres {
			t.Errorf("On DeltaIP(%s, %s) expected %d got %d", tt.ipaddr, tt.decr, tt.decres, i)
		}
	}
}

func TestDecrementIPBy(t *testing.T) {
	for _, tt := range IPDeltaTests {
		ip := DecrementIPBy(tt.ipaddr, tt.intval)
		x := CompareIPs(ip, tt.decr)
		if x != 0 {
			t.Errorf("On DecrementIPBy(%s, %d) expected %s got %s", tt.ipaddr, tt.intval, tt.decr, ip)
		}
	}
}

func TestIncrementIPBy(t *testing.T) {
	for _, tt := range IPDeltaTests {
		ip := IncrementIPBy(tt.ipaddr, tt.intval)
		x := CompareIPs(ip, tt.incr)
		if x != 0 {
			t.Errorf("On IncrementIPBy(%s, %d) expected %s got %s", tt.ipaddr, tt.intval, tt.incr, ip)
		}
	}
}

var IPDelta6Tests = []struct {
	ipaddr net.IP
	decr   net.IP
	incr   net.IP
	intval string
	incres string
	decres string
}{
	{
		net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52},
		net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 45, 3, 112, 115, 52},
		net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 47, 3, 112, 115, 52},
		"4294967296",
		"4294967296",
		"4294967296",
	},
	{
		net.IP{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0},
		net.IP{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 251, 0},
		net.IP{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		"1024",
		"255",
		"1024",
	},
	{
		net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255},
		net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 255},
		"1024",
		"1024",
		"255",
	},
}

func TestDeltaIP6(t *testing.T) {
	for _, tt := range IPDelta6Tests {
		i := DeltaIP6(tt.ipaddr, tt.incr)
		if i.String() != tt.incres {
			t.Errorf("On DeltaIP(%s, %s) expected %s got %s", tt.ipaddr, tt.incr, tt.incres, i)
		}

		i = DeltaIP6(tt.ipaddr, tt.decr)
		if i.String() != tt.decres {
			t.Errorf("On DeltaIP(%s, %s) expected %s got %s", tt.ipaddr, tt.decr, tt.decres, i)
		}
	}
}

func TestDecrementIP6By(t *testing.T) {
	for _, tt := range IPDelta6Tests {
		z := big.Int{}
		z.SetString(tt.intval, 10)
		ip := DecrementIP6By(tt.ipaddr, &z)
		x := CompareIPs(ip, tt.decr)
		if x != 0 {
			t.Errorf("On DecrementIPBy(%s, %s) expected %s got %s", tt.ipaddr, tt.intval, tt.decr, ip)
		}
	}
}

func TestIncrementIP6By(t *testing.T) {
	for _, tt := range IPDelta6Tests {
		z := big.Int{}
		z.SetString(tt.intval, 10)
		ip := IncrementIP6By(tt.ipaddr, &z)
		x := CompareIPs(ip, tt.incr)
		if x != 0 {
			t.Errorf("On IncrementIP6By(%s, %s) expected %s got %s", tt.ipaddr, tt.intval, tt.incr, ip)
		}
	}
}

var IPVersionTests = []struct {
	ipaddr   net.IP
	version  int
	eversion int
}{
	{
		net.IP{0, 0, 0, 0},
		4,
		4,
	},
	{
		net.IP{192, 168, 1, 1},
		4,
		4,
	},
	{
		net.IP{255, 255, 255, 255},
		4,
		4,
	},
	{
		net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		6,
		6,
	},
	{
		net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52},
		6,
		6,
	},
	{
		net.IP{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		6,
		6,
	},
	// these are the 6-to-4 versions of the first 3 test cases
	{
		net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 0, 0, 0, 0},
		6,
		4,
	},
	{
		net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 192, 168, 1, 1},
		6,
		4,
	},
	{
		net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255},
		6,
		4,
	},
}

func Test_Version(t *testing.T) {
	for _, tt := range IPVersionTests {
		version := Version(tt.ipaddr)
		if version != tt.version {
			t.Errorf("On %v expected %d got %d", tt.ipaddr, tt.version, version)
		}
	}
}

func Test_EffectiveVersion(t *testing.T) {
	for _, tt := range IPVersionTests {
		eversion := EffectiveVersion(tt.ipaddr)
		if eversion != tt.eversion {
			t.Errorf("On %v expected %d got %d", tt.ipaddr, tt.eversion, eversion)
		}
	}
}

var NewNetTests = []struct {
	ip      net.IP
	masklen int
	out     string
}{
	{
		net.IP{192, 168, 0, 0},
		32,
		"192.168.0.0/32",
	},
	{
		net.IP{192, 168, 0, 0},
		24,
		"192.168.0.0/24",
	},
	{
		net.IP{192, 168, 0, 7},
		32,
		"192.168.0.7/32",
	},
	{
		net.IP{192, 168, 0, 7},
		24,
		"192.168.0.0/24",
	},
	{
		net.IP{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		64,
		"2001:db8::/64",
	},
}

func TestNewNet(t *testing.T) {
	for _, tt := range NewNetTests {
		xnet := NewNet(tt.ip, tt.masklen)
		_, pnet, _ := net.ParseCIDR(tt.out)
		if xnet.String() != pnet.String() {
			t.Errorf("On NewNet(%s, %d) expected %s got %s", tt.ip.String(), tt.masklen, pnet.String(), xnet.String())
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
		net.IP{192, 168, 0, 255},
		net.IP{192, 168, 2, 0},
		"192.168.1.0/24",
		true,
		nil,
	},
	{
		net.IP{192, 168, 0, 255},
		net.IP{10, 0, 0, 0},
		"",
		false,
		ErrNoValidRange,
	},
	{
		net.IP{192, 168, 0, 255},
		net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52},
		"",
		false,
		ErrNoValidRange,
	},
	{
		net.IP{192, 168, 0, 255},
		net.IP{192, 168, 0, 255},
		"",
		false,
		ErrNoValidRange,
	},
	{
		net.IP{192, 168, 0, 255},
		net.IP{192, 168, 1, 1},
		"192.168.1.0/32",
		true,
		nil,
	},
	{
		net.IP{192, 168, 1, 0},
		net.IP{192, 168, 1, 2},
		"192.168.1.1/32",
		true,
		nil,
	},
	{
		net.IP{192, 168, 0, 255},
		net.IP{192, 168, 1, 2},
		"192.168.1.0/31",
		true,
		nil,
	},
	{
		net.IP{192, 168, 0, 255},
		net.IP{192, 168, 1, 3},
		"192.168.1.0/31",
		false,
		nil,
	},
	{
		net.IP{192, 168, 1, 0},
		net.IP{192, 168, 1, 3},
		"192.168.1.1/32",
		false,
		nil,
	},
	{
		net.IP{192, 168, 0, 255},
		net.IP{192, 168, 1, 4},
		"192.168.1.0/30",
		true,
		nil,
	},
	{
		net.IP{192, 168, 0, 255},
		net.IP{192, 168, 1, 5},
		"192.168.1.0/30",
		false,
		nil,
	},
	{
		net.IP{192, 168, 0, 254},
		net.IP{192, 168, 2, 0},
		"192.168.0.255/32",
		false,
		nil,
	},
	{
		net.IP{192, 168, 0, 255},
		net.IP{192, 168, 2, 0},
		"192.168.1.0/24",
		true,
		nil,
	},
	{
		net.IP{32, 1, 13, 183, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		net.IP{32, 1, 13, 184, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		"2001:db8::/64",
		true,
		nil,
	},
}

func TestNewNetBetween(t *testing.T) {
	for _, tt := range NewNetBetweenTests {
		xnet, exact, err := NewNetBetween(tt.start, tt.end)
		if tt.err != nil {
			if tt.err != err {
				t.Errorf("On NewNetBetween(%s, %s) expected error '%v', got '%v'", tt.start, tt.end, tt.err, err)
			}
		} else {
			if xnet.String() != tt.xnet {
				t.Errorf("On NewNetBetween(%s, %s) expected '%s', got '%s'", tt.start, tt.end, tt.xnet, xnet.String())
			}
			if exact != tt.exact {
				t.Errorf("On NewNetBetween(%s, %s) expected '%t', got '%t'", tt.start, tt.end, tt.exact, exact)
			}
		}
	}
}

var NetworkTests = []struct {
	inaddrStr  string
	ipaddr     net.IP
	inaddrMask int
	network    net.IP
	netmask    net.IPMask
	wildcard   net.IPMask
	broadcast  net.IP
	firstaddr  net.IP
	lastaddr   net.IP
	version    int
	count      string // might overflow uint64
}{
	{
		"10.1.2.3/8",
		net.IP{10, 1, 2, 3},
		8,
		net.IP{10, 0, 0, 0},
		net.IPMask{255, 0, 0, 0},
		net.IPMask{0, 255, 255, 255},
		net.IP{10, 255, 255, 255},
		net.IP{10, 0, 0, 1},
		net.IP{10, 255, 255, 254},
		4,
		"16777214",
	},
	{
		"192.168.1.1/23",
		net.IP{192, 168, 1, 1},
		23,
		net.IP{192, 168, 0, 0},
		net.IPMask{255, 255, 254, 0},
		net.IPMask{0, 0, 1, 255},
		net.IP{192, 168, 1, 255},
		net.IP{192, 168, 0, 1},
		net.IP{192, 168, 1, 254},
		4,
		"510",
	},
	{
		"192.168.1.61/26",
		net.IP{192, 168, 1, 61},
		26,
		net.IP{192, 168, 1, 0},
		net.IPMask{255, 255, 255, 192},
		net.IPMask{0, 0, 0, 63},
		net.IP{192, 168, 1, 63},
		net.IP{192, 168, 1, 1},
		net.IP{192, 168, 1, 62},
		4,
		"62",
	},
	{
		"192.168.1.66/26",
		net.IP{192, 168, 1, 66},
		26,
		net.IP{192, 168, 1, 64},
		net.IPMask{255, 255, 255, 192},
		net.IPMask{0, 0, 0, 63},
		net.IP{192, 168, 1, 127},
		net.IP{192, 168, 1, 65},
		net.IP{192, 168, 1, 126},
		4,
		"62",
	},
	{
		"192.168.1.1/30",
		net.IP{192, 168, 1, 1},
		30,
		net.IP{192, 168, 1, 0},
		net.IPMask{255, 255, 255, 252},
		net.IPMask{0, 0, 0, 3},
		net.IP{192, 168, 1, 3},
		net.IP{192, 168, 1, 1},
		net.IP{192, 168, 1, 2},
		4,
		"2",
	},
	{
		"192.168.1.1/31",
		net.IP{192, 168, 1, 1},
		31,
		net.IP{192, 168, 1, 0},
		net.IPMask{255, 255, 255, 254},
		net.IPMask{0, 0, 0, 1},
		net.IP{192, 168, 1, 1},
		net.IP{192, 168, 1, 0},
		net.IP{192, 168, 1, 1},
		4,
		"0",
	},
	{
		"192.168.1.15/32",
		net.IP{192, 168, 1, 15},
		32,
		net.IP{192, 168, 1, 15},
		net.IPMask{255, 255, 255, 255},
		net.IPMask{0, 0, 0, 0},
		net.IP{192, 168, 1, 15},
		net.IP{192, 168, 1, 15},
		net.IP{192, 168, 1, 15},
		4,
		"1",
	},
	{
		"2001:db8::/64",
		net.IP{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		64,
		net.IP{},
		net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0},
		net.IPMask{},
		net.IP{},
		net.IP{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		net.IP{32, 1, 13, 184, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255},
		6,
		"18446744073709551616",
	},
}

// ParseCIDR wraps net.ParseCIDR so it's redundant to test it except to make sure the wildcard is correct
func TestParseCIDR(t *testing.T) {
	for _, tt := range NetworkTests {
		if tt.version == 6 {
			continue
		}
		_, ipn, _ := ParseCIDR(tt.inaddrStr)
		if ipn.Wildcard().String() != tt.wildcard.String() {
			t.Errorf("On %s got Network.Wildcard == %v, want %v", tt.inaddrStr, ipn.Wildcard(), tt.wildcard)
		}
	}
}

func TestNet_BroadcastAddress(t *testing.T) {
	for _, tt := range NetworkTests {
		if tt.version == 6 {
			continue
		}
		_, ipn, _ := ParseCIDR(tt.inaddrStr)
		if addr := ipn.BroadcastAddress(); !tt.broadcast.Equal(addr) {
			t.Errorf("On %s got Network.Broadcast == %v, want %v", tt.inaddrStr, addr, tt.broadcast)
		}
	}
}

func TestNet_Version(t *testing.T) {
	for _, tt := range NetworkTests {
		_, ipnp, _ := ParseCIDR(tt.inaddrStr)
		ipnn := NewNet(tt.ipaddr, tt.inaddrMask)
		if ipnp.Version() != tt.version {
			t.Errorf("From ParseCIDR %s got Network.Version == %d, expect %d", tt.inaddrStr, ipnp.Version(), tt.version)
		}
		if ipnn.Version() != tt.version {
			t.Errorf("From NewNet %s got Network.Version == %d, want %d", tt.inaddrStr, ipnn.Version(), tt.version)
		}
	}
}

func TestNet_Count(t *testing.T) {
	for _, tt := range NetworkTests {
		if tt.version == 6 {
			continue
		}
		_, ipn, _ := ParseCIDR(tt.inaddrStr)
		count, _ := strconv.Atoi(tt.count)
		if ipn.Count() != uint32(count) {
			t.Errorf("On %s got Network.Count == %d, want %d", tt.inaddrStr, ipn.Count(), count)
		}
	}
}

func TestNet_Count4(t *testing.T) {
	for _, tt := range NetworkTests {
		if tt.version == 6 {
			continue
		}
		_, ipn, _ := ParseCIDR(tt.inaddrStr)
		count, _ := strconv.Atoi(tt.count)
		if ipn.Count() != uint32(count) {
			t.Errorf("On %s got Network.Count4 == %d, want %d", tt.inaddrStr, ipn.Count4(), count)
		}
	}
}

func TestNet_Count6(t *testing.T) {
	for _, tt := range NetworkTests {
		_, ipn, _ := ParseCIDR(tt.inaddrStr)
		count, _ := new(big.Int).SetString(tt.count, 10)
		res := ipn.Count6().Cmp(count)
		if res != 0 {
			t.Errorf("On %s got Network.Count6 == %s, want %s", tt.inaddrStr, ipn.Count6().String(), count.String())
		}
	}
}

func TestNet_FirstAddress(t *testing.T) {
	for _, tt := range NetworkTests {
		_, ipn, _ := ParseCIDR(tt.inaddrStr)
		if addr := ipn.FirstAddress(); !tt.firstaddr.Equal(addr) {
			t.Errorf("On %s got Network.FirstAddress == %v, want %v", tt.inaddrStr, addr, tt.firstaddr)
		}
	}
}

func TestNet_finalAddress(t *testing.T) {
	for _, tt := range NetworkTests {
		if tt.version == 6 {
			continue
		}
		_, ipn, _ := ParseCIDR(tt.inaddrStr)
		if addr, ones := ipn.finalAddress(); !tt.broadcast.Equal(addr) {
			t.Errorf("On %s got Network.finalAddress == %v, want %v mask length %d)", tt.inaddrStr, addr, tt.broadcast, ones)
		}
	}
}

func TestNet_LastAddress(t *testing.T) {
	for _, tt := range NetworkTests {
		_, ipn, _ := ParseCIDR(tt.inaddrStr)
		if addr := ipn.LastAddress(); !tt.lastaddr.Equal(addr) {
			t.Errorf("On %s got Network.LastAddress == %v, want %v", tt.inaddrStr, addr, tt.lastaddr)
		}
	}
}

func TestNet_NetworkAddress(t *testing.T) {
	for _, tt := range NetworkTests {
		if tt.version == 6 {
			continue
		}
		_, ipn, _ := ParseCIDR(tt.inaddrStr)
		if addr := ipn.NetworkAddress(); !tt.network.Equal(addr) {
			t.Errorf("On %s got Network.NetworkAddress == %v, want %v", tt.inaddrStr, addr, tt.network)
		}
	}
}

var enumerateTests = []struct {
	inaddr string
	total  int
	last   net.IP
}{
	{"192.168.0.0/22", 1022, net.IP{192, 168, 3, 254}},
	{"192.168.0.0/23", 510, net.IP{192, 168, 1, 254}},
	{"192.168.0.0/24", 254, net.IP{192, 168, 0, 254}},
	{"192.168.0.0/25", 126, net.IP{192, 168, 0, 126}},
	{"192.168.0.0/26", 62, net.IP{192, 168, 0, 62}},
	{"192.168.0.0/27", 30, net.IP{192, 168, 0, 30}},
	{"192.168.0.0/28", 14, net.IP{192, 168, 0, 14}},
	{"192.168.0.0/29", 6, net.IP{192, 168, 0, 6}},
	{"192.168.0.0/30", 2, net.IP{192, 168, 0, 2}},
	{"192.168.0.0/31", 2, net.IP{192, 168, 0, 1}},
	{"192.168.0.0/32", 1, net.IP{192, 168, 0, 0}},
}

func TestNet_Enumerate(t *testing.T) {
	for _, tt := range enumerateTests {
		_, ipn, _ := ParseCIDR(tt.inaddr)
		addrlist := ipn.Enumerate(0, 0)
		if len(addrlist) != tt.total {
			t.Errorf("On %s Network.Enumerate(0,0) got size %d, want %d", tt.inaddr, len(addrlist), tt.total)
		}
		x := CompareIPs(tt.last, addrlist[tt.total-1])
		if x != 0 {
			t.Errorf("On %s Network.Enumerate(0,0) got last member %+v, want %+v", tt.inaddr, addrlist[tt.total-1], tt.last)
		}

	}
}

var enumerateVariableTests = []struct {
	offset uint32
	size   uint32
	total  int
	first  net.IP
	last   net.IP
}{
	{0, 0, 1022, net.IP{192, 168, 0, 1}, net.IP{192, 168, 3, 254}},
	{1, 0, 1021, net.IP{192, 168, 0, 2}, net.IP{192, 168, 3, 254}},
	{256, 0, 766, net.IP{192, 168, 1, 1}, net.IP{192, 168, 3, 254}},
	{0, 128, 128, net.IP{192, 168, 0, 1}, net.IP{192, 168, 0, 128}},
	{20, 128, 128, net.IP{192, 168, 0, 21}, net.IP{192, 168, 0, 148}},
	{1000, 100, 22, net.IP{192, 168, 3, 233}, net.IP{192, 168, 3, 254}},
}

func TestNet_EnumerateWithVariables(t *testing.T) {
	_, ipn, _ := ParseCIDR("192.168.0.0/22")
	for _, tt := range enumerateVariableTests {
		addrlist := ipn.Enumerate(tt.size, tt.offset)
		if len(addrlist) != tt.total {
			t.Errorf("On Network.Enumerate(%d,%d) got size %d, want %d", tt.size, tt.offset, len(addrlist), tt.total)
		}
		x := CompareIPs(tt.first, addrlist[0])
		if x != 0 {
			t.Errorf("On Network.Enumerate(%d,%d) got first member %+v, want %+v", tt.size, tt.offset, addrlist[0], tt.first)
		}
		y := CompareIPs(tt.last, addrlist[len(addrlist)-1])
		if y != 0 {
			t.Errorf("On Network.Enumerate(%d,%d) got last member %+v, want %+v", tt.size, tt.offset, addrlist[len(addrlist)-1], tt.last)
		}

	}
}

var incrTests = []struct {
	inaddr   string
	ipaddr   net.IP
	nextaddr net.IP
	nexterr  error
}{
	{
		"192.168.1.0/23",
		net.IP{192, 168, 1, 0},
		net.IP{192, 168, 1, 1},
		nil,
	},
	{
		"192.168.1.0/24",
		net.IP{192, 168, 1, 254},
		net.IP{192, 168, 1, 255},
		ErrBroadcastAddress,
	},
	{
		"192.168.2.0/24",
		net.IP{192, 168, 2, 1},
		net.IP{192, 168, 2, 2},
		nil,
	},
	{
		"192.168.3.0/24",
		net.IP{192, 168, 3, 0},
		net.IP{192, 168, 3, 1},
		nil,
	},
	{
		"192.168.4.0/24",
		net.IP{192, 168, 5, 1},
		net.IP{},
		ErrAddressOutOfRange,
	},
	{
		"192.168.1.0/31",
		net.IP{192, 168, 1, 0},
		net.IP{192, 168, 1, 1},
		ErrBroadcastAddress,
	},
	{
		"192.168.1.0/32",
		net.IP{192, 168, 1, 0},
		net.IP{},
		ErrAddressAtEndOfRange,
	},
}

func TestNet_NextIP(t *testing.T) {
	for _, tt := range incrTests {
		_, ipn, _ := ParseCIDR(tt.inaddr)
		addr, err := ipn.NextIP(tt.ipaddr)
		if !addr.Equal(tt.nextaddr) {
			t.Errorf("For %s expected %v, got %v", tt.inaddr, tt.nextaddr, addr)
		}
		if err != tt.nexterr {
			t.Errorf("For %s expected \"%v\", got \"%v\"", tt.inaddr, tt.nexterr, err)
		}
	}
}

var decrTests = []struct {
	inaddr   string
	ipaddr   net.IP
	prevaddr net.IP
	preverr  error
}{
	{
		"192.168.1.0/23",
		net.IP{192, 168, 1, 0},
		net.IP{192, 168, 0, 255},
		nil,
	},
	{
		"192.168.1.0/24",
		net.IP{192, 168, 1, 254},
		net.IP{192, 168, 1, 253},
		nil,
	},
	{
		"192.168.2.0/24",
		net.IP{192, 168, 2, 1},
		net.IP{192, 168, 2, 0},
		ErrNetworkAddress,
	},
	{
		"192.168.3.0/24",
		net.IP{192, 168, 3, 0},
		net.IP{},
		ErrAddressAtEndOfRange,
	},
	{
		"192.168.4.0/24",
		net.IP{192, 168, 5, 1},
		net.IP{},
		ErrAddressOutOfRange,
	},
	{
		"192.168.1.1/31",
		net.IP{192, 168, 1, 1},
		net.IP{192, 168, 1, 0},
		ErrNetworkAddress,
	},
	{
		"192.168.1.0/32",
		net.IP{192, 168, 1, 0},
		net.IP{},
		ErrAddressAtEndOfRange,
	},
}

func TestNet_PreviousIP(t *testing.T) {
	for _, tt := range decrTests {
		_, ipn, _ := ParseCIDR(tt.inaddr)
		addr, err := ipn.PreviousIP(tt.ipaddr)
		if !addr.Equal(tt.prevaddr) {
			t.Errorf("For %s expected %v, got %v", tt.inaddr, tt.prevaddr, addr)
		}
		if err != tt.preverr {
			t.Errorf("For %s expected \"%v\", got \"%v\"", tt.inaddr, tt.preverr, err)
		}
	}
}

var supernetTests = []struct {
	in      string
	masklen int
	out     string
}{
	{
		"192.168.1.0/24",
		25,
		"192.168.1.0/24",
	},
	{
		"192.168.1.0/24",
		23,
		"192.168.0.0/23",
	},
	{
		"192.168.1.0/24",
		0,
		"192.168.0.0/23",
	},
	{
		"192.168.1.0/24",
		22,
		"192.168.0.0/22",
	},
	{
		"192.168.1.4/30",
		24,
		"192.168.1.0/24",
	},
}

func TestNet_Supernet(t *testing.T) {
	for _, tt := range supernetTests {
		_, inet, _ := ParseCIDR(tt.in)
		onet := inet.Supernet(tt.masklen)
		if onet.String() != tt.out {
			t.Errorf("On Net{%s}.Supernet(%d) expected %s got %s", tt.in, tt.masklen, tt.out, onet.String())
		}
	}
}

var compareIPTests = []struct {
	pos    int
	ipaddr net.IP
	status int
}{
	{8, net.IP{192, 168, 2, 3}, -1},
	{1, net.IP{10, 0, 0, 3}, 1},
	{0, net.IP{10, 0, 0, 1}, 1},
	{10, net.IP{192, 168, 3, 255}, -1},
	{9, net.IP{192, 168, 3, 1}, -1},
	{2, net.IP{10, 0, 1, 0}, 1},
	{7, net.IP{192, 168, 1, 1}, -1},
	{3, net.IP{44, 0, 0, 1}, 1},
	{4, net.IP{44, 0, 1, 0}, 0},
	{5, net.IP{44, 1, 0, 0}, -1},
	{6, net.IP{170, 1, 12, 1}, -1},
}

func TestCompareIPs(t *testing.T) {
	a := compareIPTests[8]
	a1 := []net.IP{}
	for _, b := range compareIPTests {
		a1 = append(a1, b.ipaddr)
		val := CompareIPs(a.ipaddr, b.ipaddr)
		if val != b.status {
			t.Errorf("For %s expected %d got %d", b.ipaddr, b.status, val)
		}
	}
	sort.Sort(ByIP(a1))
	for _, b := range compareIPTests {
		if a1[b.pos].String() != b.ipaddr.String() {
			t.Errorf("Expected %s at position %d, but found %s", b.ipaddr, b.pos, a1[b.pos])
		}
	}
}

var compareNetworks = map[int]string{
	0: "192.168.0.0/16",
	1: "192.168.0.0/23",
	2: "192.168.1.0/24",
	3: "192.168.1.0/24",
	4: "192.168.3.0/26",
	5: "192.168.3.64/26",
	6: "192.168.3.128/25",
	7: "192.168.4.0/24",
}

func TestCompareNets(t *testing.T) {
	a := ByNet{}
	for _, v := range compareNetworks {
		_, ipn, _ := ParseCIDR(v)
		a = append(a, ipn)
	}
	sort.Sort(ByNet(a))
	for k, v := range compareNetworks {
		if a[k].String() != v {
			t.Errorf("Subnet %s not at expected position %d. Got %s instead", v, k, a[k].String())
		}

	}

}

var compareCIDR = []struct {
	network string
	subnet  string
	result  bool
}{
	{"192.168.0.0/16", "192.168.45.0/24", true},
	{"192.168.45.0/24", "192.168.45.0/26", true},
	{"192.168.45.0/24", "192.168.46.0/26", false},
	{"10.1.1.1/24", "10.0.0.0/8", false},
}

func TestNet_ContainsNetwork(t *testing.T) {
	for _, cidr := range compareCIDR {
		_, ipn, _ := ParseCIDR(cidr.network)
		_, sub, _ := ParseCIDR(cidr.subnet)
		result := ipn.ContainsNet(sub)
		if result != cidr.result {
			t.Errorf("For \"%s contains %s\" expected %v got %v", cidr.network, cidr.subnet, cidr.result, result)
		}
	}
}
