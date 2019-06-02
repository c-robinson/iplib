package iplib

import (
	"math/big"
	"net"
	"testing"
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
	count := big.NewInt(1)
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
	count := big.NewInt(1)
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
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		 _ = n.Count()
	}
}

func BenchmarkNet_Count6(b *testing.B) {
	_, n, _ := ParseCIDR("2001:db8::/98")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = n.Count()
	}
}

func BenchmarkNet_Subnet_v4(b *testing.B) {
	_, n, _ := ParseCIDR("192.168.0.0/24")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = n.Subnet(25)
	}
}

func BenchmarkNet_Subnet_v6(b *testing.B) {
	_, n, _ := ParseCIDR("2001:db8::/98")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = n.Subnet(99)
	}
}

func BenchmarkNet_PreviousNet_v4(b *testing.B) {
	_, n, _ := ParseCIDR("192.168.0.0/24")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = n.PreviousNet(24)
	}
}

func BenchmarkNet_PreviousNet_v6(b *testing.B) {
	_, n, _ := ParseCIDR("2001:db8::/98")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = n.PreviousNet(24)
	}
}

func BenchmarkNet_NextNet_v4(b *testing.B) {
	_, n, _ := ParseCIDR("192.168.0.0/24")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = n.NextNet(24)
	}
}

func BenchmarkNet_NextNet_v6(b *testing.B) {
	_, n, _ := ParseCIDR("2001:db8::/98")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = n.NextNet(24)
	}
}

func BenchmarkNewNetBetween_v4(b *testing.B) {
	ipa := net.IP{10, 0, 0, 0}
	ipb := net.IP{10, 0, 0, 255}
	for i := 0; i < b.N; i++ {
		_, _, _ = NewNetBetween(ipa, ipb)
	}
}

// Sorry for  abusing the benchmark suite here, i just think it's kind of neat
// to see how quickly one can allocate the entire v4 space in a Go application
func BenchmarkNextIP_EntireV4Space(b *testing.B) {
	xip := net.IP{0, 0, 0, 0}
	b.N = 4294967294
	b.StartTimer()
	for i := 0; i <= b.N; i++ {
		xip = NextIP(xip)
	}
	b.StopTimer()
}
