//go:build amd64
// +build amd64

/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package json

import (
	"encoding/json"
	"github.com/bytedance/sonic"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"log"
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
	jsonAPI  = sonic.ConfigFastest
)

func init() {
	log.Info("Use sonic JSON library")
}

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
		return jsonAPI.Marshal(m)
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
		return jsonAPI.Unmarshal(data, m)
	}
}
