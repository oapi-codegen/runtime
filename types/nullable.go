package types

import "encoding/json"

type Nullable[T any] struct {
	Value T
	Set   bool
	Null  bool
}

func (t *Nullable[T]) UnmarshalJSON(data []byte) error {
	t.Set = true
	return json.Unmarshal(data, &t.Value)
}

func (t Nullable[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Nullable[T]) IsNullDefined() bool {
	return t.Set && t.Value == nil
}

func (t *Nullable[T]) HasValue() bool {
	return t.Set && t.Value != nil
}
