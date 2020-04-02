package iplib

import (
	"fmt"
	"net"
)

func ExampleNet4_Contains() {
	n := NewNet4(net.ParseIP("192.168.1.0"), 24)
	fmt.Println(n.Contains(net.ParseIP("192.168.1.111")))
	fmt.Println(n.Contains(net.ParseIP("10.14.0.1")))
	// Output:
	// true
	// false
}

func ExampleNet4_ContainsNet() {
	n1 := NewNet4(net.ParseIP("192.168.0.0"), 16)
	n2 := NewNet4(net.ParseIP("192.168.1.0"), 24)
	fmt.Println(n1.ContainsNet(n2))
	fmt.Println(n2.ContainsNet(n1))
	// Output:
	// true
	// false
}

func ExampleNet4_Count() {
	n := NewNet4(net.ParseIP("192.168.0.0"), 16)
	fmt.Println(n.Count())
	// Output: 65534
}

func ExampleNet4_Enumerate() {
	n := NewNet4(net.ParseIP("192.168.0.0"), 16)
	fmt.Println(n.Enumerate(2, 100))
	// Output: [192.168.0.101 192.168.0.102]
}

func ExampleNet4_NextNet() {
	n := NewNet4(net.ParseIP("192.168.1.0"), 24)
	fmt.Println(n.NextNet(24))
	// Output: 192.168.2.0/24
}

func ExampleNet4_PreviousNet() {
	n := NewNet4(net.ParseIP("192.168.1.0"), 24)
	fmt.Println(n.PreviousNet(24))
	// Output: 192.168.0.0/24
}

func ExampleNet4_Subnet() {
	n := NewNet4(net.ParseIP("192.168.0.0"), 16)
	sub, _ := n.Subnet(17)
	fmt.Println(sub)
	// Output: [192.168.0.0/17 192.168.128.0/17]
}

func ExampleNet4_Supernet() {
	n := NewNet4(net.ParseIP("192.168.1.0"), 24)
	n2, _ := n.Supernet(22)
	fmt.Println(n2)
	// Output: 192.168.0.0/22
}

func ExampleNet4_Wildcard() {
	n := NewNet4(net.ParseIP("192.168.0.0"), 16)
	fmt.Println(n.Wildcard())
	// Output: 0000ffff
}
