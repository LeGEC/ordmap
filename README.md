# Ordmap - serializable ordered maps in go

This package contains an ordered `Map` type, which keeps track of the order in which the keys have been created in it, and which is compatible with `encoding/json` and `gopkg.in/yaml.v3` for unmarshalling/marshalling.

It also contains an `Any` type, which can serve as a generic placeholder to unmarshal json or yaml data, and keeping the keys ordered for objects nested at any level in the payload.
