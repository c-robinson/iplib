/*
Package iplib provides enhanced tools for working with IP networks and
addresses. These tools are built upon and extend the generic functionality
found in the Go "net" package.

IPLib comes in two basic parts: a set of utility features for working with
net.IP (sort, increment, decrement, delta, compare; convert to hex-string or
integer) and an enhancement of net.IPNet (iplib.Net) that can calculate the
broadcast, first and last IP addresses in its block, as well as enumerating
the block into a []net.IP, and incrementing or decrementing within the
boundaries of the block.

For the most part IPLib tries to ensure that v4 and v6 addresses are treated
equally and managed transparently. The one exception is those functions which
return or require a total as an integer: for these a version-independent
function is provided and limited to uint32, but there will also be a v4 and
v6 variant, the v6 function will always take *big.Int and be able to access
the entire v6 address space. In all cases the version-independent function is
simply a router between the v4 and v6 functions that internally converts
uint32 to big.Int when necessary.

For functions where it is possible to exceed the address-space the rule is
that underflows return the version-appropriate all-zeroes address while
overflows return the all-ones.

A special note about IPv4 blocks with one host bit (/31): the only reason to
allocate such a subnet is for use as an RFC 3021 point-to-point network and
IPLib assumes this use-case and acts accordingly:

 - Count() will report 2 addresses instead of 0
 - FirstAddress() and NetworkAddress() will be equivalent
 - LastAddress() and BroadcastAddress() will be equivalent

 */
package iplib

import (
	"errors"
	"net"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math/big"
)

const (
	// MaxIPv4 is the max size of a uint32, also the IPv4 address space
	MaxIPv4 = 1<<32 - 1
)

var (
	ErrAddressAtEndOfRange = errors.New("proposed operation would cause address to exit block")
	ErrAddressOutOfRange   = errors.New("the given IP address is not a part of this netblock")
	ErrBroadcastAddress    = errors.New("address is the broadcast address of this netblock (and not considered usable)")
	ErrNetworkAddress      = errors.New("address is the network address of this netblock (and not considered usable)")
	ErrNoValidRange        = errors.New("no netblock can be found between the supplied values")
)

// ByIP implements sort.Interface for net.IP addresses
type ByIP []net.IP

// Len implements sort.interface Len(), returning the length of the
// ByIP array
func (bi ByIP) Len() int {
	return len(bi)
}

// Swap implements sort.interface Swap(), swapping two elements in our array
func (bi ByIP) Swap(a, b int) {
	bi[a], bi[b] = bi[b], bi[a]
}

// Less implements sort.interface Less(), given two elements in the array it
// returns true if the LHS should sort before the RHS. For details on the
// implementation, see CompareIPs()
func (bi ByIP) Less(a, b int) bool {
	val := CompareIPs(bi[a], bi[b])
	if val == -1 {
		return true
	}
	return false
}

// ByNet implements sort.Interface for iplib.Net based on the
// starting address of the netblock, with the netmask as a tie breaker. So if
// two Networks are submitted and one is a subset of the other, the enclosing
// network will be returned first.
type ByNet []Net

// Len implements sort.interface Len(), returning the length of the
// ByNetwork array
func (bn ByNet) Len() int {
	return len(bn)
}

// Swap implements sort.interface Swap(), swapping two elements in our array
func (bn ByNet) Swap(a, b int) {
	bn[a], bn[b] = bn[b], bn[a]
}

// Less implements sort.interface Less(), given two elements in the array it
// returns true if the LHS should sort before the RHS. For details on the
// implementation, see CompareNets()
func (bn ByNet) Less(a, b int) bool {
	val := CompareNets(bn[a], bn[b])
	if val == -1 {
		return true
	}
	return false
}

// BigintToIP6 converts a big.Int to an ip6 address and returns it as a net.IP
func BigintToIP6(z *big.Int) net.IP {
	b := z.Bytes()
	if len(b) > 16 {
		return generateNetLimits(6, 255)
	}
	if v := z.Sign(); v <= 0 {
		return generateNetLimits(6, 0)
	}

	// for cases where the resulting []byte isn't long enough
	if len(b) < 16 {
		for i := 15 - len(b); i >=0; i-- {
			b = append([]byte{0}, b...)
		}
	}
	return net.IP(b)
}

// CompareIPs is here for completeness. This is a good way to compare two
// IP objects. Since it uses bytes.Compare the return value is identical:
// 0 if a==b, -1 if a<b, 1 if a>b.
func CompareIPs(a, b net.IP) int {
	return bytes.Compare(a.To16(), b.To16())
}

// CompareNets compares two iplib.Net objects by evaluating their network
// address (the first address in a CIDR range) and, if they're equal,
// comparing their netmasks (smallest wins). This means that if a network is
// compared to one of its subnets, the enclosing network sorts first.
func CompareNets(a, b Net) int {
	val := bytes.Compare(a.NetworkAddress(), b.NetworkAddress())
	if val != 0 {
		return val
	}

	am, _ := a.Mask.Size()
	bm, _ := b.Mask.Size()

	if am == bm {
		return 0
	}
	if am < bm {
		return -1
	}
	return 1
}

// DecrementIPBy returns a net.IP that is lower than the supplied net.IP by
// the supplied integer value. If you underflow the IP space it will return
// the zero address.
func DecrementIPBy(ip net.IP, count uint32) net.IP {
	if EffectiveVersion(ip) == 4 {
		return DecrementIP4By(ip, count)
	}
	z := big.NewInt(int64(count))
	return DecrementIP6By(ip, z)
}

// DecrementIP4By returns a v4 net.IP that is lower than the supplied net.IP
// by the supplied integer value. If you underflow the IP space it will return
// 0.0.0.0
func DecrementIP4By(ip net.IP, count uint32) net.IP {
	i := IP4ToUint32(ip)
	d := i - count

	// check for underflow
	if d > i {
		return generateNetLimits(4, 0)
	}
	return Uint32ToIP4(d)
}

// DecrementIP6By returns a net.IP that is lower than the supplied net.IP by
// the supplied integer value. If you underflow the IP space it will return
// ::
func DecrementIP6By(ip net.IP, count *big.Int) net.IP {
	z := IPToBigint(ip)
	z.Sub(z, count)
	return BigintToIP6(z)
}

// DeltaIP takes two net.IP's as input and returns the difference between them
// up to the limit of uint32.
func DeltaIP(a, b net.IP) uint32 {
	if EffectiveVersion(a) == 4 && EffectiveVersion(b) == 4 {
		return DeltaIP4(a, b)
	}
	m := big.NewInt(int64(MaxIPv4))
	z := DeltaIP6(a, b)
	if v := z.Cmp(m); v > 0 {
		return MaxIPv4
	}
	return uint32(z.Uint64())
}

// DeltaIP4 takes two net.IP's as input and returns a total of the number of
// addresses between them, up to the limit of uint32.
func DeltaIP4(a, b net.IP) uint32 {
	ai := IP4ToUint32(a)
	bi := IP4ToUint32(b)

	if ai > bi {
		return ai - bi
	}
	return bi - ai
}

// DeltaIP6 takes two net.IP's as input and returns a total of the number of
// addressed between them as a big.Int. It will technically work on v4 as well
// but is considerably slower than DeltaIP4.
func DeltaIP6(a, b net.IP) *big.Int {
	ai := IPToBigint(a)
	bi := IPToBigint(b)
	i := big.NewInt(0)

	if v := ai.Cmp(bi); v >= 0 {
		return i.Sub(ai, bi)
	}
	return i.Sub(bi, ai)
}

// EffectiveVersion returns 4 if the net.IP either contains a v4 address or if
// it contains the v4-encapsulating v6 address range ::ffff
func EffectiveVersion(ip net.IP) int {
	a := ip.To4()
	if a == nil {
		return 6
	}
	return 4
}

// ForceIP4 takes a net.IP containing a 6to4 address and returns only the
// encapsulated v4 address.
func ForceIP4(ip net.IP) net.IP {
	if len(ip) == 16 {
		return ip[12:]
	}
	return ip
}

// IncrementIPBy returns a net.IP that is greater than the supplied net.IP by
// the supplied integer value. If you overflow the IP space it will return
// the all-ones address
func IncrementIPBy(ip net.IP, count uint32) net.IP {
	if Version(ip) == 4 {
		return IncrementIP4By(ip, count)
	}
	z := big.NewInt(int64(count))
	return IncrementIP6By(ip, z)
}

// IncrementIP4By returns a v4 net.IP that is greater than the supplied
// net.IP by the supplied integer value. If you overflow the IP space it
// will return 255.255.255.255
func IncrementIP4By(ip net.IP, count uint32) net.IP {
	i := IP4ToUint32(ip)
	d := i + count

	// check for overflow
	if d < i {
		return generateNetLimits(4, 255)
	}
	return Uint32ToIP4(d)
}

// IncrementIP6By returns a net.IP that is greater than the supplied net.IP by
// the supplied integer value. If you overflow the IP space it will return the
// (meaningless in this context) all-ones address
func IncrementIP6By(ip net.IP, count *big.Int) net.IP {
	z := IPToBigint(ip)
	z.Add(z, count)
	return BigintToIP6(z)
}

// IPToHexString returns the given net.IP as a hexadecimal string. This is the
// default stringer format for v6 net.IP
func IPToHexString(ip net.IP) string {
	if EffectiveVersion(ip) == 4 {
		return hex.EncodeToString(ForceIP4(ip))
	}
	return ip.String()
}

// IP4ToUint32 converts a net.IPv4 to a uint32.
func IP4ToUint32(ip net.IP) uint32 {
	if EffectiveVersion(ip) != 4 {
		return 0
	}

	i := binary.BigEndian.Uint32(ForceIP4(ip))
	return i
}

// IPToBigint converts a net.IP to big.Int.
func IPToBigint(ip net.IP) *big.Int {
	z := new(big.Int)
	z.SetBytes(ip)
	return z
}

// NextIP returns a net.IP incremented by one from the input address. This
// function is roughly as fast for v4 as IncrementIP4By(1) but is consistently
// 4x faster on v6 than IncrementIP6By(1). The bundled tests provide
// benchmarks doing so, as well as iterating over the entire v4 address space.
func NextIP(ip net.IP) net.IP {
	var ipn []byte
	if Version(ip) == 4 {
		ipn = make([]byte, 4)
		copy(ipn, ip)
	} else {
		ipn = make([]byte, 16)
		copy(ipn, ip)
	}

	for i := len(ipn) - 1; i >= 0; i-- {
		ipn[i]++
		if ipn[i] > 0 {
			return ipn
		}
	}
	return ip // if we're already at the end of range, don't wrap
}

// PrevIP returns a net.IP decremented by one from the input address. This
// function is roughly as fast for v4 as DecrementIP4By(1) but is consistently
// 4x faster on v6 than DecrementIP6By(1). The bundled tests provide
// benchmarks doing so, as well as iterating over the entire v4 address space.
func PrevIP(ip net.IP) net.IP {
	var ipn []byte
	if Version(ip) == 4 {
		ipn = make([]byte, 4)
		copy(ipn, ip.To4())
	} else {
		ipn = make([]byte, 16)
		copy(ipn, ip)
	}

	for i := len(ipn) - 1; i >= 0; i-- {
		ipn[i]--
		if ipn[i] != 255 {
			return ipn
		}
	}
	return ip // if we're already at beginning of range, don't wrap
}

// Uint32ToIP4 converts a uint32 to an ip4 address and returns it as a net.IP
func Uint32ToIP4(i uint32) net.IP {
	ip := make([]byte, 4)
	binary.BigEndian.PutUint32(ip, i)
	return ip
}

// Version returns 4 if the net.IP contains a v4 address. It will return 6 for
// any v6 address, including the v4-encapsulating v6 address range ::ffff
func Version(ip net.IP) int {
	a := ip.To4()
	if a == nil || len(ip) == 16 {
		return 6
	}
	return 4
}

func generateNetLimits(version int, filler byte) net.IP {
	var b []byte
	if version == 6 {
		version = 16
	}
	b = make([]byte, version)
	for i := range b {
		b[i] = filler
	}
	return b
}