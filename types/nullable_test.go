package types

import (
	"encoding/json"
	"fmt"
	"github.com/oapi-codegen/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
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
		assert func(obj SimpleIntNullableOptional, t *testing.T)
	}
	tests := []testCase{
		{
			name: "when not provided",
			json: []byte(`{}`),
			assert: func(obj SimpleIntNullableOptional, t *testing.T) {
				t.Helper()

				assert.Equalf(t, false, obj.ReplicaCount.IsSpecified(), "replica count should not be set")
				assert.Equalf(t, false, obj.ReplicaCount.IsNull(), "replica count should not be null")
			},
		},

		{
			name: "when explicitly set to zero value",
			json: []byte(`{"replica_count":0}`),
			assert: func(obj SimpleIntNullableOptional, t *testing.T) {
				t.Helper()

				assert.Equalf(t, true, obj.ReplicaCount.IsSpecified(), "replica count should be set")
				assert.Equalf(t, false, obj.ReplicaCount.IsNull(), "replica count should not be null")

				val, err := obj.ReplicaCount.Get()
				require.NoError(t, err)
				assert.Equalf(t, 0, val, "replica count value should be 0")

			},
		},

		{
			name: "when explicitly set to null value",
			json: []byte(`{"replica_count":null}`),
			assert: func(obj SimpleIntNullableOptional, t *testing.T) {
				t.Helper()

				assert.Equalf(t, true, obj.ReplicaCount.IsSpecified(), "replica count should be set")
				assert.Equalf(t, true, obj.ReplicaCount.IsNull(), "replica count should be null")
			},
		},

		{
			name: "when explicitly set to a specific value",
			json: []byte(`{"replica_count":5}`),
			assert: func(obj SimpleIntNullableOptional, t *testing.T) {
				t.Helper()

				assert.Equalf(t, true, obj.ReplicaCount.IsSpecified(), "replica count should not be null")
				assert.Equalf(t, false, obj.ReplicaCount.IsNull(), "replica count should not be null")

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
			tt.assert(obj, t)
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
		assert func(obj SimpleIntNullableRequired, t *testing.T)
	}
	tests := []testCase{
		{
			// For Nullable and required field there will not be a case when it is not
			// provided in json payload but this test case exists just to explain
			// the behaviour
			name: "when not provided",
			json: []byte(`{}`),
			assert: func(obj SimpleIntNullableRequired, t *testing.T) {
				t.Helper()

				assert.Equalf(t, false, obj.ReplicaCount.IsNull(), "replica count should not be null")
				assert.Equalf(t, false, obj.ReplicaCount.IsSpecified(), "replica count should not be set")
			},
		},

		{
			name: "when explicitly set to zero value",
			json: []byte(`{"replica_count":0}`),
			assert: func(obj SimpleIntNullableRequired, t *testing.T) {
				t.Helper()

				assert.Equalf(t, false, obj.ReplicaCount.IsNull(), "replica count should not be null")
				assert.Equalf(t, true, obj.ReplicaCount.IsSpecified(), "replica count should be set")

				val, err := obj.ReplicaCount.Get()
				require.NoError(t, err)
				assert.Equalf(t, 0, val, "replica count value should be 0")
			},
		},

		{
			name: "when explicitly set to null value",
			json: []byte(`{"replica_count":null}`),
			assert: func(obj SimpleIntNullableRequired, t *testing.T) {
				t.Helper()

				assert.Equalf(t, true, obj.ReplicaCount.IsSpecified(), "replica count should be set")
				assert.Equalf(t, true, obj.ReplicaCount.IsNull(), "replica count should be null")
			},
		},

		{
			name: "when explicitly set to a specific value",
			json: []byte(`{"replica_count":5}`),
			assert: func(obj SimpleIntNullableRequired, t *testing.T) {
				t.Helper()

				assert.Equalf(t, false, obj.ReplicaCount.IsNull(), "replica count should not be null")
				assert.Equalf(t, true, obj.ReplicaCount.IsSpecified(), "replica count should be set")

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
			tt.assert(obj, t)
		})
	}
}

// Idempotency tests for nullable and optional
type StringNullableOptional struct {
	ID      nullable.Nullable[string] `json:"id,omitempty"`
	Name    nullable.Nullable[string] `json:"name,omitempty"`
	Address nullable.Nullable[string] `json:"address,omitempty"`
}

func TestNullableOptionalUnmarshalIdempotency(t *testing.T) {
	var obj1 StringNullableOptional
	originalJSON1 := []byte(`{}`)
	err := json.Unmarshal(originalJSON1, &obj1)
	require.NoError(t, err)
	newJSON1, err := json.Marshal(obj1)
	require.Equal(t, originalJSON1, newJSON1)

	var obj2 StringNullableOptional
	originalJSON2 := []byte(`{"id":"12esd412"}`)
	err2 := json.Unmarshal(originalJSON2, &obj2)
	require.NoError(t, err2)
	newJSON2, err2 := json.Marshal(obj2)
	require.Equal(t, originalJSON2, newJSON2)
}

func TestNullableOptionalMarshalIdempotency(t *testing.T) {
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
	gotJSON1, err1 := json.Marshal(obj1)
	require.NoError(t, err1)
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
	gotJSON2, err2 := json.Marshal(obj2)
	require.NoError(t, err2)
	require.Equal(t, expectedJSON2, gotJSON2)
}
