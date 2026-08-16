package runtime

import (
	"net/url"
	"testing"

	"github.com/oapi-codegen/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Binding a value into an `any` destination with a declared multi-type union
// picks the first member that parses, in specificity order (boolean,
// integer, number, string) — not schema declaration order, where the
// always-succeeding string member would shadow the rest.
func TestBindStringToObject_UnionMemberSelection(t *testing.T) {
	testCases := []struct {
		name    string
		src     string
		opts    BindStringToObjectOptions
		want    any
		wantErr string
	}{
		{
			name: "integer member wins over string for integer token",
			src:  "123",
			opts: BindStringToObjectOptions{Types: []string{"string", "integer"}},
			want: int64(123),
		},
		{
			name: "number member takes fractional token",
			src:  "1.5",
			opts: BindStringToObjectOptions{Types: []string{"string", "number"}},
			want: float64(1.5),
		},
		{
			name: "integer wins over number for integer token",
			src:  "123",
			opts: BindStringToObjectOptions{Types: []string{"number", "integer", "string"}},
			want: int64(123),
		},
		{
			name: "number takes exponent token when integer present",
			src:  "1e2",
			opts: BindStringToObjectOptions{Types: []string{"integer", "number", "string"}},
			want: float64(100),
		},
		{
			name: "boolean member wins for JSON boolean literal",
			src:  "true",
			opts: BindStringToObjectOptions{Types: []string{"string", "boolean"}},
			want: true,
		},
		{
			name: "non-token falls to string member",
			src:  "abc",
			opts: BindStringToObjectOptions{Types: []string{"integer", "string"}},
			want: "abc",
		},
		{
			name: "empty string binds to string member",
			src:  "",
			opts: BindStringToObjectOptions{Types: []string{"integer", "string"}},
			want: "",
		},
		{
			name: "negative integer token",
			src:  "-42",
			opts: BindStringToObjectOptions{Types: []string{"integer", "string"}},
			want: int64(-42),
		},
		// JSON grammar, not strconv leniency: these are not JSON numbers,
		// so they must not be silently reinterpreted.
		{
			name: "leading zeros stay a string",
			src:  "007",
			opts: BindStringToObjectOptions{Types: []string{"integer", "string"}},
			want: "007",
		},
		{
			name: "plus sign stays a string",
			src:  "+1",
			opts: BindStringToObjectOptions{Types: []string{"integer", "number", "string"}},
			want: "+1",
		},
		{
			name: "surrounding space stays a string",
			src:  " 1",
			opts: BindStringToObjectOptions{Types: []string{"integer", "string"}},
			want: " 1",
		},
		{
			name: "uppercase TRUE is not a JSON boolean",
			src:  "TRUE",
			opts: BindStringToObjectOptions{Types: []string{"boolean", "string"}},
			want: "TRUE",
		},
		{
			name: "int64 overflow falls through to string",
			src:  "99999999999999999999",
			opts: BindStringToObjectOptions{Types: []string{"integer", "string"}},
			want: "99999999999999999999",
		},
		{
			name: "int64 overflow falls through to number",
			src:  "99999999999999999999",
			opts: BindStringToObjectOptions{Types: []string{"integer", "number"}},
			want: float64(99999999999999999999),
		},
		// Width formats are annotation-only: the dynamic type surface stays
		// bool/int64/float64/string/[]byte no matter what format says, so a
		// spec edit to `format` can never break a running type switch.
		{
			name: "format int32 does not narrow the integer member",
			src:  "123",
			opts: BindStringToObjectOptions{Types: []string{"integer", "string"}, Format: "int32"},
			want: int64(123),
		},
		{
			name: "format int32 does not reject values beyond int32 range",
			src:  "3000000000",
			opts: BindStringToObjectOptions{Types: []string{"integer", "string"}, Format: "int32"},
			want: int64(3000000000),
		},
		{
			name: "format float does not narrow the number member",
			src:  "1.5",
			opts: BindStringToObjectOptions{Types: []string{"number", "string"}, Format: "float"},
			want: float64(1.5),
		},
		{
			name: "format byte base64-decodes the string member",
			src:  "MTIz",
			opts: BindStringToObjectOptions{Types: []string{"integer", "string"}, Format: "byte"},
			want: []byte("123"),
		},
		{
			name: "format byte does not touch the integer member",
			src:  "123",
			opts: BindStringToObjectOptions{Types: []string{"integer", "string"}, Format: "byte"},
			want: int64(123),
		},
		{
			name:    "format byte with invalid base64 is an error",
			src:     "not!!base64",
			opts:    BindStringToObjectOptions{Types: []string{"string"}, Format: "byte"},
			wantErr: "failed to base64-decode",
		},
		{
			name: "annotation-only format is ignored",
			src:  "2024-01-01T00:00:00Z",
			opts: BindStringToObjectOptions{Types: []string{"string", "number"}, Format: "date-time"},
			want: "2024-01-01T00:00:00Z",
		},
		// Defensive handling of member names.
		{
			name: "null marker is ignored even if not stripped",
			src:  "123",
			opts: BindStringToObjectOptions{Types: []string{"null", "string", "integer"}},
			want: int64(123),
		},
		{
			name: "unstripped nullable single type binds as that type",
			src:  "abc",
			opts: BindStringToObjectOptions{Types: []string{"string", "null"}},
			want: "abc",
		},
		{
			name:    "null alone is not a bindable member",
			src:     "abc",
			opts:    BindStringToObjectOptions{Types: []string{"null"}},
			wantErr: "type union has no bindable members (declared [null])",
		},
		{
			name: "non-scalar members are skipped",
			src:  "a,b,c",
			opts: BindStringToObjectOptions{Types: []string{"array", "string"}},
			want: "a,b,c",
		},
		{
			name:    "no member parses",
			src:     "abc",
			opts:    BindStringToObjectOptions{Types: []string{"integer", "boolean"}},
			wantErr: "does not match any member of type union",
		},
		{
			name:    "only non-scalar members",
			src:     "abc",
			opts:    BindStringToObjectOptions{Types: []string{"array", "object"}},
			wantErr: "does not match any member of type union",
		},
		{
			name:    "error names only bindable members, null filtered out",
			src:     "abc",
			opts:    BindStringToObjectOptions{Types: []string{"null", "integer"}},
			wantErr: "type union [integer]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var dest any
			err := BindStringToObjectWithOptions(tc.src, &dest, tc.opts)
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, dest)
		})
	}
}

// An interface destination without a declared union keeps the historical
// error, so nothing changes for existing generated code.
func TestBindStringToObject_InterfaceWithoutTypesStillErrors(t *testing.T) {
	var dest any
	err := BindStringToObject("123", &dest)
	assert.ErrorContains(t, err, "can not bind to destination of type: interface")

	err = BindStringToObjectWithOptions("123", &dest, BindStringToObjectOptions{Type: "string"})
	assert.ErrorContains(t, err, "can not bind to destination of type: interface")
}

// Types is only consulted for interface destinations: a concrete destination
// keeps the reflection-driven path even when Types is set.
func TestBindStringToObject_ConcreteDestinationIgnoresTypes(t *testing.T) {
	var s string
	err := BindStringToObjectWithOptions("123", &s, BindStringToObjectOptions{Types: []string{"string", "integer"}})
	require.NoError(t, err)
	assert.Equal(t, "123", s)
}

// A nullable wrapper around `any` (nullable.Nullable[any] is map[bool]any)
// routes through the bool-keyed-map path and binds the inner value via the
// union walk.
func TestBindStringToObject_NullableAnyUnion(t *testing.T) {
	var dest nullable.Nullable[any]
	err := BindStringToObjectWithOptions("123", &dest, BindStringToObjectOptions{Types: []string{"string", "integer"}})
	require.NoError(t, err)
	got, err := dest.Get()
	require.NoError(t, err)
	assert.Equal(t, int64(123), got)
}

// End-to-end through the styled binder, as generated server code calls it
// for path and header parameters.
func TestBindStyledParameterWithOptions_Union(t *testing.T) {
	opts := BindStyledParameterOptions{
		ParamLocation:    ParamLocationPath,
		Required:         true,
		Types:            []string{"string", "integer"},
		ValueIsUnescaped: true,
	}

	var id any
	require.NoError(t, BindStyledParameterWithOptions("simple", "id", "42", &id, opts))
	assert.Equal(t, int64(42), id)

	require.NoError(t, BindStyledParameterWithOptions("simple", "id", "abc", &id, opts))
	assert.Equal(t, "abc", id)

	// label style strips its prefix before binding.
	require.NoError(t, BindStyledParameterWithOptions("label", "id", ".42", &id, opts))
	assert.Equal(t, int64(42), id)
}

// End-to-end through the query binder, exploded and non-exploded, required
// and optional.
func TestBindQueryParameterWithOptions_Union(t *testing.T) {
	opts := BindQueryParameterOptions{Types: []string{"string", "number", "boolean"}}
	queryParams := url.Values{
		"filter": {"1.5"},
		"flag":   {"true"},
		"note":   {"hello"},
	}

	var filter any
	require.NoError(t, BindQueryParameterWithOptions("form", true, false, "filter", queryParams, &filter, opts))
	assert.Equal(t, float64(1.5), filter)

	var flag any
	require.NoError(t, BindQueryParameterWithOptions("form", true, true, "flag", queryParams, &flag, opts))
	assert.Equal(t, true, flag)

	var note any
	require.NoError(t, BindQueryParameterWithOptions("form", false, false, "note", queryParams, &note, opts))
	assert.Equal(t, "hello", note)

	// An absent optional parameter leaves the destination untouched.
	var missing any
	require.NoError(t, BindQueryParameterWithOptions("form", true, false, "absent", queryParams, &missing, opts))
	assert.Nil(t, missing)

	// An absent required parameter is still an error.
	var absent any
	err := BindQueryParameterWithOptions("form", true, true, "absent", queryParams, &absent, opts)
	var requiredErr *RequiredParameterError
	assert.ErrorAs(t, err, &requiredErr)
}

// Opting in via the NarrowUnionNumericFormats package variable makes width
// formats load-bearing: int32/float narrow the produced dynamic type, and
// values outside the narrowed range fall through to the next member. Not
// t.Parallel(): mutates package-global state, like DefaultQueryEncoder.
func TestBindStringToObject_UnionNarrowingOptIn(t *testing.T) {
	NarrowUnionNumericFormats = true
	t.Cleanup(func() { NarrowUnionNumericFormats = false })

	testCases := []struct {
		name string
		src  string
		opts BindStringToObjectOptions
		want any
	}{
		{
			name: "format int32 narrows the integer member",
			src:  "123",
			opts: BindStringToObjectOptions{Types: []string{"integer", "string"}, Format: "int32"},
			want: int32(123),
		},
		{
			name: "int32 overflow falls through to string",
			src:  "3000000000",
			opts: BindStringToObjectOptions{Types: []string{"integer", "string"}, Format: "int32"},
			want: "3000000000",
		},
		{
			name: "int32 overflow falls through to number",
			src:  "3000000000",
			opts: BindStringToObjectOptions{Types: []string{"integer", "number"}, Format: "int32"},
			want: float64(3000000000),
		},
		{
			name: "format float narrows the number member",
			src:  "1.5",
			opts: BindStringToObjectOptions{Types: []string{"number", "string"}, Format: "float"},
			want: float32(1.5),
		},
		{
			name: "float32 overflow falls through to string",
			src:  "1e40",
			opts: BindStringToObjectOptions{Types: []string{"number", "string"}, Format: "float"},
			want: "1e40",
		},
		{
			name: "format int64 names the default width",
			src:  "123",
			opts: BindStringToObjectOptions{Types: []string{"integer", "string"}, Format: "int64"},
			want: int64(123),
		},
		{
			name: "format double names the default width",
			src:  "1.5",
			opts: BindStringToObjectOptions{Types: []string{"number", "string"}, Format: "double"},
			want: float64(1.5),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var dest any
			err := BindStringToObjectWithOptions(tc.src, &dest, tc.opts)
			require.NoError(t, err)
			assert.Equal(t, tc.want, dest)
		})
	}
}

// A non-empty interface destination never takes the union path, with or
// without Types: only empty interfaces (`any`) qualify.
func TestBindStringToObject_TypedInterfaceWithTypesStillErrors(t *testing.T) {
	var dest interface{ Foo() }
	err := BindStringToObjectWithOptions("42", &dest, BindStringToObjectOptions{Types: []string{"string", "integer"}})
	assert.ErrorContains(t, err, "can not bind to destination of type: interface")
}

// A named empty interface is still an empty interface: the NumMethod()==0
// guard admits it to the union walk exactly like a literal `any`.
func TestBindStringToObject_NamedEmptyInterfaceUnion(t *testing.T) {
	type myAny interface{}
	var dest myAny
	err := BindStringToObjectWithOptions("42", &dest, BindStringToObjectOptions{Types: []string{"string", "integer"}})
	require.NoError(t, err)
	assert.Equal(t, int64(42), dest)
}

// Non-interface kinds that flow through the Array -> Struct -> Interface ->
// Map -> default fallthrough chain keep their historical errors even when
// Types is set. This locks the chain ordering in bindstring.go's switch.
func TestBindStringToObject_NonInterfaceKindsWithTypesStillError(t *testing.T) {
	opts := BindStringToObjectOptions{Types: []string{"string", "integer"}}

	var arr [2]int
	assert.ErrorContains(t, BindStringToObjectWithOptions("42", &arr, opts),
		"can not bind to destination of type: array")

	var st struct{ X int }
	assert.ErrorContains(t, BindStringToObjectWithOptions("42", &st, opts),
		"can not bind to destination of type: struct")

	var m map[string]string
	assert.ErrorContains(t, BindStringToObjectWithOptions("42", &m, opts),
		"can not bind to destination of type: map")
}

// Styled binding still unescapes values by default (ValueIsUnescaped=false),
// and unescaping happens before member selection — percent-encoding can flip
// which member wins.
func TestBindStyledParameterWithOptions_UnionUnescapesFirst(t *testing.T) {
	opts := BindStyledParameterOptions{
		ParamLocation: ParamLocationPath,
		Required:      true,
		Types:         []string{"string", "integer"},
	}

	// "4%32" path-unescapes to "42", so the integer member wins.
	var v any
	require.NoError(t, BindStyledParameterWithOptions("simple", "id", "4%32", &v, opts))
	assert.Equal(t, int64(42), v)

	// Headers are never unescaped: "4%32" stays a string.
	hOpts := opts
	hOpts.ParamLocation = ParamLocationHeader
	require.NoError(t, BindStyledParameterWithOptions("simple", "X-Id", "4%32", &v, hOpts))
	assert.Equal(t, "4%32", v)
}

// Matrix style strips its ;name= prefix before the union walk, like label
// strips its dot.
func TestBindStyledParameterWithOptions_UnionMatrixStyle(t *testing.T) {
	opts := BindStyledParameterOptions{
		ParamLocation:    ParamLocationPath,
		Required:         true,
		Types:            []string{"string", "integer"},
		ValueIsUnescaped: true,
	}
	var v any
	require.NoError(t, BindStyledParameterWithOptions("matrix", "id", ";id=42", &v, opts))
	assert.Equal(t, int64(42), v)
}

// Optional query parameters generated without x-go-type-skip-optional-pointer
// arrive as **any (pointer to the struct's *any field). The extra-indirect
// path must allocate and bind through it.
func TestBindQueryParameterWithOptions_UnionOptionalPointer(t *testing.T) {
	opts := BindQueryParameterOptions{Types: []string{"string", "integer"}}
	queryParams := url.Values{"filter": {"42"}}

	var filter *any
	require.NoError(t, BindQueryParameterWithOptions("form", true, false, "filter", queryParams, &filter, opts))
	require.NotNil(t, filter)
	assert.Equal(t, int64(42), *filter)

	// Absent optional parameter leaves the pointer nil.
	var absent *any
	require.NoError(t, BindQueryParameterWithOptions("form", true, false, "absent", queryParams, &absent, opts))
	assert.Nil(t, absent)
}

// A parameter present with an empty value binds the empty string when a
// string member exists, and errors when it does not.
func TestBindQueryParameterWithOptions_UnionEmptyValue(t *testing.T) {
	queryParams := url.Values{"p": {""}}

	var withString any
	require.NoError(t, BindQueryParameterWithOptions("form", true, true, "p", queryParams, &withString,
		BindQueryParameterOptions{Types: []string{"integer", "string"}}))
	assert.Equal(t, "", withString)

	var withoutString any
	err := BindQueryParameterWithOptions("form", true, false, "p", queryParams, &withoutString,
		BindQueryParameterOptions{Types: []string{"integer"}})
	assert.ErrorContains(t, err, "does not match any member of type union")
}

func TestIsJSONNumber(t *testing.T) {
	valid := []string{"0", "-0", "1", "-1", "123", "1.5", "-1.5", "0.5", "1e2", "1E2", "1e+2", "1e-2", "1.5e10", "0.0"}
	for _, s := range valid {
		assert.True(t, isJSONNumber(s), "expected %q to be a JSON number", s)
	}
	invalid := []string{"", "-", "+1", "01", "007", "-01", ".5", "1.", "1e", "1e+", "0x10", "1 ", " 1", "1,000", "NaN", "Infinity", "--1", "1..5", "1.5.5"}
	for _, s := range invalid {
		assert.False(t, isJSONNumber(s), "expected %q to NOT be a JSON number", s)
	}
}

func TestIsJSONInteger(t *testing.T) {
	valid := []string{"0", "-0", "1", "-1", "123", "-123"}
	for _, s := range valid {
		assert.True(t, isJSONInteger(s), "expected %q to be a JSON integer", s)
	}
	invalid := []string{"", "1.5", "1e2", "01", "+1", "1.0", "-", "abc"}
	for _, s := range invalid {
		assert.False(t, isJSONInteger(s), "expected %q to NOT be a JSON integer", s)
	}
}
