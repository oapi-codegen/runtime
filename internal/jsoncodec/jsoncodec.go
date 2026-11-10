// Package jsoncodec holds the process-wide JSON codec.
//
// It is internal so that both the runtime package and runtime/types can share
// one codec: runtime imports types, so types cannot import runtime.
package jsoncodec

import "encoding/json"

// Codec mirrors runtime.JSONCodec.
type Codec interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

// codec is nil while the encoding/json default is in use.
var codec Codec

// Set replaces the codec. nil restores the encoding/json default.
func Set(c Codec) {
	codec = c
}

func Marshal(v any) ([]byte, error) {
	if codec == nil {
		return json.Marshal(v)
	}
	return codec.Marshal(v)
}

func Unmarshal(data []byte, v any) error {
	if codec == nil {
		return json.Unmarshal(data, v)
	}
	return codec.Unmarshal(data, v)
}
