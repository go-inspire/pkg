/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package json

import (
	"encoding/json"
	"github.com/bytedance/sonic"
	"github.com/go-inspire/pkg/encoding/json/testdata"
	gojson "github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
	"testing"
)

type payload = testdata.LargePayload

var payloadString = testdata.LargeFixture

func init() {
	var data payload
	_ = json.Unmarshal(testdata.SmallFixture, &data)
}

func Benchmark_JSON_STD_Marshal(b *testing.B) {
	var data payload
	json.Marshal(&data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(&data)
	}
}

func Benchmark_JSON_STD_Unmarshal(b *testing.B) {
	var data payload
	json.Unmarshal(payloadString, &data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Unmarshal(payloadString, &data)
	}
}

func Benchmark_JSON_jsoniter_Marshal(b *testing.B) {
	var data payload
	jsoniter.Marshal(&data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jsoniter.Marshal(&data)
	}
}

func Benchmark_JSON_jsoniter_Unmarshall(b *testing.B) {
	var data payload
	jsoniter.Unmarshal(payloadString, &data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jsoniter.Unmarshal(payloadString, &data)
	}
}

func Benchmark_JSON_gojson_Marshal(b *testing.B) {
	var data payload
	gojson.Marshal(&data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gojson.Marshal(&data)
	}
}

func Benchmark_JSON_gojson_Unmarshall(b *testing.B) {
	var data payload
	gojson.Unmarshal(payloadString, &data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gojson.Unmarshal(payloadString, &data)
	}
}

func Benchmark_JSON_sonic_Marshal(b *testing.B) {
	var data payload
	sonic.Marshal(&data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sonic.Marshal(&data)
	}
}

func Benchmark_JSON_sonic_Unmarshall(b *testing.B) {
	var data payload
	sonic.Unmarshal(payloadString, &data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sonic.Unmarshal(payloadString, &data)
	}
}
