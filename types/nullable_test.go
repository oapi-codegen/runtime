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

func TestSimpleString_IsDefined(t *testing.T) {
	type testCase struct {
		name      string
		jsonInput []byte
		wantNull  bool
	}
	tests := []testCase{
		{
			name:      "simple object: set name to some non null value",
			jsonInput: []byte(`{"name":"yolo"}`),
			wantNull:  false,
		},

		{
			name:      "simple object: set name to empty string value",
			jsonInput: []byte(`{"name":""}`),
			// since name field is present in JSON, want defined to be true
			wantNull: false,
		},

		{
			name:      "simple object: set name to null value",
			jsonInput: []byte(`{"name":null}`),
			// since name field is present in JSON, want defined to be true
			wantNull: true,
		},
		/*
			Note that it is not possible to differentiate b/w `{"name":""}` and `{"name":null}`
			as both will result in defined to be true but the value will always be the zero
			value and hence cannot tell which one was null
		*/
		{
			name:      "simple object: do not provide name in json data",
			jsonInput: []byte(`{}`),
			// since name field is present in JSON, want defined to be false
			wantNull: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			var obj SimpleString
			err := json.Unmarshal(tt.jsonInput, &obj)
			assert.NoError(t, err)
			assert.Equalf(t, tt.wantNull, obj.Name.IsNull(), "IsNull()")
			fmt.Println(obj.Name.Get())
		})
	}
}

type SimpleInt struct {
	// cannot decide if it was provided with `null` value in json
	ReplicaCount Nullable[int] `json:"replicaCount"`
}

func TestSimpleInt_IsDefined(t *testing.T) {
	type testCase struct {
		name      string
		jsonInput []byte
		wantNull  bool
	}
	tests := []testCase{
		{
			name:      "simple object: set name to some non null value",
			jsonInput: []byte(`{"replicaCount":1}`),
			wantNull:  false,
		},

		{
			name:      "simple object: set name to empty value",
			jsonInput: []byte(`{"replicaCount":0}`),
			wantNull:  false,
		},

		{
			name:      "simple object: set name to null value",
			jsonInput: []byte(`{"replicaCount":null}`),
			wantNull:  true,
		},

		{
			name:      "simple object: do not provide name in json data",
			jsonInput: []byte(`{}`),
			wantNull:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			var obj SimpleInt
			err := json.Unmarshal(tt.jsonInput, &obj)
			assert.NoError(t, err)
			assert.Equalf(t, tt.wantNull, obj.ReplicaCount.IsNull(), "IsNull()")
		})
	}
}

type SimplePointerInt struct {
	// cannot decide if it was provided with `null` value in json
	ReplicaCount Nullable[*int] `json:"replicaCount"`
}

func TestSimplePointerInt_IsDefined(t *testing.T) {
	type testCase struct {
		name      string
		jsonInput []byte
		wantNull  bool
	}
	tests := []testCase{
		{
			name:      "simple object: set name to some non null value",
			jsonInput: []byte(`{"replicaCount":1}`),
			wantNull:  false,
		},

		{
			name:      "simple object: set name to empty value",
			jsonInput: []byte(`{"replicaCount":0}`),
			wantNull:  false,
		},

		{
			name:      "simple object: set name to null value",
			jsonInput: []byte(`{"replicaCount":null}`),
			wantNull:  true,
		},

		{
			name:      "simple object: do not provide name in json data",
			jsonInput: []byte(`{}`),
			wantNull:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			var obj SimplePointerInt
			err := json.Unmarshal(tt.jsonInput, &obj)
			assert.NoError(t, err)
			assert.Equalf(t, tt.wantNull, obj.ReplicaCount.IsNull(), "IsNull()")
		})
	}
}
