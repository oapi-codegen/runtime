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
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net"
	"net/netip"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/oapi-codegen/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oapi-codegen/runtime/types"
)

// CustomStringBinder is a string type that implements Binder but not TextUnmarshaler
// This tests that Binder interface is checked for primitive types
type CustomStringBinder string

func (c *CustomStringBinder) Bind(src string) error {
	// Custom binding logic: add a prefix to demonstrate the Bind method was called
	*c = CustomStringBinder("CUSTOM:" + src)
	return nil
}

func TestBindStringToObject(t *testing.T) {
	var i int
	assert.NoError(t, BindStringToObject("5", &i))
	assert.Equal(t, 5, i)

	// Let's make sure we error out on things that can't be the correct
	// type. Since we're using reflect package setters, we'll have similar
	// unassignable type errors.
	assert.Error(t, BindStringToObject("5.7", &i))
	assert.Error(t, BindStringToObject("foo", &i))
	assert.Error(t, BindStringToObject("1,2,3", &i))

	var i8 int8
	assert.NoError(t, BindStringToObject("12", &i8))
	assert.Equal(t, int8(12), i8)
	assert.NoError(t, BindStringToObject("-12", &i8))
	assert.Equal(t, int8(-12), i8)

	assert.Error(t, BindStringToObject("5.7", &i8))
	assert.Error(t, BindStringToObject("foo", &i8))
	assert.Error(t, BindStringToObject("1,2,3", &i8))
	assert.Error(t, BindStringToObject(fmt.Sprintf("%d", math.MinInt8-1), &i8))
	assert.Error(t, BindStringToObject(fmt.Sprintf("%d", math.MaxInt8+1), &i8))

	var i16 int16
	assert.NoError(t, BindStringToObject("12", &i16))
	assert.Equal(t, int16(12), i16)
	assert.NoError(t, BindStringToObject("-12", &i16))
	assert.Equal(t, int16(-12), i16)

	assert.Error(t, BindStringToObject("5.7", &i16))
	assert.Error(t, BindStringToObject("foo", &i16))
	assert.Error(t, BindStringToObject("1,2,3", &i16))
	assert.Error(t, BindStringToObject(fmt.Sprintf("%d", math.MinInt16-1), &i16))
	assert.Error(t, BindStringToObject(fmt.Sprintf("%d", math.MaxInt16+1), &i16))

	var i32 int32
	assert.NoError(t, BindStringToObject("12", &i32))
	assert.Equal(t, int32(12), i32)
	assert.NoError(t, BindStringToObject("-12", &i32))
	assert.Equal(t, int32(-12), i32)

	assert.Error(t, BindStringToObject("5.7", &i32))
	assert.Error(t, BindStringToObject("foo", &i32))
	assert.Error(t, BindStringToObject("1,2,3", &i32))
	assert.Error(t, BindStringToObject(fmt.Sprintf("%d", math.MinInt32-1), &i32))
	assert.Error(t, BindStringToObject(fmt.Sprintf("%d", math.MaxInt32+1), &i32))

	var i64 int64
	assert.NoError(t, BindStringToObject("124", &i64))
	assert.Equal(t, int64(124), i64)
	assert.NoError(t, BindStringToObject("-124", &i64))
	assert.Equal(t, int64(-124), i64)

	assert.Error(t, BindStringToObject("5.7", &i64))
	assert.Error(t, BindStringToObject("foo", &i64))
	assert.Error(t, BindStringToObject("1,2,3", &i64))
	assert.Error(t, BindStringToObject("-9223372036854775809", &i64))
	assert.Error(t, BindStringToObject("9223372036854775808", &i64)) // 1<<63

	var u uint
	assert.NoError(t, BindStringToObject("5", &u))
	assert.Equal(t, uint(5), u)

	assert.Error(t, BindStringToObject("-5", &u))
	assert.Error(t, BindStringToObject("5.7", &u))
	assert.Error(t, BindStringToObject("foo", &u))
	assert.Error(t, BindStringToObject("1,2,3", &u))

	var u8 uint8
	assert.NoError(t, BindStringToObject("12", &u8))
	assert.Equal(t, uint8(12), u8)

	assert.Error(t, BindStringToObject("-5", &u8))
	assert.Error(t, BindStringToObject("5.7", &u8))
	assert.Error(t, BindStringToObject("foo", &u8))
	assert.Error(t, BindStringToObject("1,2,3", &u8))
	assert.Error(t, BindStringToObject(fmt.Sprintf("%d", math.MaxUint8+1), &i8))

	var u16 uint16
	assert.NoError(t, BindStringToObject("12", &u16))
	assert.Equal(t, uint16(12), u16)

	assert.Error(t, BindStringToObject("-5", &u16))
	assert.Error(t, BindStringToObject("5.7", &u16))
	assert.Error(t, BindStringToObject("foo", &u16))
	assert.Error(t, BindStringToObject("1,2,3", &u16))
	assert.Error(t, BindStringToObject(fmt.Sprintf("%d", math.MaxUint16+1), &i16))

	var u32 uint32
	assert.NoError(t, BindStringToObject("12", &u32))
	assert.Equal(t, uint32(12), u32)

	assert.Error(t, BindStringToObject("-5", &u32))
	assert.Error(t, BindStringToObject("5.7", &u32))
	assert.Error(t, BindStringToObject("foo", &u32))
	assert.Error(t, BindStringToObject("1,2,3", &u32))
	assert.Error(t, BindStringToObject(fmt.Sprintf("%d", math.MaxUint32+1), &i32))

	var u64 uint64
	assert.NoError(t, BindStringToObject("124", &u64))
	assert.Equal(t, uint64(124), u64)

	assert.Error(t, BindStringToObject("-5", &u64))
	assert.Error(t, BindStringToObject("5.7", &u64))
	assert.Error(t, BindStringToObject("foo", &u64))
	assert.Error(t, BindStringToObject("1,2,3", &u64))
	assert.Error(t, BindStringToObject("18446744073709551616", &i64)) // 1<<64

	var b bool
	assert.NoError(t, BindStringToObject("True", &b))
	assert.Equal(t, true, b)
	assert.NoError(t, BindStringToObject("true", &b))
	assert.Equal(t, true, b)
	assert.NoError(t, BindStringToObject("1", &b))
	assert.Equal(t, true, b)
	assert.NoError(t, BindStringToObject("0", &b))
	assert.Equal(t, false, b)
	assert.Error(t, BindStringToObject("-1", &b))
	assert.Error(t, BindStringToObject("hello", &b))

	var f64 float64
	assert.NoError(t, BindStringToObject("1.25", &f64))
	assert.Equal(t, float64(1.25), f64)

	assert.Error(t, BindStringToObject("foo", &f64))
	assert.Error(t, BindStringToObject("1,2,3", &f64))

	var f32 float32
	assert.NoError(t, BindStringToObject("3.125", &f32))
	assert.Equal(t, float32(3.125), f32)

	assert.Error(t, BindStringToObject("foo", &f32))
	assert.Error(t, BindStringToObject("1,2,3", &f32))
	assert.NoError(t, BindStringToObject(fmt.Sprintf("%f", math.MaxFloat32), &f32))
	assert.Error(t, BindStringToObject(fmt.Sprintf("%f", math.MaxFloat32*2), &f32))

	// This checks whether binding works through a type alias.
	type SomeType int
	var st SomeType
	assert.NoError(t, BindStringToObject("5", &st))
	assert.Equal(t, SomeType(5), st)

	// Check time binding
	now := time.Now().UTC()
	strTime := now.Format(time.RFC3339Nano)
	var parsedTime time.Time
	assert.NoError(t, BindStringToObject(strTime, &parsedTime))
	parsedTime = parsedTime.UTC()
	assert.EqualValues(t, now, parsedTime)

	now = now.Truncate(time.Second)
	strTime = now.Format(time.RFC3339)
	assert.NoError(t, BindStringToObject(strTime, &parsedTime))
	parsedTime = parsedTime.UTC()
	assert.EqualValues(t, now, parsedTime)

	// Checks whether time binding works through a type alias.
	type AliasedTime time.Time
	var aliasedTime AliasedTime
	assert.NoError(t, BindStringToObject(strTime, &aliasedTime))
	assert.EqualValues(t, now, aliasedTime)

	// Checks whether date binding works directly and through an alias.
	dateString := "2020-11-05"
	var dstDate types.Date
	assert.NoError(t, BindStringToObject(dateString, &dstDate))
	type AliasedDate types.Date
	var dstAliasedDate AliasedDate
	assert.NoError(t, BindStringToObject(dateString, &dstAliasedDate))

	// Checks whether a mock binder works and embedded types
	var mockBinder MockBinder
	assert.NoError(t, BindStringToObject(dateString, &mockBinder))
	assert.EqualValues(t, dateString, mockBinder.Format("2006-01-02"))
	var dstEmbeddedMockBinder EmbeddedMockBinder
	assert.NoError(t, BindStringToObject(dateString, &dstEmbeddedMockBinder))
	assert.EqualValues(t, dateString, dstEmbeddedMockBinder.Format("2006-01-02"))

	// Checks UUID binding
	uuidString := "bbca1470-5e1f-4c64-ba99-fa7a6d2687b0"
	var dstUUID types.UUID
	assert.NoError(t, BindStringToObject(uuidString, &dstUUID))
	assert.Equal(t, dstUUID.String(), uuidString)

	// Checks that primitive types implementing Binder are respected
	// This tests the fix for ensuring Binder is checked before type reflection
	var customString CustomStringBinder
	assert.NoError(t, BindStringToObject("hello", &customString))
	assert.Equal(t, CustomStringBinder("CUSTOM:hello"), customString)

}

// TestBindStringToObject_ByteSlice tests that BindStringToObject correctly handles
// *[]byte destinations by base64-decoding the input string, rather than failing
// or treating []byte as a generic slice.
// See: https://github.com/oapi-codegen/runtime/issues/97
func TestBindStringToObject_ByteSlice(t *testing.T) {
	opts := BindStringToObjectOptions{Type: "string", Format: "byte"}

	t.Run("valid base64 with padding", func(t *testing.T) {
		var dest []byte
		err := BindStringToObjectWithOptions("dGVzdA==", &dest, opts)
		require.NoError(t, err)
		assert.Equal(t, []byte("test"), dest)
	})

	t.Run("valid base64 without padding", func(t *testing.T) {
		var dest []byte
		err := BindStringToObjectWithOptions("dGVzdA", &dest, opts)
		require.NoError(t, err)
		assert.Equal(t, []byte("test"), dest)
	})

	t.Run("URL-safe base64", func(t *testing.T) {
		// "<<??>>+" in standard base64 is "PDw/Pz4+" but URL-safe uses "PDw_Pz4-"
		input := "PDw_Pz4-"
		var dest []byte
		err := BindStringToObjectWithOptions(input, &dest, opts)
		require.NoError(t, err)
		expected, decErr := base64.RawURLEncoding.DecodeString("PDw_Pz4-")
		require.NoError(t, decErr)
		assert.Equal(t, expected, dest)
	})

	t.Run("empty string", func(t *testing.T) {
		var dest []byte
		err := BindStringToObjectWithOptions("", &dest, opts)
		require.NoError(t, err)
		assert.Equal(t, []byte{}, dest)
	})

	t.Run("invalid base64", func(t *testing.T) {
		var dest []byte
		err := BindStringToObjectWithOptions("!!!not-base64!!!", &dest, opts)
		assert.Error(t, err)
	})
}

// A slice destination that isn't a base64-decoded []byte has no string
// representation to parse, so it must report an unhandled-type error. It used
// to fall through into the integer case instead, which meant any source that
// parsed as an integer reached v.OverflowInt on a slice value and panicked,
// while a non-numeric source happened to error out in ParseInt first.
// See: https://github.com/oapi-codegen/runtime/issues/155
func TestBindStringToObject_NonByteSliceDestination(t *testing.T) {
	const wantErr = "can not bind to destination of type: slice"

	t.Run("numeric source", func(t *testing.T) {
		var dest []string
		require.NotPanics(t, func() {
			err := BindStringToObject("123", &dest)
			require.Error(t, err)
			assert.Contains(t, err.Error(), wantErr)
		})
		assert.Nil(t, dest)
	})

	t.Run("non-numeric source", func(t *testing.T) {
		var dest []string
		err := BindStringToObject("abc", &dest)
		require.Error(t, err)
		assert.Contains(t, err.Error(), wantErr)
	})

	t.Run("integer element type", func(t *testing.T) {
		var dest []int
		require.NotPanics(t, func() {
			err := BindStringToObject("123", &dest)
			require.Error(t, err)
			assert.Contains(t, err.Error(), wantErr)
		})
	})

	t.Run("byte slice without byte format", func(t *testing.T) {
		// The base64 path is only taken for Format "byte"; without it a
		// []byte is just another slice.
		var dest []byte
		require.NotPanics(t, func() {
			err := BindStringToObject("123", &dest)
			require.Error(t, err)
			assert.Contains(t, err.Error(), wantErr)
		})
	})
}

// The panic above was reachable from generated server code: a nullable slice
// query parameter with the default form/explode serialization binds through
// the primitive path, and the nullable wrapper then binds the raw value into a
// fresh slice. `?p=123` panicked while `?p=abc` returned a binding error; both
// must now be errors.
// See: https://github.com/oapi-codegen/runtime/issues/155
func TestBindQueryParameter_NullableSliceDestination(t *testing.T) {
	for _, src := range []string{"123", "abc"} {
		t.Run(src, func(t *testing.T) {
			require.NotPanics(t, func() {
				var dest nullable.Nullable[[]string]
				err := BindQueryParameter("form", true, true, "p", url.Values{"p": {src}}, &dest)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "can not bind to destination of type: slice")
			})

			require.NotPanics(t, func() {
				var dest nullable.Nullable[[]string]
				err := BindRawQueryParameter("form", true, true, "p", "p="+src, &dest)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "can not bind to destination of type: slice")
			})
		})
	}
}

// selfParsingSlice, selfParsingStruct and selfParsingMap implement Binder on
// kinds which have no built in conversion from a string.
type selfParsingSlice []string

func (s *selfParsingSlice) Bind(src string) error {
	if src == "" {
		return errors.New("empty")
	}
	*s = strings.Split(src, "|")
	return nil
}

type selfParsingStruct struct{ Value string }

func (s *selfParsingStruct) Bind(src string) error {
	s.Value = src
	return nil
}

type selfParsingMap map[string]string

func (m *selfParsingMap) Bind(src string) error {
	k, v, _ := strings.Cut(src, "=")
	*m = selfParsingMap{k: v}
	return nil
}

// textSlice and textStruct implement encoding.TextUnmarshaler instead.
type textSlice []string

func (s *textSlice) UnmarshalText(text []byte) error {
	*s = strings.Fields(string(text))
	return nil
}

type textStruct struct{ X, Y string }

func (s *textStruct) UnmarshalText(text []byte) error {
	x, y, ok := strings.Cut(string(text), ":")
	if !ok {
		return errors.New("expected x:y")
	}
	s.X, s.Y = x, y
	return nil
}

// Types which no built in conversion can handle parse themselves, through
// Binder or encoding.TextUnmarshaler, whether or not they sit behind the
// pointer of an optional field.
// https://github.com/oapi-codegen/runtime/issues/196
func TestBindStringToObject_SelfParsingTypes(t *testing.T) {
	t.Run("Binder slice", func(t *testing.T) {
		var dst selfParsingSlice
		require.NoError(t, BindStringToObject("a|b", &dst))
		assert.Equal(t, selfParsingSlice{"a", "b"}, dst)

		var opt *selfParsingSlice
		require.NoError(t, BindStringToObject("a|b", &opt))
		require.NotNil(t, opt)
		assert.Equal(t, selfParsingSlice{"a", "b"}, *opt)
	})

	t.Run("Binder struct", func(t *testing.T) {
		var dst selfParsingStruct
		require.NoError(t, BindStringToObject("x", &dst))
		assert.Equal(t, selfParsingStruct{Value: "x"}, dst)

		var opt *selfParsingStruct
		require.NoError(t, BindStringToObject("x", &opt))
		require.NotNil(t, opt)
		assert.Equal(t, selfParsingStruct{Value: "x"}, *opt)
	})

	t.Run("Binder map", func(t *testing.T) {
		var opt *selfParsingMap
		require.NoError(t, BindStringToObject("k=v", &opt))
		require.NotNil(t, opt)
		assert.Equal(t, selfParsingMap{"k": "v"}, *opt)
	})

	t.Run("Binder error is returned", func(t *testing.T) {
		var opt *selfParsingSlice
		assert.EqualError(t, BindStringToObject("", &opt), "empty")
	})

	t.Run("TextUnmarshaler slice", func(t *testing.T) {
		var dst textSlice
		require.NoError(t, BindStringToObject("a b", &dst))
		assert.Equal(t, textSlice{"a", "b"}, dst)

		var opt *textSlice
		require.NoError(t, BindStringToObject("a b", &opt))
		require.NotNil(t, opt)
		assert.Equal(t, textSlice{"a", "b"}, *opt)
	})

	t.Run("TextUnmarshaler struct", func(t *testing.T) {
		var dst textStruct
		require.NoError(t, BindStringToObject("1:2", &dst))
		assert.Equal(t, textStruct{X: "1", Y: "2"}, dst)

		var opt *textStruct
		require.NoError(t, BindStringToObject("1:2", &opt))
		require.NotNil(t, opt)
		assert.Equal(t, textStruct{X: "1", Y: "2"}, *opt)
	})

	t.Run("TextUnmarshaler error is returned", func(t *testing.T) {
		var dst textStruct
		err := BindStringToObject("nope", &dst)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected x:y")
	})

	t.Run("standard library types", func(t *testing.T) {
		var ip net.IP
		require.NoError(t, BindStringToObject("192.0.2.1", &ip))
		assert.True(t, net.IPv4(192, 0, 2, 1).Equal(ip))

		var optIP *net.IP
		require.NoError(t, BindStringToObject("2001:db8::1", &optIP))
		require.NotNil(t, optIP)
		assert.True(t, net.ParseIP("2001:db8::1").Equal(*optIP))

		var addr netip.Addr
		require.NoError(t, BindStringToObject("192.0.2.1", &addr))
		assert.Equal(t, netip.MustParseAddr("192.0.2.1"), addr)

		var n big.Int
		require.NoError(t, BindStringToObject("123456789012345678901234567890", &n))
		assert.Equal(t, "123456789012345678901234567890", n.String())
	})

	t.Run("optional TextUnmarshaler array", func(t *testing.T) {
		var opt *uuid.UUID
		require.NoError(t, BindStringToObject("5b3a3bc5-1c9e-4e55-9d0e-1c7d5a3b6f10", &opt))
		require.NotNil(t, opt)
		assert.Equal(t, uuid.MustParse("5b3a3bc5-1c9e-4e55-9d0e-1c7d5a3b6f10"), *opt)
	})
}

// The self parsing fallback is a last resort, so everything which had a built
// in conversion keeps it, even where the type implements TextUnmarshaler too.
func TestBindStringToObject_SelfParsingIsLastResort(t *testing.T) {
	t.Run("time.Time still accepts a bare date", func(t *testing.T) {
		// time.Time.UnmarshalText only accepts RFC 3339.
		want := time.Date(2020, 11, 6, 0, 0, 0, 0, time.UTC)

		var dst time.Time
		require.NoError(t, BindStringToObject("2020-11-06", &dst))
		assert.True(t, want.Equal(dst))

		var opt *time.Time
		require.NoError(t, BindStringToObject("2020-11-06", &opt))
		require.NotNil(t, opt)
		assert.True(t, want.Equal(*opt))
	})

	t.Run("time.Time still ignores an empty string", func(t *testing.T) {
		var dst time.Time
		require.NoError(t, BindStringToObject("", &dst))
		assert.True(t, dst.IsZero())
	})

	t.Run("byte format still decodes base64 into a TextUnmarshaler", func(t *testing.T) {
		var ip net.IP
		err := BindStringToObjectWithOptions("wAACAQ==", &ip, BindStringToObjectOptions{Format: "byte"})
		require.NoError(t, err)
		assert.Equal(t, net.IP{192, 0, 2, 1}, ip)
	})

	t.Run("plain slices are still an error", func(t *testing.T) {
		var dst []string
		err := BindStringToObject("abc", &dst)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "can not bind to destination of type: slice")
	})

	t.Run("plain structs are still an error", func(t *testing.T) {
		var dst struct{ A string }
		err := BindStringToObject("abc", &dst)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "can not bind to destination of type: struct")
	})
}
