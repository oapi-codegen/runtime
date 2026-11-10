package runtime

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// noEscapeCodec overrides Marshal to leave HTML unescaped, and counts the
// calls that reach either method.
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

// streamingCodec implements StreamingJSONCodec, and counts the calls that
// reach each method.
type streamingCodec struct {
	noEscapeCodec
	writes, reads *int
}

func (c streamingCodec) MarshalWrite(w io.Writer, v any) error {
	*c.writes++
	_, err := io.WriteString(w, "streamed")
	return err
}

func (c streamingCodec) UnmarshalRead(r io.Reader, v any) error {
	*c.reads++
	*v.(*string) = "streamed"
	return nil
}

var _ StreamingJSONCodec = streamingCodec{}

// errReader fails every Read with err.
type errReader struct{ err error }

func (r errReader) Read([]byte) (int, error) { return 0, r.err }

// TestSetJSONCodec mutates the process-wide codec, so it must not run in
// parallel with other tests (they run after all sequential tests finish).
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

func TestJSONCodec_Streams(t *testing.T) {
	t.Cleanup(func() { SetJSONCodec(nil) })

	const in = "<a&b>"

	// encoded is in as the default encodes a body: HTML escaped, trailing newline.
	var encoded bytes.Buffer
	require.NoError(t, json.NewEncoder(&encoded).Encode(in))

	t.Run("default is encoding/json's Encoder and Decoder", func(t *testing.T) {
		var buf bytes.Buffer
		require.NoError(t, JSONMarshalWrite(&buf, in))
		assert.Equal(t, encoded.String(), buf.String())

		var out string
		require.NoError(t, JSONUnmarshalRead(strings.NewReader(`"x"`), &out))
		assert.Equal(t, "x", out)
		assert.ErrorIs(t, JSONUnmarshalRead(strings.NewReader(""), &out), io.EOF)
	})

	t.Run("StdJSONCodec is the default too", func(t *testing.T) {
		SetJSONCodec(StdJSONCodec{})
		var buf bytes.Buffer
		require.NoError(t, JSONMarshalWrite(&buf, in))
		assert.Equal(t, encoded.String(), buf.String(), "not the buffered path, which has no newline")
	})

	t.Run("StdJSONCodec itself doesn't stream", func(t *testing.T) {
		_, ok := any(StdJSONCodec{}).(StreamingJSONCodec)
		assert.False(t, ok, "embedding it must not promote stream methods past an overridden Marshal")
	})

	t.Run("a codec without stream methods is buffered through them", func(t *testing.T) {
		var marshals, unmarshals int
		SetJSONCodec(noEscapeCodec{marshals: &marshals, unmarshals: &unmarshals})

		var buf bytes.Buffer
		require.NoError(t, JSONMarshalWrite(&buf, in))
		assert.Equal(t, `"<a&b>"`, buf.String(), "the overridden Marshal, not encoding/json's Encoder")
		assert.Equal(t, 1, marshals)

		var out string
		require.NoError(t, JSONUnmarshalRead(strings.NewReader(`"x"`), &out))
		assert.Equal(t, "x", out)
		assert.Equal(t, 1, unmarshals)
	})

	t.Run("a codec with stream methods streams", func(t *testing.T) {
		var marshals, unmarshals, writes, reads int
		SetJSONCodec(streamingCodec{
			noEscapeCodec: noEscapeCodec{marshals: &marshals, unmarshals: &unmarshals},
			writes:        &writes,
			reads:         &reads,
		})

		var buf bytes.Buffer
		require.NoError(t, JSONMarshalWrite(&buf, in))
		assert.Equal(t, "streamed", buf.String())

		var out string
		require.NoError(t, JSONUnmarshalRead(strings.NewReader(`"x"`), &out))
		assert.Equal(t, "streamed", out)

		assert.Equal(t, 1, writes)
		assert.Equal(t, 1, reads)
		assert.Equal(t, 0, marshals, "the byte methods are not used for streams")
		assert.Equal(t, 0, unmarshals)
	})
}

func TestJSONUnmarshalRead_EmptyInput(t *testing.T) {
	t.Cleanup(func() { SetJSONCodec(nil) })

	t.Run("no value is io.EOF whichever codec is in use", func(t *testing.T) {
		var unmarshals, reads int
		for _, codec := range []JSONCodec{
			nil,
			noEscapeCodec{marshals: new(int), unmarshals: &unmarshals},
			streamingCodec{noEscapeCodec: noEscapeCodec{marshals: new(int), unmarshals: &unmarshals}, writes: new(int), reads: &reads},
		} {
			SetJSONCodec(codec)
			var out string
			assert.ErrorIs(t, JSONUnmarshalRead(nil, &out), io.EOF)
			assert.ErrorIs(t, JSONUnmarshalRead(strings.NewReader(""), &out), io.EOF)
			assert.ErrorIs(t, JSONUnmarshalRead(strings.NewReader(" \r\n\t"), &out), io.EOF)
		}
		assert.Equal(t, 0, unmarshals, "the codec never saw the empty inputs")
		assert.Equal(t, 0, reads)
	})

	t.Run("a read error is returned", func(t *testing.T) {
		SetJSONCodec(nil)
		readErr := errors.New("read failed")
		var out string
		assert.ErrorIs(t, JSONUnmarshalRead(errReader{readErr}, &out), readErr)
	})

	t.Run("the input reaches the codec unchanged", func(t *testing.T) {
		SetJSONCodec(nil)
		const in = "   {x"
		var want, got map[string]any
		wantErr := json.NewDecoder(strings.NewReader(in)).Decode(&want)
		gotErr := JSONUnmarshalRead(strings.NewReader(in), &got)
		var wantSyntax, gotSyntax *json.SyntaxError
		require.ErrorAs(t, wantErr, &wantSyntax)
		require.ErrorAs(t, gotErr, &gotSyntax)
		assert.Equal(t, wantSyntax.Offset, gotSyntax.Offset, "leading white space is replayed, not dropped")
		assert.Equal(t, wantErr.Error(), gotErr.Error())

		require.NoError(t, JSONUnmarshalRead(strings.NewReader("\n\t {\"a\":1}"), &got))
		assert.Equal(t, map[string]any{"a": 1.0}, got)
	})
}

func TestBindJSONBody_Codec(t *testing.T) {
	t.Cleanup(func() { SetJSONCodec(nil) })

	newRequest := func(body string) *http.Request {
		r, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		require.NoError(t, err)
		return r
	}

	var marshals, unmarshals int
	SetJSONCodec(noEscapeCodec{marshals: &marshals, unmarshals: &unmarshals})

	t.Run("decodes through the codec", func(t *testing.T) {
		var got map[string]string
		require.NoError(t, BindJSONBody(newRequest(` {"name":"bob"}`), &got))
		assert.Equal(t, map[string]string{"name": "bob"}, got)
		assert.Equal(t, 1, unmarshals)
	})

	t.Run("an empty body is io.EOF whichever codec is in use", func(t *testing.T) {
		var got map[string]string
		assert.ErrorIs(t, BindJSONBody(newRequest(""), &got), io.EOF)
		assert.ErrorIs(t, BindJSONBody(newRequest(" \r\n\t"), &got), io.EOF)
		assert.ErrorIs(t, BindJSONBody(&http.Request{}, &got), io.EOF)
		assert.ErrorIs(t, BindJSONBody(nil, &got), io.EOF)
		assert.Equal(t, 1, unmarshals, "the codec never saw the empty bodies")
	})
}
