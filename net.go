package iplib

import (
	"math"
	"math/big"
	"net"
	"strings"
)

// Net extends net.IPNet adding a few useful features along the way
type Net struct {
	net.IPNet
	version int
	length  int
}

// NewNet returns a new Net object containing ip at the specified masklen.
func NewNet(ip net.IP, masklen int) Net {
	var maskMax, length int
	version := EffectiveVersion(ip)
	if version == 6 {
		maskMax = 128
		length = 16
	} else {
		maskMax = 32
		length = 4
	}
	mask := net.CIDRMask(masklen, maskMax)
	n := net.IPNet{IP: ip.Mask(mask), Mask: mask}

	return Net{IPNet: n, version: version, length: length}
}

// NewNetBetween takes two net.IP's as input and will return the largest
// netblock that can fit between them (exclusive of the IP's themselves).
// If there is an exact fit it will set a boolean to true, otherwise the bool
// will be false. If no fit can be found (probably because a >= b) an
// ErrNoValidRange will be returned.
func NewNetBetween(a, b net.IP) (Net, bool, error) {
	var exact = false
	v := CompareIPs(a, b)
	if v != -1 {
		return Net{}, exact, ErrNoValidRange
	}

	if Version(a) != Version(b) {
		return Net{}, exact, ErrNoValidRange
	}

	maskMax := 128
	if EffectiveVersion(a) == 4 {
		maskMax = 32
	}

	ipa := NextIP(a)
	ipb := PreviousIP(b)
	for i := 1; i <= maskMax; i++ {
		xnet := NewNet(ipa, i)

		va := CompareIPs(xnet.NetworkAddress(), ipa)
		vb := CompareIPs(xnet.BroadcastAddress(), ipb)
		if va >= 0 && vb <= 0 {
			if va == 0 && vb == 0 {
				exact = true
			}
			return xnet, exact, nil
		}
	}
	return Net{}, exact, ErrNoValidRange
}

// ParseCIDR returns a new Net object. It is a passthrough to net.ParseCIDR
// and will return any error it generates to the caller. There is one major
// difference between how net.IPNet manages addresses and how ipnet.Net does,
// and this function exposes it: net.ParseCIDR *always* returns an IPv6
// address; if given a v4 address it returns the RFC4291 IPv4-mapped IPv6
// address internally, but treats it like v4 in practice. In contrast
// iplib.ParseCIDR will re-encode it as a v4
func ParseCIDR(s string) (net.IP, Net, error) {
	ip, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return ip, Net{}, err
	}
	if strings.Contains(s, ".") {
		masklen, _ := ipnet.Mask.Size()
		return ip[12:15], NewNet(ip.To4(), masklen), err
	}

	return ip, Net{*ipnet, 6, 16}, err
}

// BroadcastAddress returns the broadcast address for the represented network.
// In the context of IPv6 broadcast is meaningless and the value will be
// equivalent to LastAddress().
func (n Net) BroadcastAddress() net.IP {
	a, _ := n.finalAddress()
	return a
}

// ContainsNet returns true if the given Net is contained within the
// represented block
func (n Net) ContainsNet(network Net) bool {
	l1, _ := n.Mask.Size()
	l2, _ := network.Mask.Size()
	return l1 <= l2 && n.Contains(network.NetworkAddress())
}

// Count returns the total number of usable IP addresses in the represented
// network. It will work for IPv6 in cases where the return value can be
// contained in a uint32. For an IPv6-friendly implementation see Count6().
func (n Net) Count() uint32 {
	if n.version == 4 {
		return n.Count4()
	}
	m := big.NewInt(int64(MaxIPv4))
	z := n.Count6()
	if v := z.Cmp(m); v > 0 {
		return MaxIPv4
	}
	return uint32(z.Uint64())
}

// Count4 returns the total number of addresses in the represented netblock
// up to the limit of uin32. It is intended for IPv4 networks.
func (n Net) Count4() uint32 {
	ones, all := n.Mask.Size()
	exp := all - ones
	if exp == 1 {
		return uint32(0) // special handling for /31
	}
	if exp == 0 {
		return uint32(1) // special handling for /32
	}
	return uint32(math.Pow(2, float64(exp))) - 2
}

// Count6 returns the total number of usable IP addresses in the represented
// network. It works fine for v4 and v6 address ranges but is slower than
// Count().
func (n Net) Count6() *big.Int {
	ones, all := n.Mask.Size()
	exp := all - ones
	if exp == 1 {
		return big.NewInt(0)
	}
	if exp == 0 {
		return big.NewInt(1)
	}
	var z, e = big.NewInt(2), big.NewInt(int64(exp))
	if n.version == 6 {
		return z.Exp(z, e, nil)
	}
	return z.Sub(z.Exp(z, e, nil), big.NewInt(2))
}

// Enumerate generates an array of all usable addresses in Net up to the
// given size, starting at the given offset up to a maximum of the max-size
// of uint32. This is sufficient to return the entire v4 space but places an
// arbitrary constraint on v6 netblocks. If size=0 the entire block is
// enumerated.
//
// NOTE: RFC3021 (IPv4) and RFC6164 (IPv6) define a use case for netblocks of
// /31 (for IPv4) and /127 (for IPv6) for use in point-to-point links. For
// this reason enumerating networks at these lengths will return 2 a 2-element
// array even though in the v4 case it would naturally return none.
//
// For consistency, enumerating an IPv4 /32 will return the IP in a 1 element
// array.
func (n Net) Enumerate(size, offset uint32) []net.IP {
	count := n.Count()

	// offset exceeds total, return an empty array
	if offset > count {
		return []net.IP{}
	}

	// size is greater than the number of addresses that can be returned,
	// adjust the size of the slice but keep going
	if size > (count-offset) || size == 0 {
		size = count - offset
	}

	// Handle edge-case mask sizes
	if count == 1 { // Count() returns 1 if host-bits == 0
		return []net.IP{n.IP}

	}
	if count == 0 { // Count() returns 0 if host-bits == 1
		addrList := []net.IP{
			n.NetworkAddress(),
			n.BroadcastAddress(),
		}

		return addrList[offset:]
	}

	netu := IP4ToUint32(n.FirstAddress())
	netu += offset

	addrList := make([]net.IP, size)

	addrList[0] = Uint32ToIP4(netu)
	for i := uint32(1); i <= size-1; i++ {
		addrList[i] = NextIP(addrList[i-1])
	}
	return addrList
}

// FirstAddress returns the first usable address for the represented network
func (n Net) FirstAddress() net.IP {
	if n.version == 6 {
		return n.IP
	}
	i, j := n.Mask.Size()
	if i+2 > j {
		return n.IP
	}
	return NextIP(n.IP)
}

// LastAddress returns the last usable address for the represented network.
// For v6 this is the last address in the block; for v4 it is generally the
// next-to-last address, unless the block is a /31 or /32.
func (n Net) LastAddress() net.IP {
	a, ones := n.finalAddress()

	// if it's v6 return the last address
	if n.version == 6 {
		return a
	}

	// if it's v4 and either a single IP or RFC 3021, return the last address
	if ones >= 31 && n.version == 4 {
		return a
	}

	return PreviousIP(a)
}

// NetworkAddress returns the network address for the represented network, e.g.
// the lowest IP address in the given block
func (n Net) NetworkAddress() net.IP {
	return n.IP
}

// NextIP takes a net.IP as an argument and attempts to increment it by
// one. If the input is outside of the range of the represented network it will
// return an empty net.IP and an ErrAddressOutOfRange. If the resulting address
// is out of range it will return an empty net.IP and an ErrAddressAtEndOfRange.
// If the result is the broadcast address, the address _will_ be returned, but
// so will an ErrBroadcastAddress, to indicate that the address is technically
// outside the usable scope
func (n Net) NextIP(ip net.IP) (net.IP, error) {
	if !n.Contains(ip) {
		return net.IP{}, ErrAddressOutOfRange
	}
	xip := NextIP(ip)
	if !n.Contains(xip) {
		return net.IP{}, ErrAddressAtEndOfRange
	}
	// if this is the broadcast address, return it but warn the caller via error
	if n.BroadcastAddress().Equal(xip) && n.version == 4 {
		return xip, ErrBroadcastAddress
	}
	return xip, nil
}

// NextNet takes a CIDR mask-size as an argument and attempts to create a new
// Net object just after the current Net, at the requested mask length
func (n Net) NextNet(masklen int) Net {
	return NewNet(NextIP(n.BroadcastAddress()), masklen)
}

// PreviousIP takes a net.IP as an argument and attempts to decrement it by
// one. If the input is outside of the range of the represented network it will
// return an empty net.IP and an ErrAddressOutOfRange. If the resulting address
// is out of range it will return an empty net.IP and ErrAddressAtEndOfRange.
// If the result is the network address, the address _will_ be returned, but
// so will an ErrNetworkAddress, to indicate that the address is technically
// outside the usable scope
func (n Net) PreviousIP(ip net.IP) (net.IP, error) {
	if !n.Contains(ip) {
		return net.IP{}, ErrAddressOutOfRange
	}
	xip := PreviousIP(ip)
	if !n.Contains(xip) {
		return net.IP{}, ErrAddressAtEndOfRange
	}
	// if this is the network address, return it but warn the caller via error
	if n.NetworkAddress().Equal(xip) && n.version == 4 {
		return xip, ErrNetworkAddress
	}
	return xip, nil
}

// PreviousNet takes a CIDR mask-size as an argument and creates a new Net
// object just before the current one, at the requested mask length. If the
// specified mask is for a larger network than the current one then the new
// network may encompass the current one, e.g.:
//
// iplib.Net{192.168.4.0/22}.Subnet(21) -> 192.168.0.0/21
//
// In the above case 192.168.4.0/22 is part of 192.168.0.0/21
func (n Net) PreviousNet(masklen int) Net {
	return NewNet(PreviousIP(n.NetworkAddress()), masklen)
}

// Subnet takes a CIDR mask-size as an argument and carves the current Net
// object into subnets of that size, returning them as a []Net. The mask
// provided must be a larger-integer than the current mask. If set to 0 Subnet
// will carve the network in half
//
// Examples:
// Net{192.168.1.0/24}.Subnet(0)  -> []Net{192.168.1.0/25, 192.168.1.128/25}
// Net{192.168.1.0/24}.Subnet(26) -> []Net{192.168.1.0/26, 192.168.1.64/26, 192.168.1.128/26, 192.168.1.192/26}
func (n Net) Subnet(masklen int) ([]Net, error) {
	ones, all := n.Mask.Size()
	if ones > masklen {
		return nil, ErrBadMaskLength
	}

	if masklen == 0 {
		masklen = ones + 1
	}

	mask := net.CIDRMask(masklen, all)
	netlist := []Net{{net.IPNet{n.NetworkAddress(), mask}, n.version, n.length}}

	for CompareIPs(netlist[len(netlist)-1].BroadcastAddress(), n.BroadcastAddress()) == -1 {
		ng := net.IPNet{IP: NextIP(netlist[len(netlist)-1].BroadcastAddress()), Mask: mask}
		netlist = append(netlist, Net{ng, n.version, n.length})
	}
	return netlist, nil
}

// Supernet takes a CIDR mask-size as an argument and returns a Net object
// containing the supernet of the current Net at the requested mask length.
// The mask provided must be a smaller-integer than the current mask. If set
// to 0 Supernet will return the next-largest network
//
// Examples:
// Net{192.168.1.0/24}.Supernet(0)  -> Net{192.168.0.0/23}
// Net{192.168.1.0/24}.Supernet(22) -> Net{Net{192.168.0.0/22}
func (n Net) Supernet(masklen int) (Net, error) {
	ones, all := n.Mask.Size()
	if ones < masklen {
		return Net{}, ErrBadMaskLength
	}

	if masklen == 0 {
		masklen = ones - 1
	}

	mask := net.CIDRMask(masklen, all)
	ng := net.IPNet{IP: n.IP.Mask(mask), Mask: mask}
	return Net{ng, n.version, n.length}, nil
}

// Version returns the version of IP for the enclosed netblock, Either 4 or 6.
func (n Net) Version() int {
	return n.version
}

// Wildcard will return the wildcard mask for a given netmask
func (n Net) Wildcard() net.IPMask {
	wc := make([]byte, len(n.Mask))
	for pos, b := range n.Mask {
		wc[pos] = 0xff - b
	}
	return wc
}

// finalAddress returns the last address in the network. It is private
// because both LastAddress() and BroadcastAddress() rely on it, and both use
// it differently. It returns the last address in the block as well as the
// number of masked bits as an int.
func (n Net) finalAddress() (net.IP, int) {
	a := make([]byte, len(n.IP))
	ones, _ := n.Mask.Size()

	// apply wildcard to network, byte by byte
	wc := n.Wildcard()
	for pos, b := range []byte(n.IP) {
		a[pos] = b + wc[pos]
	}
	return a, ones
}
