// Copyright 2019 DeepMap, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultQueryEncoder(t *testing.T) {
	// The default encoder must not have been changed by any other test.
	assert.IsType(t, NetURLQueryEncoder{}, DefaultQueryEncoder,
		"DefaultQueryEncoder should default to NetURLQueryEncoder")

	// withQueryEncoder temporarily installs enc, restoring the previous
	// encoder when the test finishes.
	withQueryEncoder := func(t *testing.T, enc QueryEncoder) {
		t.Helper()
		prev := DefaultQueryEncoder
		DefaultQueryEncoder = enc
		t.Cleanup(func() { DefaultQueryEncoder = prev })
	}

	t.Run("net/url encoder keeps + for spaces (default)", func(t *testing.T) {
		withQueryEncoder(t, NetURLQueryEncoder{})

		result, err := StyleParamWithLocation("form", false, "filter", ParamLocationQuery, "name eq 'x'")
		assert.NoError(t, err)
		assert.EqualValues(t, "filter=name+eq+%27x%27", result)
	})

	t.Run("RFC 3986 encoder uses %20 for spaces", func(t *testing.T) {
		withQueryEncoder(t, RFC3986QueryEncoder{})

		result, err := StyleParamWithLocation("form", false, "filter", ParamLocationQuery, "name eq 'x'")
		assert.NoError(t, err)
		assert.EqualValues(t, "filter=name%20eq%20%27x%27", result)
	})

	t.Run("RFC 3986 encoder leaves literal + intact", func(t *testing.T) {
		withQueryEncoder(t, RFC3986QueryEncoder{})

		// A real '+' becomes %2B; only spaces become %20. This confirms the
		// +->%20 rewrite does not corrupt genuine plus characters.
		result, err := StyleParamWithLocation("form", false, "ts", ParamLocationQuery, "2020-01-01 22:00:00+02:00")
		assert.NoError(t, err)
		assert.EqualValues(t, "ts=2020-01-01%2022%3A00%3A00%2B02%3A00", result)
	})

	t.Run("encoder applies to exploded array values", func(t *testing.T) {
		withQueryEncoder(t, RFC3986QueryEncoder{})

		result, err := StyleParamWithLocation("form", true, "q", ParamLocationQuery, []string{"a b", "c d"})
		assert.NoError(t, err)
		assert.EqualValues(t, "q=a%20b&q=c%20d", result)
	})

	t.Run("encoder applies to param names", func(t *testing.T) {
		withQueryEncoder(t, RFC3986QueryEncoder{})

		result, err := StyleParamWithLocation("form", false, "my filter", ParamLocationQuery, "v")
		assert.NoError(t, err)
		assert.EqualValues(t, "my%20filter=v", result)
	})

	t.Run("allowReserved is RFC 3986 regardless of encoder", func(t *testing.T) {
		// The allowReserved path already emits %20 for spaces, so both
		// encoders agree.
		for _, enc := range []QueryEncoder{NetURLQueryEncoder{}, RFC3986QueryEncoder{}} {
			withQueryEncoder(t, enc)
			result, err := StyleParamWithOptions("form", false, "q", "hello world", StyleParamOptions{
				ParamLocation: ParamLocationQuery,
				AllowReserved: true,
			})
			assert.NoError(t, err)
			assert.EqualValues(t, "q=hello%20world", result)
		}
	})

	t.Run("encoder does not affect path params", func(t *testing.T) {
		withQueryEncoder(t, RFC3986QueryEncoder{})

		// Path params always use url.PathEscape, which already encodes a
		// space as %20 and is unaffected by the query encoder.
		result, err := StyleParamWithLocation("simple", false, "id", ParamLocationPath, "a b")
		assert.NoError(t, err)
		assert.EqualValues(t, "a%20b", result)
	})
}

// TestDeepObject_QueryEncoder verifies that MarshalDeepObject routes value,
// key, and param-name escaping through DefaultQueryEncoder, so the RFC 3986
// encoder produces %20 for spaces instead of '+'.
func TestDeepObject_QueryEncoder(t *testing.T) {
	prev := DefaultQueryEncoder
	DefaultQueryEncoder = RFC3986QueryEncoder{}
	t.Cleanup(func() { DefaultQueryEncoder = prev })

	src := map[string]interface{}{
		"full name": "Ada Lovelace",
	}

	marshaled, err := MarshalDeepObject(src, "my filter")
	require.NoError(t, err)

	// Param name, map key, and value all use %20 for their spaces; the
	// structural '[' and ']' delimiters remain literal.
	assert.Equal(t, "my%20filter[full%20name]=Ada%20Lovelace", marshaled)
}
