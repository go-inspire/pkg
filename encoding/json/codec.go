//go:build !amd64
// +build !amd64

/*
 * Copyright 2022 The tpool Authors. All rights reserved.
 * Use of this source code is governed by a GNU-style
 * license that can be found in the LICENSE file.
 */

package json

import (
	"encoding/json"
	gojson "github.com/goccy/go-json"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"reflect"
)

var (
	// MarshalOptions is a configurable JSON format marshaller.
	MarshalOptions = protojson.MarshalOptions{
		EmitUnpopulated: true,
	}
	// UnmarshalOptions is a configurable JSON format parser.
	UnmarshalOptions = protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}

	instance = codec{}
)

func MarshalFunc() func(v interface{}) ([]byte, error) {
	return func(v interface{}) ([]byte, error) {
		return instance.Marshal(v)
	}
}

func UnmarshalFunc() func(data []byte, v interface{}) error {
	return func(data []byte, v interface{}) error {
		return instance.Unmarshal(data, v)
	}
}

// codec is a Codec implementation with json.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case json.Marshaler:
		return m.MarshalJSON()
	case proto.Message:
		return MarshalOptions.Marshal(m)
	default:
		return gojson.Marshal(m)
	}
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	switch m := v.(type) {
	case json.Unmarshaler:
		return m.UnmarshalJSON(data)
	case proto.Message:
		return UnmarshalOptions.Unmarshal(data, m)
	default:
		rv := reflect.ValueOf(v)
		for rv := rv; rv.Kind() == reflect.Ptr; {
			if rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}
		if m, ok := reflect.Indirect(rv).Interface().(proto.Message); ok {
			return UnmarshalOptions.Unmarshal(data, m)
		}
		return gojson.Unmarshal(data, m)
	}
}
