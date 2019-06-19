## Ideas for future release
This document is an informal roadmap for `iplib`; completing these probably
gets this library to `1.0.0`

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

#### Tests have gotten out of hand
Always added to, never refactored. They need to be refactored. Also should try
to standardize what they want to test (beginning, middle, end of normal range
and beginning, end of address space).

#### IPv4 vs IPv6

the standard library uses a one-size-fits-all solution to the handling of
IP addresses and netblocks. There are some problems with this approach but
mostly it works fine. This library has followed suit but I think that may
be a mistake: at `iplib`s level of detail v4 and v6 have different concerns
and treating them the same isn't really working. For example:

- on v4 networks the first (network) and last (broadcast) addresses are
  special and are treated different. This is not true for v6.

- on v4 all 32bits of address are important to routing. For v6 this is not
  the case, and several documents from the IETF, RIPE and IANA propose
  splitting v6 addresses into a routing portion and an identity portion,
  for most practical intents, based on these prescriptions, a /64 *is*
  an address, at least from the perspective of an allocation manager which
  is what `iplib.Net` is

- on v4 it is important to have dynamically-sized subnets. The limited address
  space punishes inefficiency. Even taking the approach discussed above 64bits
  of address space is 18 quintillion addresses, but the size of the blocks
  present other challenges. Because of this, even address distributions make
  far more sense. [RIPE "Best Current Operational Practice for Operators"](https://www.ripe.net/publications/docs/ripe-690#4-2--prefix-assignment-options)
  section 4.2 recommends having all subnets conform to nibble boundaries for
  legibility and to make DNS easier to manage

- v4 fits in a `uint32` and, again using the above approach, an IPv6 space
  where only the first 64bits were relevant would fit into a `uint64` and
  remove a dependency on `math/big` that doesn't do anything useful (except
  in ip-to-integer conversion).

For these reasons I think there should be multiple implementations of
`iplib.Net`, one for v4 and another for v6. Mostly  this would mean removing
the v6 functions from the existing `Net` and creating a new `Net6` that
followed the RIPE BCOP guidelines.

#### RFC1918
The most important address-space on the (IPv4) internet is the RFC1918 private
address block designation that is effectively on the inside of every home and
institutional network in the world. It might be a good idea to have a specific
function in `iana` that returns true just for those addresses.
