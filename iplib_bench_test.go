package iplib

import (
	"net"
	"testing"

	"lukechampine.com/uint128"
)

func BenchmarkParseCIDR4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseCIDR("10.0.0.0/24")
	}
}

func BenchmarkParseCIDR6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseCIDR("2001:db8::/98")
	}
}

func BenchmarkNewNet(b *testing.B) {
	xip := net.IP{10, 0, 0, 0}
	for i := 0; i < b.N; i++ {
		NewNet(xip, 24)
	}
}

func Benchmark_DeltaIP4(b *testing.B) {
	var xip = net.IP{10, 255, 255, 255}
	var zip = net.IP{192, 168, 23, 5}
	for i := 0; i < b.N; i++ {
		_ = DeltaIP4(xip, zip)
	}
}

func Benchmark_DeltaIP6(b *testing.B) {
	var xip = net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52}
	var zip = net.IP{32, 1, 13, 184, 133, 255, 0, 0, 0, 10, 0, 15, 0, 0, 19, 0}
	for i := 0; i < b.N; i++ {
		_ = DeltaIP6(xip, zip)
	}
}

func BenchmarkPreviousIP4(b *testing.B) {
	var xip = net.IP{10, 255, 255, 255}
	for i := 0; i < b.N; i++ {
		xip = PreviousIP(xip)
	}
}

func BenchmarkDecrementIP4By(b *testing.B) {
	var xip = net.IP{10, 255, 255, 255}
	for i := 0; i < b.N; i++ {
		xip = DecrementIP4By(xip, 1)
	}
}

func BenchmarkDecrementIPBy_v4(b *testing.B) {
	var xip = net.IP{10, 255, 255, 255}
	for i := 0; i < b.N; i++ {
		xip = DecrementIPBy(xip, 1)
	}
}

func BenchmarkPreviousIP6(b *testing.B) {
	var xip = net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52}
	for i := 0; i < b.N; i++ {
		xip = PreviousIP(xip)
	}
}

func BenchmarkDecrementIP6By(b *testing.B) {
	var xip = net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52}
	count := uint128.From64(1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		xip = DecrementIP6By(xip, count)
	}
}

func BenchmarkDecrementIPBy_v6(b *testing.B) {
	var xip = net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52}
	for i := 0; i < b.N; i++ {
		xip = DecrementIPBy(xip, 1)
	}
}

func BenchmarkNextIP4(b *testing.B) {
	var xip = net.IP{10, 0, 0, 0}
	for i := 0; i < b.N; i++ {
		xip = NextIP(xip)
	}
}

func BenchmarkIncrementIP4By(b *testing.B) {
	var xip = net.IP{10, 0, 0, 0}
	for i := 0; i < b.N; i++ {
		xip = IncrementIP4By(xip, 1)
	}
}

func BenchmarkIncrementIPBy_v4(b *testing.B) {
	var xip = net.IP{10, 0, 0, 0}
	for i := 0; i < b.N; i++ {
		xip = IncrementIPBy(xip, 1)
	}
}

func BenchmarkNextIP6(b *testing.B) {
	var xip = net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52}
	for i := 0; i < b.N; i++ {
		xip = NextIP(xip)
	}
}

func BenchmarkIncrementIP6By(b *testing.B) {
	var xip = net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52}
	count := uint128.From64(1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		xip = IncrementIP6By(xip, count)
	}
}

func BenchmarkIncrementIPBy_v6(b *testing.B) {
	var xip = net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52}
	for i := 0; i < b.N; i++ {
		xip = IncrementIPBy(xip, 1)
	}
}

func BenchmarkNet_Count4(b *testing.B) {
	_, n, _ := ParseCIDR("192.168.0.0/24")
	n4 := n.(Net4)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = n4.Count()
	}
}

func BenchmarkNet_Count6(b *testing.B) {
	_, n, _ := ParseCIDR("2001:db8::/98")
	n6 := n.(Net6)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = n6.Count()
	}
}

func BenchmarkNet_Subnet_v4(b *testing.B) {
	_, n, _ := ParseCIDR("192.168.0.0/24")
	n4 := n.(Net4)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = n4.Subnet(25)
	}
}

func BenchmarkNet_Subnet_v6(b *testing.B) {
	_, n, _ := ParseCIDR("2001:db8::/98")
	n6 := n.(Net6)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = n6.Subnet(99, 0)
	}
<<<<<<< HEAD
=======

>>>>>>> 8809338 (Change from *big.Int to uint128.Uint128)
}

func BenchmarkNet_PreviousNet_v4(b *testing.B) {
	_, n, _ := ParseCIDR("192.168.0.0/24")
	n4 := n.(Net4)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = n4.PreviousNet(24)
	}
}

func BenchmarkNet_PreviousNet_v6(b *testing.B) {
	_, n, _ := ParseCIDR("2001:db8::/98")
	n6 := n.(Net6)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = n6.PreviousNet(24)
	}
}

func BenchmarkNet_NextNet_v4(b *testing.B) {
	_, n, _ := ParseCIDR("192.168.0.0/24")
	n4 := n.(Net4)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = n4.NextNet(24)
	}
}

func BenchmarkNet_NextNet_v6(b *testing.B) {
	_, n, _ := ParseCIDR("2001:db8::/98")
	n6 := n.(Net6)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = n6.NextNet(24)
	}
}

func BenchmarkNewNetBetween_v4(b *testing.B) {
	ipa := net.IP{10, 0, 0, 0}
	ipb := net.IP{10, 0, 0, 255}
	for i := 0; i < b.N; i++ {
		_, _, _ = NewNetBetween(ipa, ipb)
	}
}

func BenchmarkNewNetBetween_v6(b *testing.B) {
	ipa, _, _ := net.ParseCIDR("::")
	ipb, _, _ := net.ParseCIDR("ffff::")
	for i := 0; i < b.N; i++ {
		_, _, _ = NewNetBetween(ipa, ipb)
	}
}

func BenchmarkNet6_DecrementIP6WithinHostmask(b *testing.B) {
	var xip = net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52}
	count := uint128.From64(1)
	hm := NewHostMask(8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		xip, _ = DecrementIP6WithinHostmask(xip, hm, count)
	}
}

func BenchmarkNet6_PreviousIPWithinHostmask(b *testing.B) {
	var xip = net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52}
	hm := NewHostMask(8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		xip, _ = PreviousIP6WithinHostmask(xip, hm)
	}
}

func BenchmarkNet6_IncrementIP6WithinHostmask(b *testing.B) {
	var xip = net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52}
	count := uint128.From64(1)
	hm := NewHostMask(8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		xip, _ = DecrementIP6WithinHostmask(xip, hm, count)
	}
}

func BenchmarkNet6_NextIPWithinHostmask(b *testing.B) {
	var xip = net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52}
	hm := NewHostMask(8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		xip, _ = NextIP6WithinHostmask(xip, hm)
	}
<<<<<<< HEAD
=======
}

func BenchmarkNet6_IncrementIP6WithinHostmask(b *testing.B) {
	var xip = net.IP{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52}
	count := uint128.From64(1)
	hm := NewHostMask(8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		xip, _ = IncrementIP6WithinHostmask(xip, hm, count)
	}
>>>>>>> 8809338 (Change from *big.Int to uint128.Uint128)
}
