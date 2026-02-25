package runtime

import (
	"encoding/base64"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsByteSlice(t *testing.T) {
	assert.True(t, isByteSlice(reflect.TypeOf([]byte{})))
	assert.True(t, isByteSlice(reflect.TypeOf([]uint8{})))

	assert.False(t, isByteSlice(reflect.TypeOf([]int{})))
	assert.False(t, isByteSlice(reflect.TypeOf([]string{})))
	assert.False(t, isByteSlice(reflect.TypeOf("")))
	assert.False(t, isByteSlice(reflect.TypeOf(0)))

	// Type alias for []byte should also be detected
	type MyBytes []byte
	assert.True(t, isByteSlice(reflect.TypeOf(MyBytes{})))
}

func TestBase64Decode(t *testing.T) {
	t.Run("standard encoding with padding", func(t *testing.T) {
		// "test" → StdEncoding → "dGVzdA=="
		b, err := base64Decode("dGVzdA==")
		require.NoError(t, err)
		assert.Equal(t, []byte("test"), b)
	})

	t.Run("standard encoding without padding", func(t *testing.T) {
		// "test" → RawStdEncoding → "dGVzdA"
		b, err := base64Decode("dGVzdA")
		require.NoError(t, err)
		assert.Equal(t, []byte("test"), b)
	})

	t.Run("URL-safe encoding without padding", func(t *testing.T) {
		// Contains '_' and '-' → dispatches to RawURLEncoding
		input := "PDw_Pz4-"
		b, err := base64Decode(input)
		require.NoError(t, err)
		expected, decErr := base64.RawURLEncoding.DecodeString(input)
		require.NoError(t, decErr)
		assert.Equal(t, expected, b)
	})

	t.Run("URL-safe encoding with padding", func(t *testing.T) {
		// Contains '_' and '=' → dispatches to URLEncoding (padded)
		input := base64.URLEncoding.EncodeToString([]byte("<<??>>"))
		b, err := base64Decode(input)
		require.NoError(t, err)
		assert.Equal(t, []byte("<<??>>"), b)
	})

	t.Run("empty string", func(t *testing.T) {
		b, err := base64Decode("")
		require.NoError(t, err)
		assert.Equal(t, []byte{}, b)
	})

	t.Run("invalid base64", func(t *testing.T) {
		_, err := base64Decode("!!!not-base64!!!")
		assert.Error(t, err)
	})

	t.Run("padded input not corrupted by wrong decoder", func(t *testing.T) {
		// This is the key correctness test. With the old blind cascade,
		// RawStdEncoding would "succeed" on padded input but produce
		// garbage bytes (treating '=' as data). The new logic dispatches
		// padded input to StdEncoding directly.
		input := base64.StdEncoding.EncodeToString([]byte("hello world"))
		b, err := base64Decode(input)
		require.NoError(t, err)
		assert.Equal(t, []byte("hello world"), b)
	})
}
