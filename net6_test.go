package iplib

import (
	"net"
	"sort"
	"testing"
)

var NewNet6Tests = []struct {
	s           string
	addr        net.IP
	netmasklen  int
	hostmasklen int
	isEmpty     bool
}{
	{"2001:db8::/32", net.ParseIP("2001:db8::"), 32, 0, false},
	{"2001:db8::/32", net.ParseIP("2001:db8::"), 32, 16, false},
	{"", net.ParseIP("2001:db8::"), 33, 96, true},
	{"", net.ParseIP("2001:db8::"), 112, 17, true},
	{"2001:db8::/112", net.ParseIP("2001:db8::"), 112, 15, false},
	{"2001:db8::/127", net.ParseIP("2001:db8::"), 127, 0, false},
	{"2001:db8::/128", net.ParseIP("2001:db8::"), 128, 0, false},
}

func TestNewNet6(t *testing.T) {
	for i, tt := range NewNet6Tests {
		ipn := NewNet6(tt.addr, tt.netmasklen, tt.hostmasklen)
		if (tt.isEmpty == true && ipn.IP() != nil) || (tt.isEmpty == false && ipn.IP() == nil) {
			t.Errorf("[%d] expect isEmpty == %t, but is not", i, tt.isEmpty)
		} else if tt.isEmpty == false {
			if ipn.String() != tt.s {
				t.Errorf("[%d] Net6 want %s got %s", i, tt.s, ipn.String())
			}
			hostmasklen, _ := ipn.Hostmask.Size()
			netmasklen, _ := ipn.Mask().Size()
			if tt.hostmasklen != hostmasklen {
				t.Errorf("[%d] want hostmask size %d got %d", i, tt.hostmasklen, hostmasklen)
			}
			if tt.netmasklen != netmasklen {
				t.Errorf("[%d] want netmask size %d got %d", i, tt.netmasklen, netmasklen)
			}
		}
	}
}

var Net6FromStrTests = []struct {
	ins     string
	outs    string
	isEmpty bool
}{
	{"2001:db8::/64", "2001:db8::/64", false},
	{"notanaddress!!", "", true},
	{"::ffff:c0a8:0000/16", "", true},
}

func TestNet6FromStr(t *testing.T) {
	for i, tt := range Net6FromStrTests {
		ipn := Net6FromStr(tt.ins)
		if (tt.isEmpty == true && ipn.IP() != nil) || (tt.isEmpty == false && ipn.IP() == nil) {
			t.Errorf("[%d] expect isEmpty == %t, but is not", i, tt.isEmpty)
		} else if tt.isEmpty == false {
			if tt.outs != ipn.String() {
				t.Errorf("[%d] want %s got %s", i, tt.outs, ipn.String())
			}
		}
	}
}

var Net6Tests = []struct {
	ip          string
	firstaddr   string
	lastaddr    string
	hostmask    int
	hostmaskpos int
	netmasklen  int
	count       string // converted to uint128.Uint128
}{
	// 0-7 hostmask applied on byte boundaries
	{
		"2001:0db8:1234:5678:9abc:def0:1234:5678",
		"2001:db8:1234:5678:9a00::",
		"2001:0db8:1234:5678:9aff:ffff:ffff:ffff",
		0, -1, 72, "72057594037927936",
	},
	{
		"2001:0db8:1234:5678:9abc:def0:1234:5678",
		"2001:0db8:1234:5678::",
		"2001:0db8:1234:5678:ffff:ffff:ffff:ff00",
		8, 15, 64, "72057594037927936",
	},
	{
		"2001:0db8:1234:5678:9abc:def0:1234:5678",
		"2001:0db8:1234:5678::",
		"2001:0db8:1234:5678:ffff:ffff:ffff:0",
		16, 14, 64, "281474976710656",
	},
	{
		"2001:0db8:1234:5678:9abc:def0:1234:5678",
		"2001:0db8:1234:5678::",
		"2001:0db8:1234:5678:ffff:ffff:ff00:0",
		24, 13, 64, "1099511627776",
	},
	{
		"2001:0db8:1234:5678:9abc:def0:1234:5678",
		"2001:0db8:1234:5678::",
		"2001:0db8:1234:5678:ffff:ffff::",
		32, 12, 64, "4294967296",
	},
	{
		"2001:0db8:1234:5678:9abc:def0:1234:5678",
		"2001:0db8:1234:5678::",
		"2001:0db8:1234:5678:ffff:ff00::",
		40, 11, 64, "16777216",
	},
	{
		"2001:0db8:1234:5678:9abc:def0:1234:5678",
		"2001:0db8:1234:5678::",
		"2001:0db8:1234:5678:ffff::",
		48, 10, 64, "65536",
	},
	{
		"2001:0db8:1234:5678:9abc:def0:1234:5678",
		"2001:0db8:1234:5678::",
		"2001:0db8:1234:5678:ff00::",
		56, 9, 64, "256",
	},

	// 7-15: hostmask applied within a byte
	{
		"2001:0db8:1234:5678:9abc:def0:1234:5678",
		"2001:0db8:1234:5678::",
		"2001:0db8:1234:5678:7f00::",
		57, 8, 64, "128",
	},
	{
		"2001:0db8:1234:5678:9abc:def0:1234:5678",
		"2001:0db8:1234:5678::",
		"2001:0db8:1234:5678:3f00::",
		58, 8, 64, "64",
	},
	{
		"2001:0db8:1234:5678:9abc:def0:1234:5678",
		"2001:0db8:1234:5678::",
		"2001:0db8:1234:5678:1f00::",
		59, 8, 64, "32",
	},
	{
		"2001:0db8:1234:5678:9abc:def0:1234:5678",
		"2001:0db8:1234:5678::",
		"2001:0db8:1234:5678:0f00::",
		60, 8, 64, "16",
	},
	{
		"2001:0db8:1234:5678:9abc:def0:1234:5678",
		"2001:0db8:1234:5678::",
		"2001:0db8:1234:5678:0700::",
		61, 8, 64, "8",
	},
	{
		"2001:0db8:1234:5678:9abc:def0:1234:5678",
		"2001:0db8:1234:5678::",
		"2001:0db8:1234:5678:0300::",
		62, 8, 64, "4",
	},
	{
		"2001:0db8:1234:5678:9abc:def0:1234:5678",
		"2001:0db8:1234:5678::",
		"2001:0db8:1234:5678:0100::",
		63, 8, 64, "2",
	},
	{
		"2001:0db8:1234:5678:9abc:def0:1234:5678",
		"2001:0db8:1234:5678::",
		"",
		64, -1, 64, "0",
	}, // entire address masked

	// 16-19: hostmask and netmask applied on same byte
	{
		"2001:0db8:1234:5678:9abc:def0::",
		"2001:db8:1234:5678:9abc:8000::",
		"2001:0db8:1234:5678:9abc:ff00::",
		40, 11, 81, "128",
	},
	{
		"2001:0db8:1234:5678:9abc:def0::",
		"2001:db8:1234:5678:9abc:c000::",
		"2001:0db8:1234:5678:9abc:7f00::",
		41, 10, 82, "32",
	},
	{
		"2001:0db8:1234:5678:9abc:def0::",
		"2001:db8:1234:5678:9abc:c000::",
		"2001:0db8:1234:5678:9abc:1f00::",
		42, 10, 83, "8",
	},
	{
		"2001:0db8:1234:5678:9abc:def0::",
		"2001:db8:1234:5678:9abc:d000::",
		"2001:0db8:1234:5678:9abc:ff00::",
		43, 10, 84, "2",
	},
	// 20-21, address with /128 and /127 respectively
	{
		"2001:0db8:1234:5678:9abc:def0:1234:5678",
		"2001:db8:1234:5678:9abc:def0:1234:5678",
		"2001:0db8:1234:5678:9abc:def0:1234:5678",
		0, -1, 128, "1",
	},
	{
		"2001:0db8:1234:5678:9abc:def0:1234:5678",
		"2001:db8:1234:5678:9abc:def0:1234:5678",
		"2001:0db8:1234:5678:9abc:def0:1234:5679",
		0, -1, 127, "2",
	},
	// 22, address with no netmask or hostmask
	{
		"2001:0db8::",
		"::",
		"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
		0, -1, 0, "340282366920938463463374607431768211455",
	},
	{
		"::",
		"::",
		"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
		0, -1, 0, "340282366920938463463374607431768211455",
	},
}

func TestNet6_Version(t *testing.T) {
	for i, tt := range Net6Tests {
		ipn := NewNet6(net.ParseIP(tt.ip), tt.netmasklen, tt.hostmask)
		if ipn.Version() != IP6Version {
			t.Errorf("[%d] want version 6, got %d", i, ipn.Version())
		}
	}
}

func TestNet6_Count(t *testing.T) {
	for i, tt := range Net6Tests {
		ipn := NewNet6(net.ParseIP(tt.ip), tt.netmasklen, tt.hostmask)

		if ipn.IPNet.IP == nil {
			if tt.count != "0" {
				t.Fatalf("[%d] produced nil Net6{}, but should not have", i)
			}
			continue
		}

		if tt.count != ipn.Count().String() {
			t.Errorf("[%d] count: want %s got %s", i, tt.count, ipn.Count().String())
		}
	}
}

func TestNet6_FirstAddress(t *testing.T) {
	for i, tt := range Net6Tests {
		firstAddr := net.ParseIP(tt.firstaddr)
		ipn := NewNet6(net.ParseIP(tt.ip), tt.netmasklen, tt.hostmask)

		if ipn.IPNet.IP == nil {
			if tt.count != "0" {
				t.Fatalf("[%d] produced nil Net6{}, but should not have", i)
			}
			continue
		}

		if v := CompareIPs(firstAddr, ipn.IP()); v != 0 {
			t.Errorf("[%d] network address: want %s got %s", i, firstAddr, ipn.IP())
		}

		if v := CompareIPs(firstAddr, ipn.FirstAddress()); v != 0 {
			t.Errorf("[%d] first address: want %s got %s", i, firstAddr, ipn.FirstAddress())
		}
	}
}

func TestNet6_LastAddress(t *testing.T) {
	for i, tt := range Net6Tests {
		lastAddr := net.ParseIP(tt.lastaddr)
		ipn := NewNet6(net.ParseIP(tt.ip), tt.netmasklen, tt.hostmask)

		la := ipn.LastAddress()

		if v := CompareIPs(lastAddr, la); v != 0 {
			t.Errorf("[%d] last address: want %s got %s", i, lastAddr, la)
		}
	}
}

func TestNet6_BoundaryByte(t *testing.T) {
	for i, tt := range Net6Tests {
		ipn := NewNet6(net.ParseIP(tt.ip), tt.netmasklen, tt.hostmask)
		_, bpos := ipn.Hostmask.BoundaryByte()
		if bpos != tt.hostmaskpos {
			t.Errorf("[%d] boundary position: want %d got %d", i, tt.hostmaskpos, bpos)
		}
	}
}

func TestNewNet6WrongVersion(t *testing.T) {
	n := NewNet6(ForceIP4(net.ParseIP("10.0.0.0")), 8, 0)
	if v := CompareIPs(n.IP(), nil); v != 0 {
		t.Errorf("Expected empty Net6, got %s", n.IP())
	}
}

var enumerate6Tests = []struct {
	inaddr      net.IP
	hostmasklen int
	netmasklen  int
	total       int
	last        net.IP
}{
	{ // one element network returns itself
		net.ParseIP("2001:db8:1000:2000:3000:4000::"),
		0, 128, 1,
		net.ParseIP("2001:db8:1000:2000:3000:4000::"),
	},
	{ // no address in list
		net.ParseIP("2001:db8:1000:2000:3000:4000::"),
		64, 64, 0,
		net.ParseIP("::"),
	},
	{ // RFC6164
		net.ParseIP("2001:db8:1000:2000:3000:4000::"),
		0, 127, 2,
		net.ParseIP("2001:db8:1000:2000:3000:4000::1"),
	},
	{
		net.ParseIP("2001:db8:1000:2000:3000:4000::"),
		48, 64, 65536,
		net.ParseIP("2001:db8:1000:2000:ffff::"),
	},
}

func TestNet6_Enumerate(t *testing.T) {
	for i, tt := range enumerate6Tests {
		n := NewNet6(tt.inaddr, tt.netmasklen, tt.hostmasklen)
		addrlist := n.Enumerate(0, 0)
		if len(addrlist) != tt.total {
			t.Errorf("[%d] total want %d got %d", i, tt.total, len(addrlist))
		}
		if len(addrlist) > 0 {
			if v := CompareIPs(addrlist[len(addrlist)-1], tt.last); v != 0 {
				t.Errorf("[%d] last address: want %s got %s", i, tt.last, addrlist[len(addrlist)-1])
			}
		}
		for ii, a := range addrlist {
			if a == nil {
				t.Errorf("[%d] address %d is nil", i, ii)
			}
		}
	}
}

var enumerate6VariableTests = []struct {
	hostmasklen int
	netmasklen  int
	offset      int
	size        int
	total       int
	first       net.IP
	last        net.IP
}{
	{ // no offset, enumerate entire block
		56, 56, 0, 0, 65536,
		net.ParseIP("2001:db8:1000:2000::"),
		net.ParseIP("2001:db8:1000:20ff:ff00::"),
	},
	{ // enumerate the entire back half
		56, 56, 32768, 0, 32768,
		net.ParseIP("2001:db8:1000:2080::"),
		net.ParseIP("2001:db8:1000:20ff:ff00::"),
	},
	{ // enumerate half of the back half
		56, 56, 32768, 16384, 16384,
		net.ParseIP("2001:db8:1000:2080::"),
		net.ParseIP("2001:db8:1000:20bf:ff00::"),
	},
	{ // enumerate past the boundary
		56, 56, 65000, 5000, 536,
		net.ParseIP("2001:db8:1000:20fd:e800::"),
		net.ParseIP("2001:db8:1000:20ff:ff00::"),
	},
	{ // enumerate starting after the boundary
		56, 56, 65537, 16, 0,
		net.ParseIP("2001:db8:1000:2080::"),
		net.ParseIP("2001:db8:1000:20bf:ff00::"),
	},
}

func TestNet6_EnumerateWithVariables(t *testing.T) {
	ip := net.ParseIP("2001:db8:1000:2000:3000:4000::")
	for i, tt := range enumerate6VariableTests {
		n := NewNet6(ip, tt.netmasklen, tt.hostmasklen)
		addrlist := n.Enumerate(tt.size, tt.offset)
		if len(addrlist) != tt.total {
			t.Errorf("[%d] size: want %d got %d", i, tt.total, len(addrlist))
		}
		if len(addrlist) > 0 {
			x := CompareIPs(tt.first, addrlist[0])
			if x != 0 {
				t.Errorf("[%d] first member: want %s got %s", i, tt.first, addrlist[0])
			}
			y := CompareIPs(tt.last, addrlist[len(addrlist)-1])
			if y != 0 {
				t.Errorf("[%d] last member: want %s got %s", i, tt.last, addrlist[len(addrlist)-1])
			}
		}
	}
}

var incr6Tests = []struct {
	netmask  int
	hostmask int
	thisaddr net.IP
	nextaddr net.IP
	err      error
}{
	{ // address not in the netblock
		64, 0,
		net.ParseIP("2001:db8:123:4567::"),
		nil,
		ErrAddressOutOfRange,
	},
	{ // address outside of hostmask
		64, 56,
		net.ParseIP("2001:db8:1234:5678:ff::1"),
		nil,
		ErrAddressOutOfRange,
	},
	{ // increment from 1st to second address, no hostmask
		64, 0,
		net.ParseIP("2001:db8:1234:5678::"),
		net.ParseIP("2001:db8:1234:5678::1"),
		nil,
	},
	{ // increment from 1st to second address, with hostmask
		64, 56,
		net.ParseIP("2001:db8:1234:5678::"),
		net.ParseIP("2001:db8:1234:5678:100::"),
		nil,
	},
	{ // increment from last address, no hostmask
		64, 0,
		net.ParseIP("2001:db8:1234:5678:ffff:ffff:ffff:ffff"),
		nil,
		ErrAddressOutOfRange,
	},
	{ // increment from last address, no hostmask
		64, 56,
		net.ParseIP("2001:db8:1234:5678:ff00::"),
		nil,
		ErrAddressOutOfRange,
	},
}

func TestNet6_NextIP(t *testing.T) {
	netaddr := net.ParseIP("2001:db8:1234:5678::")
	for i, tt := range incr6Tests {
		ipn := NewNet6(netaddr, tt.netmask, tt.hostmask)
		nextaddr, err := ipn.NextIP(tt.thisaddr)
		if e := compareErrors(err, tt.err); len(e) > 0 {
			t.Errorf("[%d] %s", i, e)
		} else {
			if !nextaddr.Equal(tt.nextaddr) {
				t.Errorf("[%d] want %s got %s", i, tt.nextaddr, nextaddr)
			}
		}
	}
}

func TestNet6_NextIPBadStartAddress(t *testing.T) {
	ipn := NewNet6(net.ParseIP("2001:db8:1234:5678::"), 64, 56)
	ip, err := ipn.NextIP(net.ParseIP("2001:db8:1234:5678::12"))
	if e := compareErrors(err, ErrAddressOutOfRange); len(e) > 0 {
		t.Errorf("expected out-of-range error, got IP '%s', error '%s'", ip, err)
	}
	ip, err = ipn.NextIP(net.ParseIP("2001:db8:1234:5677::"))
	if e := compareErrors(err, ErrAddressOutOfRange); len(e) > 0 {
		t.Errorf("expected out-of-range error, got IP '%s', error '%s'", ip, err)
	}
}

var incr6SubnetTests = []struct {
	netmasklen int
	// hostmasklen int
	next Net6
}{
	{64, NewNet6(net.ParseIP("2001:db8:1234:5679::"), 64, 56)},
	{0, NewNet6(net.ParseIP("2001:db8:1234:5679::"), 64, 56)},
	{48, NewNet6(net.ParseIP("2001:db8:1235::"), 48, 56)},
	{1, NewNet6(net.ParseIP("8000::"), 1, 56)},
	{72, Net6{}}, // netmask + hostnask == 128
}

func TestNet6_NextNet(t *testing.T) {
	ipn := NewNet6(net.ParseIP("2001:db8:1234:5678::"), 64, 63)
	for i, tt := range incr6SubnetTests {
		next := ipn.NextNet(tt.netmasklen)
		if v := CompareNets(next, tt.next); v != 0 { // WHY IS COMPARE NOT WORKING?!
			t.Errorf("[%d] want %s got %s", i, tt.next, next)
		}
	}
}

var decr6Tests = []struct {
	netmask  int
	hostmask int
	thisaddr net.IP
	prevaddr net.IP
	err      error
}{
	{ // address not in the netblock
		64, 0,
		net.ParseIP("2001:db8:123:4567::"),
		nil,
		ErrAddressOutOfRange,
	},
	{ // address outside of hostmask
		64, 56,
		net.ParseIP("2001:db8:1234:5678:ff::10"),
		nil,
		ErrAddressOutOfRange,
	},
	{ // decrement from last address, no hostmask
		64, 0,
		net.ParseIP("2001:db8:1234:5678:ffff:ffff:ffff:ffff"),
		net.ParseIP("2001:db8:1234:5678:ffff:ffff:ffff:fffe"),
		nil,
	},
	{ // decrement from last address, with hostmask
		64, 56,
		net.ParseIP("2001:db8:1234:5678:ff00::"),
		net.ParseIP("2001:db8:1234:5678:fe00::"),
		nil,
	},
	{ // decrement from first address, no hostmask
		64, 0,
		net.ParseIP("2001:db8:1234:5678::"),
		nil,
		ErrAddressOutOfRange,
	},
	{ // decrement from first address, no hostmask
		64, 56,
		net.ParseIP("2001:db8:1234:5678::"),
		nil,
		ErrAddressOutOfRange,
	},
}

func TestNet6_PreviousIP(t *testing.T) {
	netaddr := net.ParseIP("2001:db8:1234:5678::")
	for i, tt := range decr6Tests {
		ipn := NewNet6(netaddr, tt.netmask, tt.hostmask)
		prevaddr, err := ipn.PreviousIP(tt.thisaddr)
		if e := compareErrors(err, tt.err); len(e) > 0 {
			t.Errorf("[%d] %s", i, e)
		} else {
			if !prevaddr.Equal(tt.prevaddr) {
				t.Errorf("[%d] wamt %s got %s", i, tt.prevaddr, prevaddr)
			}
		}
	}
}

func TestNet6_PreviousIPBadStartAddress(t *testing.T) {
	ipn := NewNet6(net.ParseIP("2001:db8:1234:5678::"), 64, 56)
	ip, err := ipn.PreviousIP(net.ParseIP("2001:db8:1234:5678::12"))
	if e := compareErrors(err, ErrAddressOutOfRange); len(e) > 0 {
		t.Errorf("expected out-of-range error, got IP '%s', error '%s'", ip, err)
	}
	ip, err = ipn.PreviousIP(net.ParseIP("2001:db8:1234:5677::"))
	if e := compareErrors(err, ErrAddressOutOfRange); len(e) > 0 {
		t.Errorf("expected out-of-range error, got IP '%s', error '%s'", ip, err)
	}
}

var decr6SubnetTests = []struct {
	netmasklen int
	prev       Net6
}{
	{64, NewNet6(net.ParseIP("2001:db8:1234:5677::"), 64, 56)},
	{0, NewNet6(net.ParseIP("2001:db8:1234:5677::"), 64, 56)},
	{48, NewNet6(net.ParseIP("2001:db8:1233::"), 48, 56)},
	{5, NewNet6(net.ParseIP("1800::"), 5, 56)},
	{72, Net6{}}, // netmask + hostnask == 128
}

func TestNet6_PreviousNet(t *testing.T) {
	ipn := NewNet6(net.ParseIP("2001:db8:1234:5678::"), 64, 56)
	for i, tt := range decr6SubnetTests {
		prev := ipn.PreviousNet(tt.netmasklen)
		if v := CompareNets(prev, tt.prev); v != 0 {
			t.Errorf("[%d] want %v got %v", i, tt.prev, prev)
		}
	}
}

var subnet6Tests = []struct {
	netmasklen  int
	hostmasklen int
	subnets     []Net6
	err         error
}{
	{
		0, 0,
		[]Net6{
			NewNet6(net.ParseIP("2001:db8:1234:5678::"), 65, 0),
			NewNet6(net.ParseIP("2001:db8:1234:5678:8000::"), 65, 0),
		},
		nil,
	},
	{
		68, 61,
		[]Net6{},
		ErrBadMaskLength,
	},
	{
		65, 0,
		[]Net6{
			NewNet6(net.ParseIP("2001:db8:1234:5678::"), 65, 0),
			NewNet6(net.ParseIP("2001:db8:1234:5678:8000::"), 65, 0),
		},
		nil,
	},
	{
		66, 0,
		[]Net6{
			NewNet6(net.ParseIP("2001:db8:1234:5678::"), 66, 0),
			NewNet6(net.ParseIP("2001:db8:1234:5678:4000::"), 66, 0),
			NewNet6(net.ParseIP("2001:db8:1234:5678:8000::"), 66, 0),
			NewNet6(net.ParseIP("2001:db8:1234:5678:c000::"), 66, 0),
		},
		nil,
	},
	{
		63, 0,
		nil,
		ErrBadMaskLength,
	},
}

func TestNet6_Subnet(t *testing.T) {
	for i, tt := range subnet6Tests {
		ipn := NewNet6(net.ParseIP("2001:db8:1234:5678::"), 64, tt.hostmasklen)
		subnets, err := ipn.Subnet(tt.netmasklen, tt.hostmasklen)
		if e := compareErrors(err, tt.err); len(e) > 0 {
			t.Errorf("[%d] %s", i, e)
		} else {
			if v := compareNet6Arrays(subnets, tt.subnets); v == false {
				t.Errorf("[%d] want len %d got %d: %v", i, len(tt.subnets), len(subnets), subnets)
			}
		}
	}
}

var supernet6Tests = []struct {
	in         Net6
	netmasklen int
	out        Net6
	err        error
}{
	{
		Net6FromStr("2001:db8:1234:5678::/64"),
		60,
		Net6FromStr("2001:db8:1234:5670::/60"),
		nil,
	},
	{
		Net6FromStr("2001:db8:1234:5671::/64"),
		63,
		Net6FromStr("2001:db8:1234:5670::/63"),
		nil,
	},
	{
		Net6FromStr("2001:db8:1234:5671::/64"),
		0,
		Net6FromStr("2001:db8:1234:5670::/63"),
		nil,
	},
	{
		Net6FromStr("2001:db8:1234:5678::/64"),
		65,
		Net6{},
		ErrBadMaskLength,
	},
	{
		Net6FromStr("::/0"),
		0,
		Net6{},
		nil,
	},
}

func TestNet6_Supernet(t *testing.T) {
	for i, tt := range supernet6Tests {
		out, err := tt.in.Supernet(tt.netmasklen, 0)
		if e := compareErrors(err, tt.err); len(e) > 0 {
			t.Errorf("[%d] %s", i, e)
		} else {
			if v := CompareNets(out, tt.out); v != 0 {
				t.Errorf("[%d] want %s got %s", i, tt.out, out)
			}
		}
	}
}

func TestCompareNets6(t *testing.T) {
	net6map := map[int]Net6{
		0: Net6FromStr("::/0"),
		1: Net6FromStr("::/128"),
		2: Net6FromStr("2001:db8::/96"),
		3: Net6FromStr("2001:db8:12::/64"),
		4: Net6FromStr("2001:db8:1234::/64"),
		5: Net6FromStr("2001:db8:1234:5678::/63"),
		6: Net6FromStr("2001:db8:1234:5678::/64"),
		7: Net6FromStr("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/128"),
	}

	net6list := ByNet{}
	for _, ipn := range net6map {
		net6list = append(net6list, ipn)
	}
	sort.Sort(ByNet(net6list))
	for pos, ipn := range net6map {
		if v := CompareNets(net6list[pos], ipn); v != 0 {
			for i, aipn := range net6list {
				if v := CompareNets(ipn, aipn); v == 0 {
					t.Errorf("subnet %s want position %d got %d", ipn, pos, i)
					break
				}
			}
		}
	}
}

var containsNet6Tests = []struct {
	netblock1 string
	netblock2 string
	result    bool
}{
	{"2001:db8:1000:2000::/64", "2001:db8:1000:2000:3000::/72", true},
	{"2001:db8:1000:2000::/64", "2001:db8:1000:2000:3000:4000:5000:6000/127", true},
	{"2001:db8:1000:2000:3000::/72", "2001:db8:1000:2000::/64", false},
	{"2001:db8:1000:2000::/64", "2001:db8:1000:3000::/64", false},
	{"2001:db8:1000:2000::/64", "2001:db8:1000:2000::/64", true},
}

func TestNet6_ContainsNet(t *testing.T) {
	for i, tt := range containsNet6Tests {
		_, ipn, _ := ParseCIDR(tt.netblock1)
		_, sub, _ := ParseCIDR(tt.netblock2)
		result := ipn.ContainsNet(sub)
		if result != tt.result {
			t.Errorf("[%d] For \"%s contains %s\" want %v got %v", i, tt.netblock1, tt.netblock2, tt.result, result)
		}
	}
}

func TestNet6_RandomIP(t *testing.T) {
	for i, tt := range containsNet6Tests {
		_, ipn, _ := ParseCIDR(tt.netblock1)
		rip := ipn.(Net6).RandomIP()
		if !ipn.Contains(rip) {
			t.Errorf("[%d] address %s not in %s", i, rip, ipn)
		}
	}
}

var controlsTests = []struct {
	ipn   Net6
	addrs map[string]bool
}{
	{
		NewNet6(net.ParseIP("2001:db8:1::"), 56, 64),
		map[string]bool{
			"2001:db8:1:1::":    true,
			"2001:db8:2::":      false,
			"2001:db8:1:ff:1::": false,
		},
	},
}

func TestNet6_Controls(t *testing.T) {
	for _, tt := range controlsTests {
		for ip, v := range tt.addrs {
			if tt.ipn.Controls(net.ParseIP(ip)) != v {
				t.Errorf("Net6 '%s' ip '%s' want %t got %t", tt.ipn, ip, v, tt.ipn.Controls(net.ParseIP(ip)))
			}
		}
	}
}

func compareNet6Arrays(a []Net6, b []Net6) bool {
	if len(a) != len(b) {
		return false
	}

	for i, n := range a {
		if v := CompareNets(n, b[i]); v != 0 {
			return false
		}
	}

	return true
}
