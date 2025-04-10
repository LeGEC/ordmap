# Ordmap - serializable ordered maps in go

This package contains an ordered `Map` type, which keeps track of the order in which the keys have been created in it, and which is compatible with `encoding/json` and `gopkg.in/yaml.v3` for unmarshalling/marshalling.

It also contains an `Any` type, which can serve as a generic placeholder to unmarshal json or yaml data, and keeping the keys ordered for objects nested at any level in the payload.

Check [the doc](https://pkg.go.dev/github.com/LeGEC/ordmap) for more details, and examples.

## Goals and non goals

This package aims at being a correct target to unmarshal/marshal json and yaml, and tries to follow as close as possible the behavior of the standard
`json` package, as well as the `gopkg.in/yaml.v3` package.

It was not designed with performance in mind, so it may compare quite badly with other packages performance wise.
