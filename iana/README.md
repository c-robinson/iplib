# iana
[![Documentation](https://godoc.org/github.com/c-robinson/iplib?status.svg)](http://godoc.org/github.com/c-robinson/iplib/v2/iana)
[![Go Report Card](https://goreportcard.com/badge/github.com/c-robinson/iplib)](https://goreportcard.com/report/github.com/c-robinson/iplib)
[![Tests](https://img.shields.io/github/actions/workflow/status/c-robinson/iplib/test.yml?branch=main&longCache=true&label=Test&logo=github%20actions&logoColor=fff)](https://github.com/c-robinson/iplib/actions?query=workflow%3ATest)

This package imports the [Internet Assigned Number Authority (IANA)](https://www.iana.org/)
Special IP Address Registry for [IPv4](https://www.iana.org/assignments/iana-ipv4-special-registry/iana-ipv4-special-registry.xhtml)
and [IPv6](https://www.iana.org/assignments/iana-ipv6-special-registry/iana-ipv6-special-registry.xhtml)
and exposes it as a data structure. Functions allow a caller to compare the
registry against `net.IP` and `iplib.Net` objects to see if they contain or
are contained within an reserved IP address block.

## Installing

```sh
go get -u github.com/c-robinson/iplib/v2
```

## Using IANA

Here are examples comparing against both an address and a network. Note that in
the network case it is entirely possible for a broad-enough network to contain
multiple reservations. If this is the case all reservations will be returned.

```go
package main

import (
	"fmt"
	"net"
	
	"github.com/c-robinson/iplib/v2"
	"github.com/c-robinson/iplib/v2/iana"
)

func main() {
	ipa := net.ParseIP("144.21.21.21")
	ipb := net.ParseIP("192.168.12.5")
	
	res := iana.GetReservationsForIP(ipa)
	fmt.Println(len(res))                 // 0
	
	res = iana.GetReservationsForIP(ipb)
	fmt.Println(len(res))                 // 1
	fmt.Println(res[0].Title)             // "Private-Use"
	fmt.Println(res[0].RFC)               // ["RFC1918"]
	
	_, neta, _ := iplib.ParseCIDR("2001::/16")
	
	res = iana.GetReservationsForNetwork(neta)
	fmt.Println(len(res))                     // 10
	fmt.Println(iana.IsForwardable(neta))     // false
	fmt.Println(iana.IsGlobal(neta))          // false
	fmt.Println(iana.IsReserved(neta))        // true
	fmt.Println(iana.GetRFCsForNetwork(neta)) // all relevant RFCs, in this case: 
	                                          // [RFC1752,RFC2928,RFC3849,RFC4380,
	                                          //  RFC5180,RFC7343,RFC7450,RFC7535,
	                                          //  RFC7723,RFC7954,RFC8155,RFC8190]
}
```
