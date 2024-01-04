package types_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/* ...... */
// Adapted from https://www.calhoun.io/how-to-determine-if-a-json-key-has-been-set-to-null-or-not-provided/
type Nullable[T any] struct {
	// Value contains the underlying value of the field. If `Set` is true, and `Null` is false, **??**
	Value *T
	// Set will be true if the field was sent
	Set bool
	// Valid will be true if the value is a valid type - either a value of T or as an explicit `null`
	Valid bool
}

func (t *Nullable[T]) UnmarshalJSON(data []byte) error {
	// If this method is called, there was a value explicitly sent, which was either <nil> or a value of `T`
	t.Set = true

	// we received an explicit value of null
	if string(data) == "null" {
		// which is deemed valid, because we allow either <nil> or a value of `T`
		t.Valid = true
		t.Value = nil
		return nil
	}

	// we received a value of `T`
	var temp T
	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("couldn't unmarshal JSON: %w", err)
	}
	t.Value = &temp
	t.Valid = true
	return nil
}

/* </adapted from https://www.calhoun.io/how-to-determine-if-a-json-key-has-been-set-to-null-or-not-provided/> */
/* ...... */

func (t Nullable[T]) IsSet() bool {
	return t.Set
}

func (t Nullable[T]) IsValid() bool {
	return t.Valid
}

func TestNullable2(t *testing.T) {
	obj := struct {
		ID Nullable[int] `json:"id"`
	}{}

	// when unset
	obj.ID = Nullable[int]{}
	assert.False(t, obj.ID.IsSet())
	assert.False(t, obj.ID.IsValid())
	assert.Nil(t, obj.ID.Value)

	// when empty body
	obj.ID = Nullable[int]{}
	err := json.Unmarshal([]byte("{}"), &obj)
	require.NoError(t, err)
	fmt.Printf("empty\t%v: %v %v\n", obj.ID, obj.ID.IsSet(), obj.ID.IsValid())

	assert.False(t, obj.ID.IsSet())
	assert.False(t, obj.ID.IsValid())
	assert.Nil(t, obj.ID.Value)

	// when explicit null body
	obj.ID = Nullable[int]{}
	err = json.Unmarshal([]byte(`{"id": null}`), &obj)
	require.NoError(t, err)
	fmt.Printf("null\t%v: %v %v\n", obj.ID, obj.ID.IsSet(), obj.ID.IsValid())

	assert.True(t, obj.ID.IsSet())
	assert.True(t, obj.ID.IsValid())
	assert.Nil(t, obj.ID.Value)

	// when explicit zero value
	obj.ID = Nullable[int]{}
	err = json.Unmarshal([]byte(`{"id": 0}`), &obj)
	require.NoError(t, err)
	fmt.Printf("zero\t%v: %v %v\n", obj.ID, obj.ID.IsSet(), obj.ID.IsValid())
	assert.True(t, obj.ID.IsSet())
	assert.True(t, obj.ID.IsValid())
	if assert.NotNil(t, obj.ID.Value) {
		assert.Equal(t, 0, *obj.ID.Value)
	}

	// when explicit value
	obj.ID = Nullable[int]{}
	err = json.Unmarshal([]byte(`{"id": 1230}`), &obj)
	require.NoError(t, err)
	fmt.Printf("val\t%v: %v %v\n", obj.ID, obj.ID.IsSet(), obj.ID.IsValid())
	assert.True(t, obj.ID.IsSet())
	assert.True(t, obj.ID.IsValid())
	if assert.NotNil(t, obj.ID.Value) {
		assert.Equal(t, 1230, *obj.ID.Value)
	}
	assert.True(t, obj.ID.Valid)
}
