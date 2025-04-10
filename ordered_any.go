package ordmap

// Any wraps an 'any' value.
//
// Its purpose is to be used as the target for `json.Umarshal()` or `yaml.Unmarshal()`,
// and to create a go structure which keeps track of the order of the keys as they
// appeared in the initial document.
//
// See the implementation of `Any.UnmarshalJSON()` and `Any.UnmarshalYAML()`
type Any struct {
	v any
}

func (x *Any) V() any {
	return x.v
}
