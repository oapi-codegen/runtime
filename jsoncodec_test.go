package runtime

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// noEscapeCodec overrides Marshal only, and counts the calls that reach it.
type noEscapeCodec struct {
	StdJSONCodec
	marshals, unmarshals *int
}

func (c noEscapeCodec) Marshal(v any) ([]byte, error) {
	*c.marshals++
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), nil
}

func (c noEscapeCodec) Unmarshal(data []byte, v any) error {
	*c.unmarshals++
	return c.StdJSONCodec.Unmarshal(data, v)
}

func TestSetJSONCodec(t *testing.T) {
	t.Cleanup(func() { SetJSONCodec(nil) })

	const in = "<a&b>"

	// Default: encoding/json, which escapes HTML.
	escaped, err := json.Marshal(in)
	require.NoError(t, err)
	require.NotEqual(t, `"<a&b>"`, string(escaped))

	got, err := JSONMarshal(in)
	require.NoError(t, err)
	assert.Equal(t, string(escaped), string(got))

	var marshals, unmarshals int
	SetJSONCodec(noEscapeCodec{marshals: &marshals, unmarshals: &unmarshals})

	got, err = JSONMarshal(in)
	require.NoError(t, err)
	assert.Equal(t, `"<a&b>"`, string(got))
	assert.Equal(t, 1, marshals)

	var out string
	require.NoError(t, JSONUnmarshal(got, &out))
	assert.Equal(t, in, out)
	assert.Equal(t, 1, unmarshals)

	// The runtime's own encoding goes through the codec too.
	_, err = MarshalDeepObject(map[string]string{"k": "v"}, "p")
	require.NoError(t, err)
	assert.Equal(t, 2, marshals)

	// So does runtime/types, which cannot import this package.
	d := types.Date{Time: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)}
	got, err = d.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, `"2020-01-02"`, string(got))
	assert.Equal(t, 3, marshals)
	require.NoError(t, d.UnmarshalJSON(got))
	assert.Equal(t, 2, unmarshals)

	// nil restores the default.
	SetJSONCodec(nil)
	got, err = JSONMarshal(in)
	require.NoError(t, err)
	assert.Equal(t, string(escaped), string(got))
	assert.Equal(t, 3, marshals)
}
