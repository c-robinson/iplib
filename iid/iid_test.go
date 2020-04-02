package iid

import (
	"crypto"
	_ "crypto/sha512"
	"net"
	"testing"

	"github.com/c-robinson/iplib"
)

var RFC7217AddrTests = []struct {
	netid   string
	secret  string
	counter int64
	htype   crypto.Hash
	scope   Scope
	out     string
	err     error
}{
	{
		"01234567",
		"secret",
		1,
		crypto.SHA256,
		ScopeGlobal,
		"2001:db8::c6fa:ba02:41ab:282c",
		nil,
	},
	{
		"01234567",
		"secret",
		1,
		crypto.SHA384,
		ScopeGlobal,
		"2001:db8::51b3:c6b0:4e14:3519",
		nil,
	},
	{
		"76543210",
		"secret",
		1,
		crypto.SHA384,
		ScopeGlobal,
		"2001:db8::703d:9ce9:741a:80f1",
		nil,
	},
	{
		"01234567",
		"terces",
		1,
		crypto.SHA384,
		ScopeGlobal,
		"2001:db8::97dd:7cba:9f02:c412",
		nil,
	},
	{
		"01234567",
		"secret",
		1,
		crypto.SHA384,
		ScopeLocal,
		"2001:db8::51b3:c6b0:4e14:3519",
		nil,
	},
	{
		"01234567",
		"secret",
		2,
		crypto.SHA384,
		ScopeGlobal,
		"2001:db8::606a:57c0:dacf:706",
		nil,
	},
}

func TestGenerateRFC7217Addr(t *testing.T) {
	ip := net.ParseIP("2001:db8::")
	hw, _ := net.ParseMAC("77:88:99:aa:bb:cc")
	for i, tt := range RFC7217AddrTests {
		out, err := GenerateRFC7217Addr(ip, hw, tt.counter, []byte(tt.netid), []byte(tt.secret), tt.htype, tt.scope)
		if tt.err == nil && err != nil {
			t.Errorf("[%d] got unexpected error: %s", i, err.Error())
		} else if tt.err != nil && err == nil {
			t.Errorf("[%d] expected error, got none", i)
		} else {
			ttout := net.ParseIP(tt.out)
			v := iplib.CompareIPs(ttout, out)
			if v != 0 {
				t.Errorf("[%d] wrong address. Expected '%s' got '%s'", i, ttout, out)
			}
		}
	}
}

var IPTests = []struct {
	name    string
	address string
	res     bool
	rfc     string
}{
	{
		"Broken",
		"192.168.1.1",
		false,
		"",
	},
	{
		"NotReserved",
		"25:100:200::0200:5EFF:FF00:521A",
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
		"aaaa:bbbb:cccc:dddd:0200:5EFF:FE00:5211",
		true,
		"RFC4291",
	},
	{
		"ReservedProxyMobile",
		"aaaa:bbbb:cccc:dddd:0200:5EFF:FE00:5213",
		true,
		"RFC6543",
	},
}

func TestGetReservationsForIP(t *testing.T) {
	for _, tt := range IPTests {
		ip := net.ParseIP(tt.address)
		r := GetReservationsForIP(ip)
		if tt.res == false {
			if r != nil {
				t.Errorf("%s: expected no results, got '%s'", tt.name, r.Title)
			}
		} else {
			if r == nil {
				t.Errorf("%s: got no result but one was expected", tt.name)
			} else {
				if r.RFC != tt.rfc {
					t.Errorf("%s got wrong reservation, expected '%s' got %s", tt.name, tt.rfc, r.RFC)
				}
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
		"192.168.1.1",
		"bb:aa:cc:dd:ee:ff",
		"",
		"",
		"",
		"",
	},
	{
		"2001:db8:1111:2222::",
		"bb:aa:cc:dd:ee",
		"",
		"",
		"",
		"",
	},
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
		"2001:db8:1111:2222:b9aa:ccdd:ddcc:aabb",
		"2001:db8:1111:2222:bbaa:ccdd:ddcc:aabb",
	},
	{
		"2001:db8:1111:2222::",
		"99:88:77:66:55:44:33:22",
		"2001:db8:1111:2222:9b88:7766:5544:3322",
		"2001:db8:1111:2222:9988:7766:5544:3322",
		"2001:db8:1111:2222:9b88:7766:5544:3322",
		"2001:db8:1111:2222:9988:7766:5544:3322",
	},
}

func TestMakeEUI64Addr(t *testing.T) {
	for i, tt := range EUI64Tests {
		inaddr := net.ParseIP(tt.inaddr)
		hwaddr, _ := net.ParseMAC(tt.hwaddr)
		if iplib.EffectiveVersion(inaddr) == 4 || len(hwaddr) < 4 {
			out := MakeEUI64Addr(inaddr, hwaddr, ScopeGlobal)
			if out != nil {
				t.Errorf("[%d] expected <nil> got '%s'", i, out)
			}
			continue
		}

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

var OpaqueAddrTests = []struct {
	netid   string
	secret  string
	counter int64
	out     string
}{
	{
		"01234567",
		"secret",
		1,
		"2001:db8::c6fa:ba02:41ab:282c",
	},
	{
		"76543210",
		"secret",
		1,
		"2001:db8::8377:6cc:7e00:a088",
	},
	{
		"",
		"secret",
		1,
		"2001:db8::f67e:7072:5572:d4bc",
	},
	{
		"01234567",
		"terces",
		1,
		"2001:db8::5a42:5d26:73bc:28a6",
	},
	{
		"01234567",
		"secret",
		2,
		"2001:db8::ab77:a9d1:2391:5994",
	},
}

func TestMakeOpaqueAddr(t *testing.T) {
	ip := net.ParseIP("2001:db8::")
	hw, _ := net.ParseMAC("77:88:99:aa:bb:cc")
	for i, tt := range OpaqueAddrTests {
		out, err := MakeOpaqueAddr(ip, hw, tt.counter, []byte(tt.netid), []byte(tt.secret))

		if err != nil {
			t.Errorf("[%d] got  unexpected error: %s", i, err)
		}

		ttout := net.ParseIP(tt.out)
		v := iplib.CompareIPs(ttout, out)
		if v != 0 {
			t.Errorf("[%d] wrong address. Expected '%s' got '%s'", i, ttout, out)
		}
	}
}
