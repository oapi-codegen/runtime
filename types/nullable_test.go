package types

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleNullable_foo() {
	obj := struct {
		ID *Nullable[int] `json:"id,omitempty"`
	}{}
	fmt.Printf("obj.ID.IsNull(): %v\n", obj.ID.IsNull())
	fmt.Printf("obj.ID.IsSet(): %v\n", obj.ID.IsSet())
	// Output:
}

func ExampleNullable_marshal() {
	obj := struct {
		ID *Nullable[int] `json:"id,omitempty"`
	}{}

	// when it's not set
	b, err := json.Marshal(obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf(`JSON: %s`+"\n", b)
	fmt.Println("---")

	// when it's set explicitly to nil
	obj.ID = &Nullable[int]{}
	obj.ID.Value = nil
	obj.ID.Set = true

	b, err = json.Marshal(obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf(`JSON: %s`+"\n", b)
	fmt.Println("---")

	// when it's set explicitly to the zero value
	var v int
	obj.ID.Value = &v
	obj.ID.Set = true

	b, err = json.Marshal(obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf(`JSON: %s`+"\n", b)
	fmt.Println("---")

	// when it's set explicitly to a specific value
	v = 12345
	obj.ID.Value = &v
	obj.ID.Set = true

	b, err = json.Marshal(obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf(`JSON: %s`+"\n", b)
	fmt.Println("---")

	// Output:
	// JSON: {}
	// ---
	// JSON: {"id":null}
	// ---
	// JSON: {"id":0}
	// ---
	// JSON: {"id":12345}
	// ---
}

func ExampleNullable_unmarshal() {
	obj := struct {
		Name Nullable[string] `json:"name"`
	}{}

	// when it's not set
	err := json.Unmarshal([]byte(`
		{
		}
		`), &obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("obj.Name.Set: %v\n", obj.Name.Set)
	fmt.Printf("obj.Name.Value: %v\n", obj.Name.Value)
	fmt.Println("---")

	// when it's set explicitly to nil
	err = json.Unmarshal([]byte(`
		{
		"name": null
		}
		`), &obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("obj.Name.Set: %v\n", obj.Name.Set)
	fmt.Printf("obj.Name.Value: %v\n", obj.Name.Value)
	fmt.Println("---")

	// when it's set explicitly to the zero value
	err = json.Unmarshal([]byte(`
		{
		"name": ""
		}
		`), &obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("obj.Name.Set: %v\n", obj.Name.Set)
	if obj.Name.Value == nil {
		fmt.Println("Error: expected obj.Name.Value to have a value, but was <nil>")
		return
	}
	fmt.Printf("obj.Name.Value: %#v\n", *obj.Name.Value)
	fmt.Println("---")

	// when it's set explicitly to a specific value
	err = json.Unmarshal([]byte(`
		{
		"name": "foo"
		}
		`), &obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("obj.Name.Set: %v\n", obj.Name.Set)
	if obj.Name.Value == nil {
		fmt.Println("Error: expected obj.Name.Value to have a value, but was <nil>")
		return
	}
	fmt.Printf("obj.Name.Value: %#v\n", *obj.Name.Value)
	fmt.Println("---")

	// Output:
	// obj.Name.Set: false
	// obj.Name.Value: <nil>
	// ---
	// obj.Name.Set: true
	// obj.Name.Value: <nil>
	// ---
	// obj.Name.Set: true
	// obj.Name.Value: ""
	// ---
	// obj.Name.Set: true
	// obj.Name.Value: "foo"
	// ---
}

type SimpleString struct {
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
			assert.Equalf(t, tt.wantSet, obj.Name.Set, "Set")
			fmt.Println(obj.Name.Get())
		})
	}
}

type SimpleInt struct {
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
			assert.Equalf(t, tt.wantSet, obj.ReplicaCount.Set, "Set")
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
			assert.Equalf(t, tt.wantSet, obj.ReplicaCount.Set, "Set")
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
				require.NotNilf(t, obj.SimpleInt.Value, "NN")
				require.NotNilf(t, obj.SimpleString.Value, "NN")
				require.NotNilf(t, obj.StringList.Value, "NN")
				assert.Equalf(t, false, obj.SimpleInt.Value.ReplicaCount.Set, "replica count should not be set")
				assert.Equalf(t, false, obj.SimpleInt.Value.ReplicaCount.IsNull(), "replica count should not be null")
				assert.Equalf(t, false, obj.SimpleString.Value.Name.Set, "name should not be set")
				assert.Equalf(t, false, obj.SimpleString.Value.Name.IsNull(), "name should not be null")
				assert.Equalf(t, false, obj.StringList.Set, "string list should not be set")
				assert.Equalf(t, false, obj.StringList.IsNull(), "string list should not be null")
			},
		},

		{
			name:      "replica count having non null value",
			jsonInput: []byte(`{"simple_int":{"replicaCount":1}}`),
			assert: func(obj TestComplex, t *testing.T) {
				require.NotNilf(t, obj.SimpleInt.Value, "NN")
				require.NotNilf(t, obj.SimpleString.Value, "NN")
				assert.Equalf(t, false, obj.SimpleInt.Value.ReplicaCount.IsNull(), "replica count should NOT be null")
				assert.Equalf(t, true, obj.SimpleInt.Value.ReplicaCount.Set, "replica count should be set")
				assert.Equalf(t, false, obj.SimpleString.Value.Name.Set, "name should NOT be set")
				assert.Equalf(t, false, obj.SimpleString.Value.Name.IsNull(), "name should NOT be null")
				gotValue, isSet := obj.SimpleInt.Value.ReplicaCount.Get()
				assert.Equalf(t, true, isSet, "replica count should NOT be null")
				assert.Equalf(t, 1, *gotValue, "replica count should be 1")
			},
		},

		{
			name:      "string list having null value",
			jsonInput: []byte(`{"string_list": null}`),
			assert: func(obj TestComplex, t *testing.T) {
				require.NotNilf(t, obj.StringList.Value, "NN")
				assert.Equalf(t, true, obj.StringList.Set, "string_list should be set")
				assert.Equalf(t, true, obj.StringList.IsNull(), "string_list should be null")
			},
		},

		{
			name:      "string list having non null value",
			jsonInput: []byte(`{"string_list": ["foo", "bar"]}`),
			assert: func(obj TestComplex, t *testing.T) {
				assert.Equalf(t, true, obj.StringList.Set, "string_list should be set")
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
				assert.Equalf(t, true, obj.StringList.Set, "string_list should be set")
				assert.Equalf(t, false, obj.StringList.IsNull(), "string_list should not be null")
				gotStringList, isNull := obj.StringList.Get()
				assert.Equalf(t, false, isNull, "string_list should not be null")
				assert.Equalf(t, []string{}, gotStringList, "string_list should have the values as provided in the jSON")

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj TestComplex
			err := json.Unmarshal(tt.jsonInput, &obj)
			assert.NoError(t, err)
			tt.assert(obj, t)
		})
	}
}
