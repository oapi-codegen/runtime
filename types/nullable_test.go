package types

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SimpleString struct {
	// cannot decide if it was provided with `null` value in json
	Name Nullable[string] `json:"name"`
}

func TestSimpleString(t *testing.T) {
	type testCase struct {
		name      string
		jsonInput []byte
		wantNull  bool
		wantSet   bool
	}
	tests := []testCase{
		{
			name:      "simple object: set name to some non null value",
			jsonInput: []byte(`{"name":"yolo"}`),
			wantNull:  false,
			wantSet:   true,
		},

		{
			name:      "simple object: set name to empty string value",
			jsonInput: []byte(`{"name":""}`),
			wantNull:  false,
			wantSet:   true,
		},

		{
			name:      "simple object: set name to null value",
			jsonInput: []byte(`{"name":null}`),
			wantNull:  true,
			wantSet:   true,
		},

		{
			name:      "simple object: do not provide name in json data",
			jsonInput: []byte(`{}`),
			wantNull:  false,
			wantSet:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			var obj SimpleString
			err := json.Unmarshal(tt.jsonInput, &obj)
			assert.NoError(t, err)
			assert.Equalf(t, tt.wantNull, obj.Name.IsNull(), "IsNull()")
			assert.Equalf(t, tt.wantSet, obj.Name.IsSet(), "IsSet()")
			fmt.Println(obj.Name.Get())
		})
	}
}

type SimpleInt struct {
	// cannot decide if it was provided with `null` value in json
	ReplicaCount Nullable[int] `json:"replicaCount"`
}

func TestSimpleInt(t *testing.T) {
	type testCase struct {
		name      string
		jsonInput []byte
		wantNull  bool
		wantSet   bool
	}
	tests := []testCase{
		{
			name:      "simple object: set name to some non null value",
			jsonInput: []byte(`{"replicaCount":1}`),
			wantNull:  false,
			wantSet:   true,
		},

		{
			name:      "simple object: set name to empty value",
			jsonInput: []byte(`{"replicaCount":0}`),
			wantNull:  false,
			wantSet:   true,
		},

		{
			name:      "simple object: set name to null value",
			jsonInput: []byte(`{"replicaCount":null}`),
			wantNull:  true,
			wantSet:   true,
		},

		{
			name:      "simple object: do not provide name in json data",
			jsonInput: []byte(`{}`),
			wantNull:  false,
			wantSet:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			var obj SimpleInt
			err := json.Unmarshal(tt.jsonInput, &obj)
			assert.NoError(t, err)
			assert.Equalf(t, tt.wantNull, obj.ReplicaCount.IsNull(), "IsNull()")
			assert.Equalf(t, tt.wantSet, obj.ReplicaCount.IsSet(), "IsSet()")
		})
	}
}

type SimplePointerInt struct {
	// cannot decide if it was provided with `null` value in json
	ReplicaCount Nullable[*int] `json:"replicaCount"`
}

func TestSimplePointerInt(t *testing.T) {
	type testCase struct {
		name      string
		jsonInput []byte
		wantNull  bool
		wantSet   bool
	}
	tests := []testCase{
		{
			name:      "simple object: set name to some non null value",
			jsonInput: []byte(`{"replicaCount":1}`),
			wantNull:  false,
			wantSet:   true,
		},

		{
			name:      "simple object: set name to empty value",
			jsonInput: []byte(`{"replicaCount":0}`),
			wantNull:  false,
			wantSet:   true,
		},

		{
			name:      "simple object: set name to null value",
			jsonInput: []byte(`{"replicaCount":null}`),
			wantNull:  true,
			wantSet:   true,
		},

		{
			name:      "simple object: do not provide name in json data",
			jsonInput: []byte(`{}`),
			wantNull:  false,
			wantSet:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			var obj SimplePointerInt
			err := json.Unmarshal(tt.jsonInput, &obj)
			assert.NoError(t, err)
			assert.Equalf(t, tt.wantNull, obj.ReplicaCount.IsNull(), "IsNull()")
			assert.Equalf(t, tt.wantSet, obj.ReplicaCount.IsSet(), "IsSet()")
		})
	}
}

type TestComplex struct {
	SimpleInt    Nullable[SimpleInt]    `json:"simple_int"`
	SimpleString Nullable[SimpleString] `json:"simple_string"`
	StringList   Nullable[[]string]     `json:"string_list"`
}

func TestMixed(t *testing.T) {
	type testCase struct {
		name      string
		jsonInput []byte
		assert    func(obj TestComplex, t *testing.T)
	}
	tests := []testCase{
		{
			name:      "empty json input",
			jsonInput: []byte(`{}`),
			assert: func(obj TestComplex, t *testing.T) {
				assert.Equalf(t, false, obj.SimpleInt.Value.ReplicaCount.IsSet(), "replica count should not be set")
				assert.Equalf(t, false, obj.SimpleInt.Value.ReplicaCount.IsNull(), "replica count should not be null")
				assert.Equalf(t, false, obj.SimpleString.Value.Name.IsSet(), "name should not be set")
				assert.Equalf(t, false, obj.SimpleString.Value.Name.IsNull(), "name should not be null")
				assert.Equalf(t, false, obj.StringList.IsSet(), "string list should not be set")
				assert.Equalf(t, false, obj.StringList.IsNull(), "string list should not be null")
			},
		},

		{
			name:      "replica count having non null value",
			jsonInput: []byte(`{"simple_int":{"replicaCount":1}}`),
			assert: func(obj TestComplex, t *testing.T) {
				assert.Equalf(t, false, obj.SimpleInt.Value.ReplicaCount.IsNull(), "replica count should NOT be null")
				assert.Equalf(t, true, obj.SimpleInt.Value.ReplicaCount.IsSet(), "replica count should be set")
				assert.Equalf(t, false, obj.SimpleString.Value.Name.IsSet(), "name should NOT be set")
				assert.Equalf(t, false, obj.SimpleString.Value.Name.IsNull(), "name should NOT be null")
				gotValue, isNull := obj.SimpleInt.Value.ReplicaCount.Get()
				assert.Equalf(t, false, isNull, "replica count should NOT be null")
				assert.Equalf(t, 1, gotValue, "replica count should be 1")
			},
		},

		{
			name:      "string list having null value",
			jsonInput: []byte(`{"string_list": null}`),
			assert: func(obj TestComplex, t *testing.T) {
				assert.Equalf(t, true, obj.StringList.IsSet(), "string_list should be set")
				assert.Equalf(t, true, obj.StringList.IsNull(), "string_list should be null")
			},
		},

		{
			name:      "string list having non null value",
			jsonInput: []byte(`{"string_list": ["foo", "bar"]}`),
			assert: func(obj TestComplex, t *testing.T) {
				assert.Equalf(t, true, obj.StringList.IsSet(), "string_list should be set")
				assert.Equalf(t, false, obj.StringList.IsNull(), "string_list should not be null")
				gotStringList, isNull := obj.StringList.Get()
				assert.Equalf(t, false, isNull, "string_list should not be null")
				assert.Equalf(t, []string{"foo", "bar"}, gotStringList, "string_list should have the values as provided in the jSON")

			},
		},

		{
			name:      "set string list having empty value",
			jsonInput: []byte(`{"string_list":[]}`),
			assert: func(obj TestComplex, t *testing.T) {
				assert.Equalf(t, true, obj.StringList.IsSet(), "string_list should be set")
				assert.Equalf(t, false, obj.StringList.IsNull(), "string_list should not be null")
				gotStringList, isNull := obj.StringList.Get()
				assert.Equalf(t, false, isNull, "string_list should not be null")
				assert.Equalf(t, []string{}, gotStringList, "string_list should have the values as provided in the jSON")

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			var obj TestComplex
			err := json.Unmarshal(tt.jsonInput, &obj)
			assert.NoError(t, err)
			tt.assert(obj, t)
		})
	}
}
