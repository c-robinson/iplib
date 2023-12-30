# iid
[![Documentation](https://godoc.org/github.com/c-robinson/iplib?status.svg)](http://godoc.org/github.com/c-robinson/iplib/v2/iid)
[![Go Report Card](https://goreportcard.com/badge/github.com/c-robinson/iplib)](https://goreportcard.com/report/github.com/c-robinson/iplib)
[![Tests](https://img.shields.io/github/actions/workflow/status/c-robinson/iplib/test.yml?branch=main&longCache=true&label=Test&logo=github%20actions&logoColor=fff)](https://github.com/c-robinson/iplib/actions?query=workflow%3ATest)

This package implements functions for generating and validating IPv6 Interface
Identifiers (IID's) for use in link-local, global unicast and Stateless Address
Autoconfiguration (SLAAC). For the purposes of this module an IID is an IPv6
address constructed, somehow, from information which uniquely identifies a
given interface on a network, and is unique _within_ that network.

## Installing

```sh
go get -u github.com/c-robinson/iplib/v2
```

## Using IID

This library contains functions for uniting `net.IP` and `net.HardwareAddr`
addresses in order to generate globally unique IPv6 addresses. The simplest of
which is the "Modified EUI-64 address" described in [RFC4291 section 2.5.1](https://tools.ietf.org/html/rfc4291#section-2.5.1)

```go
package main

import (
	"fmt"
	"net"
	
	"github.com/c-robinson/iplib/v2/iid"
)

func main() {
	ip := net.ParseIP("2001:db8:1111:2222::")
	hw, _ := net.ParseMAC("99:88:77:66:55:44")
	myiid := iid.MakeEUI64Addr(ip, hw, iid.ScopeGlobal)
	fmt.Println(myiid) // will be "2001:db8:1111:2222:9b88:77ff:fe66:5544"
}
```

EUI64 is fine for a local subnet, but since it is tied to a hardware address
and guessable by design it is a privacy nightmare as outlined in [RFC4941](https://tools.ietf.org/html/rfc4941).

[RFC7217](https://tools.ietf.org/html/rfc7217) defines an algorithm to create
"semantically opaque" IID's based on the local interface by hashing the address
with a secret key, a counter, and some optional additional data. The resulting
IID is pseudo-random (the same inputs will result in the same outputs) so care
must be taken while generating it. This function has some requirements:

- `secret` a `[]byte` that is a closely-held secret key
- `counter` an `int64`, this is what provides the address its mutability. The
  RFC specifies that this counter should be incremented every time the same
  ipaddr/hwaddr pair is used as input and should be stored in non-volatile
  memory to preserve it
- `netid` is an optional parameter to improve the privacy of the results, it
  is suggested that this be some other bit of information from the local
  network such as an 802.11 SSID.

NOTE: it is possible, though very unlikely, that an address generated this way
might collide with the [IANA Reserved Interface Identifier List](https://www.iana.org/assignments/ipv6-interface-ids/ipv6-interface-ids.xhtml),
if this happens an `iid.ErrIIDAddressCollision` will be returned. If so
`counter` should be incremented and the function re-run.

```go
package main

import (
	"fmt"
	"net"
	
	"github.com/c-robinson/iplib/v2/iid"
)

func main() {
	ip      := net.ParseIP("2001:db8::")
	hw, _   := net.ParseMAC("77:88:99:aa:bb:cc")
	counter := int64(1)
	netid   := []byte("01234567")
	secret  := []byte("secret")
	
	myiid, err := iid.MakeOpaqueAddr(ip, hw, counter, netid, secret)
	if err != nil {
		fmt.Println("a very unlikely collision occurred!")
	}
	fmt.Println(myiid) // will be "2001:db8::c6fa:ba02:41ab:282c"
}
```

`MakeOpaqueIID()` is an implementation of the RFC's specified function
using its' preferred [SHA256](https://golang.org/pkg/crypto/sha256/)
hashing algorithm and a `iid.ScopeGlobal` scope. If either of these is not to
your liking you can roll your own by calling the underlying function.

NOTE: if you use any hashing algorithm other than SHA224 or SHA256 you will
need to import both `"crypto"` _and_ the crypto submodule with your specific
implementation first (e.g. `_ "golang.org/x/crypto/blake2s"`. Also note that
the RFC _specifically prohibits_ MD5 as being too insecure for use. Here's an
example using [SHA512](https://golang.org/pkg/crypto/sha512/)

```go
package main

import (
	"crypto"
	_ "crypto/sha512"
	"fmt"
	"net"
	
	"github.com/c-robinson/iplib/v2/iid"
)

func main() {
	ip      := net.ParseIP("2001:db8::")
	hw, _   := net.ParseMAC("77:88:99:aa:bb:cc")
	counter := int64(1)
	netid   := []byte("01234567")
	secret  := []byte("secret")
	
	myiid, err := iid.GenerateRFC7217Addr(ip, hw, counter, netid, secret, crypto.SHA384, iid.ScopeGlobal)
	if err != nil {
		fmt.Println("a very unlikely collision occurred!")
	}
	fmt.Println(myiid) // will be "2001:db8::51b3:c6b0:4e14:3519"
}
```

Finally, to be entirely RFC7217-compliant a function _should_ check it's
results to make sure they don't collide with the IANA Reserved Interface
Identifier List. In the name of "using every part of the buffalo" the function
is exposed for the extremely unlikely case where anyone needs it:

```go
package main

import (
	"fmt"
	"net"
	
	"github.com/c-robinson/iplib/v2/iid"
)

func main() {
	ip := net.ParseIP("2001:db8::0200:5EFF:FE00:5211")
	res := iid.GetReservationsForIP(ip)
	fmt.Println(res.RFC) // will be "RFC4291"
	fmt.Println(res.Title) // "Reserved IPv6 Interface Identifiers corresponding to the IANA Ethernet Block"
}
```
