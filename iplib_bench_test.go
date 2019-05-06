package iplib

import (
	"testing"
	"net"
	"math/big"
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
	xip := net.IP{10,0,0,0}
	for i := 0; i < b.N; i++ {
		NewNet(xip, 24)
	}
}

func BenchmarkNextIP4(b *testing.B) {
	var xip = net.IP{10,0,0,0}
	for i := 0; i < b.N; i++ {
		xip = NextIP(xip)
	}
}

func BenchmarkIncrementIP4By(b *testing.B) {
	var xip = net.IP{10,0,0,0}
	for i := 0; i < b.N; i++ {
		xip = IncrementIP4By(xip, 1)
	}
}

func BenchmarkIncrementIPBy_v4(b *testing.B) {
	var xip = net.IP{10,0,0,0}
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

// Sorry for  abusing the benchmark suite here, i just think it's kind of neat
// to see how quickly one can allocate the entire v4 space in a Go application
func BenchmarkNextIP_EntireV4Space(b *testing.B) {
	xip := net.IP{0,0,0,0}
	b.N = 4294967294
	b.StartTimer()
	for i := 0; i <= b.N; i++ {
		xip = NextIP(xip)
	}
	b.StopTimer()
}
