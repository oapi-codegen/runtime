package types

import "encoding/json"

// Optional type which can help distinguish between if a value was explicitly
// provided in JSON or not
type Optional[T any] struct {
	// Value is the actual value of the field.
	Value T
	// Defined indicates that the field was provided in JSON if it is true.
	// If a field is not provided in JSON, then `Defined` is false and `Value`
	// contains the `zero-value` of the field type e.g "" for string,
	// 0 for int, nil for pointer etc
	Defined bool
}

// UnmarshalJSON implements the Unmarshaler interface.
func (t *Optional[T]) UnmarshalJSON(data []byte) error {
	t.Defined = true
	return json.Unmarshal(data, &t.Value)
}

// MarshalJSON implements the  Marshaler interface.
func (t Optional[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

// IsDefined returns true if the value is explicitly provided in json
func (t *Optional[T]) IsDefined() bool {
	return t.Defined
}
