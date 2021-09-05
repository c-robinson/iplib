package iplib

import (
	"fmt"
	"math/big"
	"net"
)

func ExampleBigintToIP6() {
	z := big.Int{}
	z.SetString("42540766452641154071740215577757643572", 10)
	fmt.Println(BigintToIP6(&z))
	// Output: 2001:db8:85a3::8a2e:370:7334
}

func ExampleCompareIPs() {
	fmt.Println(CompareIPs(net.ParseIP("192.168.1.0"), net.ParseIP("192.168.1.1")))
	fmt.Println(CompareIPs(net.ParseIP("10.0.0.0"), net.ParseIP("10.0.0.0")))
	fmt.Println(CompareIPs(net.ParseIP("2001:db8::100"), net.ParseIP("2001:db8::99")))
	// Output:
	// -1
	// 0
	// 1
}

func ExampleDecrementIP4By() {
	ip := net.ParseIP("192.168.2.0")
	fmt.Println(DecrementIP4By(ip, 255))
	// Output: 192.168.1.1
}

func ExampleDecrementIP6By() {
	z := big.NewInt(16777215)
	ip := net.ParseIP("2001:db8::ffff:ffff")
	fmt.Println(DecrementIP6By(ip, z))
	// Output: 2001:db8::ff00:0
}

func ExampleDecrementIP6WithinHostmask() {
	ip := net.ParseIP("2001:db8:1000::")
	ip1, _ := DecrementIP6WithinHostmask(ip, NewHostMask(0), big.NewInt(1))
	ip2, _ := DecrementIP6WithinHostmask(ip, NewHostMask(56), big.NewInt(1))
	fmt.Println(ip1)
	fmt.Println(ip2)
	// Output:
	// 2001:db8:fff:ffff:ffff:ffff:ffff:ffff
	// 2001:db8:fff:ffff:ff00::
}

func ExampleDeltaIP4() {
	ipa := net.ParseIP("192.168.1.1")
	ipb := net.ParseIP("192.168.2.0")
	fmt.Println(DeltaIP4(ipa, ipb))
	// Output: 255
}

func ExampleDeltaIP6() {
	ipa := net.ParseIP("2001:db8::ffff:ffff")
	ipb := net.ParseIP("2001:db8::ff00:0")
	fmt.Println(DeltaIP6(ipa, ipb))
	// Output: 16777215
}

func ExampleEffectiveVersion() {
	fmt.Println(EffectiveVersion(net.ParseIP("192.168.1.1")))
	fmt.Println(EffectiveVersion(net.ParseIP("::ffff:c0a8:101")))
	fmt.Println(EffectiveVersion(net.ParseIP("2001:db8::c0a8:101")))
	// Output:
	// 4
	// 4
	// 6
}

func ExampleExpandIP6() {
	fmt.Println(ExpandIP6(net.ParseIP("2001:db8::1")))
	// Output: 2001:0db8:0000:0000:0000:0000:0000:0001
}

func ExampleForceIP4() {
	fmt.Println(len(ForceIP4(net.ParseIP("::ffff:c0a8:101"))))
	// Output: 4
}

func ExampleHexStringToIP() {
	ip := HexStringToIP("c0a80101")
	fmt.Println(ip.String())
	// Output: 192.168.1.1
}

func ExampleIncrementIP4By() {
	ip := net.ParseIP("192.168.1.1")
	fmt.Println(IncrementIP4By(ip, 255))
	// Output: 192.168.2.0
}

func ExampleIncrementIP6By() {
	z := big.NewInt(16777215)
	ip := net.ParseIP("2001:db8::ff00:0")
	fmt.Println(IncrementIP6By(ip, z))
	// Output: 2001:db8::ffff:ffff
}

func ExampleIncrementIP6WithinHostmask() {
	ip := net.ParseIP("2001:db8:1000::")
	ip1, _ := IncrementIP6WithinHostmask(ip, NewHostMask(0), big.NewInt(1))
	ip2, _ := IncrementIP6WithinHostmask(ip, NewHostMask(56), big.NewInt(1))
	fmt.Println(ip1)
	fmt.Println(ip2)
	// Output:
	// 2001:db8:1000::1
	// 2001:db8:1000:0:100::
}

func ExampleIP4ToARPA() {
	fmt.Println(IP4ToARPA(net.ParseIP("192.168.1.1")))
	// Output: 1.1.168.192.in-addr.arpa
}

func ExampleIP6ToARPA() {
	fmt.Println(IP6ToARPA(net.ParseIP("2001:db8::ffff:ffff")))
	// Output: f.f.f.f.f.f.f.f.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa
}

func ExampleIP4ToUint32() {
	fmt.Println(IP4ToUint32(net.ParseIP("192.168.1.1")))
	// Output: 3232235777
}

func ExampleIPToBinaryString() {
	fmt.Println(IPToBinaryString(net.ParseIP("192.168.1.1")))
	fmt.Println(IPToBinaryString(net.ParseIP("2001:db8::ffff:ffff")))
	// Output:
	// 11000000.10101000.00000001.00000001
	// 00100000.00000001.00001101.10111000.00000000.00000000.00000000.00000000.00000000.00000000.00000000.00000000.11111111.11111111.11111111.11111111
}

func ExampleIPToHexString() {
	fmt.Println(IPToHexString(net.ParseIP("192.168.1.1")))
	// Output:
	// c0a80101
}

func ExampleNextIP() {
	fmt.Println(NextIP(net.ParseIP("192.168.1.1")))
	fmt.Println(NextIP(net.ParseIP("2001:db8::ffff:fffe")))
	// Output:
	// 192.168.1.2
	// 2001:db8::ffff:ffff
}

func ExampleNextIP6WithinHostmask() {
	ip, _ := NextIP6WithinHostmask(net.ParseIP("2001:db8:1234:5678::"), NewHostMask(56))
	fmt.Println(ip)
	// Output: 2001:db8:1234:5678:100::
}

func ExamplePreviousIP() {
	fmt.Println(PreviousIP(net.ParseIP("192.168.1.2")))
	fmt.Println(PreviousIP(net.ParseIP("2001:db8::ffff:ffff")))
	// Output:
	// 192.168.1.1
	// 2001:db8::ffff:fffe
}

func ExamplePreviousIP6WithinHostmask() {
	ip, _ := PreviousIP6WithinHostmask(net.ParseIP("2001:db8:1234:5678::"), NewHostMask(56))
	fmt.Println(ip)
	// Output: 2001:db8:1234:5677:ff00::
}

func ExampleUint32ToIP4() {
	fmt.Println(Uint32ToIP4(3232235777))
	// Output: 192.168.1.1
}

func ExampleVersion() {
	fmt.Println(Version(ForceIP4(net.ParseIP("192.168.1.1"))))
	fmt.Println(Version(net.ParseIP("::ffff:c0a8:101")))
	fmt.Println(Version(net.ParseIP("2001:db8::c0a8:101")))
	// Output:
	// 4
	// 6
	// 6
}
