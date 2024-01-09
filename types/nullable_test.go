package types

import (
	"encoding/json"
	"testing"

	"github.com/oapi-codegen/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type SimpleStringNullableRequired struct {
	// A required field must be `nullable.Nullable`
	Name nullable.Nullable[string] `json:"name"`
}

func TestSimpleStringNullableRequired(t *testing.T) {
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
			wantNull:  false,
		},

		{
			name:      "simple object: set name to null value",
			jsonInput: []byte(`{"name":null}`),
			wantNull:  true,
		},

		{
			// For Nullable and required field there will not be a case when it is not
			// provided in json payload but this test case exists just to explain
			// the behaviour
			name:      "simple object: do not provide name in json data",
			jsonInput: []byte(`{}`),
			wantNull:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj SimpleStringNullableRequired
			err := json.Unmarshal(tt.jsonInput, &obj)
			require.NoError(t, err)
			assert.Equalf(t, tt.wantNull, obj.Name.IsNull(), "IsNull()")
		})
	}
}

type SimpleStringNullableOptional struct {
	// An optional field must be `nullable.Nullable` and have the JSON tag `omitempty`
	Name nullable.Nullable[string] `json:"name,omitempty"`
}

func TestSimpleStringNullableOptional(t *testing.T) {
	type testCase struct {
		name          string
		jsonInput     []byte
		wantNull      bool
		wantSpecified bool
	}
	tests := []testCase{
		{
			name:          "simple object: set name to some non null value",
			jsonInput:     []byte(`{"name":"yolo"}`),
			wantNull:      false,
			wantSpecified: true,
		},

		{
			name:          "simple object: set name to empty string value",
			jsonInput:     []byte(`{"name":""}`),
			wantNull:      false,
			wantSpecified: true,
		},

		{
			name:          "simple object: set name to null value",
			jsonInput:     []byte(`{"name":null}`),
			wantNull:      true,
			wantSpecified: true,
		},

		{
			name:          "simple object: do not provide name in json data",
			jsonInput:     []byte(`{}`),
			wantNull:      false,
			wantSpecified: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj SimpleStringNullableOptional
			err := json.Unmarshal(tt.jsonInput, &obj)
			require.NoError(t, err)
			assert.Equalf(t, tt.wantNull, obj.Name.IsNull(), "IsNull()")
			assert.Equalf(t, tt.wantSpecified, obj.Name.IsSpecified(), "IsSpecified")
		})
	}
}

type SimpleIntNullableRequired struct {
	// Nullable type should be used for  `nullable and required` fields.
	ReplicaCount nullable.Nullable[int] `json:"replica_count"`
}

func TestSimpleIntNullableRequired(t *testing.T) {
	type testCase struct {
		name      string
		jsonInput []byte
		wantNull  bool
	}
	tests := []testCase{
		{
			name:      "simple object: set name to some non null value",
			jsonInput: []byte(`{"replica_count":1}`),
			wantNull:  false,
		},

		{
			name:      "simple object: set name to empty value",
			jsonInput: []byte(`{"replica_count":0}`),
			wantNull:  false,
		},

		{
			name:      "simple object: set name to null value",
			jsonInput: []byte(`{"replica_count":null}`),
			wantNull:  true,
		},

		{
			// For Nullable and required field there will not be a case when it is not
			// provided in json payload but this test case exists just to explain
			// the behaviour
			name:      "simple object: do not provide name in json data",
			jsonInput: []byte(`{}`),
			wantNull:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj SimpleIntNullableRequired
			err := json.Unmarshal(tt.jsonInput, &obj)
			require.NoError(t, err)
			assert.Equalf(t, tt.wantNull, obj.ReplicaCount.IsNull(), "IsNull()")
		})
	}
}

type SimpleIntNullableOptional struct {
	ReplicaCount nullable.Nullable[int] `json:"replica_count,omitempty"`
}

func TestSimpleIntNullableOptional(t *testing.T) {
	type testCase struct {
		name          string
		jsonInput     []byte
		wantNull      bool
		wantSpecified bool
	}
	tests := []testCase{
		{
			name:          "simple object: set name to some non null value",
			jsonInput:     []byte(`{"replica_count":1}`),
			wantNull:      false,
			wantSpecified: true,
		},

		{
			name:          "simple object: set name to empty value",
			jsonInput:     []byte(`{"replica_count":0}`),
			wantNull:      false,
			wantSpecified: true,
		},

		{
			name:          "simple object: set name to null value",
			jsonInput:     []byte(`{"replica_count":null}`),
			wantNull:      true,
			wantSpecified: true,
		},

		{
			name:          "simple object: do not provide name in json data",
			jsonInput:     []byte(`{}`),
			wantNull:      false,
			wantSpecified: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj SimpleIntNullableOptional
			err := json.Unmarshal(tt.jsonInput, &obj)
			require.NoError(t, err)

			assert.Equalf(t, tt.wantNull, obj.ReplicaCount.IsNull(), "IsNull()")
			assert.Equalf(t, tt.wantSpecified, obj.ReplicaCount.IsSpecified(), "IsSpecified")
		})
	}
}

func TestNullableOptional_MarshalJSON(t *testing.T) {
	type testCase struct {
		name string
		obj  SimpleIntNullableOptional
		want []byte
	}
	tests := []testCase{
		{
			// When object is not set, and it is an optional type, it does not marshal
			name: "when obj is explicitly not set",
			want: []byte(`{}`),
		},
		{
			name: "when obj is explicitly set to a specific value",
			obj: SimpleIntNullableOptional{
				ReplicaCount: nullable.NewNullableWithValue(5),
			},
			want: []byte(`{"replica_count":5}`),
		},
		{
			name: "when obj is explicitly set to zero value",
			obj: SimpleIntNullableOptional{
				ReplicaCount: nullable.NewNullableWithValue(0),
			},
			want: []byte(`{"replica_count":0}`),
		},
		{
			name: "when obj is explicitly set to null value",
			obj: SimpleIntNullableOptional{
				ReplicaCount: nullable.NewNullNullable[int](),
			},
			want: []byte(`{"replica_count":null}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.obj)
			require.NoError(t, err)
			assert.Equalf(t, tt.want, got, "MarshalJSON()")
		})
	}
}

func TestNullableOptional_UnmarshalJSON(t *testing.T) {
	type testCase struct {
		name   string
		json   []byte
		assert func(t *testing.T, obj SimpleIntNullableOptional)
	}
	tests := []testCase{
		{
			name: "when not provided",
			json: []byte(`{}`),
			assert: func(t *testing.T, obj SimpleIntNullableOptional) {
				t.Helper()

				assert.Falsef(t, obj.ReplicaCount.IsSpecified(), "replica count should not be specified")
				assert.Falsef(t, obj.ReplicaCount.IsNull(), "replica count should not be null")
			},
		},

		{
			name: "when explicitly set to zero value",
			json: []byte(`{"replica_count":0}`),
			assert: func(t *testing.T, obj SimpleIntNullableOptional) {
				t.Helper()

				assert.Truef(t, obj.ReplicaCount.IsSpecified(), "replica count should be specified")
				assert.Falsef(t, obj.ReplicaCount.IsNull(), "replica count should not be null")

				actual, err := obj.ReplicaCount.Get()
				require.NoError(t, err)

				assert.Equalf(t, 0, actual, "replica count value should be 0")
			},
		},

		{
			name: "when explicitly set to null value",
			json: []byte(`{"replica_count":null}`),
			assert: func(t *testing.T, obj SimpleIntNullableOptional) {
				t.Helper()

				assert.Truef(t, obj.ReplicaCount.IsSpecified(), "replica count should be set")
				assert.Truef(t, obj.ReplicaCount.IsNull(), "replica count should be null")
			},
		},

		{
			name: "when explicitly set to a specific value",
			json: []byte(`{"replica_count":5}`),
			assert: func(t *testing.T, obj SimpleIntNullableOptional) {
				t.Helper()

				assert.Truef(t, obj.ReplicaCount.IsSpecified(), "replica count should not be null")
				assert.Falsef(t, obj.ReplicaCount.IsNull(), "replica count should not be null")

				actual, err := obj.ReplicaCount.Get()
				require.NoError(t, err)

				assert.Equalf(t, 5, actual, "replica count value should be 5")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj SimpleIntNullableOptional
			err := json.Unmarshal(tt.json, &obj)
			require.NoError(t, err)
			tt.assert(t, obj)
		})
	}
}

func TestNullableRequired_MarshalJSON(t *testing.T) {
	type testCase struct {
		name string
		obj  SimpleIntNullableRequired
		want []byte
	}
	tests := []testCase{
		{
			// When object is not set ( not provided in json )
			// it marshals to zero value.
			// For Nullable and required field there will not be a case when it is not
			// provided in json payload but this test case exists just to explain
			// the behaviour
			name: "when obj is explicitly not set",
			want: []byte(`{"replica_count":0}`),
		},
		{
			name: "when obj is explicitly set to a specific value",
			obj: SimpleIntNullableRequired{
				ReplicaCount: nullable.NewNullableWithValue(5),
			},
			want: []byte(`{"replica_count":5}`),
		},
		{
			name: "when obj is explicitly set to zero value",
			obj: SimpleIntNullableRequired{
				ReplicaCount: nullable.NewNullableWithValue(0),
			},
			want: []byte(`{"replica_count":0}`),
		},
		{
			name: "when obj is explicitly set to null value",
			obj: SimpleIntNullableRequired{
				ReplicaCount: nullable.NewNullNullable[int](),
			},
			want: []byte(`{"replica_count":null}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.obj)
			require.NoError(t, err)
			assert.Equalf(t, tt.want, got, "MarshalJSON()")
		})
	}
}

func TestNullableRequired_UnmarshalJSON(t *testing.T) {
	type testCase struct {
		name   string
		json   []byte
		assert func(t *testing.T, obj SimpleIntNullableRequired)
	}
	tests := []testCase{
		{
			// For Nullable and required field there will not be a case when it is not
			// provided in json payload but this test case exists just to explain
			// the behaviour
			name: "when not provided",
			json: []byte(`{}`),
			assert: func(t *testing.T, obj SimpleIntNullableRequired) {
				t.Helper()

				assert.Falsef(t, obj.ReplicaCount.IsNull(), "replica count should not be null")
			},
		},

		{
			name: "when explicitly set to zero value",
			json: []byte(`{"replica_count":0}`),
			assert: func(t *testing.T, obj SimpleIntNullableRequired) {
				t.Helper()

				assert.Falsef(t, obj.ReplicaCount.IsNull(), "replica count should not be null")

				actual, err := obj.ReplicaCount.Get()
				require.NoError(t, err)

				assert.Equalf(t, 0, actual, "replica count value should be 0")
			},
		},

		{
			name: "when explicitly set to null value",
			json: []byte(`{"replica_count":null}`),
			assert: func(t *testing.T, obj SimpleIntNullableRequired) {
				t.Helper()

				assert.Truef(t, obj.ReplicaCount.IsNull(), "replica count should be null")
			},
		},

		{
			name: "when explicitly set to a specific value",
			json: []byte(`{"replica_count":5}`),
			assert: func(t *testing.T, obj SimpleIntNullableRequired) {
				t.Helper()

				assert.Falsef(t, obj.ReplicaCount.IsNull(), "replica count should not be null")

				actual, err := obj.ReplicaCount.Get()
				require.NoError(t, err)

				assert.Equalf(t, 5, actual, "replica count value should be 5")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj SimpleIntNullableRequired
			err := json.Unmarshal(tt.json, &obj)
			require.NoError(t, err)

			tt.assert(t, obj)
		})
	}
}
