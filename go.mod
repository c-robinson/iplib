module github.com/c-robinson/iplib/v2

go 1.20

replace (
	github.com/c-robinson/iplib/iana/v2 => ./iana
	github.com/c-robinson/iplib/iid/v2 => ./iid
)

require lukechampine.com/uint128 v1.3.0
