# IPLib

A hopefully useful, asperationally full-featured library for working with IP
addresses and networks.

IPLib is an extension of the `net.IP` utilities and is intended to make working
with IP addresses a little bit easier by providing tools to manage blocks of
addresses. Tools include:

[![Documentation](https://godoc.org/github.com/<username>/<library>?status.svg)](http://godoc.org/github.com/c-robinson/iplib)
[![Go Report Card](https://goreportcard.com/badge/github.com/<username>/<library>)](https://goreportcard.com/report/github.com/c-robinson/iplib)

##### net.IP tools

Some simple tools for performing common tasks against IP objects:

- Compare two addresses
- Get the delta between two addresses
- Sort
- Decrement or increment addresses
- Print v4 as a hexadecimal string
- Print v6 in fully expanded form
- Convert between net.IP, integer and hexadecimal
- Get the version of a v4 address or force a 6to4 address to be a v4 address

##### iplib.IPNet

An enhancement of `net.IPNet` providing features such as:

- Retrieve the wildcard mask
- Get the network, broadcast, first and last usable addresses
- Increment or decrement an address within the boundaries of a netblock
- Enumerate all or part of a netblock to `[]net.IP`
- Allocate subnets
- Find free space between allocated subnets
- Expand subnets if space allows

#### Using IPNet

There are a series of functions for working with v4 or v6 `net.IP` objects:

```Go
import "github.com/c-robinson/iplib"

main() {
    ipa := net.IP{192,168,1,1}
    ipb := net.IP{192,168,1,5}

    c := iplib.CompareIPs(ipa, ipb)
    d := iplib.DeltaIP(ipa, ipb)
    
    ipc := iplib.IncrementIPBy(ipb, 20)
    fmt.Println(IPToHextString(ipc))

    iplist := []net.IP{ ipb, ipc, ipa }
    sorted := sort.Sort(iplib.ByAddr(iplist))

    vers := iplib.Version(ipc)

    asInt := iplib.IP4ToUint32(ipc)
    ipd := iplib.Uint32ToIP4(asInt+20)

}
```

Addresses that require or return a count default to using `uint32`, which is
sufficient for working with the entire IPv4 space. As a rule these functions
are just lowest-common wrappers around IPv4- or IPv6-specific functions. The
IPv6-specific variants use `big.Int` so they can access the entire v6 space:


```Go
main() {
    ip6a := net.IP{32,1,13,184,133,163,0,0,0,0,138,46,3,112,115,52}

    m := big.NewInt(int64(iplib.MaxUint32)) // MaxUint32 is 4,294,967,296
    ip6b := iplib.IncrementIP6By(ip6a, m)

    d := iplib.DeltaIP6(ip6a, ip6b)

    asBigint := iplib.IP6ToBigint(ip6b)
    ip6c := iplib.BigintToIP6(asBigint.Add(m))
}

```

To work with networks simply create an `iplib.IPNet` object:

```Go
import "path/to/iplib"

main() {
    ipa, ipna, erra := iplib.NewNet("192.168.1.0/22")
    if erra != nil {
        // become sad!
    }

    ipb, ipnb, errb := iplib.NewNet("192.168.2.0/24")
    if errb != nil {
        // still sad!
    }

    fmt.Println(ipna.Count()) // there's a Count6() for IPv6

    iplib.CompareNets(ipna, ipnb)

    ipnlist := []iplib.Net{ipnb, ipna}
    sorted := sort.Sort(ByNet(ipnlist))
    
    v := ipna.ContainsNet(ipnb)

    ipc, errc := ipna.NextIP(ipna.NetworkAddress())
    if errc != nil {
        // be sad
    }

    ipd, errd := ipna.PreviousIP(ipna.NetworkAddress())
    if errd != nil {
        // rejoice, for this will error out (underflowing network boundary)
    }

    iplist := ipna.Enumerate(0, 0)

```
