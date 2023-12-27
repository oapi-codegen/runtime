package types

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// nullBytes is a JSON null literal
var nullBytes = []byte("null")

// Nullable type which can help distinguish between if a value was explicitly
// provided `null` in JSON or not
type Nullable[T any] struct {
	Value T
	Set   bool
	Null  bool
}

// UnmarshalJSON implements the Unmarshaler interface.
func (t *Nullable[T]) UnmarshalJSON(data []byte) error {
	t.Set = true
	if bytes.Equal(data, nullBytes) {
		t.Null = true
		return nil
	}
	if err := json.Unmarshal(data, &t.Value); err != nil {
		return fmt.Errorf("couldn't unmarshal JSON: %w", err)
	}
	t.Null = false
	return nil
}

// MarshalJSON implements the  Marshaler interface.
func (t Nullable[T]) MarshalJSON() ([]byte, error) {
	if t.IsNull() {
		return nullBytes, nil
	}
	return json.Marshal(t.Value)
}

// IsNull returns true if the value is explicitly provided `null` in json
func (t *Nullable[T]) IsNull() bool {
	return t.Null
}

// IsSet returns true if the value is provided in json
func (t *Nullable[T]) IsSet() bool {
	return t.Set
}

func (t *Nullable[T]) Get() (value T, null bool) {
	return t.Value, t.IsNull()
}
