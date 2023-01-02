package iplib

import (
	"math/big"
	"net"
	"reflect"
	"sort"
	"testing"
)

func TestCopyIP(t *testing.T) {
	ipa := net.ParseIP("192.168.23.5")
	ipb := CopyIP(ipa)
	if reflect.ValueOf(ipa).Pointer() == reflect.ValueOf(ipb).Pointer() {
		t.Errorf("CopyIP() failed to copy (copied IP shares pointer with original)!")
	}
	if CompareIPs(ipa, ipb) != 0 {
		t.Errorf("CopyIP() failed to copy (value of copied IP does not match original)!")
	}
}

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
	for i, tt := range IPTests {
		x := CompareIPs(tt.next, NextIP(tt.ipaddr))
		if x != 0 {
			t.Errorf("[%d] want %s, got %s", i, tt.next, NextIP(tt.ipaddr))
		}
	}
}

func TestPrevIP(t *testing.T) {
	for i, tt := range IPTests {
		x := CompareIPs(tt.prev, PreviousIP(tt.ipaddr))
		if x != 0 {
			t.Errorf("[%d] want %s got %s", i, tt.prev, PreviousIP(tt.ipaddr))
		}
	}
}

func TestIP4ToUint32(t *testing.T) {
	for i, tt := range IPTests {
		z := IP4ToUint32(tt.ipaddr)
		if z != tt.intval {
			t.Errorf("[%d] want %d got %d", i, tt.intval, z)
		}
	}
}

func TestIPToHexString(t *testing.T) {
	for i, tt := range IPTests {
		s := IPToHexString(tt.ipaddr)
		if s != tt.hexval {
			t.Errorf("[%d] want %s got %s", i, tt.hexval, s)
		}
	}
}

func TestIPToBinaryString(t *testing.T) {
	for i, tt := range IPTests {
		s := IPToBinaryString(tt.ipaddr)
		if s != tt.binval {
			t.Errorf("[%d] expected %s, got %s", i, tt.binval, s)
		}
	}
}

func TestHexStringToIP(t *testing.T) {
	for i, tt := range IPTests {
		ip := HexStringToIP(tt.hexval)
		x := CompareIPs(tt.ipaddr, ip)
		if x != 0 {
			t.Errorf("[%d] want %s got %s", i, tt.ipaddr, ip)
		}
	}
}

func TestHexStringToIPBadVals(t *testing.T) {
	ip := HexStringToIP("placebo")
	if ip != nil {
		t.Errorf("non-ip word should return nil")
	}
	ip = HexStringToIP("2001:db8::/24")
	if ip != nil {
		t.Errorf("network address should return nil")
	}
}

func TestUint32ToIP4(t *testing.T) {
	for i, tt := range IPTests {
		ip := Uint32ToIP4(tt.intval)
		x := CompareIPs(ip, tt.ipaddr)
		if x != 0 {
			t.Errorf("[%d] want %s got %s", i, tt.ipaddr, ip)
		}
	}
}

func TestIP4ToARPA(t *testing.T) {
	for i, tt := range IPTests {
		s := IPToARPA(tt.ipaddr)
		if s != tt.inarpa {
			t.Errorf("[%d] want %s got %s", i, tt.inarpa, s)
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
	for i, tt := range IP6Tests {
		x := CompareIPs(net.ParseIP(tt.next), NextIP(net.ParseIP(tt.ipaddr)))
		if x != 0 {
			t.Errorf("[%d] want %s got %s", i, tt.next, NextIP(net.ParseIP(tt.ipaddr)))
		}
	}
}

func TestPrevIP6(t *testing.T) {
	for i, tt := range IP6Tests {
		x := CompareIPs(net.ParseIP(tt.prev), PreviousIP(net.ParseIP(tt.ipaddr)))
		if x != 0 {
			t.Errorf("[%d] want %s got %s", i, tt.prev, PreviousIP(net.ParseIP(tt.ipaddr)))
		}
	}
}

func TestIP6ToBigint(t *testing.T) {
	for _, tt := range IP6Tests {
		i := IPToBigint(net.ParseIP(tt.ipaddr))
		if i.String() != tt.bigintval {
			t.Errorf("[%d] want %s got %s", i, tt.bigintval, i.String())
		}
	}
}

func TestIP6ToUint64(t *testing.T) {
	for i, tt := range IP6Tests {
		z := IP6ToUint64(net.ParseIP(tt.ipaddr))
		if z != tt.int64val {
			t.Errorf("[%d] want %d got %d", i, tt.int64val, z)
		}
	}
}

func TestUint64ToIP6(t *testing.T) {
	for i, tt := range IP6Tests {
		ip := Uint64ToIP6(tt.int64val)
		x := CompareIPs(ip, net.ParseIP(tt.hostbits))
		if x != 0 {
			t.Errorf("[%d] want %s got %s", i, tt.hostbits, ip)
		}
	}
}

func TestIP6ToBinaryString(t *testing.T) {
	for i, tt := range IP6Tests {
		s := IPToBinaryString(net.ParseIP(tt.ipaddr))
		if s != tt.binval {
			t.Errorf("[%d] want %s, got %s", i, tt.binval, s)
		}
	}
}

func TestIP6ToHexString(t *testing.T) {
	for i, tt := range IP6Tests {
		s := IPToHexString(net.ParseIP(tt.ipaddr))
		if s != tt.hexval {
			t.Errorf("[%d] want %s got %s", i, tt.hexval, s)
		}
	}
}

func TestBigintToIP6(t *testing.T) {
	for i, tt := range IP6Tests {
		z := big.Int{}
		z.SetString(tt.bigintval, 10)
		ip := BigintToIP6(&z)
		x := CompareIPs(ip, net.ParseIP(tt.ipaddr))
		if x != 0 {
			t.Errorf("[%d] want %s got %s", i, tt.ipaddr, ip)
		}
	}
}

func TestExpandIP6(t *testing.T) {
	for i, tt := range IP6Tests {
		s := ExpandIP6(net.ParseIP(tt.ipaddr))
		if s != tt.expand {
			t.Errorf("[%d] want %s got %s", i, tt.expand, s)
		}
	}
}

func TestIP6ToARPA(t *testing.T) {
	for i, tt := range IP6Tests {
		s := IPToARPA(net.ParseIP(tt.ipaddr))
		if s != tt.inarpa {
			t.Errorf("[%d] want %s got %s", i, tt.inarpa, s)
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
	for i, tt := range IPDeltaTests {
		z := DeltaIP(tt.ipaddr, tt.incr)
		if z != tt.incres {
			t.Errorf("[%d] on increment: want %d got %d", i, tt.incres, z)
		}

		z = DeltaIP(tt.ipaddr, tt.decr)
		if z != tt.decres {
			t.Errorf("[%d]on decrement: want %d got %d", i, tt.decres, z)
		}
	}
}

func TestDeltaIPMaxValue(t *testing.T) {
	i := DeltaIP(net.ParseIP("2001:db8::"), net.ParseIP("2001:db8:1234:5678::"))
	if i != MaxIPv4 {
		t.Errorf("want %d got %d", MaxIPv4, i)
	}
}

func TestDecrementIPBy(t *testing.T) {
	for i, tt := range IPDeltaTests {
		ip := DecrementIPBy(tt.ipaddr, tt.intval)
		x := CompareIPs(ip, tt.decr)
		if x != 0 {
			t.Errorf("[%d] want %s got %s", i, tt.decr, ip)
		}
	}
}

func TestIncrementIPBy(t *testing.T) {
	for i, tt := range IPDeltaTests {
		ip := IncrementIPBy(tt.ipaddr, tt.intval)
		x := CompareIPs(ip, tt.incr)
		if x != 0 {
			t.Errorf("[%d] want %s got %s", i, tt.incr, ip)
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
	for i, tt := range IPDelta6Tests {
		z := DeltaIP6(tt.ipaddr, tt.incr)
		if z.String() != tt.incres {
			t.Errorf("[%d] on increment: want %s got %s", i, tt.incres, z)
		}

		z = DeltaIP6(tt.ipaddr, tt.decr)
		if z.String() != tt.decres {
			t.Errorf("[%d] on decrement: want %s got %s", i, tt.decres, z)
		}
	}
}

func TestDecrementIP6By(t *testing.T) {
	for i, tt := range IPDelta6Tests {
		z := big.Int{}
		z.SetString(tt.intval, 10)
		ip := DecrementIP6By(tt.ipaddr, &z)
		x := CompareIPs(ip, tt.decr)
		if x != 0 {
			t.Errorf("[%d] want %s got %s", i, tt.decr, ip)
		}
	}
}

func TestIncrementIP6By(t *testing.T) {
	for i, tt := range IPDelta6Tests {
		z := big.Int{}
		z.SetString(tt.intval, 10)
		ip := IncrementIP6By(tt.ipaddr, &z)
		x := CompareIPs(ip, tt.incr)
		if x != 0 {
			t.Errorf("[%d] want %s got %s", i, tt.incr, ip)
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
	{net.IP{0, 0, 0, 0}, 4, 4},
	{net.IP{192, 168, 1, 1}, 4, 4},
	{net.IP{255, 255, 255, 255}, 4, 4},
	{net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 6, 6},
	{net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52}, 6, 6},
	{net.IP{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, 6, 6},
	// these are the 4-in-6 versions of the first 3 test cases
	{net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 0, 0, 0, 0}, 6, 4},
	{net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 192, 168, 1, 1}, 6, 4},
	{net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255}, 6, 4},
}

func Test_Version(t *testing.T) {
	for i, tt := range IPVersionTests {
		version := Version(tt.ipaddr)
		if version != tt.version {
			t.Errorf("[%d] want %d got %d", i, tt.version, version)
		}
	}
}

func Test_EffectiveVersion(t *testing.T) {
	for i, tt := range IPVersionTests {
		eversion := EffectiveVersion(tt.ipaddr)
		if eversion != tt.eversion {
			t.Errorf("[%d] want %d got %d", i, tt.eversion, eversion)
		}
	}
}

func Test_EffectiveVersionNil(t *testing.T) {
	eversion := EffectiveVersion(nil)
	if eversion != 0 {
		t.Errorf("want 0, got %d", eversion)
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
	for i, b := range compareIPTests {
		a1 = append(a1, b.ipaddr)
		val := CompareIPs(a.ipaddr, b.ipaddr)
		if val != b.status {
			t.Errorf("[%d] want %d got %d", i, b.status, val)
		}
	}
	sort.Sort(ByIP(a1))
	for i, b := range compareIPTests {
		if a1[b.pos].String() != b.ipaddr.String() {
			t.Errorf("[%d] want %s at position %d got %s", i, b.ipaddr, b.pos, a1[b.pos])
		}
	}
}

var isAllTests = []struct {
	ipaddr net.IP
	isones bool
	iszero bool
	is4in6 bool
}{
	{net.IP{0, 0, 0, 0}, false, true, false},
	{net.IP{255, 255, 255, 255}, true, false, false},
	{net.IP{192, 168, 1, 1}, false, false, false},
	{net.ParseIP("::ffff:0:0"), false, true, true},
	{net.ParseIP("::ffff:ffff:ffff"), true, false, true},
	{net.ParseIP("::ffff:c0a8:0101"), false, false, true},
	{net.ParseIP("::"), false, true, false},
	{net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"), true, false, false},
	{net.ParseIP("2001:db8::1"), false, false, false},
}

func TestIs4in6(t *testing.T) {
	for i, tt := range isAllTests {
		v := Is4in6(tt.ipaddr)
		if v != tt.is4in6 {
			t.Errorf("[%d] want %t got %t", i, tt.is4in6, v)
		}
	}
}

func TestIsAllOnes(t *testing.T) {
	for i, tt := range isAllTests {
		v := IsAllOnes(tt.ipaddr)
		if v != tt.isones {
			t.Errorf("[%d] want %t got %t", i, tt.isones, v)
		}
	}
}

func TestIsAllZeroes(t *testing.T) {
	for i, tt := range isAllTests {
		v := IsAllZeroes(tt.ipaddr)
		if v != tt.iszero {
			t.Errorf("[%d] want %t got %t", i, tt.iszero, v)
		}
	}
}
