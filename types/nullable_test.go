package types

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/oapi-codegen/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type SimpleStringNullableRequired struct {
	// A required field must be `nullable.Nullable`.
	Name nullable.Nullable[string] `json:"name"`
}

func TestSimpleStringNullableRequired(t *testing.T) {
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
			// For Nullable and required field there will not be a case when it is not
			// provided in json payload but this test case exists just to explain
			// the behaviour
			name:          "simple object: do not provide name in json data",
			jsonInput:     []byte(`{}`),
			wantNull:      false,
			wantSpecified: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj SimpleStringNullableRequired
			err := json.Unmarshal(tt.jsonInput, &obj)
			require.NoError(t, err)
			assert.Equalf(t, tt.wantNull, obj.Name.IsNull(), "IsNull()")
			assert.Equalf(t, tt.wantSpecified, obj.Name.IsSpecified(), "IsSpecified()")
		})
	}
}

type SimpleStringNullableOptional struct {
	// An optional field must be `nullable.Nullable` and have the JSON tag `omitempty`.
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
			assert.Equalf(t, tt.wantSpecified, obj.Name.IsSpecified(), "IsSpecified()")
		})
	}
}

type SimpleIntNullableRequired struct {
	// A required field must be `nullable.Nullable`.
	ReplicaCount nullable.Nullable[int] `json:"replica_count"`
}

func TestSimpleIntNullableRequired(t *testing.T) {
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
			// For Nullable and required field there will not be a case when it is not
			// provided in json payload but this test case exists just to explain
			// the behaviour
			name:          "simple object: do not provide name in json data",
			jsonInput:     []byte(`{}`),
			wantNull:      false,
			wantSpecified: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj SimpleIntNullableRequired
			err := json.Unmarshal(tt.jsonInput, &obj)
			require.NoError(t, err)
			assert.Equalf(t, tt.wantNull, obj.ReplicaCount.IsNull(), "IsNull()")
			assert.Equalf(t, tt.wantSpecified, obj.ReplicaCount.IsSpecified(), "IsSpecified()")
		})
	}
}

type SimpleIntNullableOptional struct {
	// An optional field must be `nullable.Nullable` and have the JSON tag `omitempty`
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
			assert.Equalf(t, tt.wantSpecified, obj.ReplicaCount.IsSpecified(), "IsSpecified()")
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
			// When object is not set ( not provided in json )
			// it marshals to zero value.
			name: "when obj is explicitly not set",
			want: []byte(`{}`),
		},
		{
			name: "when obj is explicitly set to a specific value",
			obj: SimpleIntNullableOptional{
				ReplicaCount: nullable.Nullable[int]{
					true: 5,
				},
			},
			want: []byte(`{"replica_count":5}`),
		},
		{
			name: "when obj is explicitly set to zero value",
			obj: SimpleIntNullableOptional{
				ReplicaCount: nullable.Nullable[int]{
					true: 0,
				},
			},
			want: []byte(`{"replica_count":0}`),
		},
		{
			name: "when obj is explicitly set to null value",
			obj: SimpleIntNullableOptional{
				ReplicaCount: nullable.Nullable[int]{
					false: 0,
				},
			},
			want: []byte(`{"replica_count":null}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.obj)
			fmt.Println(string(got))
			require.NoError(t, err)
			assert.Equalf(t, string(tt.want), string(got), "MarshalJSON()")
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

				assert.Falsef(t, obj.ReplicaCount.IsSpecified(), "replica count should not be set")
				assert.Falsef(t, obj.ReplicaCount.IsNull(), "replica count should not be null")
			},
		},

		{
			name: "when explicitly set to zero value",
			json: []byte(`{"replica_count":0}`),
			assert: func(t *testing.T, obj SimpleIntNullableOptional) {
				t.Helper()

				assert.Truef(t, obj.ReplicaCount.IsSpecified(), "replica count should be set")
				assert.Falsef(t, obj.ReplicaCount.IsNull(), "replica count should not be null")

				val, err := obj.ReplicaCount.Get()
				require.NoError(t, err)
				assert.Equalf(t, 0, val, "replica count value should be 0")
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

				assert.Truef(t, obj.ReplicaCount.IsSpecified(), "replica count should be set")
				assert.Falsef(t, obj.ReplicaCount.IsNull(), "replica count should not be null")

				val, err := obj.ReplicaCount.Get()
				require.NoError(t, err)
				assert.Equalf(t, 5, val, "replica count value should be 5")
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
				ReplicaCount: nullable.Nullable[int]{
					true: 5,
				},
			},
			want: []byte(`{"replica_count":5}`),
		},
		{
			name: "when obj is explicitly set to zero value",
			obj: SimpleIntNullableRequired{
				ReplicaCount: nullable.Nullable[int]{
					true: 0,
				},
			},
			want: []byte(`{"replica_count":0}`),
		},
		{
			name: "when obj is explicitly set to null value",
			obj: SimpleIntNullableRequired{
				ReplicaCount: nullable.Nullable[int]{
					false: 0,
				},
			},
			want: []byte(`{"replica_count":null}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.obj)
			fmt.Println(string(got))
			require.NoError(t, err)
			assert.Equalf(t, string(tt.want), string(got), "MarshalJSON()")
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
				assert.Falsef(t, obj.ReplicaCount.IsSpecified(), "replica count should not be set")
			},
		},

		{
			name: "when explicitly set to zero value",
			json: []byte(`{"replica_count":0}`),
			assert: func(t *testing.T, obj SimpleIntNullableRequired) {
				t.Helper()

				assert.Falsef(t, obj.ReplicaCount.IsNull(), "replica count should not be null")
				assert.Truef(t, obj.ReplicaCount.IsSpecified(), "replica count should be set")

				val, err := obj.ReplicaCount.Get()
				require.NoError(t, err)
				assert.Equalf(t, 0, val, "replica count value should be 0")
			},
		},

		{
			name: "when explicitly set to null value",
			json: []byte(`{"replica_count":null}`),
			assert: func(t *testing.T, obj SimpleIntNullableRequired) {
				t.Helper()

				assert.Truef(t, obj.ReplicaCount.IsSpecified(), "replica count should be set")
				assert.Truef(t, obj.ReplicaCount.IsNull(), "replica count should be null")
			},
		},

		{
			name: "when explicitly set to a specific value",
			json: []byte(`{"replica_count":5}`),
			assert: func(t *testing.T, obj SimpleIntNullableRequired) {
				t.Helper()

				assert.Falsef(t, obj.ReplicaCount.IsNull(), "replica count should not be null")
				assert.Truef(t, obj.ReplicaCount.IsSpecified(), "replica count should be set")

				val, err := obj.ReplicaCount.Get()
				require.NoError(t, err)
				assert.Equalf(t, 5, val, "replica count value should be 5")
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

type ComplexNullable struct {
	Config    nullable.Nullable[Config] `json:"config,omitempty"`
	Location  nullable.Nullable[string] `json:"location"`
	NodeCount nullable.Nullable[int]    `json:"node_count,omitempty"`
}

type Config struct {
	CPU nullable.Nullable[string] `json:"cpu,omitempty"`
	RAM nullable.Nullable[string] `json:"ram,omitempty"`
}

func TestComplexNullable(t *testing.T) {
	type testCase struct {
		name      string
		jsonInput []byte
		assert    func(t *testing.T, obj ComplexNullable)
	}
	tests := []testCase{
		{
			name:      "complex object: empty value",
			jsonInput: []byte(`{}`),
			assert: func(t *testing.T, obj ComplexNullable) {
				t.Helper()

				assert.Falsef(t, obj.Config.IsSpecified(), "config should not be set")
				assert.Falsef(t, obj.Config.IsNull(), "config should not be null")

				assert.Falsef(t, obj.NodeCount.IsSpecified(), "node count should not be set")
				assert.Falsef(t, obj.NodeCount.IsNull(), "node count should not be null")

				assert.Falsef(t, obj.Location.IsSpecified(), "location should not be set")
				assert.Falsef(t, obj.Location.IsNull(), "location should not be null")
			},
		},
		{
			name:      "complex object: empty config value",
			jsonInput: []byte(`{"config":{}}`),
			assert: func(t *testing.T, obj ComplexNullable) {
				t.Helper()

				assert.Truef(t, obj.Config.IsSpecified(), "config should be set")
				assert.Falsef(t, obj.Config.IsNull(), "config should not be null")

				gotConfig, err := obj.Config.Get()
				require.NoError(t, err)

				assert.Falsef(t, gotConfig.CPU.IsSpecified(), "cpu should not be set")
				assert.Falsef(t, gotConfig.CPU.IsNull(), "cpu should not be null")

				assert.Falsef(t, gotConfig.RAM.IsSpecified(), "ram should not be set")
				assert.Falsef(t, gotConfig.RAM.IsNull(), "ram should not be null")
			},
		},

		{
			name:      "complex object: setting only cpu config value",
			jsonInput: []byte(`{"config":{"cpu":"500"}}`),
			assert: func(t *testing.T, obj ComplexNullable) {
				t.Helper()

				assert.Truef(t, obj.Config.IsSpecified(), "config should be set")
				assert.Falsef(t, obj.Config.IsNull(), "config should not be null")

				gotConfig, err := obj.Config.Get()
				require.NoError(t, err)

				assert.Truef(t, gotConfig.CPU.IsSpecified(), "cpu should be set")
				assert.Falsef(t, gotConfig.CPU.IsNull(), "cpu should not be null")

				assert.Falsef(t, gotConfig.RAM.IsSpecified(), "ram should not be set")
				assert.Falsef(t, gotConfig.RAM.IsNull(), "ram should not be null")
			},
		},

		{
			name:      "complex object: setting only ram config value",
			jsonInput: []byte(`{"config":{"ram":"1024"}}`),
			assert: func(t *testing.T, obj ComplexNullable) {
				t.Helper()

				assert.Truef(t, obj.Config.IsSpecified(), "config should be set")
				assert.Falsef(t, obj.Config.IsNull(), "config should not be null")

				gotConfig, err := obj.Config.Get()
				require.NoError(t, err)

				assert.Falsef(t, gotConfig.CPU.IsSpecified(), "cpu should not be set")
				assert.Falsef(t, gotConfig.CPU.IsNull(), "cpu should not be null")

				assert.Truef(t, gotConfig.RAM.IsSpecified(), "ram should be set")
				assert.Falsef(t, gotConfig.RAM.IsNull(), "ram should not be null")
			},
		},

		{
			name:      "complex object: setting config to null",
			jsonInput: []byte(`{"config":null}`),
			assert: func(t *testing.T, obj ComplexNullable) {
				t.Helper()

				assert.Truef(t, obj.Config.IsSpecified(), "config should be set")
				assert.Truef(t, obj.Config.IsNull(), "config should not be null")
			},
		},

		{
			name:      "complex object: setting only cpu config to null",
			jsonInput: []byte(`{"config":{"cpu":null}}`),
			assert: func(t *testing.T, obj ComplexNullable) {
				t.Helper()

				assert.Truef(t, obj.Config.IsSpecified(), "config should be set")
				assert.Falsef(t, obj.Config.IsNull(), "config should not be null")

				gotConfig, err := obj.Config.Get()
				require.NoError(t, err)

				assert.Truef(t, gotConfig.CPU.IsSpecified(), "cpu should be set")
				assert.Truef(t, gotConfig.CPU.IsNull(), "cpu should be null")

				assert.Falsef(t, gotConfig.RAM.IsSpecified(), "ram should not be set")
				assert.Falsef(t, gotConfig.RAM.IsNull(), "ram should not be null")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj ComplexNullable
			err := json.Unmarshal(tt.jsonInput, &obj)
			require.NoError(t, err)

			tt.assert(t, obj)
		})
	}
}

// Tests to validate that repeated unmarshal and marshal should not change
// the JSON for optional and nullable.
type StringNullableOptional struct {
	ID      nullable.Nullable[string] `json:"id,omitempty"`
	Name    nullable.Nullable[string] `json:"name,omitempty"`
	Address nullable.Nullable[string] `json:"address,omitempty"`
}

func TestNullableOptionalUnmarshal(t *testing.T) {
	var obj1 StringNullableOptional
	originalJSON1 := []byte(`{}`)

	err := json.Unmarshal(originalJSON1, &obj1)
	require.NoError(t, err)
	newJSON1, err := json.Marshal(obj1)
	require.NoError(t, err)

	require.Equal(t, originalJSON1, newJSON1)

	var obj2 StringNullableOptional
	originalJSON2 := []byte(`{"id":"12esd412"}`)

	err = json.Unmarshal(originalJSON2, &obj2)
	require.NoError(t, err)
	newJSON2, err := json.Marshal(obj2)
	require.NoError(t, err)

	require.Equal(t, originalJSON2, newJSON2)
}

func TestNullableOptionalMarshal(t *testing.T) {
	obj1 := StringNullableOptional{
		ID: nullable.Nullable[string]{
			true: "id-1",
		},
		Name: nullable.Nullable[string]{
			true: "",
		},
		Address: nullable.Nullable[string]{
			false: "",
		},
	}
	expectedJSON1 := []byte(`{"id":"id-1","name":"","address":null}`)
	gotJSON1, err := json.Marshal(obj1)
	require.NoError(t, err)
	require.Equal(t, expectedJSON1, gotJSON1)

	obj2 := StringNullableOptional{
		ID: nullable.Nullable[string]{
			true: "id-1",
		},
		Address: nullable.Nullable[string]{
			false: "",
		},
	}
	expectedJSON2 := []byte(`{"id":"id-1","address":null}`)
	gotJSON2, err := json.Marshal(obj2)
	require.NoError(t, err)
	require.Equal(t, expectedJSON2, gotJSON2)
}
