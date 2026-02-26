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
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oapi-codegen/runtime/types"
)

// TestBindStyledParameter_ByteSlice tests that BindStyledParameterWithOptions
// correctly handles *[]byte destinations by base64-decoding the parameter value,
// rather than treating []byte as a generic slice and splitting on commas.
// See: https://github.com/oapi-codegen/runtime/issues/97
func TestBindStyledParameter_ByteSlice(t *testing.T) {
	expected := []byte("test")

	tests := []struct {
		name    string
		style   string
		explode bool
		value   string
	}{
		{"simple/no-explode", "simple", false, "dGVzdA=="},
		{"simple/explode", "simple", true, "dGVzdA=="},
		{"label/no-explode", "label", false, ".dGVzdA=="},
		{"label/explode", "label", true, ".dGVzdA=="},
		{"matrix/no-explode", "matrix", false, ";data=dGVzdA=="},
		{"matrix/explode", "matrix", true, ";data=dGVzdA=="},
		{"form/no-explode", "form", false, "data=dGVzdA=="},
		{"form/explode", "form", true, "data=dGVzdA=="},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var dest []byte
			err := BindStyledParameterWithOptions(tc.style, "data", tc.value, &dest, BindStyledParameterOptions{
				ParamLocation: ParamLocationUndefined,
				Explode:       tc.explode,
				Required:      true,
				Type:          "string",
				Format:        "byte",
			})
			require.NoError(t, err)
			assert.Equal(t, expected, dest)
		})
	}
}

// TestBindQueryParameter_ByteSlice tests that BindQueryParameter correctly
// handles *[]byte destinations by base64-decoding the query parameter value.
// See: https://github.com/oapi-codegen/runtime/issues/97
func TestBindQueryParameter_ByteSlice(t *testing.T) {
	expected := []byte("test")

	opts := BindQueryParameterOptions{Type: "string", Format: "byte"}

	t.Run("form/explode/required", func(t *testing.T) {
		var dest []byte
		queryParams := url.Values{"data": {"dGVzdA=="}}
		err := BindQueryParameterWithOptions("form", true, true, "data", queryParams, &dest, opts)
		require.NoError(t, err)
		assert.Equal(t, expected, dest)
	})

	t.Run("form/no-explode/required", func(t *testing.T) {
		var dest []byte
		queryParams := url.Values{"data": {"dGVzdA=="}}
		err := BindQueryParameterWithOptions("form", false, true, "data", queryParams, &dest, opts)
		require.NoError(t, err)
		assert.Equal(t, expected, dest)
	})

	t.Run("form/explode/optional/present", func(t *testing.T) {
		var dest *[]byte
		queryParams := url.Values{"data": {"dGVzdA=="}}
		err := BindQueryParameterWithOptions("form", true, false, "data", queryParams, &dest, opts)
		require.NoError(t, err)
		require.NotNil(t, dest)
		assert.Equal(t, expected, *dest)
	})

	t.Run("form/explode/optional/absent", func(t *testing.T) {
		var dest *[]byte
		queryParams := url.Values{}
		err := BindQueryParameterWithOptions("form", true, false, "data", queryParams, &dest, opts)
		require.NoError(t, err)
		assert.Nil(t, dest)
	})

	t.Run("form/explode/optional/empty", func(t *testing.T) {
		var dest []byte
		queryParams := url.Values{"data": {""}}
		err := BindQueryParameterWithOptions("form", true, false, "data", queryParams, &dest, opts)
		require.NoError(t, err)
		assert.Equal(t, []byte{}, dest)
	})
}

// MockBinder is just an independent version of Binder that has the Bind implemented
type MockBinder struct {
	time.Time
}

func (d *MockBinder) Bind(src string) error {
	// Don't fail on empty string.
	if src == "" {
		return nil
	}
	parsedTime, err := time.Parse(types.DateFormat, src)
	if err != nil {
		return fmt.Errorf("error parsing '%s' as date: %s", src, err)
	}
	d.Time = parsedTime
	return nil
}

func (d MockBinder) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(types.DateFormat))
}

func (d *MockBinder) UnmarshalJSON(data []byte) error {
	var dateStr string
	err := json.Unmarshal(data, &dateStr)
	if err != nil {
		return err
	}
	parsed, err := time.Parse(types.DateFormat, dateStr)
	if err != nil {
		return err
	}
	d.Time = parsed
	return nil
}

// EmbeddedMockBinder has an embedded MockBinder and so keeps the Binder Method from MockBinder.
type EmbeddedMockBinder struct {
	MockBinder
}

// AnotherMockBinder is an entirely new type we have to create a bind method with it to implement Binder as well.
type AnotherMockBinder MockBinder

func (b *AnotherMockBinder) Bind(src string) error {
	// Don't fail on empty string.
	if src == "" {
		return nil
	}
	parsedTime, err := time.Parse(types.DateFormat, src)
	if err != nil {
		return fmt.Errorf("error parsing '%s' as date: %s", src, err)
	}
	b.Time = parsedTime
	return nil
}

func TestSplitParameter(t *testing.T) {
	// Please see the parameter serialization docs to understand these test
	// cases

	expectedPrimitive := []string{"5"}
	expectedArray := []string{"3", "4", "5"}
	expectedObject := []string{"role", "admin", "firstName", "Alex"}
	expectedExplodedObject := []string{"role=admin", "firstName=Alex"}

	var result []string
	var err error
	//  ------------------------ simple style ---------------------------------
	result, err = splitStyledParameter("simple",
		false,
		false,
		"id",
		"5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedPrimitive, result)

	result, err = splitStyledParameter("simple",
		false,
		false,
		"id",
		"3,4,5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedArray, result)

	result, err = splitStyledParameter("simple",
		false,
		true,
		"id",
		"role,admin,firstName,Alex")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedObject, result)

	result, err = splitStyledParameter("simple",
		true,
		false,
		"id",
		"5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedPrimitive, result)

	result, err = splitStyledParameter("simple",
		true,
		false,
		"id",
		"3,4,5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedArray, result)

	result, err = splitStyledParameter("simple",
		true,
		true,
		"id",
		"role=admin,firstName=Alex")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedExplodedObject, result)

	//  ------------------------ label style ---------------------------------
	result, err = splitStyledParameter("label",
		false,
		false,
		"id",
		".5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedPrimitive, result)

	result, err = splitStyledParameter("label",
		false,
		false,
		"id",
		".3,4,5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedArray, result)

	result, err = splitStyledParameter("label",
		false,
		true,
		"id",
		".role,admin,firstName,Alex")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedObject, result)

	result, err = splitStyledParameter("label",
		true,
		false,
		"id",
		".5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedPrimitive, result)

	result, err = splitStyledParameter("label",
		true,
		false,
		"id",
		".3.4.5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedArray, result)

	result, err = splitStyledParameter("label",
		true,
		true,
		"id",
		".role=admin.firstName=Alex")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedExplodedObject, result)

	//  ------------------------ matrix style ---------------------------------
	result, err = splitStyledParameter("matrix",
		false,
		false,
		"id",
		";id=5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedPrimitive, result)

	result, err = splitStyledParameter("matrix",
		false,
		false,
		"id",
		";id=3,4,5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedArray, result)

	result, err = splitStyledParameter("matrix",
		false,
		true,
		"id",
		";id=role,admin,firstName,Alex")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedObject, result)

	result, err = splitStyledParameter("matrix",
		true,
		false,
		"id",
		";id=5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedPrimitive, result)

	result, err = splitStyledParameter("matrix",
		true,
		false,
		"id",
		";id=3;id=4;id=5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedArray, result)

	result, err = splitStyledParameter("matrix",
		true,
		true,
		"id",
		";role=admin;firstName=Alex")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedExplodedObject, result)

	// ------------------------ form style ---------------------------------
	result, err = splitStyledParameter("form",
		false,
		false,
		"id",
		"id=5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedPrimitive, result)

	result, err = splitStyledParameter("form",
		false,
		false,
		"id",
		"id=3,4,5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedArray, result)

	result, err = splitStyledParameter("form",
		false,
		true,
		"id",
		"id=role,admin,firstName,Alex")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedObject, result)

	result, err = splitStyledParameter("form",
		true,
		false,
		"id",
		"id=5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedPrimitive, result)

	result, err = splitStyledParameter("form",
		true,
		false,
		"id",
		"id=3&id=4&id=5")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedArray, result)

	result, err = splitStyledParameter("form",
		true,
		true,
		"id",
		"role=admin&firstName=Alex")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedExplodedObject, result)
}

func TestBindQueryParameter(t *testing.T) {
	t.Run("deepObject", func(t *testing.T) {
		type Object struct {
			Count int `json:"count"`
		}
		type Nested struct {
			Object  Object   `json:"object"`
			Objects []Object `json:"objects"`
		}
		type ID struct {
			FirstName *string     `json:"firstName"`
			LastName  *string     `json:"lastName"`
			Role      string      `json:"role"`
			Birthday  *types.Date `json:"birthday"`
			Married   *MockBinder `json:"married"`
			Nested    Nested      `json:"nested"`
		}

		expectedName := "Alex"
		expectedDeepObject := &ID{
			FirstName: &expectedName,
			Role:      "admin",
			Birthday:  &types.Date{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			Married:   &MockBinder{time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC)},
			Nested: Nested{
				Object:  Object{Count: 123},
				Objects: []Object{{Count: 1}, {Count: 2}},
			},
		}

		actual := new(ID)
		paramName := "id"
		queryParams := url.Values{
			"id[firstName]":                 {"Alex"},
			"id[role]":                      {"admin"},
			"foo":                           {"bar"},
			"id[birthday]":                  {"2020-01-01"},
			"id[married]":                   {"2020-02-02"},
			"id[nested][object][count]":     {"123"},
			"id[nested][objects][0][count]": {"1"},
			"id[nested][objects][1][count]": {"2"},
		}

		err := BindQueryParameter("deepObject", true, false, paramName, queryParams, &actual)
		assert.NoError(t, err)
		assert.Equal(t, expectedDeepObject, actual)
	})

	t.Run("form", func(t *testing.T) {
		expected := &MockBinder{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}
		birthday := &MockBinder{}
		queryParams := url.Values{
			"birthday": {"2020-01-01"},
		}
		err := BindQueryParameter("form", true, false, "birthday", queryParams, &birthday)
		assert.NoError(t, err)
		assert.Equal(t, expected, birthday)
	})

	t.Run("optional", func(t *testing.T) {
		queryParams := url.Values{
			"time":   {"2020-12-09T16:09:53+00:00"},
			"number": {"100"},
			"text":   {"loremipsum"},
		}
		// An optional time will be a pointer to a time in a parameter object
		var optionalTime *time.Time
		err := BindQueryParameter("form", true, false, "notfound", queryParams, &optionalTime)
		require.NoError(t, err)
		assert.Nil(t, optionalTime)

		var optionalNumber *int
		err = BindQueryParameter("form", true, false, "notfound", queryParams, &optionalNumber)
		require.NoError(t, err)
		assert.Nil(t, optionalNumber)

		var optionalNonPointerText = ""
		err = BindQueryParameter("form", true, false, "notfound", queryParams, &optionalNonPointerText)
		require.NoError(t, err)
		assert.Zero(t, optionalNonPointerText)

		err = BindQueryParameter("form", true, false, "text", queryParams, &optionalNonPointerText)
		require.NoError(t, err)
		assert.Equal(t, "loremipsum", optionalNonPointerText)

		// If we require values, we require errors when they're not present.
		err = BindQueryParameter("form", true, true, "notfound", queryParams, &optionalTime)
		assert.Error(t, err)
		err = BindQueryParameter("form", true, true, "notfound", queryParams, &optionalNumber)
		assert.Error(t, err)

		var optionalPointerText *string
		err = BindQueryParameter("form", true, false, "notfound", queryParams, &optionalPointerText)
		require.NoError(t, err)
		assert.Zero(t, optionalPointerText)

		err = BindQueryParameter("form", true, false, "text", queryParams, &optionalPointerText)
		require.NoError(t, err)
		require.NotNil(t, optionalPointerText)
		assert.Equal(t, "loremipsum", *optionalPointerText)
	})
}

func TestBindParameterViaAlias(t *testing.T) {
	// We don't need to check every parameter format type here, since the binding
	// code is identical irrespective of parameter type, buy we do want to test
	// a bunch of types.
	type AString string
	type Aint int
	type Afloat float64
	type Atime = time.Time
	type Adate = MockBinder

	type AliasTortureTest struct {
		S  AString  `json:"s"`
		Sp *AString `json:"sp,omitempty"`
		I  Aint     `json:"i"`
		Ip *Aint    `json:"ip,omitempty"`
		F  Afloat   `json:"f"`
		Fp *Afloat  `json:"fp,omitempty"`
		T  Atime    `json:"t"`
		Tp *Atime   `json:"tp,omitempty"`
		D  Adate    `json:"d"`
		Dp *Adate   `json:"dp,omitempty"`
	}

	now := time.Now().UTC()
	later := now.Add(time.Hour)

	queryParams := url.Values{
		"alias[s]":  {"str"},
		"alias[sp]": {"strp"},
		"alias[i]":  {"1"},
		"alias[ip]": {"2"},
		"alias[f]":  {"3.5"},
		"alias[fp]": {"4.5"},
		"alias[t]":  {now.Format(time.RFC3339Nano)},
		"alias[tp]": {later.Format(time.RFC3339Nano)},
		"alias[d]":  {"2020-11-06"},
		"alias[dp]": {"2020-11-07"},
	}

	dst := new(AliasTortureTest)

	err := BindQueryParameter("deepObject", true, false, "alias", queryParams, &dst)
	require.NoError(t, err)

	var sp AString = "strp"
	var ip Aint = 2
	var fp Afloat = 4.5
	dp := Adate{Time: time.Date(2020, 11, 7, 0, 0, 0, 0, time.UTC)}

	expected := AliasTortureTest{
		S:  "str",
		Sp: &sp,
		I:  1,
		Ip: &ip,
		F:  3.5,
		Fp: &fp,
		T:  now,
		Tp: &later,
		D:  Adate{Time: time.Date(2020, 11, 6, 0, 0, 0, 0, time.UTC)},
		Dp: &dp,
	}

	// Compare field by field, makes errors easier to track.
	assert.EqualValues(t, expected.S, dst.S)
	assert.EqualValues(t, expected.Sp, dst.Sp)
	assert.EqualValues(t, expected.I, dst.I)
	assert.EqualValues(t, expected.Ip, dst.Ip)
	assert.EqualValues(t, expected.F, dst.F)
	assert.EqualValues(t, expected.Fp, dst.Fp)
	assert.EqualValues(t, expected.T, dst.T)
	assert.EqualValues(t, expected.Tp, dst.Tp)
	assert.EqualValues(t, expected.D, dst.D)
	assert.EqualValues(t, expected.Dp, dst.Dp)
}

// bindParamsToExplodedObject has to special case some types. Make sure that
// these non-object types are handled correctly. The other parts of the functionality
// are tested via more generic code above.
func TestBindParamsToExplodedObject(t *testing.T) {
	now := time.Now().UTC()
	values := url.Values{
		"time": {now.Format(time.RFC3339Nano)},
		"date": {"2020-11-06"},
	}

	var dstTime time.Time
	fieldsPresent, err := bindParamsToExplodedObject("time", values, &dstTime)
	assert.NoError(t, err)
	assert.True(t, fieldsPresent)
	assert.EqualValues(t, now, dstTime)

	type AliasedTime time.Time
	var aDstTime AliasedTime
	fieldsPresent, err = bindParamsToExplodedObject("time", values, &aDstTime)
	assert.NoError(t, err)
	assert.True(t, fieldsPresent)
	assert.EqualValues(t, now, aDstTime)

	expectedDate := MockBinder{Time: time.Date(2020, 11, 6, 0, 0, 0, 0, time.UTC)}

	var dstDate MockBinder
	fieldsPresent, err = bindParamsToExplodedObject("date", values, &dstDate)
	assert.NoError(t, err)
	assert.True(t, fieldsPresent)
	assert.EqualValues(t, expectedDate, dstDate)

	var eDstDate EmbeddedMockBinder
	fieldsPresent, err = bindParamsToExplodedObject("date", values, &eDstDate)
	assert.NoError(t, err)
	assert.True(t, fieldsPresent)
	assert.EqualValues(t, expectedDate, dstDate)

	var nTDstDate AnotherMockBinder
	fieldsPresent, err = bindParamsToExplodedObject("date", values, &nTDstDate)
	assert.NoError(t, err)
	assert.True(t, fieldsPresent)
	assert.EqualValues(t, expectedDate, nTDstDate)

	type ObjectWithOptional struct {
		Time *time.Time `json:"time,omitempty"`
	}

	var optDstTime ObjectWithOptional
	fieldsPresent, err = bindParamsToExplodedObject("explodedObject", values, &optDstTime)
	assert.NoError(t, err)
	assert.True(t, fieldsPresent)
	assert.EqualValues(t, &now, optDstTime.Time)
}

func TestFindRawQueryParam(t *testing.T) {
	tests := []struct {
		name       string
		rawQuery   string
		paramName  string
		wantValues []string
		wantFound  bool
	}{
		{
			name:       "simple value",
			rawQuery:   "color=red",
			paramName:  "color",
			wantValues: []string{"red"},
			wantFound:  true,
		},
		{
			name:       "not found",
			rawQuery:   "color=red",
			paramName:  "size",
			wantValues: nil,
			wantFound:  false,
		},
		{
			name:       "empty query",
			rawQuery:   "",
			paramName:  "color",
			wantValues: nil,
			wantFound:  false,
		},
		{
			name:       "multiple values (exploded)",
			rawQuery:   "color=red&color=blue&color=green",
			paramName:  "color",
			wantValues: []string{"red", "blue", "green"},
			wantFound:  true,
		},
		{
			name:       "comma in value stays encoded",
			rawQuery:   "color=a%2Cb",
			paramName:  "color",
			wantValues: []string{"a%2Cb"},
			wantFound:  true,
		},
		{
			name:       "empty value",
			rawQuery:   "color=",
			paramName:  "color",
			wantValues: []string{""},
			wantFound:  true,
		},
		{
			name:       "no equals sign",
			rawQuery:   "color",
			paramName:  "color",
			wantValues: []string{""},
			wantFound:  true,
		},
		{
			name:       "encoded key",
			rawQuery:   "my%20color=red",
			paramName:  "my color",
			wantValues: []string{"red"},
			wantFound:  true,
		},
		{
			name:       "mixed params",
			rawQuery:   "size=large&color=red&shape=round",
			paramName:  "color",
			wantValues: []string{"red"},
			wantFound:  true,
		},
		{
			name:       "value with special chars",
			rawQuery:   "color=red%26blue%3Dgreen",
			paramName:  "color",
			wantValues: []string{"red%26blue%3Dgreen"},
			wantFound:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, found := findRawQueryParam(tt.rawQuery, tt.paramName)
			assert.Equal(t, tt.wantFound, found)
			assert.Equal(t, tt.wantValues, values)
		})
	}
}

func TestBindRawQueryParameter(t *testing.T) {
	type TestObject struct {
		FirstName string `json:"firstName"`
		Role      string `json:"role"`
	}

	t.Run("form/explode=false", func(t *testing.T) {
		t.Run("string slice simple", func(t *testing.T) {
			var dest []string
			err := BindRawQueryParameter("form", false, true, "color", "color=red,green,blue", &dest)
			require.NoError(t, err)
			assert.Equal(t, []string{"red", "green", "blue"}, dest)
		})

		t.Run("string slice with comma in value", func(t *testing.T) {
			var dest []string
			err := BindRawQueryParameter("form", false, true, "color", "color=a,b,c%2Cd", &dest)
			require.NoError(t, err)
			assert.Equal(t, []string{"a", "b", "c,d"}, dest)
		})

		t.Run("string slice with multiple special chars", func(t *testing.T) {
			var dest []string
			err := BindRawQueryParameter("form", false, true, "color", "color=a%2Cb,c%26d,e%3Df", &dest)
			require.NoError(t, err)
			assert.Equal(t, []string{"a,b", "c&d", "e=f"}, dest)
		})

		t.Run("int slice", func(t *testing.T) {
			var dest []int
			err := BindRawQueryParameter("form", false, true, "ids", "ids=1,2,3", &dest)
			require.NoError(t, err)
			assert.Equal(t, []int{1, 2, 3}, dest)
		})

		t.Run("primitive string", func(t *testing.T) {
			var dest string
			err := BindRawQueryParameter("form", false, true, "color", "color=red", &dest)
			require.NoError(t, err)
			assert.Equal(t, "red", dest)
		})

		t.Run("primitive int", func(t *testing.T) {
			var dest int
			err := BindRawQueryParameter("form", false, true, "count", "count=42", &dest)
			require.NoError(t, err)
			assert.Equal(t, 42, dest)
		})

		t.Run("struct (object)", func(t *testing.T) {
			var dest TestObject
			err := BindRawQueryParameter("form", false, true, "id", "id=firstName,Alex,role,admin", &dest)
			require.NoError(t, err)
			assert.Equal(t, TestObject{FirstName: "Alex", Role: "admin"}, dest)
		})

		t.Run("struct with encoded comma in value", func(t *testing.T) {
			var dest TestObject
			err := BindRawQueryParameter("form", false, true, "id", "id=firstName,Alex%2CBob,role,admin", &dest)
			require.NoError(t, err)
			assert.Equal(t, TestObject{FirstName: "Alex,Bob", Role: "admin"}, dest)
		})

		t.Run("required missing", func(t *testing.T) {
			var dest []string
			err := BindRawQueryParameter("form", false, true, "color", "other=red", &dest)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "required")
		})

		t.Run("optional missing returns nil", func(t *testing.T) {
			var dest *[]string
			err := BindRawQueryParameter("form", false, false, "color", "other=red", &dest)
			require.NoError(t, err)
			assert.Nil(t, dest)
		})

		t.Run("optional present is populated", func(t *testing.T) {
			var dest *string
			err := BindRawQueryParameter("form", false, false, "color", "color=red", &dest)
			require.NoError(t, err)
			require.NotNil(t, dest)
			assert.Equal(t, "red", *dest)
		})

		t.Run("duplicate param errors", func(t *testing.T) {
			var dest []string
			err := BindRawQueryParameter("form", false, true, "color", "color=red&color=blue", &dest)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "not exploded")
		})
	})

	t.Run("form/explode=true", func(t *testing.T) {
		t.Run("string slice", func(t *testing.T) {
			var dest []string
			err := BindRawQueryParameter("form", true, true, "color", "color=red&color=green&color=blue", &dest)
			require.NoError(t, err)
			assert.Equal(t, []string{"red", "green", "blue"}, dest)
		})

		t.Run("int slice", func(t *testing.T) {
			var dest []int
			err := BindRawQueryParameter("form", true, true, "ids", "ids=1&ids=2&ids=3", &dest)
			require.NoError(t, err)
			assert.Equal(t, []int{1, 2, 3}, dest)
		})

		t.Run("primitive", func(t *testing.T) {
			var dest string
			err := BindRawQueryParameter("form", true, true, "color", "color=red", &dest)
			require.NoError(t, err)
			assert.Equal(t, "red", dest)
		})

		t.Run("struct", func(t *testing.T) {
			var dest TestObject
			err := BindRawQueryParameter("form", true, true, "id", "firstName=Alex&role=admin", &dest)
			require.NoError(t, err)
			assert.Equal(t, TestObject{FirstName: "Alex", Role: "admin"}, dest)
		})

		t.Run("required missing", func(t *testing.T) {
			var dest []string
			err := BindRawQueryParameter("form", true, true, "color", "other=red", &dest)
			assert.Error(t, err)
		})

		t.Run("optional missing", func(t *testing.T) {
			var dest *string
			err := BindRawQueryParameter("form", true, false, "color", "other=red", &dest)
			require.NoError(t, err)
			assert.Nil(t, dest)
		})
	})

	t.Run("deepObject/explode=true", func(t *testing.T) {
		type ID struct {
			FirstName *string `json:"firstName"`
			Role      string  `json:"role"`
		}
		var dest ID
		err := BindRawQueryParameter("deepObject", true, false, "id", "id%5BfirstName%5D=Alex&id%5Brole%5D=admin", &dest)
		require.NoError(t, err)
		expectedName := "Alex"
		assert.Equal(t, ID{FirstName: &expectedName, Role: "admin"}, dest)
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("deepObject explode=false", func(t *testing.T) {
			var dest TestObject
			err := BindRawQueryParameter("deepObject", false, true, "id", "id=foo", &dest)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "exploded")
		})

		t.Run("spaceDelimited", func(t *testing.T) {
			var dest []string
			err := BindRawQueryParameter("spaceDelimited", false, true, "color", "color=a%20b%20c", &dest)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "spaceDelimited")
		})

		t.Run("pipeDelimited", func(t *testing.T) {
			var dest []string
			err := BindRawQueryParameter("pipeDelimited", false, true, "color", "color=a|b|c", &dest)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "pipeDelimited")
		})

		t.Run("unknown style", func(t *testing.T) {
			var dest string
			err := BindRawQueryParameter("unknown", false, true, "color", "color=red", &dest)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid")
		})
	})
}

func TestRoundTripQueryParameter(t *testing.T) {
	type TestObject struct {
		FirstName string `json:"firstName"`
		Role      string `json:"role"`
	}

	tests := []struct {
		name      string
		style     string
		explode   bool
		paramName string
		value     interface{}
		dest      interface{} // pointer to zero value of dest type
		expected  interface{} // expected value after round-trip
	}{
		{
			name:      "form/false string slice",
			style:     "form",
			explode:   false,
			paramName: "color",
			value:     []string{"red", "green", "blue"},
			dest:      &[]string{},
			expected:  []string{"red", "green", "blue"},
		},
		{
			name:      "form/false string slice with commas",
			style:     "form",
			explode:   false,
			paramName: "color",
			value:     []string{"a,b", "c", "d,e,f"},
			dest:      &[]string{},
			expected:  []string{"a,b", "c", "d,e,f"},
		},
		{
			name:      "form/false int slice",
			style:     "form",
			explode:   false,
			paramName: "ids",
			value:     []int{1, 2, 3},
			dest:      &[]int{},
			expected:  []int{1, 2, 3},
		},
		{
			name:      "form/false primitive string",
			style:     "form",
			explode:   false,
			paramName: "color",
			value:     "red",
			dest:      new(string),
			expected:  "red",
		},
		{
			name:      "form/false struct",
			style:     "form",
			explode:   false,
			paramName: "id",
			value:     TestObject{FirstName: "Alex", Role: "admin"},
			dest:      &TestObject{},
			expected:  TestObject{FirstName: "Alex", Role: "admin"},
		},
		{
			name:      "form/true string slice",
			style:     "form",
			explode:   true,
			paramName: "color",
			value:     []string{"red", "green", "blue"},
			dest:      &[]string{},
			expected:  []string{"red", "green", "blue"},
		},
		{
			name:      "form/true int slice",
			style:     "form",
			explode:   true,
			paramName: "ids",
			value:     []int{1, 2, 3},
			dest:      &[]int{},
			expected:  []int{1, 2, 3},
		},
		{
			name:      "form/true struct",
			style:     "form",
			explode:   true,
			paramName: "id",
			value:     TestObject{FirstName: "Alex", Role: "admin"},
			dest:      &TestObject{},
			expected:  TestObject{FirstName: "Alex", Role: "admin"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize
			raw, err := StyleParamWithLocation(tt.style, tt.explode, tt.paramName, ParamLocationQuery, tt.value)
			require.NoError(t, err, "StyleParamWithLocation failed")

			// Deserialize
			err = BindRawQueryParameter(tt.style, tt.explode, true, tt.paramName, raw, tt.dest)
			require.NoError(t, err, "BindRawQueryParameter failed for raw=%q", raw)

			// Compare
			actual := reflect.ValueOf(tt.dest).Elem().Interface()
			assert.Equal(t, tt.expected, actual)
		})
	}

	t.Run("deepObject/true struct", func(t *testing.T) {
		original := TestObject{FirstName: "Alex", Role: "admin"}

		raw, err := StyleParamWithLocation("deepObject", true, "id", ParamLocationQuery, original)
		require.NoError(t, err)

		var dest TestObject
		err = BindRawQueryParameter("deepObject", true, true, "id", raw, &dest)
		require.NoError(t, err)
		assert.Equal(t, original, dest)
	})
}

func TestBindStyledParameterWithLocation(t *testing.T) {
	t.Run("bigNumber", func(t *testing.T) {
		expectedBig := big.NewInt(12345678910)
		var dstBigNumber big.Int

		err := BindStyledParameterWithOptions("simple", "id", "12345678910", &dstBigNumber, BindStyledParameterOptions{
			ParamLocation: ParamLocationUndefined,
			Explode:       false,
			Required:      false,
		})
		assert.NoError(t, err)
		assert.Equal(t, *expectedBig, dstBigNumber)
	})

	t.Run("object", func(t *testing.T) {
		type Object struct {
			Key1 string `json:"key1"`
			Key2 string `json:"key2"`
		}
		expectedObject := Object{
			Key1: "value1",
			Key2: "42",
		}
		var dstObject Object

		err := BindStyledParameterWithOptions("simple", "map", "key1,value1,key2,42", &dstObject, BindStyledParameterOptions{
			ParamLocation: ParamLocationUndefined,
			Explode:       false,
			Required:      false,
		})
		assert.NoError(t, err)
		assert.EqualValues(t, expectedObject, dstObject)
	})

	t.Run("map", func(t *testing.T) {
		expectedMap := map[string]any{
			"key1": "value1",
			"key2": "42",
		}
		var dstMap map[string]any

		err := BindStyledParameterWithOptions("simple", "map", "key1,value1,key2,42", &dstMap, BindStyledParameterOptions{
			ParamLocation: ParamLocationUndefined,
			Explode:       false,
			Required:      false,
		})
		assert.NoError(t, err)
		assert.EqualValues(t, expectedMap, dstMap)
	})
}
