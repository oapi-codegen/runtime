package runtime

import (
	"io"

	"github.com/oapi-codegen/runtime/internal/jsoncodec"
)

// JSONCodec encodes and decodes JSON for this package and for the code
// oapi-codegen generates. It has two methods:
//
//	Marshal(v any) ([]byte, error)
//	Unmarshal(data []byte, v any) error
//
// A codec must produce and accept standard JSON: it selects an implementation
// and its options, such as encoding/json/v2 or an encoder that doesn't escape
// HTML, not a different format. A codec that also implements
// StreamingJSONCodec encodes and decodes request and response bodies as
// streams.
type JSONCodec = jsoncodec.Codec

// StreamingJSONCodec is optionally implemented by a JSONCodec that encodes to
// and decodes from a stream, which request and response bodies use:
//
//	MarshalWrite(w io.Writer, v any) error
//	UnmarshalRead(r io.Reader, v any) error
//
// Without it, JSONMarshalWrite and JSONUnmarshalRead buffer the whole body
// through Marshal and Unmarshal. Whether MarshalWrite ends the value with a
// newline, and whether UnmarshalRead rejects data after the value, is up to
// the implementation: encoding/json's Encoder and Decoder do and don't,
// respectively. Neither this package nor generated code relies on either.
type StreamingJSONCodec = jsoncodec.Streaming

// StdJSONCodec is encoding/json's json.Marshal and json.Unmarshal. Embed it in
// a codec to replace one method.
//
// It doesn't implement StreamingJSONCodec, so a codec that embeds it handles
// bodies through its own Marshal and Unmarshal. The default codec, which
// SetJSONCodec(nil) restores, adds encoding/json's Encoder and Decoder for
// them.
type StdJSONCodec = jsoncodec.Std

// SetJSONCodec replaces the codec used by this package and by the code
// oapi-codegen generates, process-wide. Call it during initialisation; it is
// not synchronised. nil restores the default.
func SetJSONCodec(c JSONCodec) { jsoncodec.Set(c) }

// JSONMarshal encodes v with the current JSONCodec.
func JSONMarshal(v any) ([]byte, error) { return jsoncodec.Marshal(v) }

// JSONUnmarshal decodes data into v with the current JSONCodec.
func JSONUnmarshal(data []byte, v any) error { return jsoncodec.Unmarshal(data, v) }

// JSONMarshalWrite encodes v to w with the current JSONCodec: with its
// MarshalWrite when it implements StreamingJSONCodec, otherwise with Marshal
// followed by a single Write.
func JSONMarshalWrite(w io.Writer, v any) error { return jsoncodec.MarshalWrite(w, v) }

// JSONUnmarshalRead decodes v from r with the current JSONCodec: with its
// UnmarshalRead when it implements StreamingJSONCodec, otherwise by reading r
// to its end and calling Unmarshal.
func JSONUnmarshalRead(r io.Reader, v any) error { return jsoncodec.UnmarshalRead(r, v) }
