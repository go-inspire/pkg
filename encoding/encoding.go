/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

// Package encoding 提供了数据编码和解码的通用接口和工具函数。
package encoding

import (
	"github.com/go-inspire/pkg/encoding/json"
	"io"
)

// Codec 定义了用于编码和解码消息的接口。
// 此接口的实现必须是线程安全的，因为 Codec 的方法可能会被并发的 goroutine 调用。
type Codec interface {
	// Marshal 将对象编码为二进制格式。
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal 将二进制数据解析为对象。
	Unmarshal(data []byte, v interface{}) error
}

// MarshalJSON 将对象编码为 JSON 格式。
func MarshalJSON(v interface{}) ([]byte, error) {
	return json.MarshalFunc()(v)
}

// UnmarshalJSON 将 JSON 数据解析为对象。
func UnmarshalJSON(data []byte, v interface{}) error {
	return json.UnmarshalFunc()(data, v)
}

// EncoderFunc 是一个编码器函数类型，用于将对象编码到输出流。
type EncoderFunc func(v interface{}) error

// NewEncoderFunc 创建一个新的编码器函数，用于将对象编码到指定的输出流。
func NewEncoderFunc(writer io.Writer) EncoderFunc {
	return json.NewEncoderFunc(writer)
}
