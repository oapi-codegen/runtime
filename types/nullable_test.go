package types

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SimpleString struct {
	// cannot decide if it was provided with `null` value in json
	Name Optional[string] `json:"name"`
}

func TestSimpleString_IsDefined(t *testing.T) {
	type testCase struct {
		name        string
		jsonInput   []byte
		wantNull    bool
		wantDefined bool
	}
	tests := []testCase{
		{
			name:      "simple object: set name to some non null value",
			jsonInput: []byte(`{"name":"yolo"}`),
			// since name field is present in JSON, want defined to be true
			wantDefined: true,
		},

		{
			name:      "simple object: set name to empty string value",
			jsonInput: []byte(`{"name":""}`),
			// since name field is present in JSON, want defined to be true
			wantDefined: true,
		},

		{
			name:      "simple object: set name to null value",
			jsonInput: []byte(`{"name":null}`),
			// since name field is present in JSON, want defined to be true
			wantDefined: true,
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
			wantDefined: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			var obj SimpleString
			err := json.Unmarshal(tt.jsonInput, &obj)
			assert.NoError(t, err)
			assert.Equalf(t, tt.wantDefined, obj.Name.IsDefined(), "IsDefined()")

		})
	}
}

type SimpleStringPointer struct {
	// can decide if it was provided with `null` value in json
	Name Optional[*string] `json:"name"`
}

func TestSimpleStringPointer_IsDefined(t *testing.T) {
	type testCase struct {
		name        string
		jsonInput   []byte
		wantDefined bool
		wantNull    bool
	}
	tests := []testCase{
		{
			name:      "simple object: set name to some non null value",
			jsonInput: []byte(`{"name":"yolo"}`),
			// since name field is present in JSON, want defined to be true
			wantDefined: true,
		},

		{
			name:      "simple object: set name to empty string value",
			jsonInput: []byte(`{"name":""}`),
			// since name field is present in JSON, want defined to be true
			wantDefined: true,
		},

		{
			name:      "simple object: set name to null value",
			jsonInput: []byte(`{"name":null}`),
			// since name field is present in JSON, want defined to be true
			wantDefined: true,
			wantNull:    true,
		},
		/*
				Note that it is possible to differentiate b/w `{"name":""}` and `{"name":null}`
				as both will result in defined to be true but the value will always be zero
				value for `{"name":""}` and nil for `{"name":null}`.
			    We could tell which one was null because of (pointer) Nullable[*string]
		*/

		{
			name:      "simple object: do not provide name in json data",
			jsonInput: []byte(`{}`),
			// since name field is present in JSON, want defined to be false
			wantDefined: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			var obj SimpleStringPointer
			err := json.Unmarshal(tt.jsonInput, &obj)
			assert.NoError(t, err)
			assert.Equalf(t, tt.wantDefined, obj.Name.IsDefined(), "IsDefined()")
			gotNull := false
			if obj.Name.IsDefined() && obj.Name.Value == nil {
				gotNull = true
			}
			assert.Equalf(t, tt.wantNull, gotNull, "Null Check")
		})
	}
}

type SimpleInt struct {
	// cannot decide if it was provided with `null` value in json
	ReplicaCount Optional[int] `json:"replicaCount"`
}

func TestSimpleInt_IsDefined(t *testing.T) {
	type testCase struct {
		name        string
		jsonInput   []byte
		wantDefined bool
	}
	tests := []testCase{
		{
			name:        "simple object: set name to some non null value",
			jsonInput:   []byte(`{"replicaCount":1}`),
			wantDefined: true,
		},

		{
			name:      "simple object: set name to empty value",
			jsonInput: []byte(`{"replicaCount":0}`),
			// since name field is present in JSON want defined to be true
			wantDefined: true,
		},

		{
			name:      "simple object: set name to null value",
			jsonInput: []byte(`{"replicaCount":null}`),
			// since name field is present in JSON want defined to be true
			wantDefined: true,
		},

		{
			name:      "simple object: do not provide name in json data",
			jsonInput: []byte(`{}`),
			// since name field is NOT present in JSON want defined to be false
			wantDefined: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			var obj SimpleInt
			err := json.Unmarshal(tt.jsonInput, &obj)
			assert.NoError(t, err)
			assert.Equalf(t, tt.wantDefined, obj.ReplicaCount.IsDefined(), "IsDefined()")
		})
	}
}

type SimpleIntPointer struct {
	// can decide if it was provided with `null` value in json
	ReplicaCount Optional[*int] `json:"replicaCount"`
}

func TestSimpleIntPointer_IsDefined(t *testing.T) {
	type testCase struct {
		name        string
		jsonInput   []byte
		wantDefined bool
		wantNull    bool
	}
	tests := []testCase{
		{
			name:      "simple object: set name to some non null value",
			jsonInput: []byte(`{"replicaCount":1}`),
			// since replicaCount field is present in JSON want defined to be true
			wantDefined: true,
		},

		{
			name:      "simple object: set name to empty value",
			jsonInput: []byte(`{"replicaCount":0}`),
			// since replicaCount field is present in JSON want defined to be true
			wantDefined: true,
		},

		{
			name:      "simple object: set name to null value",
			jsonInput: []byte(`{"replicaCount":null}`),
			// since replicaCount field is present in JSON want defined to be true
			wantDefined: true,
			wantNull:    true,
		},

		{
			name:      "simple object: do not provide name in json data",
			jsonInput: []byte(`{}`),
			// since replicaCount field is NOT present in JSON want defined to be false
			wantDefined: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			var obj SimpleIntPointer
			err := json.Unmarshal(tt.jsonInput, &obj)
			assert.NoError(t, err)
			assert.Equalf(t, tt.wantDefined, obj.ReplicaCount.IsDefined(), "IsDefined()")
			gotNull := false
			if obj.ReplicaCount.IsDefined() && obj.ReplicaCount.Value == nil {
				gotNull = true
			}
			assert.Equalf(t, tt.wantNull, gotNull, "Null Check")
		})
	}
}
