# iid
[![Documentation](https://godoc.org/github.com/c-robinson/iplib?status.svg)](http://godoc.org/github.com/c-robinson/iplib/iana)
[![CircleCI](https://circleci.com/gh/c-robinson/iplib/tree/master.svg?style=svg)](https://circleci.com/gh/c-robinson/iplib/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/c-robinson/iplib)](https://goreportcard.com/report/github.com/c-robinson/iplib)
[![Coverage Status](https://coveralls.io/repos/github/c-robinson/iplib/badge.svg?branch=master)](https://coveralls.io/github/c-robinson/iplib?branch=master)

This package implements functions for generating and validating IPv6 Interface
Identifiers (IID's). For the purposes of this module an IID is an IPv6 address
constructed, somehow, from information which uniquely identifies a given
interface on a network, and is unique _within_ that network.

