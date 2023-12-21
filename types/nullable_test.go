package types

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SimpleString struct {
	Name Nullable[string] `json:"name"`
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
			// since name field is present in JSON and is NOT null, want null to be false
			wantNull: false,
			// since name field is present in JSON, want defined to be true
			wantDefined: true,
		},

		{
			name:      "simple object: set name to empty string value",
			jsonInput: []byte(`{"name":""}`),
			// since name field is present in JSON and is NOT null, want null to be false
			wantNull: false,
			// since name field is present in JSON, want defined to be true
			wantDefined: true,
		},

		{
			name:      "simple object: set name to null value",
			jsonInput: []byte(`{"name":null}`),
			// since name field is present in JSON and is  null, want null to be true
			wantNull: true,
			// since name field is present in JSON, want defined to be true
			wantDefined: true,
		},

		{
			name:      "simple object: do not provide name in json data",
			jsonInput: []byte(`{}`),
			// since name field is NOT present in JSON, want null to be false
			wantNull: false,
			// since name field is present in JSON, want defined to be false
			wantDefined: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			var obj SimpleString
			err := json.Unmarshal(tt.jsonInput, &obj)
			assert.NoError(t, err)
			assert.Equalf(t, tt.wantNull, obj.Name.IsNull(), "IsNull()")
			assert.Equalf(t, tt.wantDefined, obj.Name.IsDefined(), "IsDefined()")
		})
	}
}

type SimpleStringPointer struct {
	Name Nullable[*string] `json:"name"`
}

func TestSimpleStringPointer_IsDefined(t *testing.T) {
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
			// since name field is present in JSON and is NOT null, want null to be false
			wantNull: false,
			// since name field is present in JSON, want defined to be true
			wantDefined: true,
		},

		{
			name:      "simple object: set name to empty string value",
			jsonInput: []byte(`{"name":""}`),
			// since name field is present in JSON and is NOT null, want null to be false
			wantNull: false,
			// since name field is present in JSON, want defined to be true
			wantDefined: true,
		},

		{
			name:      "simple object: set name to null value",
			jsonInput: []byte(`{"name":null}`),
			// since name field is present in JSON and is  null, want null to be true
			wantNull: true,
			// since name field is present in JSON, want defined to be true
			wantDefined: true,
		},

		{
			name:      "simple object: do not provide name in json data",
			jsonInput: []byte(`{}`),
			// since name field is NOT present in JSON, want null to be false
			wantNull: false,
			// since name field is present in JSON, want defined to be false
			wantDefined: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			var obj SimpleStringPointer
			err := json.Unmarshal(tt.jsonInput, &obj)
			assert.NoError(t, err)
			assert.Equalf(t, tt.wantNull, obj.Name.IsNull(), "IsNull()")
			assert.Equalf(t, tt.wantDefined, obj.Name.IsDefined(), "IsDefined()")
		})
	}
}

type SimpleInt struct {
	ReplicaCount Nullable[int] `json:"replicaCount"`
}

func TestSimpleInt_IsDefined(t *testing.T) {
	type testCase struct {
		name        string
		jsonInput   []byte
		wantNull    bool
		wantDefined bool
	}
	tests := []testCase{
		{
			name:      "simple object: set name to some non null value",
			jsonInput: []byte(`{"replicaCount":1}`),
			// since replicaCount field is present in JSON but is NOT null want null to be false
			wantNull: false,
			// since name field is present in JSON want defined to be true
			wantDefined: true,
		},

		{
			name:      "simple object: set name to empty value",
			jsonInput: []byte(`{"replicaCount":0}`),
			// since replicaCount field is present in JSON but is NOT null want null to be false
			wantNull: false,
			// since name field is present in JSON want defined to be true
			wantDefined: true,
		},

		{
			name:      "simple object: set name to null value",
			jsonInput: []byte(`{"replicaCount":null}`),
			// since replicaCount field is present in JSON and is null, want null to be true
			wantNull: true,
			// since name field is present in JSON want defined to be true
			wantDefined: true,
		},

		{
			name:      "simple object: do not provide name in json data",
			jsonInput: []byte(`{}`),
			// since name field is NOT present in JSON, want null to be false
			wantNull: false,
			// since name field is NOT present in JSON want defined to be false
			wantDefined: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			var obj SimpleInt
			err := json.Unmarshal(tt.jsonInput, &obj)
			assert.NoError(t, err)
			assert.Equalf(t, tt.wantNull, obj.ReplicaCount.IsNull(), "IsNull()")
			assert.Equalf(t, tt.wantDefined, obj.ReplicaCount.IsDefined(), "IsDefined()")
		})
	}
}

type SimpleIntPointer struct {
	ReplicaCount Nullable[*int] `json:"replicaCount"`
}

func TestSimpleIntPointer_IsDefined(t *testing.T) {
	type testCase struct {
		name        string
		jsonInput   []byte
		wantNull    bool
		wantDefined bool
	}
	tests := []testCase{
		{
			name:      "simple object: set name to some non null value",
			jsonInput: []byte(`{"replicaCount":1}`),
			// since replicaCount field is present in JSON but is NOT null, want null false
			wantNull: false,
			// since replicaCount field is present in JSON want defined to be true
			wantDefined: true,
		},

		{
			name:      "simple object: set name to empty value",
			jsonInput: []byte(`{"replicaCount":0}`),
			// since replicaCount field is present in JSON but is NOT null, want null false
			wantNull: false,
			// since replicaCount field is present in JSON want defined to be true
			wantDefined: true,
		},

		{
			name:      "simple object: set name to null value",
			jsonInput: []byte(`{"replicaCount":null}`),
			// since replicaCount field is present in JSON and is null, want null true
			wantNull: true,
			// since replicaCount field is present in JSON want defined to be true
			wantDefined: true,
		},

		{
			name:      "simple object: do not provide name in json data",
			jsonInput: []byte(`{}`),
			// since replicaCount field is NOT present in JSON, want null false
			wantNull: false,
			// since replicaCount field is NOT present in JSON want defined to be false
			wantDefined: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			var obj SimpleIntPointer
			err := json.Unmarshal(tt.jsonInput, &obj)
			assert.NoError(t, err)
			assert.Equalf(t, tt.wantNull, obj.ReplicaCount.IsNull(), "IsNull()")
		})
	}
}
