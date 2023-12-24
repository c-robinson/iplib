package iplib

import (
	"bytes"
	"net"
	"testing"

	"lukechampine.com/uint128"
)

var hostMaskTests = []struct {
	masklen int
	mask    net.IPMask
	bpos    int
	bvalue  byte
}{
	{0, net.IPMask{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, -1, 0x00},
	{1, net.IPMask{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x80}, 15, 0x80},
	{2, net.IPMask{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xc0}, 15, 0xc0},
	{3, net.IPMask{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xe0}, 15, 0xe0},
	{4, net.IPMask{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf0}, 15, 0xf0},
	{5, net.IPMask{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf8}, 15, 0xf8},
	{6, net.IPMask{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xfc}, 15, 0xfc},
	{7, net.IPMask{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xfe}, 15, 0xfe},
	{8, net.IPMask{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff}, 15, 0xff},
	{16, net.IPMask{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff}, 14, 0xff},
	{32, net.IPMask{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff}, 12, 0xff},
	{64, net.IPMask{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, 8, 0xff},
	{58, net.IPMask{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xc0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, 8, 0xc0},
	{127, net.IPMask{0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, 0, 0xfe},
	{128, net.IPMask{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, 0, 0xff},
}

func TestNewHostMask(t *testing.T) {
	for i, tt := range hostMaskTests {
		mask := NewHostMask(tt.masklen)
		v := bytes.Compare(mask, tt.mask)
		if v != 0 {
			t.Errorf("[%d] got wrong mask value for masklen %d: want %v got %v", i, tt.masklen, tt.mask, mask)
		}

		masklen, _ := mask.Size()
		if masklen != tt.masklen {
			t.Errorf("[%d] got wrong mask length: want %d got %d", i, tt.masklen, masklen)
		}

		bvalue, bpos := mask.BoundaryByte()
		if bvalue != tt.bvalue {
			t.Errorf("[%d] got wrong boundary byte value: want '%x' got '%x'", i, tt.bvalue, bvalue)
		}
		if bpos != tt.bpos {
			t.Errorf("[%d] got wrong boundary byte position: want %d got %d", i, tt.bpos, bpos)
		}
	}
}

var boundaryByteDeltaTests = []struct {
	bb        byte
	bv        byte
	count     uint128.Uint128
	decrval   byte
	decrcount uint128.Uint128
	incrval   byte
	incrcount uint128.Uint128
}{
	{0x00, 0x00, uint128.From64(0), 0x00, uint128.From64(0), 0x00, uint128.From64(0)},
	{0x00, 0x00, uint128.From64(1), 0xff, uint128.From64(1), 0x01, uint128.From64(0)},
	{0x00, 0x01, uint128.From64(0), 0x01, uint128.From64(0), 0x01, uint128.From64(0)},
	{0x00, 0x01, uint128.From64(1), 0x00, uint128.From64(0), 0x02, uint128.From64(0)},
	{0x00, 0x09, uint128.From64(1024), 0x09, uint128.From64(4), 0x09, uint128.From64(4)},
	{0x80, 0x09, uint128.From64(1024), 0x09, uint128.From64(8), 0x09, uint128.From64(8)},
	{0x40, 0x09, uint128.From64(1024), 0x89, uint128.From64(6), 0x49, uint128.From64(5)},
	{0x20, 0x09, uint128.From64(1024), 0x69, uint128.From64(5), 0x89, uint128.From64(4)},
	{0x10, 0x09, uint128.From64(1024), 0xb9, uint128.From64(5), 0x49, uint128.From64(4)},
	{0x08, 0x09, uint128.From64(1024), 0xe1, uint128.From64(5), 0x29, uint128.From64(4)},
}

func Test_decrementBoundaryByte(t *testing.T) {
	for i, tt := range boundaryByteDeltaTests {
		decrcount, decrval := decrementBoundaryByte(tt.bb, tt.bv, tt.count)
		if decrval != tt.decrval {
			t.Errorf("[%d] got wrong output byte: want 0x%02x, got 0x%02x", i, tt.decrval, decrval)
		}
		if v := decrcount.Cmp(tt.decrcount); v != 0 {
			t.Errorf("[%d] got wrong output count: want %d, got %d", i, tt.decrcount, decrcount)
		}
	}
}

func Test_incrementBoundaryByte(t *testing.T) {
	for i, tt := range boundaryByteDeltaTests {
		count, incrval := incrementBoundaryByte(tt.bb, tt.bv, tt.count)
		if incrval != tt.incrval {
			t.Errorf("[%d] got wrong output byte: want 0x%02x, got 0x%02x", i, tt.incrval, incrval)
		}
		if v := count.Cmp(tt.incrcount); v != 0 {
			t.Errorf("[%d] got wrong output count: want %d, got %d", i, tt.incrcount, count)
		}
	}
}

var unmaskedBytesDeltaTest = []struct {
	inval   []byte
	incount uint128.Uint128
	decrval []byte
	incrval []byte
}{
	{[]byte{0, 255}, uint128.From64(255), []byte{0, 0}, []byte{1, 254}},
	{[]byte{254, 1}, uint128.From64(1), []byte{254, 0}, []byte{254, 2}},
	{[]byte{0, 255}, uint128.From64(1), []byte{0, 254}, []byte{1, 0}},
}

func Test_decrementUnmaskedBytes(t *testing.T) {
	for i, tt := range unmaskedBytesDeltaTest {
		decrval := decrementUnmaskedBytes(tt.inval, tt.incount)
		if v := bytes.Compare(tt.decrval, decrval); v != 0 {
			t.Errorf("[%d] got wrong output array: want %+v, got %+v", i, tt.decrval, decrval)
		}
	}
}

func Test_incrementUnmaskedBytes(t *testing.T) {
	for i, tt := range unmaskedBytesDeltaTest {
		incrval := incrementUnmaskedBytes(tt.inval, tt.incount)
		if v := bytes.Compare(tt.incrval, incrval); v != 0 {
			t.Errorf("[%d] got wrong output array: want %+v, got %+v", i, tt.incrval, incrval)
		}
	}
}

var IPHostmaskDeltaTests = []struct {
	ipaddr   net.IP
	hostmask int
	decr     net.IP
	decrErr  error
	incr     net.IP
	incrErr  error
	prev     net.IP
	prevErr  error
	next     net.IP
	nextErr  error
}{
	{ // 0
		net.ParseIP("2001:db8:1234:5678::"), 0,
		net.ParseIP("2001:db8:1234:5677:ffff:ffff:ffff:fc18"), nil,
		net.ParseIP("2001:db8:1234:5678::3e8"), nil,
		net.ParseIP("2001:db8:1234:5677:ffff:ffff:ffff:ffff"), nil,
		net.ParseIP("2001:db8:1234:5678::1"), nil,
	}, { // 1
		net.ParseIP("2001:db8:1234:5678:9900::"), 56,
		net.ParseIP("2001:db8:1234:5674:b100::"), nil,
		net.ParseIP("2001:db8:1234:567c:8100::"), nil,
		net.ParseIP("2001:db8:1234:5678:9800::"), nil,
		net.ParseIP("2001:db8:1234:5678:9a00::"), nil,
	}, { // 2
		net.ParseIP("2001:db8:1234:5678:ff00::"), 56,
		net.ParseIP("2001:db8:1234:5675:1700::"), nil,
		net.ParseIP("2001:db8:1234:567c:e700::"), nil,
		net.ParseIP("2001:db8:1234:5678:fe00::"), nil,
		net.ParseIP("2001:db8:1234:5679::"), nil,
	}, { // 3
		net.ParseIP("::"), 56,
		net.ParseIP(""), ErrAddressOutOfRange,
		net.ParseIP("::3:e800:0:0:0"), nil,
		net.IP{}, ErrAddressOutOfRange,
		net.ParseIP("::100:0:0:0"), nil,
	}, { // 4
		net.ParseIP("ffff:ffff:ffff:ffff:ff00::"), 56,
		net.ParseIP("ffff:ffff:ffff:fffc:1700::"), nil,
		net.IP{}, ErrAddressOutOfRange,
		net.ParseIP("ffff:ffff:ffff:ffff:fe00::"), nil,
		net.IP{}, ErrAddressOutOfRange,
	}, { // 5
		net.ParseIP("2001:db8:1234:5678:9906::"), 53,
		net.ParseIP("2001:db8:1234:5678:1c06::"), nil,
		net.ParseIP("2001:db8:1234:5679:1606::"), nil,
		net.ParseIP("2001:db8:1234:5678:9905::"), nil,
		net.ParseIP("2001:db8:1234:5678:9907::"), nil,
	}, { // 6
		net.ParseIP("2001:db8:1234:5678:9907::"), 53,
		net.ParseIP("2001:db8:1234:5678:1c07::"), nil,
		net.ParseIP("2001:db8:1234:5679:1607::"), nil,
		net.ParseIP("2001:db8:1234:5678:9906::"), nil,
		net.ParseIP("2001:db8:1234:5678:9a00::"), nil,
	}, { // 7
		net.ParseIP("2001:db8:1234:5678:9908::"), 53,
		net.ParseIP(""), ErrAddressOutOfRange,
		net.ParseIP("2001:db8:1234:5679:1700::"), nil,
		net.ParseIP("2001:db8:1234:5678:9907::"), nil,
		net.ParseIP("2001:db8:1234:5678:9909::"), nil,
	}, { // 8
		net.ParseIP("2001:db8:1234:5678:ff::"), 56,
		net.ParseIP(""), ErrAddressOutOfRange,
		net.ParseIP("2001:db8:1234:567c:e700::"), nil,
		net.ParseIP(""), ErrAddressOutOfRange,
		net.ParseIP(""), ErrAddressOutOfRange,
	},
}

func TestDecrementIP6WithinHostmask(t *testing.T) {
	for i, tt := range IPHostmaskDeltaTests {
		count := uint128.From64(1000)
		hm := NewHostMask(tt.hostmask)
		decr, err := DecrementIP6WithinHostmask(tt.ipaddr, hm, count)
		if e := compareErrors(err, tt.decrErr); len(e) > 0 {
			t.Errorf("[%d] %s (%s)", i, e, decr)
		} else {
			x := CompareIPs(decr, tt.decr)
			if x != 0 {
				t.Errorf("[%d] expected %s got %s", i, tt.decr, decr)
			}
		}
	}
}

func TestIncrementIP6WithinHostmask(t *testing.T) {
	for i, tt := range IPHostmaskDeltaTests {
		count := uint128.From64(1000)
		hm := NewHostMask(tt.hostmask)
		incr, err := IncrementIP6WithinHostmask(tt.ipaddr, hm, count)
		if e := compareErrors(err, tt.incrErr); len(e) > 0 {
			t.Errorf("[%d] %s (%s)", i, e, incr)
		} else {
			x := CompareIPs(incr, tt.incr)
			if x != 0 {
				t.Errorf("[%d] expected %s got %s", i, tt.incr, incr)
			}
		}
	}
}

func TestNextIPWithinHostmask(t *testing.T) {
	for i, tt := range IPHostmaskDeltaTests {
		next, err := NextIP6WithinHostmask(tt.ipaddr, NewHostMask(tt.hostmask))
		if e := compareErrors(err, tt.nextErr); len(e) > 0 {
			t.Errorf("[%d] %s (%s)", i, e, next)
		} else {
			x := CompareIPs(next, tt.next)
			if x != 0 {
				t.Errorf("[%d] expected %s got %s", i, tt.next, next)
			}
		}
	}
}

func TestPreviousIPWithinHostmask(t *testing.T) {
	for i, tt := range IPHostmaskDeltaTests {
		prev, err := PreviousIP6WithinHostmask(tt.ipaddr, NewHostMask(tt.hostmask))
		if e := compareErrors(err, tt.prevErr); len(e) > 0 {
			t.Errorf("[%d] %s (%s)", i, e, prev)
		} else {
			x := CompareIPs(prev, tt.prev)
			if x != 0 {
				t.Errorf("[%d] expected %s got %s", i, tt.prev, prev)
			}
		}
	}
}

func compareErrors(got, want error) string {
	if got == nil && want == nil {
		return ""
	}
	if got == nil && want != nil {
		return "wanted error, but got none"
	}
	if got != nil && want == nil {
		return "got unexpected error: " + got.Error()
	}
	if got.Error() != want.Error() {
		return "got wrong error: " + got.Error()
	}
	return ""
}
