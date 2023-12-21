package types

import "encoding/json"

// Nullable type which can distinguish between an explicit `null` vs not provided
// in JSON when un-marshaled to go type.
type Nullable[T any] struct {
	// Value is the actual value of the field.
	Value *T
	// Defined indicates that the field was provided in JSON if it is true.
	// If a field is not provided in JSON, then `Defined` is false and `Value`
	// contains the `zero-value` of the field type e.g "" for string,
	// 0 for int, nil for pointer etc
	Defined bool
}

// UnmarshalJSON implements the Unmarshaler interface.
func (t *Nullable[T]) UnmarshalJSON(data []byte) error {
	t.Defined = true
	return json.Unmarshal(data, &t.Value)
}

// MarshalJSON implements the  Marshaler interface.
func (t Nullable[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

// IsNull returns true if the value is explicitly provided `null` in json
func (t *Nullable[T]) IsNull() bool {
	return t.IsDefined() && t.Value == nil
}

// IsDefined returns true if the value is explicitly provided in json
func (t *Nullable[T]) IsDefined() bool {
	return t.Defined
}
