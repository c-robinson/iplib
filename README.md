# IPLib 
[![Documentation](https://godoc.org/github.com/c-robinson/iplib?status.svg)](http://godoc.org/github.com/c-robinson/iplib)
[![CircleCI](https://circleci.com/gh/c-robinson/iplib/tree/master.svg?style=svg)](https://circleci.com/gh/c-robinson/iplib/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/c-robinson/iplib)](https://goreportcard.com/report/github.com/c-robinson/iplib)
[![Coverage Status](https://coveralls.io/repos/github/c-robinson/iplib/badge.svg?branch=master)](https://coveralls.io/github/c-robinson/iplib?branch=master)

I really enjoy Python's [ipaddress](https://docs.python.org/3/library/ipaddress.html)
library and Ruby's [ipaddr](https://ruby-doc.org/stdlib-2.5.1/libdoc/ipaddr/rdoc/IPAddr.html),
I think you can write a lot of neat software if some of the little problems
around manipulating IP addresses and netblocks are taken care of for you, so I
set out to write something like them for my language of choice, Go. This is
what I've come up with.

[IPLib](http://godoc.org/github.com/c-robinson/iplib) is a hopefully useful,
aspirationally full-featured library built around and on top of the address
primitives found in the [net](https://golang.org/pkg/net/) package, it seeks
to make them more accessible and easier to manipulate. 

It includes:

##### net.IP tools

Some simple tools for performing common tasks against IP objects:

- Compare two addresses
- Get the delta between two addresses
- Sort
- Decrement or increment addresses
- Print v4 as a hexadecimal string
- Print v6 in fully expanded form
- Convert between net.IP, integer and hexadecimal
- Get the version of a v4 address or force a IPv4-mapped IPv6address to be a 
  v4 address

##### iplib.IPNet

An enhancement of `net.IPNet` providing features such as:

- Retrieve the wildcard mask
- Get the network, broadcast, first and last usable addresses
- Increment or decrement an address within the boundaries of a netblock
- Enumerate all or part of a netblock to `[]net.IP`
- Allocate subnets and supernets

## Sub-modules

- [iana](https://github.com/c-robinson/iplib/tree/master/iana) - a module for referencing 
  IP netblocks against the [Internet Assigned Numbers Authority's](https://www.iana.org/)
  Special IP Address Registry
- [iid](https://github.com/c-robinson/iplib/tree/master/iid) - a module for
  generating and validating IPv6 Interface Identifiers, including [RFC4291](https://tools.ietf.org/html/rfc4291)
  modified EUI64 and [RFC7217](https://tools.ietf.org/html/rfc7217)
  Semantically Opaque addresses

## Installing

```sh
go get -u github.com/c-robinson/iplib
```

## Using IPLib

There are a series of functions for working with v4 or v6 `net.IP` objects:

```go
package main

import (
	"fmt"
	"net"
	"sort"
	
	"github.com/c-robinson/iplib"
)


func main() {
	ipa := net.ParseIP("192.168.1.1")
	ipb := iplib.IncrementIPBy(ipa, 15)      // ipb is 192.168.1.16
	ipc := iplib.NextIP(ipa)                 // ipc is 192.168.1.2

	fmt.Println(iplib.CompareIPs(ipa, ipb))  // -1
    
	fmt.Println(iplib.DeltaIP(ipa, ipb))     // 15
    
	fmt.Println(iplib.IPToHexString(ipc))    // "c0a80102"

	iplist := []net.IP{ ipb, ipc, ipa }
	sort.Sort(iplib.ByIP(iplist))            // []net.IP{ipa, ipc, ipb}

	fmt.Println(iplib.IP4ToUint32(ipa))      // 3232235777
	fmt.Println(iplib.IPToBinaryString(ipa))  // 11000000.10101000.00000001.00000001
	ipd := iplib.Uint32ToIP4(iplib.IP4ToUint32(ipa)+20) // ipd is 192.168.1.21
	fmt.Println(iplib.IP4ToARPA(ipa))        // 1.1.168.192.in-addr.arpa
}
```

Addresses that require or return a count default to using `uint32`, which is
sufficient for working with the entire IPv4 space. As a rule these functions
are just lowest-common wrappers around IPv4- or IPv6-specific functions. The
IPv6-specific variants use `big.Int` so they can access the entire v6 space:


```go
package main

import (
	"fmt"
	"math/big"
	"net"
	"sort"
	
	"github.com/c-robinson/iplib"
)


func main() {
	ipa := net.ParseIP("2001:db8::1")
	ipb := iplib.IncrementIPBy(ipa, 15)      // ipb is 2001:db8::16
	ipc := iplib.NextIP(ipa)                 // ipc is 2001:db8::2

	fmt.Println(iplib.CompareIPs(ipa, ipb))  // -1
    
	fmt.Println(iplib.DeltaIP6(ipa, ipb))     // 15
    
	fmt.Println(iplib.ExpandIP6(ipa))        // "2001:0db8:0000:0000:0000:0000:0000:0001"
	fmt.Println(iplib.IPToBigint(ipa))       // 42540766411282592856903984951653826561 
	fmt.Println(iplib.IPToBinaryString(ipa)) // 00100000.00000001.00001101.10111000.00000000.00000000.00000000.00000000.00000000.00000000.00000000.00000000.00000000.00000000.00000000.00000001
    
	iplist := []net.IP{ ipb, ipc, ipa }
	sort.Sort(iplib.ByIP(iplist))            // []net.IP{ipa, ipc, ipb}

	m := big.NewInt(int64(iplib.MaxIPv4))    // e.g. 4,294,967,295
	ipd := iplib.IncrementIP6By(ipa, m)      // ipd is 2001:db8::1:0:0

	fmt.Println(iplib.DeltaIP6(ipb, ipd))    // 4294967274
	fmt.Println(iplib.IP6ToARPA(ipa))        // 1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa
}

```

To work with networks simply create an `iplib.IPNet` object:

```go
package main

import (
	"fmt"
	"net"
	"sort"
	
	"github.com/c-robinson/iplib"
)

func main() {
	// this calls net.ParseCIDR() under the hood, but returns an iplib.Net object
	_, ipna, err := iplib.ParseCIDR("192.168.1.0/22")
	if err != nil {
		// this will be an error from the net package 
	}
	
	// NewNet() wants a net.IP and is waaaaaaaaaaaaaaaaay faster
	ipb := net.ParseIP("192.168.2.0")
	ipnb := iplib.NewNet(ipb, 22)
    
	// ...works for IPv6 too
	ipc := net.ParseIP("2001:db8::1")
	ipnc := iplib.NewNet(ipc, 64)

	fmt.Println(ipna.Count())                  // 1022 -- good enough for ipv4, but...
    
	fmt.Println(ipnc.Count())                  // 4294967295 -- ...sigh
	fmt.Println(ipnc.Count6())                 // 18446744073709551616 -- yay Count6() !

	fmt.Println(iplib.CompareNets(ipna, ipnb)) // -1

	ipnlist := []iplib.Net{ipnb, ipna, ipnc}
	sort.Sort(iplib.ByNet(ipnlist))            // []iplib.Net{ ipna, ipnb, ipnc } 
    
	elist := ipna.Enumerate(0, 0)
	fmt.Println(len(elist))                    // 1022
    
	fmt.Println(ipna.ContainsNet(ipnb))        // true
    
	fmt.Println(ipna.NetworkAddress())         // 192.168.1.0
	fmt.Println(ipna.FirstAddress())           // 192.168.1.1
	fmt.Println(ipna.LastAddress())            // 192.168.3.254
	fmt.Println(ipna.BroadcastAddress())       // 192.168.3.255
    
	fmt.Println(ipnc.NetworkAddress())         // 2001:db8::1 -- meaningless in IPv6
	fmt.Println(ipnc.FirstAddress())           // 2001:db8::1
	fmt.Println(ipnc.LastAddress())            // 2001:db8::ffff:ffff:ffff:ffff
	fmt.Println(ipnc.BroadcastAddress())       // 2001:db8::ffff:ffff:ffff:ffff
    
	ipa1 := net.ParseIP("2001:db8::2")
	ipa1, err = ipna.PreviousIP(ipa1)         //  net.IP{2001:db8::1}, nil
	ipa1, err = ipna.PreviousIP(ipa1)         //  net.IP{}, ErrAddressAtEndOfRange
}
```

`iplib.IPNet` objects can be used to generate subnets and supernets:

```go
package main

import (
	"fmt"
	
	"github.com/c-robinson/iplib"
)

func main() {
    _, ipna, _ := iplib.ParseCIDR("192.168.4.0/22")
    fmt.Println(ipna.Subnet(24))   // []iplib.Net{ 192.168.4.0/24, 192.168.5.0/24, 
                                   //              192.168.6.0/24, 192.168.7.0/24 }
    ipnb, err := ipna.Supernet(21) // 192.168.0.0/21
    
    ipnc := ipna.PreviousNet(21)   // 192.168.0.0/21
    
    ipnd := ipna.NextNet(21)       // 192.168.8.0/21
}
```
