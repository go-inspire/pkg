package json

import (
	"encoding/json"
	"github.com/bytedance/sonic"
	gojson "github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
	"testing"
)

var work = []string{
	"0xe5e25e241b61e31f8e87bdbade565315fce55f40e66087f1d513cc7487dc6aa6",
	"0x5087743bd711255a8ea71ad43c8f377492ca5073be1f5fff40c076f3db7dbeb3",
	"0x000001ad7f29abcaf485787a6520ec08d23699194119a5c37387b71906614310",
}

type Data struct {
	Id      int64       `json:"id"`
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   interface{} `json:"error,omitempty"`
}

var dataJsonString = []byte(`{"id":0,"jsonrpc":"2.0","result":["0xe5e25e241b61e31f8e87bdbade565315fce55f40e66087f1d513cc7487dc6aa6","0x5087743bd711255a8ea71ad43c8f377492ca5073be1f5fff40c076f3db7dbeb3","0x000001ad7f29abcaf485787a6520ec08d23699194119a5c37387b71906614310"]}`)

func Benchmark_JSON_STD_Marshal(b *testing.B) {
	data := Data{Id: 0, Version: "2.0", Error: nil, Result: work}
	json.Marshal(&data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(&data)
	}
}

func Benchmark_JSON_STD_Unmarshal(b *testing.B) {
	var data Data
	json.Unmarshal(dataJsonString, &data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Unmarshal(dataJsonString, &data)
	}
}

func Benchmark_JSON_jsoniter_Marshal(b *testing.B) {
	data := Data{Id: 0, Version: "2.0", Error: nil, Result: work}
	jsoniter.Marshal(&data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jsoniter.Marshal(&data)
	}
}

func Benchmark_JSON_jsoniter_Unmarshall(b *testing.B) {
	var data Data
	jsoniter.Unmarshal(dataJsonString, &data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jsoniter.Unmarshal(dataJsonString, &data)
	}
}

func Benchmark_JSON_gojson_Marshal(b *testing.B) {
	data := Data{Id: 0, Version: "2.0", Error: nil, Result: work}
	gojson.Marshal(&data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gojson.Marshal(&data)
	}
}

func Benchmark_JSON_gojson_Unmarshall(b *testing.B) {
	var data Data
	gojson.Unmarshal(dataJsonString, &data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gojson.Unmarshal(dataJsonString, &data)
	}
}

func Benchmark_JSON_sonic_Marshal(b *testing.B) {
	data := Data{Id: 0, Version: "2.0", Error: nil, Result: work}
	sonic.Marshal(&data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sonic.Marshal(&data)
	}
}

func Benchmark_JSON_sonic_Unmarshall(b *testing.B) {
	var data Data
	sonic.Unmarshal(dataJsonString, &data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sonic.Unmarshal(dataJsonString, &data)
	}
}
