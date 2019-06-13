package iid

import (
	"net"
	"testing"
)

var IPTests = []struct {
	name     string
	address  string
	res      bool
}{

	{
		"NotReserved",
		"25:100:200::195:16",
		false,
	},
	{
		"ReservedAnycast",
		"::",
		true,
	},
	{
		"ReservedEthernet",
		"0200:5EFF:FE00:521a:AAAA:BBBB:CCCC:DDDD",
		true,
	},
}

func TestGetReservationsForIP(t *testing.T) {
	for _, tt := range IPTests {
		ip := net.ParseIP(tt.address)
		r := GetReservationsForIP(ip)
		if len(r) > 0 && tt.res && len(r) == 0 && tt.res == true {
			t.Errorf("%s returned wrong reservations status", tt.name)
		}
	}
}
