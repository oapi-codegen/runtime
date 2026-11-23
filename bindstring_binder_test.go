package runtime

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oapi-codegen/runtime/types"
)

// percentBinder is an integer which does its own parsing: it tolerates a
// trailing percent sign, and rejects anything over 100.
type percentBinder int

func (p *percentBinder) Bind(src string) error {
	n, err := strconv.Atoi(strings.TrimSuffix(src, "%"))
	if err != nil {
		return err
	}
	if n > 100 {
		return errors.New("over 100")
	}
	*p = percentBinder(n)
	return nil
}

// A Binder takes precedence over the built in conversions whether or not it
// sits behind the pointer of an optional field. The optional case used to
// parse the underlying primitive without calling Bind, which skipped the
// type's validation and refused its syntax.
// https://github.com/oapi-codegen/runtime/issues/196
func TestBindStringToObject_BinderBehindOptionalPointer(t *testing.T) {
	t.Run("custom syntax is accepted", func(t *testing.T) {
		var dst percentBinder
		require.NoError(t, BindStringToObject("50%", &dst))
		assert.Equal(t, percentBinder(50), dst)

		var opt *percentBinder
		require.NoError(t, BindStringToObject("50%", &opt))
		require.NotNil(t, opt)
		assert.Equal(t, percentBinder(50), *opt)
	})

	t.Run("validation is applied", func(t *testing.T) {
		var dst percentBinder
		assert.EqualError(t, BindStringToObject("150", &dst), "over 100")

		var opt *percentBinder
		assert.EqualError(t, BindStringToObject("150", &opt), "over 100")
	})

	t.Run("string kind", func(t *testing.T) {
		var opt *CustomStringBinder
		require.NoError(t, BindStringToObject("hello", &opt))
		require.NotNil(t, opt)
		assert.Equal(t, CustomStringBinder("CUSTOM:hello"), *opt)
	})

	t.Run("types.Date binds as before", func(t *testing.T) {
		var opt *types.Date
		require.NoError(t, BindStringToObject("2020-11-06", &opt))
		require.NotNil(t, opt)
		assert.True(t, time.Date(2020, 11, 6, 0, 0, 0, 0, time.UTC).Equal(opt.Time))

		var empty *types.Date
		require.NoError(t, BindStringToObject("", &empty))
		require.NotNil(t, empty)
		assert.True(t, empty.IsZero())

		var bad *types.Date
		require.Error(t, BindStringToObject("not-a-date", &bad))
	})

	t.Run("types without a Binder are unaffected", func(t *testing.T) {
		var opt *int
		require.NoError(t, BindStringToObject("150", &opt))
		require.NotNil(t, opt)
		assert.Equal(t, 150, *opt)

		var bad *int
		require.Error(t, BindStringToObject("50%", &bad))
	})

	t.Run("optional property of an exploded form object", func(t *testing.T) {
		type Filters struct {
			Required percentBinder  `json:"required,omitempty"`
			Optional *percentBinder `json:"optional,omitempty"`
		}

		var dst *Filters
		values := url.Values{"required": {"25%"}, "optional": {"50%"}}
		require.NoError(t, BindQueryParameter("form", true, false, "filters", values, &dst))
		require.NotNil(t, dst)
		assert.Equal(t, percentBinder(25), dst.Required)
		require.NotNil(t, dst.Optional)
		assert.Equal(t, percentBinder(50), *dst.Optional)

		dst = nil
		values = url.Values{"optional": {"150"}}
		err := BindQueryParameter("form", true, false, "filters", values, &dst)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "over 100")
	})
}
