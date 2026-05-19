package runtime

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type InnerArrayObject struct {
	Names []string `json:"names"`
}

type InnerObject struct {
	Name string
	ID   int
}

type InnerObject2 struct {
	Foo string
	Is  bool
}

type InnerObject3 struct {
	Foo   string
	Count *int `json:"count,omitempty"`
}

// These are all possible field types, mandatory and optional.
type AllFields struct {
	I    int              `json:"i"`
	Oi   *int             `json:"oi,omitempty"`
	Ab   *[]bool          `json:"ab,omitempty"`
	F    float32          `json:"f"`
	Of   *float32         `json:"of,omitempty"`
	B    bool             `json:"b"`
	Ob   *bool            `json:"ob,omitempty"`
	As   []string         `json:"as"`
	Oas  *[]string        `json:"oas,omitempty"`
	O    InnerObject      `json:"o"`
	Ao   []InnerObject2   `json:"ao"`
	Aop  *[]InnerObject3  `json:"aop"`
	Onas InnerArrayObject `json:"onas"`
	Oo   *InnerObject     `json:"oo,omitempty"`
	D    MockBinder       `json:"d"`
	Od   *MockBinder      `json:"od,omitempty"`
	M    map[string]int   `json:"m"`
	Om   *map[string]int  `json:"om,omitempty"`
}

func TestDeepObject(t *testing.T) {
	oi := int(5)
	of := float32(3.7)
	ob := true
	oas := []string{"foo", "bar"}
	oo := InnerObject{
		Name: "Marcin Romaszewicz",
		ID:   123,
	}
	om := map[string]int{
		"additional": 1,
	}

	d := MockBinder{Time: time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)}

	two := 2

	srcObj := AllFields{
		I:   12,
		Oi:  &oi,
		F:   4.2,
		Of:  &of,
		B:   true,
		Ob:  &ob,
		Ab:  &[]bool{true},
		As:  []string{"hello", "world"},
		Oas: &oas,
		O: InnerObject{
			Name: "Joe Schmoe",
			ID:   456,
		},
		Ao: []InnerObject2{
			{Foo: "bar", Is: true},
			{Foo: "baz"},
		},
		Aop: &[]InnerObject3{
			{Foo: "a"},
			{Foo: "b", Count: &two},
		},
		Onas: InnerArrayObject{
			Names: []string{"Bill", "Frank"},
		},
		Oo: &oo,
		D:  d,
		Od: &d,
		M:  om,
		Om: &om,
	}

	marshaled, err := MarshalDeepObject(srcObj, "p")
	require.NoError(t, err)
	// Spaces in values are encoded as '+' per url.QueryEscape (form-encoded
	// convention, matching the pre-v2.7.0 behavior of url.Values.Encode()).
	require.EqualValues(t, "p[ab][0]=true&p[ao][0][Foo]=bar&p[ao][0][Is]=true&p[ao][1][Foo]=baz&p[ao][1][Is]=false&p[aop][0][Foo]=a&p[aop][1][Foo]=b&p[aop][1][count]=2&p[as][0]=hello&p[as][1]=world&p[b]=true&p[d]=2020-02-01&p[f]=4.2&p[i]=12&p[m][additional]=1&p[o][ID]=456&p[o][Name]=Joe+Schmoe&p[oas][0]=foo&p[oas][1]=bar&p[ob]=true&p[od]=2020-02-01&p[of]=3.7&p[oi]=5&p[om][additional]=1&p[onas][names][0]=Bill&p[onas][names][1]=Frank&p[oo][ID]=123&p[oo][Name]=Marcin+Romaszewicz", marshaled)

	// Use url.ParseQuery for the round-trip — it decodes percent-escapes and
	// '+' back to literal characters, matching what a real HTTP server does.
	params, err := url.ParseQuery(marshaled)
	require.NoError(t, err)

	var dstObj AllFields
	err = UnmarshalDeepObject(&dstObj, "p", params)
	require.NoError(t, err)
	assert.EqualValues(t, srcObj, dstObj)
}

// Item represents an item object for testing array of objects
type Item struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func TestDeepObject_TimeFields(t *testing.T) {
	type TimeObject struct {
		Created time.Time `json:"created"`
	}

	t.Run("RFC3339 time parses correctly", func(t *testing.T) {
		params := url.Values{}
		params.Set("p[created]", "2024-01-15T10:30:00Z")

		var dst TimeObject
		err := UnmarshalDeepObject(&dst, "p", params)
		require.NoError(t, err)
		assert.Equal(t, time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC), dst.Created)
	})

	t.Run("date-only string parses correctly as time.Time", func(t *testing.T) {
		params := url.Values{}
		params.Set("p[created]", "2024-01-15")

		var dst TimeObject
		err := UnmarshalDeepObject(&dst, "p", params)
		require.NoError(t, err)
		assert.Equal(t, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), dst.Created)
	})

	t.Run("invalid time string returns error", func(t *testing.T) {
		params := url.Values{}
		params.Set("p[created]", "not-a-time")

		var dst TimeObject
		err := UnmarshalDeepObject(&dst, "p", params)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error parsing")
	})
}

func TestDeepObject_ArrayOfObjects(t *testing.T) {
	// Test case for:
	// name: items
	// style: deepObject
	// required: false
	// explode: true
	// schema:
	//   type: array
	//   minItems: 1
	//   items:
	//     type: object

	srcArray := []Item{
		{Name: "first", Value: "value1"},
		{Name: "second", Value: "value2"},
	}

	// Marshal the array to deepObject format
	marshaled, err := MarshalDeepObject(srcArray, "items")
	require.NoError(t, err)
	t.Log("Marshaled:", marshaled)

	// Expected format for array of objects with explode:true should be:
	// items[0][name]=first&items[0][value]=value1&items[1][name]=second&items[1][value]=value2

	// Parse the marshaled string into url.Values
	params := make(url.Values)
	marshaledParts := strings.Split(marshaled, "&")
	for _, p := range marshaledParts {
		parts := strings.Split(p, "=")
		require.Equal(t, 2, len(parts))
		params.Set(parts[0], parts[1])
	}

	// Unmarshal back to the destination array
	var dstArray []Item
	err = UnmarshalDeepObject(&dstArray, "items", params)
	require.NoError(t, err)

	// Verify the result matches the source
	assert.EqualValues(t, srcArray, dstArray)
	assert.Len(t, dstArray, 2)
	assert.Equal(t, "first", dstArray[0].Name)
	assert.Equal(t, "value1", dstArray[0].Value)
	assert.Equal(t, "second", dstArray[1].Name)
	assert.Equal(t, "value2", dstArray[1].Value)
}

func TestDeepObject_NonIndexedArray(t *testing.T) {
	t.Run("primitive string array", func(t *testing.T) {
		params := url.Values{}
		params.Add("p[vals]", "a")
		params.Add("p[vals]", "b")

		type Obj struct {
			Vals []string `json:"vals"`
		}

		var dst Obj
		err := UnmarshalDeepObject(&dst, "p", params)
		require.NoError(t, err)
		assert.Equal(t, []string{"a", "b"}, dst.Vals)
	})

	t.Run("object with mixed scalar and non-indexed array", func(t *testing.T) {
		params := url.Values{}
		params.Set("p[op]", "eq")
		params.Add("p[vals]", "a")
		params.Add("p[vals]", "b")

		type Filter struct {
			Op   string   `json:"op"`
			Vals []string `json:"vals"`
		}

		var dst Filter
		err := UnmarshalDeepObject(&dst, "p", params)
		require.NoError(t, err)
		assert.Equal(t, "eq", dst.Op)
		assert.Equal(t, []string{"a", "b"}, dst.Vals)
	})
}

// assertDeepObjectWireSafe asserts the marshaled deepObject output is safe to
// use on the wire and round-trips correctly. Three invariants:
//
//  1. Every byte is 7-bit ASCII. RFC 3986 requires non-ASCII bytes to be
//     percent-encoded; Go's net/url is lenient and won't surface the bug, so
//     we check the bytes directly.
//  2. The output survives a full URL round trip:
//     URL.RawQuery → URL.String() → url.Parse() → URL.Query()
//     This catches characters (notably '#') that ParseQuery accepts in
//     isolation but that break when the string is assigned to URL.RawQuery
//     and re-parsed as part of a full URL.
//  3. Re-parsed values + UnmarshalDeepObject reconstruct the original.
func assertDeepObjectWireSafe(t *testing.T, marshaled, paramName string, dst, src interface{}) {
	t.Helper()

	if i := strings.IndexFunc(marshaled, func(r rune) bool {
		return r > unicode.MaxASCII
	}); i >= 0 {
		t.Fatalf("non-ASCII rune at byte offset %d; deepObject output must be percent-encoded; raw=%q", i, marshaled)
	}

	u, err := url.Parse("http://example.test/path")
	require.NoError(t, err)
	u.RawQuery = marshaled
	roundtrippedURL := u.String()

	parsedURL, err := url.Parse(roundtrippedURL)
	require.NoError(t, err, "URL round-trip failed; raw=%q url=%q", marshaled, roundtrippedURL)
	require.Empty(t, parsedURL.Fragment,
		"unexpected URL fragment %q after round-trip — unescaped '#' leaked into RawQuery; raw=%q",
		parsedURL.Fragment, marshaled)

	parsedQuery := parsedURL.Query()
	for k := range parsedQuery {
		require.True(t, strings.HasPrefix(k, paramName+"["),
			"unexpected top-level key %q after URL round-trip; raw=%q", k, marshaled)
	}

	err = UnmarshalDeepObject(dst, paramName, parsedQuery)
	require.NoError(t, err, "unmarshal failed; raw=%q", marshaled)
	assert.Equal(t, src, reflect.ValueOf(dst).Elem().Interface(),
		"round-trip mismatch; raw=%q", marshaled)
}

// TestDeepObject_URLEncoding pins issue
// https://github.com/oapi-codegen/runtime/issues/131: the marshaller must
// percent-encode reserved characters and non-ASCII bytes in values, in map
// keys, and in the param name itself. Without this, an `&` in a value
// silently injects an extra query parameter and non-ASCII content produces
// invalid URIs that RFC3986-compliant servers reject with 400.
func TestDeepObject_URLEncoding(t *testing.T) {
	type Inner struct {
		Name string `json:"name"`
	}
	type Outer struct {
		Field string `json:"field"`
		Inner Inner  `json:"inner"`
	}

	// Adversarial value cases. Each case puts the bad string in Field, with a
	// safe constant value in Inner.Name so we can also detect injection (a
	// stray top-level "inner.name=safe" pair would prove the value escaped its
	// enclosing param).
	valueCases := []struct {
		name string
		val  string
	}{
		{"empty", ""},
		{"plain ascii", "Alex"},
		{"space", "a b c"},
		{"ampersand", "a&admin=true"},
		{"equals", "k=v"},
		{"question mark", "what?"},
		{"hash fragment", "before#after"},
		{"plus literal", "1+2=3"},
		{"semicolon", "a;b;c"},
		{"comma", "a,b,c"},
		{"slash", "a/b/c"},
		{"colon", "key:value"},
		{"at sign", "user@host"},
		{"percent literal", "100%"},
		{"square brackets", "[bracketed]"},
		{"curly braces", "{json}"},
		{"reserved kitchen sink", "&=?#+;,/:@!$'()*[]"},
		{"single quote", "it's"},
		{"double quote", `say "hi"`},
		{"backtick", "`code`"},
		{"backslash", `a\b`},
		{"tab", "a\tb"},
		{"newline", "a\nb"},
		{"carriage return", "a\rb"},
		{"null byte", "a\x00b"},
		{"cjk hiragana", "こんにちは"},
		{"cjk han", "你好世界"},
		{"cyrillic", "Привет"},
		{"arabic", "مرحبا"},
		{"hebrew", "שלום"},
		{"thai", "สวัสดี"},
		{"emoji 4-byte", "🚀"},
		{"emoji zwj sequence", "👨‍👩‍👧‍👦"},
		{"combining accents", "é"}, // e + combining acute
		{"surrogate-range", "\U0001F600\U0001F601"},
		{"mixed unicode and reserved", "filter&q=こんにちは"},
		{"long repeated reserved", strings.Repeat("&", 64)},
	}

	for _, tc := range valueCases {
		t.Run("value/"+tc.name, func(t *testing.T) {
			src := Outer{Field: tc.val, Inner: Inner{Name: "safe"}}

			marshaled, err := MarshalDeepObject(src, "p")
			require.NoError(t, err)

			var dst Outer
			assertDeepObjectWireSafe(t, marshaled, "p", &dst, src)
		})
	}

	// Adversarial map-key cases. JSON-encoded structs cannot produce
	// arbitrary keys, but map[string]string can, and that's a real user
	// pattern (filters[name]=... where the filter key comes from a user).
	keyCases := []struct {
		name string
		key  string
	}{
		{"plain ascii", "name"},
		{"space in key", "first name"},
		{"ampersand in key", "a&b"},
		{"equals in key", "a=b"},
		{"dollar prefix", "$eq"},
		{"non-ascii key", "名前"},
		{"emoji key", "🔑"},
		{"reserved kitchen sink key", "&=?#+;,/:@"},
	}

	for _, tc := range keyCases {
		t.Run("key/"+tc.name, func(t *testing.T) {
			src := map[string]string{tc.key: "safe"}

			marshaled, err := MarshalDeepObject(src, "p")
			require.NoError(t, err)

			var dst map[string]string
			assertDeepObjectWireSafe(t, marshaled, "p", &dst, src)
		})
	}

	// Adversarial param-name. The codegen emits the parameter name as
	// declared in the OpenAPI spec, which is usually a safe identifier, but
	// nothing in the spec forbids spaces or non-ASCII names — and the
	// marshaller prepends the name to each fragment without escaping.
	paramNameCases := []string{
		"plain",
		"with space",
		"with&amp",
		"フィルター",
		"🔥",
	}
	for _, pn := range paramNameCases {
		t.Run("paramName/"+pn, func(t *testing.T) {
			src := map[string]string{"name": "safe"}

			marshaled, err := MarshalDeepObject(src, pn)
			require.NoError(t, err)

			// Marshaled output must be 7-bit ASCII regardless of paramName
			// content. We can't use the full helper here because the helper
			// checks for `paramName+"["` literally in parsed keys; once
			// encoded, the prefix is the percent-encoded form. Just check the
			// ASCII invariant and that the output parses.
			if i := strings.IndexFunc(marshaled, func(r rune) bool {
				return r > unicode.MaxASCII
			}); i >= 0 {
				t.Fatalf("non-ASCII rune at byte offset %d in paramName-prefixed output; raw=%q", i, marshaled)
			}

			u, err := url.Parse("http://example.test/path")
			require.NoError(t, err)
			u.RawQuery = marshaled
			parsedURL, err := url.Parse(u.String())
			require.NoError(t, err, "URL round-trip failed; raw=%q", marshaled)
			require.Empty(t, parsedURL.Fragment, "unexpected fragment; raw=%q", marshaled)
		})
	}

	// Non-ASCII assertion: the issue text spells out the exact expected
	// percent-encoded form for "こんにちは". Lock that in so we don't drift to
	// some other encoder.
	t.Run("non-ascii uses utf-8 percent-encoding (RFC3986 §2.5)", func(t *testing.T) {
		src := map[string]string{"name": "こんにちは"}
		marshaled, err := MarshalDeepObject(src, "p")
		require.NoError(t, err)
		// "こんにちは" → E3 81 93 E3 82 93 E3 81 AB E3 81 A1 E3 81 AF.
		// Each byte becomes %XX; url.QueryEscape uses uppercase hex.
		assert.Contains(t, marshaled,
			"%E3%81%93%E3%82%93%E3%81%AB%E3%81%A1%E3%81%AF",
			"expected UTF-8 percent-encoded value; got %q", marshaled)
	})
}
