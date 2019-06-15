package iid

import (
	"net"
	"testing"

	"github.com/c-robinson/iplib"
)

var IPTests = []struct {
	name     string
	address  string
	res      bool
	rfc      string
}{

	{
		"NotReserved",
		"25:100:200::195:16",
		false,
		"",
	},
	{
		"ReservedAnycast",
		"::",
		true,
		"RFC4291",
	},
	{
		"ReservedEthernet",
		"0200:5EFF:FE00:521A:AAAA:BBBB:CCCC:DDDD",
		true,
		"RFC4291",
	},
	{
		"ReservedProxyMobile",
		"0200:5EFF:FE00:5213:1234::",
		true,
		"RFC6543",
	},
}

func TestGetReservationsForIP(t *testing.T) {
	for _, tt := range IPTests {
		ip := net.ParseIP(tt.address)
		r := GetReservationsForIP(ip)
		if len(r) > 0 && tt.res && len(r) == 0 && tt.res == true {
			t.Errorf("%s returned wrong reservations status", tt.name)
			if r[0].RFC != tt.rfc {
				t.Errorf("%s returned wrong reservation RFC, expected '%s' got '%s'", tt.name, tt.rfc, r[0].RFC)
			}
		}
	}
}

var EUI64Tests = []struct {
	inaddr    string
	hwaddr    string
	outGlobal string
	outLocal  string
	outInvert string
	outNone   string
}{
	{
		"2001:db8:1111:2222::",
		"bb:aa:cc:dd:ee:ff",
		"2001:db8:1111:2222:bbaa:ccff:fedd:eeff",
		"2001:db8:1111:2222:b9aa:ccff:fedd:eeff",
		"2001:db8:1111:2222:b9aa:ccff:fedd:eeff",
		"2001:db8:1111:2222:bbaa:ccff:fedd:eeff",
	},
	{
		"2001:db8:1111:2222::",
		"99:88:77:66:55:44",
		"2001:db8:1111:2222:9b88:77ff:fe66:5544",
		"2001:db8:1111:2222:9988:77ff:fe66:5544",
		"2001:db8:1111:2222:9b88:77ff:fe66:5544",
		"2001:db8:1111:2222:9988:77ff:fe66:5544",
	},
	{
		"2001:db8:1111:2222::",
		"bb:aa:cc:dd:dd:cc:aa:bb",
		"2001:db8:1111:2222:bbaa:ccdd:ddcc:aabb",
		"2001:db8:1111:2222:b9aa:ccdd:ddcc:aabb",
		"2001:db8:1111:2222:b9aa:ccdd:ddcc:aabb", // this one is 2001:db8:1111:2222:bbaa:ccdd:ddcc:aabb
		"2001:db8:1111:2222:bbaa:ccdd:ddcc:aabb",
	},
	{
		"2001:db8:1111:2222::",
		"99:88:77:66:55:44:33:22",
		"2001:db8:1111:2222:9b88:7766:5544:3322",
		"2001:db8:1111:2222:9988:7766:5544:3322",
		"2001:db8:1111:2222:9b88:7766:5544:3322",
		"2001:db8:1111:2222:9988:7766:5544:3322", // this one is 2001:db8:1111:2222:9b88:7766:5544:3322
	},
}

func TestMakeEUI64Addr(t *testing.T) {
	for i, tt := range EUI64Tests {
		inaddr := net.ParseIP(tt.inaddr)
		hwaddr, _ := net.ParseMAC(tt.hwaddr)

		out := MakeEUI64Addr(inaddr, hwaddr, ScopeGlobal)
		if val := iplib.CompareIPs(out, net.ParseIP(tt.outGlobal)); val != 0 {
			t.Errorf("[%d] '%s' outGlobal: expected %s got %s", i, tt.hwaddr, tt.outGlobal, out)
		}

		out = MakeEUI64Addr(inaddr, hwaddr, ScopeLocal)
		if val := iplib.CompareIPs(out, net.ParseIP(tt.outLocal)); val != 0 {
			t.Errorf("[%d] '%s' outLocal: expected %s got %s", i, tt.hwaddr, tt.outLocal, out)
		}

		out = MakeEUI64Addr(inaddr, hwaddr, ScopeInvert)
		if val := iplib.CompareIPs(out, net.ParseIP(tt.outInvert)); val != 0 {
			t.Errorf("[%d] '%s' outInvert: expected %s got %s", i, tt.hwaddr, tt.outInvert, out)
		}

		out = MakeEUI64Addr(inaddr, hwaddr, ScopeNone)
		if val := iplib.CompareIPs(out, net.ParseIP(tt.outNone)); val != 0 {
			t.Errorf("[%d] '%s' outNone: expected %s got %s", i, tt.hwaddr, tt.outNone, out)
		}
	}
}
