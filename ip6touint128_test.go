package iplib

import (
	"net"
	"testing"

	"lukechampine.com/uint128"
)

func TestIP6ToUint128WithNonIPv6Input(t *testing.T) {
	cases := []struct {
		name string
		ip   net.IP
	}{
		{"nil", nil},
		{"ipv4 4-byte", net.IP{192, 0, 2, 1}},
		{"ipv4 16-byte", net.ParseIP("192.0.2.1")},
		{"ipv4-mapped", net.ParseIP("::ffff:192.0.2.1")},
		{"truncated", net.IP{0x20, 0x01, 0x0d}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := IP6ToUint128(c.ip); !got.Equals(uint128.Uint128{}) {
				t.Fatalf("IP6ToUint128(%v) = %s, want the zero value", c.ip, got)
			}
		})
	}
}

func TestIP6ToUint128KeepsIPv6Values(t *testing.T) {
	ip := net.ParseIP("2001:db8:85a3::8a2e:370:7334")
	want := uint128.FromBytesBE(ip.To16())
	if got := IP6ToUint128(ip); !got.Equals(want) {
		t.Fatalf("IP6ToUint128(%v) = %s, want %s", ip, got, want)
	}
}

func TestDeltaIP6WithNilInput(t *testing.T) {
	ip := net.ParseIP("2001:db8::1")
	if got := DeltaIP6(nil, ip); !got.Equals(IP6ToUint128(ip)) {
		t.Fatalf("DeltaIP6(nil, %v) = %s", ip, got)
	}
}
