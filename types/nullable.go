package types

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// nullBytes is a JSON null literal
var nullBytes = []byte("null")

// Nullable allows defining that a
// provided `null` in JSON or not
type Nullable[T any] struct {
	// Value contains the underlying value of the field. If `Set` is true, and `Null` is false, **??**
	Value *T
	// Set will be true if the field was sent.
	Set bool
}

// Worst case we would have Unmarshal work correctly, and Marshal be broken
// https://stackoverflow.com/questions/70025330/how-to-allow-omitempty-only-unmarshal-and-not-when-marshal

func (t *Nullable[T]) IsSet() bool {
	if t == nil {
		return false
	}
	return t.Set
}

// UnmarshalJSON implements the Unmarshaler interface.
func (t *Nullable[T]) UnmarshalJSON(data []byte) error {
	fmt.Println(data)
	t.Set = true
	if bytes.Equal(data, nullBytes) {
		// t.Null = true
		return nil
	}
	// fmt.Printf("data: %v\n", data)
	// fmt.Printf("t.Value: %v\n", t.Value)
	var tt T
	if err := json.Unmarshal(data, &tt); err != nil {
		return fmt.Errorf("couldn't unmarshal JSON: %w", err)
	}
	// fmt.Printf("t.Value: %v\n", t.Value)
	t.Value = &tt
	// fmt.Printf("t.Value: %v\n", t.Value)
	// fmt.Printf("t.Value: %v\n", *t.Value)
	// t.Null = false
	return nil
}

// MarshalJSON implements the  Marshaler interface.
func (t Nullable[T]) MarshalJSON() ([]byte, error) {
	// TODO
	// TODO
	// TODO
	// if !t.Set {
	// 	// return []byte(""), nil
	// 	return nil, nil
	// }
	// TODO
	// TODO
	// TODO

	if t.IsNull() {
		return nullBytes, nil
	}
	return json.Marshal(t.Value)
}

// IsNull returns true if the value is explicitly provided `null` in json
func (t *Nullable[T]) IsNull() bool {
	if t == nil {
		return false
	}

	return t.Value == nil
}

// Get retrieves the value of underlying nullable field, and indicates whether the value was set or not.
// If `set == false`, then `value` can be ignored
// If `set == true` and `value == nil`: the field was sent explicitly with the value `null`
// If `set == true` and `value != nil`: the field was sent with the contents at `*value`
func (t *Nullable[T]) Get() (value *T, set bool) {
	return t.Value, t.Set
}
