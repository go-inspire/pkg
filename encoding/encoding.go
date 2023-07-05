/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package encoding

import (
	"github.com/go-inspire/pkg/encoding/json"
	"io"
)

// Codec defines the interface Transport uses to encode and decode messages.  Note
// that implementations of this interface must be thread safe; a Codec's
// methods can be called from concurrent goroutines.
type Codec interface {
	// Marshal returns the wire format of v.
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal parses the wire format into v.
	Unmarshal(data []byte, v interface{}) error
}

func MarshalJSON(v interface{}) ([]byte, error) {
	return json.MarshalFunc()(v)
}

func UnmarshalJSON(data []byte, v interface{}) error {
	return json.UnmarshalFunc()(data, v)
}

// EncoderFunc adapts an encoder function into Encoder
type EncoderFunc func(v interface{}) error

func NewEncoderFunc(writer io.Writer) EncoderFunc {
	return json.NewEncoderFunc(writer)
}
