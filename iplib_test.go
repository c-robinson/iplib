package iplib

import (
	"math/big"
	"net"
	"sort"
	"testing"
)

var IPTests = []struct {
	ipaddr net.IP
	next   net.IP
	prev   net.IP
	intval uint32
	hexval string
	inarpa string
	binval string
}{
	{
		net.ParseIP("10.1.2.3"),
		net.ParseIP("10.1.2.4"),
		net.ParseIP("10.1.2.2"),
		167838211,
		"0a010203",
		"3.2.1.10.in-addr.arpa",
		"00001010.00000001.00000010.00000011",
	},
	{
		net.ParseIP("10.1.2.255"),
		net.ParseIP("10.1.3.0"),
		net.ParseIP("10.1.2.254"),
		167838463,
		"0a0102ff",
		"255.2.1.10.in-addr.arpa",
		"00001010.00000001.00000010.11111111",
	},
	{
		net.ParseIP("10.1.2.0"),
		net.ParseIP("10.1.2.1"),
		net.ParseIP("10.1.1.255"),
		167838208,
		"0a010200",
		"0.2.1.10.in-addr.arpa",
		"00001010.00000001.00000010.00000000",
	},
	{
		net.ParseIP("255.255.255.255"),
		net.ParseIP("255.255.255.255"),
		net.ParseIP("255.255.255.254"),
		4294967295,
		"ffffffff",
		"255.255.255.255.in-addr.arpa",
		"11111111.11111111.11111111.11111111",
	},
	{
		net.ParseIP("0.0.0.0"),
		net.ParseIP("0.0.0.1"),
		net.ParseIP("0.0.0.0"),
		0,
		"00000000",
		"0.0.0.0.in-addr.arpa",
		"00000000.00000000.00000000.00000000",
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
		x := CompareIPs(tt.prev, PreviousIP(tt.ipaddr))
		if x != 0 {
			t.Errorf("On PreviousIP(%+v) expected %+v, got %+v", tt.ipaddr, tt.prev, PreviousIP(tt.ipaddr))
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

func TestIPToBinaryString(t *testing.T) {
	for _, tt := range IPTests {
		s := IPToBinaryString(tt.ipaddr)
		if s != tt.binval {
			t.Errorf("On IPToBinaryString(%+v) expected %s, got %s", tt.ipaddr, tt.binval, s)
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

func TestIP4ToARPA(t *testing.T) {
	for _, tt := range IPTests {
		s := IPToARPA(tt.ipaddr)
		if s != tt.inarpa {
			t.Errorf("On IP4ToARPA(%s) expected %s, got %s", tt.ipaddr, tt.inarpa, s)
		}
	}
}

var IP6Tests = []struct {
	ipaddr    string
	next      string
	prev      string
	bigintval string
	int64val  uint64
	hostbits  string
	hexval    string
	expand    string
	inarpa    string
	binval    string
}{
	{
		"2001:db8:85a3::8a2e:370:7334",
		"2001:db8:85a3::8a2e:370:7335",
		"2001:db8:85a3::8a2e:370:7333",
		"42540766452641154071740215577757643572",
		2306139570357600256,
		"2001:db8:85a3::",
		"2001:db8:85a3::8a2e:370:7334",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		"4.3.3.7.0.7.3.0.e.2.a.8.0.0.0.0.0.0.0.0.3.a.5.8.8.b.d.0.1.0.0.2.ip6.arpa",
		"00100000.00000001.00001101.10111000.10000101.10100011.00000000.00000000.00000000.00000000.10001010.00101110.00000011.01110000.01110011.00110100",
	},
	{
		"::",
		"::1",
		"::",
		"0",
		0,
		"::",
		"::",
		"0000:0000:0000:0000:0000:0000:0000:0000",
		"0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa",
		"00000000.00000000.00000000.00000000.00000000.00000000.00000000.00000000.00000000.00000000.00000000.00000000.00000000.00000000.00000000.00000000",
	},
	{
		"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
		"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
		"ffff:ffff:ffff:ffff:ffff:ffff:ffff:fffe",
		"340282366920938463463374607431768211455",
		18446744073709551615,
		"ffff:ffff:ffff:ffff::",
		"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
		"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
		"f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.ip6.arpa",
		"11111111.11111111.11111111.11111111.11111111.11111111.11111111.11111111.11111111.11111111.11111111.11111111.11111111.11111111.11111111.11111111",
	},
}

func TestNextIP6(t *testing.T) {
	for _, tt := range IP6Tests {
		x := CompareIPs(net.ParseIP(tt.next), NextIP(net.ParseIP(tt.ipaddr)))
		if x != 0 {
			t.Errorf("On IPv6 NextIP(%s) expected %s, got %s", tt.ipaddr, tt.next, NextIP(net.ParseIP(tt.ipaddr)))
		}
	}
}

func TestPrevIP6(t *testing.T) {
	for _, tt := range IP6Tests {
		x := CompareIPs(net.ParseIP(tt.prev), PreviousIP(net.ParseIP(tt.ipaddr)))
		if x != 0 {
			t.Errorf("On IPv6 PreviousIP(%s) expected %s, got %s", tt.ipaddr, tt.prev, PreviousIP(net.ParseIP(tt.ipaddr)))
		}
	}
}

func TestIP6ToBigint(t *testing.T) {
	for _, tt := range IP6Tests {
		i := IPToBigint(net.ParseIP(tt.ipaddr))
		if i.String() != tt.bigintval {
			t.Errorf("On IPToBigint(%s) expected %s, got %s", tt.ipaddr, tt.bigintval, i.String())
		}
	}
}

func TestIP6ToUint64(t *testing.T) {
	for _, tt := range IP6Tests {
		i := IP6ToUint64(net.ParseIP(tt.ipaddr))
		if i != tt.int64val {
			t.Errorf("On IP6ToUint64(%s) expected %d, got %d", tt.ipaddr, tt.int64val, i)
		}
	}
}

func TestUint64ToIP6(t *testing.T) {
	for _, tt := range IP6Tests {
		ip := Uint64ToIP6(tt.int64val)
		x := CompareIPs(ip, net.ParseIP(tt.hostbits))
		if x != 0 {
			t.Errorf("On IPv6 Uint64ToIP6(%s) expected %s, got %s", tt.ipaddr, tt.hostbits, ip)
		}
	}
}

func TestIP6ToBinaryString(t *testing.T) {
	for _, tt := range IP6Tests {
		s := IPToBinaryString(net.ParseIP(tt.ipaddr))
		if s != tt.binval {
			t.Errorf("On IPv6 IPToBinaryString(%s) expected %s, got %s", tt.ipaddr, tt.binval, s)
		}
	}
}

func TestIP6ToHexString(t *testing.T) {
	for _, tt := range IP6Tests {
		s := IPToHexString(net.ParseIP(tt.ipaddr))
		if s != tt.hexval {
			t.Errorf("On IPv6 IPToHexString(%s) expected %s, got %s", tt.ipaddr, tt.hexval, s)
		}
	}
}

func TestBigintToIP6(t *testing.T) {
	for _, tt := range IP6Tests {
		z := big.Int{}
		z.SetString(tt.bigintval, 10)
		ip := BigintToIP6(&z)
		x := CompareIPs(ip, net.ParseIP(tt.ipaddr))
		if x != 0 {
			t.Errorf("On BigintToIP6(%s) expected %s, got %s", tt.bigintval, tt.ipaddr, ip)
		}
	}
}

func TestExpandIP6(t *testing.T) {
	for _, tt := range IP6Tests {
		s := ExpandIP6(net.ParseIP(tt.ipaddr))
		if s != tt.expand {
			t.Errorf("On ExpandIP6(%s) expected '%s', got '%s'", tt.ipaddr, tt.expand, s)
		}
	}
}

func TestIP6ToARPA(t *testing.T) {
	for _, tt := range IP6Tests {
		s := IPToARPA(net.ParseIP(tt.ipaddr))
		if s != tt.inarpa {
			t.Errorf("On IP4ToARPA(%s) expected %s, got %s", tt.ipaddr, tt.inarpa, s)
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
		net.ParseIP("192.168.2.2"),
		net.ParseIP("192.168.1.1"),
		net.ParseIP("192.168.3.3"),
		257,
		257,
		257,
	},
	{
		net.ParseIP("10.0.0.0"),
		net.ParseIP("9.0.0.0"),
		net.ParseIP("11.0.0.0"),
		16777216,
		16777216,
		16777216,
	},
	{
		net.ParseIP("255.255.255.0"),
		net.ParseIP("255.255.252.0"),
		net.ParseIP("255.255.255.255"),
		768,
		255,
		768,
	},
	{
		net.ParseIP("0.0.0.255"),
		net.ParseIP("0.0.0.0"),
		net.ParseIP("0.0.3.255"),
		768,
		768,
		255,
	},
	{
		net.ParseIP("2001:db8:85a3::8a2e:370:7334"),
		net.ParseIP("2001:db8:85a3::8a2e:370:731c"),
		net.ParseIP("2001:db8:85a3::8a2e:370:734c"),
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
		// THIS IS JUST FOR TESTING, THE ONE BELOW IS FOR FIXING
		net.ParseIP("2001:db8:85a3::8a2e:370:7334"),
		net.ParseIP("2001:db8:85a3::8a2d:370:7334"),
		net.ParseIP("2001:db8:85a3::8a2f:370:7334"),
		"4294967296",
		"4294967296",
		"4294967296",
	},
	{
		net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ff00"),
		net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:fb00"),
		net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"),
		"1024",
		"255",
		"1024",
	},
	{
		net.ParseIP("::ff"),
		net.ParseIP("::"),
		net.ParseIP("::4ff"),
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
	// ParseIP() *always* returns 4-in-6 addresses, so we specify exactly what
	// we want here
	{net.IP{0, 0, 0, 0}, 4, 4 },
	{net.IP{192, 168, 1, 1}, 4, 4 },
	{net.IP{255, 255, 255, 255}, 4, 4 },
	{net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 6, 6 },
	{net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52}, 6, 6 },
	{net.IP{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, 6, 6 },
	// these are the 4-in-6 versions of the first 3 test cases
	{net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 0, 0, 0, 0}, 6, 4 },
	{net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 192, 168, 1, 1}, 6, 4 },
	{net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255}, 6, 4 },
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

var compareIPTests = []struct {
	pos    int
	ipaddr net.IP
	status int
}{
	{8, net.ParseIP("192.168.2.3"), -1},
	{1, net.ParseIP("10.0.0.3"), 1},
	{0, net.ParseIP("10.0.0.1"), 1},
	{10, net.ParseIP("192.168.3.255"), -1},
	{9, net.ParseIP("192.168.3.1"), -1},
	{2, net.ParseIP("10.0.1.0"), 1},
	{7, net.ParseIP("192.168.1.1"), -1},
	{3, net.ParseIP("44.0.0.1"), 1},
	{4, net.ParseIP("44.0.1.0"), 0},
	{5, net.ParseIP("44.1.0.0"), -1},
	{6, net.ParseIP("170.1.12.1"), -1},
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

var isAllTests = []struct{
	ipaddr net.IP
	isones bool
	iszero bool
	is4in6 bool
}{
	{ net.IP{0,0,0,0}, false, true, false },
	{ net.IP{255,255,255,255}, true, false, false },
	{ net.IP{192,168,1,1}, false, false, false },
	{ net.ParseIP("::ffff:0:0"), false, true, true },
	{ net.ParseIP("::ffff:ffff:ffff"), true, false, true },
	{ net.ParseIP("::ffff:c0a8:0101"), false, false, true },
	{ net.ParseIP("::"), false, true, false },
	{ net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"), true, false, false },
	{ net.ParseIP("2001:db8::1"), false, false, false },
}

func TestIs4in6(t *testing.T) {
	for _, tt := range isAllTests {
		v := Is4in6(tt.ipaddr)
		if v != tt.is4in6 {
			t.Errorf("%s: expected %t got %t", tt.ipaddr, tt.is4in6, v)
		}
	}
}

func TestIsAllOnes(t *testing.T) {
	for _, tt := range isAllTests {
		v := IsAllOnes(tt.ipaddr)
		if v != tt.isones {
			t.Errorf("%s: expected %t got %t", tt.ipaddr, tt.isones, v)
		}
	}
}

func TestIsAllZeroes(t *testing.T) {
	for _, tt := range isAllTests {
		v := IsAllZeroes(tt.ipaddr)
		if v != tt.iszero {
			t.Errorf("%s: expected %t got %t", tt.ipaddr, tt.iszero, v)
		}
	}
}