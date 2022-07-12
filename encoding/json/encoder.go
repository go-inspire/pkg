//go:build !amd64
// +build !amd64

/*
 * Copyright 2022 The tpool Authors. All rights reserved.
 * Use of this source code is governed by a GNU-style
 * license that can be found in the LICENSE file.
 */

package json

import (
	gojson "github.com/goccy/go-json"
	"io"
)

// An Encoder writes JSON values to an output stream.
type encoder struct {
	inner *gojson.Encoder
}

func (enc *encoder) Encode(val interface{}) (err error) {
	return enc.inner.Encode(val)
}

func NewEncoderFunc(w io.Writer) func(v interface{}) error {
	enc := encoder{
		inner: gojson.NewEncoder(w),
	}
	return func(v interface{}) error {
		return enc.Encode(v)
	}
}
