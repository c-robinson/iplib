package iplib

import (
	"fmt"
	"net"
)

func ExampleNet6_Contains() {
	n := NewNet6(net.ParseIP("2001:db8:1234:5678::"), 64, 0)
	fmt.Println(n.Contains(net.ParseIP("2001:db8:1234:5678::1")))
	fmt.Println(n.Contains(net.ParseIP("2001:db8:1234::")))
	// Output:
	// true
	// false
}

func ExampleNet6_Count() {
	// without hostmask
	n := NewNet6(net.ParseIP("2001:db8:1234:5678::"), 64, 0)
	fmt.Println(n.Count())

	// with hostmask set to 56, leaving 8 usable bytes between the two masks
	n = NewNet6(net.ParseIP("2001:db8:1234:5678::"), 64, 56)
	fmt.Println(n.Count())
	// Output:
	// 18446744073709551616
	// 256
}

func ExampleNet6_LastAddress() {
	// without hostmask
	n := NewNet6(net.ParseIP("2001:db8:1234:5678::"), 64, 0)
	fmt.Println(n.LastAddress())

	// with hostmask set to 56, leaving 8 usable bytes between the two masks
	n = NewNet6(net.ParseIP("2001:db8:1234:5678::"), 64, 56)
	fmt.Println(n.LastAddress())
	// Output:
	// 2001:db8:1234:5678:ffff:ffff:ffff:ffff
	// 2001:db8:1234:5678:ff00::
}

func ExampleNet6_NextIP() {
	// without hostmask
	n := NewNet6(net.ParseIP("2001:db8:1234:5678::"), 64, 0)
	fmt.Println(n.NextIP(net.ParseIP("2001:db8:1234:5678::")))

	// with hostmask set to 56, leaving 8 usable bytes between the two masks
	n = NewNet6(net.ParseIP("2001:db8:1234:5678::"), 64, 56)
	fmt.Println(n.NextIP(net.ParseIP("2001:db8:1234:5678::")))

	// as above, but trying to scan past the end of the netblock
	fmt.Println(n.NextIP(net.ParseIP("2001:db8:1234:5678:ff00::")))
	// Output:
	// 2001:db8:1234:5678::1 <nil>
	// 2001:db8:1234:5678:100:: <nil>
	// <nil> address is not a part of this netblock
}

func ExampleNet6_NextNet() {
	n := NewNet6(net.ParseIP("2001:db8:1234:5678::"), 64, 0)
	fmt.Println(n.NextNet(0))
	// Output: 2001:db8:1234:5679::/64
}

func ExampleNet6_PreviousIP() {
	// without hostmask
	n := NewNet6(net.ParseIP("2001:db8:1234:5678::"), 64, 0)
	fmt.Println(n.PreviousIP(net.ParseIP("2001:db8:1234:5678:ff00::")))

	// as above, but trying to scan past the end of the netblock
	fmt.Println(n.PreviousIP(net.ParseIP("2001:db8:1234:5678::")))

	// with hostmask set to 56, leaving 8 usable bytes between the two masks
	n = NewNet6(net.ParseIP("2001:db8:1234:5678::"), 64, 56)
	fmt.Println(n.PreviousIP(net.ParseIP("2001:db8:1234:5678:ff00::")))
	// Output:
	// 2001:db8:1234:5678:feff:ffff:ffff:ffff <nil>
	// <nil> address is not a part of this netblock
	// 2001:db8:1234:5678:fe00:: <nil>
}

func ExampleNet6_PreviousNet() {
	n := NewNet6(net.ParseIP("2001:db8:1234:5678::"), 64, 0)

	// at the same netmask
	fmt.Println(n.PreviousNet(0))

	// at a larger netmask (result encompasses the starting network)
	fmt.Println(n.PreviousNet(62))
	// Output:
	// 2001:db8:1234:5677::/64
	// 2001:db8:1234:5674::/62

}

func ExampleNet6_Subnet() {
	n := NewNet6(net.ParseIP("2001:db8:1234:5678::"), 64, 0)
	for _, i := range []int{65, 66} {
		sub, _ := n.Subnet(i, 0)
		fmt.Println(sub)
	}
	// Output:
	// [2001:db8:1234:5678::/65 2001:db8:1234:5678:8000::/65]
	// [2001:db8:1234:5678::/66 2001:db8:1234:5678:4000::/66 2001:db8:1234:5678:8000::/66 2001:db8:1234:5678:c000::/66]
}

func ExampleNet6_Supernet() {
	n := NewNet6(net.ParseIP("2001:db8:1234:5678::"), 64, 0)
	fmt.Println(n.Supernet(0, 0))
	// Output: 2001:db8:1234:5678::/63 <nil>
}
