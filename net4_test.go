package iplib

import (
	"net"
	"sort"
	"testing"
)

var NewNet4Tests = []struct {
	s       string
	addr    net.IP
	masklen int
	is4in6  bool
	isEmpty bool
}{
	{"192.168.0.0/16", ForceIP4(net.ParseIP("192.168.0.0")), 16, false, false},
	{"192.168.0.0/16", net.ParseIP("::ffff:c0a8:0000"), 16, true, false},
	{"", net.ParseIP("192.168.0.0"), 33, false, true},
}

func TestNewNet4(t *testing.T) {
	for i, tt := range NewNet4Tests {
		ipn := NewNet4(tt.addr, tt.masklen)
		if (tt.isEmpty == true && ipn.IP() != nil) || (tt.isEmpty == false && ipn.IP() == nil) {
			t.Errorf("[%d] expect isEmpty == %t, but is not", i, tt.isEmpty)
		} else if tt.isEmpty == false {
			if ipn.String() != tt.s {
				t.Errorf("[%d] Net4 want %s got %s", i, tt.s, ipn.String())
			}
			if ipn.is4in6 != tt.is4in6 {
				t.Errorf("[%d] is4in6 want %t got %t", i, tt.is4in6, ipn.is4in6)
			}
		}
	}
}

var Net4FromStrTests = []struct {
	ins     string
	outs    string
	isEmpty bool
}{
	{"192.168.0.0/16", "192.168.0.0/16", false},
	{"notanaddress!!", "", true},
	{"::ffff:c0a8:0000/16", "192.168.0.0/16", false},
	{"2001:db8::/32", "", true},
}

func TestNet4FromStr(t *testing.T) {
	for i, tt := range Net4FromStrTests {
		ipn := Net4FromStr(tt.ins)
		if (tt.isEmpty == true && ipn.IP() != nil) || (tt.isEmpty == false && ipn.IP() == nil) {
			t.Errorf("[%d] expect isEmpty == %t, but is not", i, tt.isEmpty)
		} else if tt.isEmpty == false {
			if tt.outs != ipn.String() {
				t.Errorf("[%d] want %s got %s", i, tt.outs, ipn.String())
			}
		}
	}
}

var Net4Tests = []struct {
	ip        net.IP
	network   net.IP
	netmask   net.IPMask
	wildcard  net.IPMask
	broadcast net.IP
	firstaddr net.IP
	lastaddr  net.IP
	masklen   int
	count     uint32
}{
	{
		net.ParseIP("10.1.2.3"),
		net.ParseIP("10.0.0.0"),
		net.IPMask{255, 0, 0, 0},
		net.IPMask{0, 255, 255, 255},
		net.ParseIP("10.255.255.255"),
		net.ParseIP("10.0.0.1"),
		net.ParseIP("10.255.255.254"),
		8, 16777214,
	},
	{
		net.ParseIP("192.168.1.1"),
		net.ParseIP("192.168.0.0"),
		net.IPMask{255, 255, 254, 0},
		net.IPMask{0, 0, 1, 255},
		net.ParseIP("192.168.1.255"),
		net.ParseIP("192.168.0.1"),
		net.ParseIP("192.168.1.254"),
		23, 510,
	},
	{
		net.ParseIP("192.168.1.61"),
		net.ParseIP("192.168.1.0"),
		net.IPMask{255, 255, 255, 192},
		net.IPMask{0, 0, 0, 63},
		net.ParseIP("192.168.1.63"),
		net.ParseIP("192.168.1.1"),
		net.ParseIP("192.168.1.62"),
		26, 62,
	},
	{
		net.ParseIP("192.168.1.66"),
		net.ParseIP("192.168.1.64"),
		net.IPMask{255, 255, 255, 192},
		net.IPMask{0, 0, 0, 63},
		net.ParseIP("192.168.1.127"),
		net.ParseIP("192.168.1.65"),
		net.ParseIP("192.168.1.126"),
		26, 62,
	},
	{
		net.ParseIP("192.168.1.1"),
		net.ParseIP("192.168.1.0"),
		net.IPMask{255, 255, 255, 252},
		net.IPMask{0, 0, 0, 3},
		net.ParseIP("192.168.1.3"),
		net.ParseIP("192.168.1.1"),
		net.ParseIP("192.168.1.2"),
		30, 2,
	},
	{
		net.ParseIP("192.168.1.1"),
		net.ParseIP("192.168.1.0"),
		net.IPMask{255, 255, 255, 254},
		net.IPMask{0, 0, 0, 1},
		net.ParseIP("192.168.1.1"),
		net.ParseIP("192.168.1.0"),
		net.ParseIP("192.168.1.1"),
		31, 2,
	},
	{
		net.ParseIP("192.168.1.15"),
		net.ParseIP("192.168.1.15"),
		net.IPMask{255, 255, 255, 255},
		net.IPMask{0, 0, 0, 0},
		net.ParseIP("192.168.1.15"),
		net.ParseIP("192.168.1.15"),
		net.ParseIP("192.168.1.15"),
		32, 1,
	},
}

func TestNet4_BroadcastAddress(t *testing.T) {
	for i, tt := range Net4Tests {
		ipn := NewNet4(tt.ip, tt.masklen)
		if addr := ipn.BroadcastAddress(); !tt.broadcast.Equal(addr) {
			t.Errorf("[%d] want %v got %v", i, tt.broadcast, addr)
		}
	}
}

func TestNet4_Version(t *testing.T) {
	for i, tt := range Net4Tests {
		ipn := NewNet4(tt.ip, tt.masklen)
		if ipn.Version() != IP4Version {
			t.Errorf("[%d] want version 4, got %d", i, ipn.Version())
		}
	}
}

func TestNet4_Count(t *testing.T) {
	for i, tt := range Net4Tests {
		ipn := NewNet4(tt.ip, tt.masklen)
		if ipn.Count() != tt.count {
			t.Errorf("[%d] want %d got %d", i, tt.count, ipn.Count())
		}
	}
}

func TestNet4_FirstAddress(t *testing.T) {
	for i, tt := range Net4Tests {
		ipn := NewNet4(tt.ip, tt.masklen)
		if addr := ipn.FirstAddress(); !tt.firstaddr.Equal(addr) {
			t.Errorf("[%d] want %s got %s", i, tt.firstaddr, addr)
		}
	}
}

func TestNet4_finalAddress(t *testing.T) {
	for i, tt := range Net4Tests {
		ipn := NewNet4(tt.ip, tt.masklen)
		if addr, ones := ipn.finalAddress(); !tt.broadcast.Equal(addr) {
			t.Errorf("[%d] want %s got %s (%d))", i, tt.broadcast, addr, ones)
		}
	}
}

func TestNet4_LastAddress(t *testing.T) {
	for i, tt := range Net4Tests {
		ipn := NewNet4(tt.ip, tt.masklen)
		if addr := ipn.LastAddress(); !tt.lastaddr.Equal(addr) {
			t.Errorf("[%d] want %s got %s", i, tt.lastaddr, addr)
		}
	}
}

func TestNet4_NetworkAddress(t *testing.T) {
	for i, tt := range Net4Tests {
		ipn := NewNet4(tt.ip, tt.masklen)
		if addr := ipn.IP(); !tt.network.Equal(addr) {
			t.Errorf("[%d] want %s got %s", i, tt.network, addr)
		}
	}
}

func TestWildcard(t *testing.T) {
	for i, tt := range Net4Tests {
		ipn := NewNet4(tt.ip, tt.masklen)
		if ipn.Wildcard().String() != tt.wildcard.String() {
			t.Errorf("[%d] want %s got %s", i, tt.wildcard, ipn.Wildcard())
		}
	}
}

var enumerate4Tests = []struct {
	incidr string
	total  int
	last   net.IP
}{
	{"192.168.0.0/22", 1022, net.ParseIP("192.168.3.254")},
	{"192.168.0.0/23", 510, net.ParseIP("192.168.1.254")},
	{"192.168.0.0/24", 254, net.ParseIP("192.168.0.254")},
	{"192.168.0.0/25", 126, net.ParseIP("192.168.0.126")},
	{"192.168.0.0/26", 62, net.ParseIP("192.168.0.62")},
	{"192.168.0.0/27", 30, net.ParseIP("192.168.0.30")},
	{"192.168.0.0/28", 14, net.ParseIP("192.168.0.14")},
	{"192.168.0.0/29", 6, net.ParseIP("192.168.0.6")},
	{"192.168.0.0/30", 2, net.ParseIP("192.168.0.2")},
	{"192.168.0.0/31", 2, net.ParseIP("192.168.0.1")},
	{"192.168.0.0/32", 1, net.ParseIP("192.168.0.0")},
}

func TestNet4_Enumerate(t *testing.T) {
	for i, tt := range enumerate4Tests {
		_, ipn, _ := ParseCIDR(tt.incidr)
		ipn4 := ipn.(Net4)
		addrlist := ipn4.Enumerate(0, 0)
		if len(addrlist) != tt.total {
			t.Errorf("[%d] want size %d got %d", i, tt.total, len(addrlist))
		}
		x := CompareIPs(tt.last, addrlist[tt.total-1])
		if x != 0 {
			t.Errorf("[%d] want last address %s, got %s", i, tt.last, addrlist[tt.total-1])
		}
	}
}

var enumerate4VariableTests = []struct {
	offset int
	size   int
	total  int
	first  net.IP
	last   net.IP
}{
	{
		0, 0, 1022,
		net.ParseIP("192.168.0.1"),
		net.ParseIP("192.168.3.254"),
	},
	{
		1, 0, 1021,
		net.ParseIP("192.168.0.2"),
		net.ParseIP("192.168.3.254"),
	},
	{
		256, 0, 766,
		net.ParseIP("192.168.1.1"),
		net.ParseIP("192.168.3.254"),
	},
	{
		0, 128, 128,
		net.ParseIP("192.168.0.1"),
		net.ParseIP("192.168.0.128"),
	},
	{
		20, 128, 128,
		net.ParseIP("192.168.0.21"),
		net.ParseIP("192.168.0.148"),
	},
	{
		1000, 100, 22,
		net.ParseIP("192.168.3.233"),
		net.ParseIP("192.168.3.254"),
	},
	{
		1023, 0, 0,
		net.ParseIP("192.168.3.233"),
		net.ParseIP("192.168.3.254"),
	},
}

func TestNet4_EnumerateWithVariables(t *testing.T) {
	_, ipn, _ := ParseCIDR("192.168.0.0/22")
	ipn4 := ipn.(Net4)
	for i, tt := range enumerate4VariableTests {
		addrlist := ipn4.Enumerate(tt.size, tt.offset)
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

func TestNet4_EnumerateRFC3021(t *testing.T) {
	ipn := NewNet4(net.ParseIP("192.168.1.0"), 31)
	addrlist := ipn.Enumerate(0, 0)
	if len(addrlist) != 2 {
		t.Errorf("want 2, got %d", len(addrlist))
	}
}

var incr4Tests = []struct {
	inaddr   string
	thisaddr net.IP
	nextaddr net.IP
	err      error
}{
	{
		"192.168.1.0/23",
		net.ParseIP("192.168.1.0"),
		net.ParseIP("192.168.1.1"),
		nil,
	},
	{
		"192.168.1.0/24",
		net.ParseIP("192.168.1.254"),
		net.ParseIP("192.168.1.255"),
		ErrBroadcastAddress,
	},
	{
		"192.168.2.0/24",
		net.ParseIP("192.168.2.1"),
		net.ParseIP("192.168.2.2"),
		nil,
	},
	{
		"192.168.3.0/24",
		net.ParseIP("192.168.3.0"),
		net.ParseIP("192.168.3.1"),
		nil,
	},
	{
		"192.168.4.0/24",
		net.ParseIP("192.168.5.1"),
		net.IP{},
		ErrAddressOutOfRange,
	},
	{
		"192.168.1.0/31",
		net.ParseIP("192.168.1.0"),
		net.ParseIP("192.168.1.1"),
		ErrBroadcastAddress,
	},
	{
		"192.168.1.0/32",
		net.ParseIP("192.168.1.0"),
		net.IP{},
		ErrAddressOutOfRange,
	},
}

func TestNet4_NextIP(t *testing.T) {
	for i, tt := range incr4Tests {
		_, ipn, _ := ParseCIDR(tt.inaddr)
		ipn4 := ipn.(Net4)
		nextaddr, err := ipn4.NextIP(tt.thisaddr)
		if e := compareErrors(err, tt.err); len(e) > 0 {
			t.Errorf("[%d] %s (%s)", i, e, tt.thisaddr)
		} else {
			if !nextaddr.Equal(tt.nextaddr) {
				t.Errorf("For %s expected %v, got %v", tt.thisaddr, tt.nextaddr, nextaddr)
			}
		}
	}
}

var incr4SubnetTests = []struct {
	netblock Net4
	netmask  int
	next     Net4
}{
	{Net4FromStr("192.168.0.0/24"), 24, Net4FromStr("192.168.1.0/24")},
	{Net4FromStr("192.168.0.0/24"), 25, Net4FromStr("192.168.1.0/25")},
	{Net4FromStr("192.168.0.0/24"), 23, Net4FromStr("192.168.0.0/23")},
	{Net4FromStr("255.255.255.0/24"), 24, Net4FromStr("255.255.255.0/24")},
}

func TestNet4_NextNet(t *testing.T) {
	for i, tt := range incr4SubnetTests {
		next := tt.netblock.NextNet(tt.netmask)
		if v := CompareNets(next, tt.next); v != 0 {
			t.Errorf("[%d] want %v got %v", i, tt.next, next)
		}
	}
}

var decr4Tests = []struct {
	inaddr   string
	thisaddr net.IP
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
		ErrAddressOutOfRange,
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
		ErrAddressOutOfRange,
	},
}

func TestNet4_PreviousIP(t *testing.T) {
	for i, tt := range decr4Tests {
		_, ipn, _ := ParseCIDR(tt.inaddr)
		ipn4 := ipn.(Net4)
		prevaddr, err := ipn4.PreviousIP(tt.thisaddr)
		if e := compareErrors(err, tt.preverr); len(e) > 0 {
			t.Errorf("[%d] %s (%s)", i, e, prevaddr)
		} else {
			if !prevaddr.Equal(tt.prevaddr) {
				t.Errorf("[%d] want %s, got %s", i, tt.prevaddr, prevaddr)
			}
		}
	}
}

var decr4SubnetTests = []struct {
	netblock Net4
	netmask  int
	prev     Net4
}{
	{Net4FromStr("192.168.1.0/24"), 24, Net4FromStr("192.168.0.0/24")},
	{Net4FromStr("192.168.1.0/24"), 25, Net4FromStr("192.168.0.128/25")},
	{Net4FromStr("192.168.1.0/24"), 23, Net4FromStr("192.168.0.0/23")},
	{Net4FromStr("0.0.0.0/24"), 24, Net4FromStr("0.0.0.0/24")},
}

func TestNet4_PreviousNet(t *testing.T) {
	for i, tt := range decr4SubnetTests {
		prev := tt.netblock.PreviousNet(tt.netmask)
		if v := CompareNets(prev, tt.prev); v != 0 {
			t.Errorf("[%d] want %s got %s", i, tt.prev, prev)
		}
	}
}

var subnet4Tests = []struct {
	netblock Net4
	netmask  int
	subnets  []string
	err      error
}{
	{
		Net4FromStr("192.168.0.0/24"), 0,
		[]string{"192.168.0.0/25", "192.168.0.128/25"},
		nil,
	},
	{
		Net4FromStr("192.168.0.0/24"), 25,
		[]string{"192.168.0.0/25", "192.168.0.128/25"},
		nil,
	},
	{
		Net4FromStr("192.168.0.0/24"), 26,
		[]string{"192.168.0.0/26", "192.168.0.64/26", "192.168.0.128/26", "192.168.0.192/26"},
		nil,
	},
	{
		Net4FromStr("192.168.0.0/24"), 23,
		[]string{},
		ErrBadMaskLength,
	},
	{
		Net4FromStr("192.168.0.0/32"), 0,
		[]string{},
		ErrBadMaskLength,
	},
}

func TestNet4_Subnet(t *testing.T) {
	for i, tt := range subnet4Tests {
		subnets, err := tt.netblock.Subnet(tt.netmask)
		if e := compareErrors(err, tt.err); len(e) > 0 {
			t.Errorf("[%d] %s", i, e)
		} else {
			v := compareNet4ArraysToStringRepresentation(subnets, tt.subnets)
			if v == false {
				t.Errorf("[%d] want len %d got %d: %v", i, len(tt.subnets), len(subnets), subnets)
			}
		}
	}
}

var supernet4Tests = []struct {
	in      Net4
	masklen int
	out     Net4
	err     error
}{
	{
		Net4FromStr("192.168.1.0/24"), 23, Net4FromStr("192.168.0.0/23"),
		nil,
	},
	{
		Net4FromStr("192.168.1.0/24"), 0, Net4FromStr("192.168.0.0/23"),
		nil,
	},
	{
		Net4FromStr("192.168.1.0/24"), 22, Net4FromStr("192.168.1.0/22"),
		nil,
	},
	{
		Net4FromStr("192.168.1.4/30"), 24, Net4FromStr("192.168.1.0/24"),
		nil,
	},
	{
		Net4FromStr("192.168.0.0/24"), 25, Net4{},
		ErrBadMaskLength,
	},
}

func TestNet4_Supernet(t *testing.T) {
	for i, tt := range supernet4Tests {
		out, err := tt.in.Supernet(tt.masklen)
		if e := compareErrors(err, tt.err); len(e) > 0 {
			t.Errorf("[%d] %s", i, e)
		} else {
			if v := CompareNets(out, tt.out); v != 0 {
				t.Errorf("[%d] want %s got %s", i, tt.out, out)
			}
		}
	}
}

func TestCompareNets(t *testing.T) {
	net4map := map[int]Net4{
		0: Net4FromStr("192.168.0.0/16"),
		1: Net4FromStr("192.168.0.0/23"),
		2: Net4FromStr("192.168.1.0/24"),
		3: Net4FromStr("192.168.1.0/24"),
		4: Net4FromStr("192.168.3.0/26"),
		5: Net4FromStr("192.168.3.64/26"),
		6: Net4FromStr("192.168.3.128/25"),
		7: Net4FromStr("192.168.4.0/24"),
	}

	net4list := ByNet{}
	for _, ipn := range net4map {
		net4list = append(net4list, ipn)
	}
	sort.Sort(ByNet(net4list))
	for pos, ipn := range net4map {
		if v := CompareNets(net4list[pos], ipn); v != 0 {
			for i, aipn := range net4list {
				if v := CompareNets(ipn, aipn); v == 0 {
					t.Errorf("subnet %s want position %d got %d", ipn, pos, i)
					break
				}
			}
		}
	}
}

var containsNet4Tests = []struct {
	ipn1   Net4
	ipn2   Net4
	result bool
}{
	{Net4FromStr("192.168.0.0/16"), Net4FromStr("192.168.45.0/24"), true},
	{Net4FromStr("192.168.45.0/24"), Net4FromStr("192.168.45.0/26"), true},
	{Net4FromStr("192.168.45.0/24"), Net4FromStr("192.168.46.0/26"), false},
	{Net4FromStr("10.1.1.1/24"), Net4FromStr("10.0.0.0/8"), false},
}

func TestNet4_ContainsNet(t *testing.T) {
	for i, tt := range containsNet4Tests {
		result := tt.ipn1.ContainsNet(tt.ipn2)
		if result != tt.result {
			t.Errorf("[%d] want %t got %t", i, tt.result, result)
		}
	}
}

func TestNet4_Is4in6(t *testing.T) {
	nf := Net4FromStr("192.168.0.0./16")
	//nf := NewNet4(ForceIP4(net.ParseIP("192.168.0.0")), 16)
	if nf.Is4in6() != false {
		t.Errorf("should be false for '192.168.0.0/16'")
	}
	nt := NewNet4(net.ParseIP("::ffff:c0a8:0000"), 16)
	if nt.Is4in6() != true {
		t.Errorf("should be true for '::ffff:c0a8:0000/16'")
	}
}

func compareNet4ArraysToStringRepresentation(a []Net4, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, n := range a {
		if n.String() != b[i] {
			return false
		}
	}

	return true
}
