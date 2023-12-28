module github.com/c-robinson/iplib/v2

go 1.20

replace (
	github.com/c-robinson/iplib/v2/iana => ./iana
	github.com/c-robinson/iplib/v2/iid => ./iid
)

require lukechampine.com/uint128 v1.3.0
