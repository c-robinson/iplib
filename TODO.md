## Ideas for future release
This document is an informal roadmap for `iplib`.

#### Standardize output
Both for `net` and `iplib` stringified outputs can vary: sometimes punctuation
or separators are used and other times they are not. Each output format should
have a canonical output format and some way to (non)prettify it.

#### Net.PreviousNetWithoutOverlap()
Depending on what `Net`s netmask is and what the proposed netmask for the new
object is, `PreviousNet()` might return a network that overlaps the current
one. This is fine but there should be an explicit function (or parameter to
the existing function) to force `PreviousNet()` to return a new network in an
entirely separate address-space, even if this means the two are not actually
adjacent.

#### Net.PreviousNetAtBestSize()
In-line with above there should be a function, or a way to tell the existing
function, to create a new object that is directly adjacent to it at whatever
the netmask has to be to do that (perhaps `Net.PreviousNet(0)`).

#### NewNetBetween is terrible
Pretty much that. If the problems with `PreviousNet()` are solved it probably
provides a fix for this as well.

#### Return errors
Following from the example set by `net` iplib sometimes returns `nil` when it
cannot do something with an IP address. This is fine for the core libraries
but probably not ok otherwise and errors should be returned instead.

#### Tests have gotten out of hand
Always added to, never refactored. They need to be refactored. Also should try
to standardize what they want to test (beginning, middle, end of normal range
and beginning, end of address space).

#### IPv6-specific functions
There's a lot of handling for v4-specific needs (broadcast address, next/prev
IP) but almost none for IPv6's new concerns and that should be fixed. Ideas
include:

- functions for subnetting, prev/next netblock that abide by nibble boundaries
  per [RIPE "Best Current Operational Practice for Operators"](https://www.ripe.net/publications/docs/ripe-690#4-2--prefix-assignment-options)
  section 4.2

- following from the above it might be nice to be able to set some allocation
  boundaries inside an `iplib.Net` object in the v6 context, such as the total
  size of allocatable space (ISPs are typically granted a /32) and the size to
  be granted to Customer Premises Equipment (CPE) (preferably a /48, but might
  change based on customer size/site count). The workflow would then allow
  iterating within the ISP netblock, assigning subnets at the given CPE size
  while respecting the nibble boundary.

- functions for allocating /64s as if they were IPs as per [RFC7934 section 6](https://tools.ietf.org/html/rfc7934#section-6)

- functions for generating interface identifiers for link-local and global
  use based on Modified IEEE EUI-64 hardware addresses as described in
  [RFC 4291](https://tools.ietf.org/html/rfc4291#section-2.5.1)

- The list of RFCs starting with RFC 4941 (and containing at least 7217 
  and 8064) describe mechanisms for enhancing the privacy of self-generated
  addresses by pseudo-randomly modifying the last 64 bits of an address.
  This might make for an interesting sub-module

#### RFC1918
The most important address-space on the (IPv4) internet is the RFC1918 private
address block designation that is effectively on the inside of every home and
institutional network in the world. It might be a good idea to have a specific
function in `iana` that returns true just for those addresses.

#### GoDoc fixes
There are some small problems with godoc rendering that need fixing.
