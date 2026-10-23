package runtime

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/oapi-codegen/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oapi-codegen/runtime/types"
)

// decodeGeneric decodes JSON the way kin-openapi's validator does: into an
// `any`, keeping numbers as json.Number.
func decodeGeneric(t *testing.T, doc string) any {
	t.Helper()
	dec := json.NewDecoder(strings.NewReader(doc))
	dec.UseNumber()
	var v any
	require.NoError(t, dec.Decode(&v))
	return v
}

type bindJSONEmbedded struct {
	Promoted string `json:"promoted"`
	Shadowed string `json:"shadowed"`
}

// BindJSONEmbeddedPtr is exported because encoding/json refuses to allocate an
// embedded nil pointer to an unexported struct type.
type BindJSONEmbeddedPtr struct {
	ViaPointer int `json:"viaPointer"`
}

type bindJSONNested struct {
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

type bindJSONEnum string

// bindJSONAll exercises every field shape oapi-codegen generates for a model.
type bindJSONAll struct {
	bindJSONEmbedded
	*BindJSONEmbeddedPtr

	Shadowed string `json:"shadowed"` // wins over the embedded one

	Bool    bool    `json:"bool"`
	Int     int     `json:"int"`
	Int8    int8    `json:"int8"`
	Int64   int64   `json:"int64"`
	Uint32  uint32  `json:"uint32"`
	Float32 float32 `json:"float32"`
	Float64 float64 `json:"float64"`
	String  string  `json:"string"`
	Enum    bindJSONEnum

	PtrString *string         `json:"ptrString,omitempty"`
	PtrInt    *int64          `json:"ptrInt,omitempty"`
	PtrStruct *bindJSONNested `json:"ptrStruct,omitempty"`

	Strings  []string                  `json:"strings"`
	Ints     []int32                   `json:"ints"`
	Structs  []bindJSONNested          `json:"structs"`
	PtrSlice *[]string                 `json:"ptrSlice,omitempty"`
	Array    [3]int                    `json:"array"`
	Bytes    []byte                    `json:"bytes"`
	Map      map[string]int            `json:"map"`
	MapAny   map[string]any            `json:"mapAny"`
	IntKeys  map[int]string            `json:"intKeys"`
	Any      any                       `json:"any"`
	Raw      json.RawMessage           `json:"raw"`
	Nested   bindJSONNested            `json:"nested"`
	Time     time.Time                 `json:"time"`
	Date     types.Date                `json:"date"`
	PtrDate  *types.Date               `json:"ptrDate,omitempty"`
	UUID     types.UUID                `json:"uuid"`
	Email    types.Email               `json:"email"`
	Nullable nullable.Nullable[string] `json:"nullable,omitempty"`
	Quoted   int                       `json:"quoted,string"`

	Ignored    string `json:"-"`
	unexported string
	NoTag      string
	Dash       string `json:"-,"` //nolint:staticcheck // SA5008: "-," deliberately names the member "-"; the binder must match encoding/json here
}

func TestBindJSONValue_MatchesEncodingJSON(t *testing.T) {
	cases := []struct {
		name string
		doc  string
	}{
		{"empty object", `{}`},
		{"scalars", `{"bool":true,"int":-7,"int8":127,"int64":9007199254740993,"uint32":4294967295,"float32":1.5,"float64":2.25e10,"string":"s","Enum":"a"}`},
		{"pointers set", `{"ptrString":"p","ptrInt":42,"ptrStruct":{"name":"n","score":0.5}}`},
		{"pointers null", `{"ptrString":null,"ptrInt":null,"ptrStruct":null,"ptrSlice":null,"ptrDate":null}`},
		{"collections", `{"strings":["a","b"],"ints":[1,2,3],"structs":[{"name":"x"},{"score":1}],"ptrSlice":["q"],"array":[1,2],"bytes":"aGVsbG8=","map":{"a":1,"b":2},"mapAny":{"n":1,"s":"x","l":[1,"y",null],"o":{"k":false}},"intKeys":{"1":"one","-2":"minus two"}}`},
		{"empty collections", `{"strings":[],"ints":[],"structs":[],"map":{},"mapAny":{},"array":[]}`},
		{"null collections", `{"strings":null,"map":null,"mapAny":null,"any":null,"raw":null,"bytes":null}`},
		{"any", `{"any":{"deep":[1,2.5,"x",true,null,{"n":1e3}]}}`},
		{"any number", `{"any":12345678901234567890}`},
		{"raw", `{"raw":{"keep":["me",1]}}`},
		{"nested", `{"nested":{"name":"inner","score":3}}`},
		{"unmarshalers", `{"time":"2024-01-02T03:04:05Z","date":"2024-01-02","ptrDate":"2023-12-31","uuid":"6ba7b810-9dad-11d1-80b4-00c04fd430c8","email":"a@example.com"}`},
		{"nullable set", `{"nullable":"v"}`},
		{"nullable null", `{"nullable":null}`},
		{"quoted", `{"quoted":"17"}`},
		{"embedded", `{"promoted":"yes","shadowed":"outer","viaPointer":5}`},
		{"ignored and untagged", `{"-":"dash","Ignored":"no","unexported":"no","NoTag":"yes"}`},
		{"unknown members", `{"nope":1,"also":{"x":[1]}}`},
		{"unicode escapes", `{"string":"é😀<>&"}`},
		{"big float into float64", `{"float64":1e308}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var want bindJSONAll
			require.NoError(t, json.Unmarshal([]byte(tc.doc), &want))

			var got bindJSONAll
			require.NoError(t, BindJSONValue(decodeGeneric(t, tc.doc), &got))

			// json.RawMessage is re-encoded by the binder, so compare it by
			// JSON value rather than by bytes.
			if len(want.Raw) > 0 || len(got.Raw) > 0 {
				assert.JSONEq(t, string(want.Raw), string(got.Raw))
				want.Raw, got.Raw = nil, nil
			}
			assert.Empty(t, got.unexported, "unexported fields must never be set")
			assert.Equal(t, want, got)
		})
	}
}

func TestBindJSONValue_ExactMemberNames(t *testing.T) {
	type user struct {
		Name string `json:"name"`
		Role string `json:"role"`
		ID   int    `json:"id,omitempty"`
	}

	// encoding/json would bind "ROLE" to Role and "ID" to ID; the binder
	// only honours the exact names the validator saw.
	var got user
	require.NoError(t, BindJSONValue(decodeGeneric(t, `{"name":"bob","role":"user","ROLE":"admin","ID":42,"Name":"BOB"}`), &got))
	assert.Equal(t, user{Name: "bob", Role: "user"}, got)

	var viaJSON user
	require.NoError(t, json.Unmarshal([]byte(`{"name":"bob","role":"user","ROLE":"admin","ID":42}`), &viaJSON))
	assert.Equal(t, "admin", viaJSON.Role, "sanity check: encoding/json is case-insensitive")
	assert.Equal(t, 42, viaJSON.ID)
}

func TestBindJSONValue_Errors(t *testing.T) {
	type target struct {
		Int    int            `json:"int"`
		Uint   uint8          `json:"uint"`
		Float  float32        `json:"float"`
		Bool   bool           `json:"bool"`
		String string         `json:"string"`
		Slice  []int          `json:"slice"`
		Array  [1]int         `json:"array"`
		Map    map[string]int `json:"map"`
		Nested struct{}       `json:"nested"`
		Date   types.Date     `json:"date"`
		Bytes  []byte         `json:"bytes"`
		Quoted int            `json:"quoted,string"`
	}
	cases := []struct {
		name string
		doc  string
	}{
		{"string into int", `{"int":"1"}`},
		{"fraction into int", `{"int":1.5}`},
		{"overflow int8-ish uint8", `{"uint":256}`},
		{"negative into uint", `{"uint":-1}`},
		{"overflow float32", `{"float":1e300}`},
		{"number into bool", `{"bool":1}`},
		{"number into string", `{"string":1}`},
		{"object into slice", `{"slice":{}}`},
		{"string into array", `{"array":"x"}`},
		{"array into map", `{"map":[]}`},
		{"array into struct", `{"nested":[]}`},
		{"bad date", `{"date":"not-a-date"}`},
		{"bad base64", `{"bytes":"%%%"}`},
		{"quoted non-string", `{"quoted":17}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var viaJSON target
			require.Error(t, json.Unmarshal([]byte(tc.doc), &viaJSON), "sanity check: encoding/json rejects it too")

			got := target{Int: 99}
			err := BindJSONValue(decodeGeneric(t, tc.doc), &got)
			require.Error(t, err)
			assert.Equal(t, target{Int: 99}, got, "destination must be untouched on error")
		})
	}

	t.Run("embedded nil pointer to unexported struct", func(t *testing.T) {
		type inner struct {
			V int `json:"v"`
		}
		type outer struct {
			*inner
		}
		doc := `{"v":1}`
		var viaJSON outer
		require.Error(t, json.Unmarshal([]byte(doc), &viaJSON))
		var got outer
		require.Error(t, BindJSONValue(decodeGeneric(t, doc), &got))
	})

	t.Run("top-level mismatch", func(t *testing.T) {
		var s []string
		assert.Error(t, BindJSONValue(decodeGeneric(t, `{"a":1}`), &s))
		var n int
		assert.Error(t, BindJSONValue(decodeGeneric(t, `"x"`), &n))
	})

	t.Run("unsupported destinations", func(t *testing.T) {
		var ch chan int
		assert.ErrorIs(t, BindJSONValue(decodeGeneric(t, `1`), &ch), errUnsupportedJSONDestination)
		var iface io.Reader
		assert.ErrorIs(t, BindJSONValue(decodeGeneric(t, `{}`), &iface), errUnsupportedJSONDestination)
		var m map[bool]int
		assert.ErrorIs(t, BindJSONValue(decodeGeneric(t, `{"true":1}`), &m), errUnsupportedJSONDestination)
	})

	t.Run("bad destination", func(t *testing.T) {
		var s string
		assert.Error(t, BindJSONValue("x", s))
		assert.Error(t, BindJSONValue("x", nil))
		assert.Error(t, BindJSONValue("x", (*string)(nil)))
	})
}

func TestBindJSONValue_TopLevelShapes(t *testing.T) {
	var s []bindJSONNested
	require.NoError(t, BindJSONValue(decodeGeneric(t, `[{"name":"a"},{"score":2}]`), &s))
	assert.Equal(t, []bindJSONNested{{Name: "a"}, {Score: 2}}, s)

	var m map[string]*bindJSONNested
	require.NoError(t, BindJSONValue(decodeGeneric(t, `{"x":{"name":"a"},"y":null}`), &m))
	assert.Equal(t, map[string]*bindJSONNested{"x": {Name: "a"}, "y": nil}, m)

	var str string
	require.NoError(t, BindJSONValue(decodeGeneric(t, `"text"`), &str))
	assert.Equal(t, "text", str)

	var id uuid.UUID
	require.NoError(t, BindJSONValue(decodeGeneric(t, `"6ba7b810-9dad-11d1-80b4-00c04fd430c8"`), &id))
	assert.Equal(t, "6ba7b810-9dad-11d1-80b4-00c04fd430c8", id.String())

	var anyV any
	require.NoError(t, BindJSONValue(decodeGeneric(t, `[1,"a"]`), &anyV))
	assert.Equal(t, []any{float64(1), "a"}, anyV)

	// Go numeric values, not just json.Number, are accepted.
	var n int64
	require.NoError(t, BindJSONValue(float64(42), &n))
	assert.Equal(t, int64(42), n)
	assert.Error(t, BindJSONValue(42.5, &n))

	// A top-level null zeroes the destination.
	ptr := &bindJSONNested{Name: "x"}
	require.NoError(t, BindJSONValue(nil, &ptr))
	assert.Nil(t, ptr)
}

func TestBindJSONValue_EmbeddedConflicts(t *testing.T) {
	type a struct {
		X string
		Y string
	}
	type b struct {
		X string
		Y string `json:"Y"`
	}
	type outer struct {
		a
		b
		Z string
	}
	// "X" is claimed by two untagged fields at the same depth and is
	// dropped; "Y" at the same depth goes to the single tagged field.
	doc := `{"X":"gone","Y":"tagged","Z":"z"}`
	var want outer
	require.NoError(t, json.Unmarshal([]byte(doc), &want))
	var got outer
	require.NoError(t, BindJSONValue(decodeGeneric(t, doc), &got))
	assert.Equal(t, want, got)
	assert.Equal(t, "", got.a.X)
	assert.Equal(t, "", got.b.X)
	assert.Equal(t, "tagged", got.b.Y)
	assert.Equal(t, "z", got.Z)
}

type recursiveNode struct {
	Value    int             `json:"value"`
	Children []recursiveNode `json:"children,omitempty"`
	Next     *recursiveNode  `json:"next,omitempty"`
}

func TestBindJSONValue_Recursive(t *testing.T) {
	doc := `{"value":1,"children":[{"value":2,"next":{"value":3}}],"next":{"value":4,"children":[]}}`
	var want recursiveNode
	require.NoError(t, json.Unmarshal([]byte(doc), &want))
	var got recursiveNode
	require.NoError(t, BindJSONValue(decodeGeneric(t, doc), &got))
	assert.Equal(t, want, got)
}

// validatedReadCloser mimics the body kin-openapi leaves behind with
// Options.KeepValidatedValue set.
type validatedReadCloser struct {
	io.ReadCloser
	value any
}

func (b validatedReadCloser) ValidatedValue() any { return b.value }

func TestBindJSONBody(t *testing.T) {
	type user struct {
		Name string `json:"name"`
		Role string `json:"role"`
	}
	const doc = `{"name":"bob","role":"user","ROLE":"admin"}`

	newRequest := func(body string) *http.Request {
		return httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	}

	t.Run("binds the validated value", func(t *testing.T) {
		r := newRequest(doc)
		r.Body = validatedReadCloser{ReadCloser: r.Body, value: decodeGeneric(t, doc)}
		var got user
		require.NoError(t, BindJSONBody(r, &got))
		assert.Equal(t, user{Name: "bob", Role: "user"}, got)

		// The body was not consumed.
		rest, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Equal(t, doc, string(rest))
	})

	t.Run("decodes the body without a validated value", func(t *testing.T) {
		r := newRequest(doc)
		var got user
		require.NoError(t, BindJSONBody(r, &got))
		assert.Equal(t, user{Name: "bob", Role: "admin"}, got, "plain encoding/json behaviour is kept")
	})

	t.Run("falls back to the body when the value does not bind", func(t *testing.T) {
		r := newRequest(doc)
		r.Body = validatedReadCloser{ReadCloser: r.Body, value: decodeGeneric(t, `{"name":1}`)}
		var got user
		require.NoError(t, BindJSONBody(r, &got))
		assert.Equal(t, user{Name: "bob", Role: "admin"}, got)
	})

	t.Run("fallback errors are encoding/json's", func(t *testing.T) {
		r := newRequest(`{"name":1}`)
		var got user
		err := BindJSONBody(r, &got)
		var typeErr *json.UnmarshalTypeError
		require.ErrorAs(t, err, &typeErr)

		r = newRequest("")
		assert.ErrorIs(t, BindJSONBody(r, &got), io.EOF)
	})

	t.Run("nil request decodes nothing", func(t *testing.T) {
		var got user
		r := &http.Request{Body: io.NopCloser(bytes.NewReader(nil))}
		assert.ErrorIs(t, BindJSONBody(r, &got), io.EOF)
	})
}

func TestBindJSONValue_FieldCacheIsStable(t *testing.T) {
	f1 := cachedJSONFields(typeOf[bindJSONAll]())
	f2 := cachedJSONFields(typeOf[bindJSONAll]())
	assert.Equal(t, f1, f2)
	_, hasIgnored := f1["Ignored"]
	assert.False(t, hasIgnored)
	_, hasUnexported := f1["unexported"]
	assert.False(t, hasUnexported)
	assert.Equal(t, []int{1, 0}, f1["viaPointer"].index)
	assert.Equal(t, "-", f1["-"].name)
	assert.True(t, f1["quoted"].quoted)
}

func typeOf[T any]() reflect.Type { return reflect.TypeFor[T]() }
