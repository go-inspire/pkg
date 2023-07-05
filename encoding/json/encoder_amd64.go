//go:build amd64
// +build amd64

/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package json

import (
	"github.com/bytedance/sonic"
	"io"
)

// An Encoder writes JSON values to an output stream.
type encoder struct {
	inner sonic.Encoder
}

func (enc *encoder) Encode(val interface{}) (err error) {
	return enc.inner.Encode(val)
}

func NewEncoderFunc(w io.Writer) func(v interface{}) error {
	enc := encoder{
		inner: jsonAPI.NewEncoder(w),
	}
	return func(v interface{}) error {
		return enc.Encode(v)
	}
}
