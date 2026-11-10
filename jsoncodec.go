package runtime

import (
	"encoding/json"

	"github.com/oapi-codegen/runtime/internal/jsoncodec"
)

// JSONCodec encodes and decodes JSON for generated code and for this package.
type JSONCodec interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

// StdJSONCodec is the default, backed by encoding/json. Embed it to replace one method.
type StdJSONCodec struct{}

func (StdJSONCodec) Marshal(v any) ([]byte, error)      { return json.Marshal(v) }
func (StdJSONCodec) Unmarshal(data []byte, v any) error { return json.Unmarshal(data, v) }

// SetJSONCodec replaces the codec used by every generated package in the process and by
// this package. Call it during initialisation; it is not synchronised. nil restores the default.
func SetJSONCodec(c JSONCodec) { jsoncodec.Set(c) }

// JSONMarshal encodes v with the current JSONCodec.
func JSONMarshal(v any) ([]byte, error) { return jsoncodec.Marshal(v) }

// JSONUnmarshal decodes data into v with the current JSONCodec.
func JSONUnmarshal(data []byte, v any) error { return jsoncodec.Unmarshal(data, v) }
