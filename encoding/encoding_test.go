/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package encoding

import (
	"reflect"
	"strings"
	"testing"
)

type testEmbed struct {
	Level1a int `json:"a"`
	Level1b int `json:"b"`
	Level1c int `json:"c"`
}

type testMessage struct {
	Field1 string     `json:"a"`
	Field2 string     `json:"b"`
	Field3 string     `json:"c"`
	Embed  *testEmbed `json:"embed,omitempty"`
}

type mock struct {
	value int
}

const (
	Unknown = iota
	Gopher
	Zebra
)

func (a *mock) UnmarshalJSON(b []byte) error {
	var s string
	if err := UnmarshalJSON(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	default:
		a.value = Unknown
	case "gopher":
		a.value = Gopher
	case "zebra":
		a.value = Zebra
	}

	return nil
}

func (a *mock) MarshalJSON() ([]byte, error) {
	var s string
	switch a.value {
	default:
		s = "unknown"
	case Gopher:
		s = "gopher"
	case Zebra:
		s = "zebra"
	}

	return MarshalJSON(s)
}

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect string
	}{
		{
			input:  &testMessage{},
			expect: `{"a":"","b":"","c":""}`,
		},
		{
			input:  &testMessage{Field1: "a", Field2: "b", Field3: "c"},
			expect: `{"a":"a","b":"b","c":"c"}`,
		},
		{
			input:  &mock{value: Gopher},
			expect: `"gopher"`,
		},
		// 添加更多测试用例
		{
			input:  &testMessage{Embed: &testEmbed{Level1a: 1, Level1b: 2, Level1c: 3}},
			expect: `{"a":"","b":"","c":"","embed":{"a":1,"b":2,"c":3}}`,
		},
		{
			input:  map[string]interface{}{"key": "value", "number": 123},
			expect: `{"key":"value","number":123}`,
		},
		{
			input:  []int{1, 2, 3},
			expect: `[1,2,3]`,
		},
		{
			input:  nil,
			expect: `null`,
		},
	}
	for _, v := range tests {
		data, err := MarshalJSON(v.input)
		if err != nil {
			t.Errorf("marshal(%#v): %s", v.input, err)
		}
		if got, want := string(data), v.expect; strings.ReplaceAll(got, " ", "") != want {
			if strings.Contains(want, "\n") {
				t.Errorf("marshal(%#v):\nHAVE:\n%s\nWANT:\n%s", v.input, got, want)
			} else {
				t.Errorf("marshal(%#v):\nhave %#q\nwant %#q", v.input, got, want)
			}
		}
	}
}

func TestUnmarshalJSON(t *testing.T) {
	p := testMessage{}
	p4 := &mock{}
	tests := []struct {
		input  string
		expect interface{}
	}{
		{
			input:  `{"a":"","b":"","c":""}`,
			expect: &testMessage{},
		},
		{
			input:  `{"a":"a","b":"b","c":"c"}`,
			expect: &p,
		},
		{
			input:  `"zebra"`,
			expect: p4,
		},
		// 添加更多测试用例
		{
			input:  `{"a":"","b":"","c":"","embed":{"a":1,"b":2,"c":3}}`,
			expect: &testMessage{Embed: &testEmbed{}},
		},
		{
			input:  `null`,
			expect: new(interface{}),
		},
	}
	for _, v := range tests {
		want := []byte(v.input)
		err := UnmarshalJSON(want, v.expect)
		if err != nil {
			t.Errorf("marshal(%#v): %s", v.input, err)
		}
		got, err := MarshalJSON(v.expect)
		if err != nil {
			t.Errorf("marshal(%#v): %s", v.input, err)
		}
		if !reflect.DeepEqual(strings.ReplaceAll(string(got), " ", ""), strings.ReplaceAll(string(want), " ", "")) {
			t.Errorf("marshal(%#v):\nhave %#q\nwant %#q", v.input, got, want)
		}
	}
}

// 添加对 EncoderFunc 的测试
func TestEncoderFunc(t *testing.T) {
	var buf strings.Builder
	encoder := NewEncoderFunc(&buf)

	// 测试编码简单对象
	testObj := &testMessage{Field1: "test", Field2: "encoder", Field3: "func"}
	err := encoder(&testObj)
	if err != nil {
		t.Errorf("EncoderFunc error: %v", err)
	}

	// 验证输出 - 注意处理换行符
	output := strings.TrimSpace(buf.String())
	expected := `{"a":"test","b":"encoder","c":"func"}`
	if strings.ReplaceAll(output, " ", "") != strings.ReplaceAll(expected, " ", "") {
		t.Errorf("EncoderFunc output:\nhave %#q\nwant %#q", output, expected)
	}

	// 重置缓冲区并测试编码数组
	buf.Reset()
	testArray := []string{"one", "two", "three"}
	err = encoder(testArray)
	if err != nil {
		t.Errorf("EncoderFunc error: %v", err)
	}

	// 验证数组输出 - 注意处理换行符
	output = strings.TrimSpace(buf.String())
	expected = `["one","two","three"]`
	if strings.ReplaceAll(output, " ", "") != strings.ReplaceAll(expected, " ", "") {
		t.Errorf("EncoderFunc output:\nhave %#q\nwant %#q", output, expected)
	}
}

// 添加错误处理测试
func TestMarshalJSONError(t *testing.T) {
	// 创建一个会导致 JSON 编码错误的对象
	badObj := map[string]interface{}{
		"func": func() {}, // 函数不能被 JSON 编码
	}

	_, err := MarshalJSON(badObj)
	if err == nil {
		t.Error("Expected error when marshaling function, got nil")
	}
}

func TestUnmarshalJSONError(t *testing.T) {
	// 测试无效的 JSON
	invalidJSON := []byte(`{"invalid": json`)
	var result map[string]interface{}

	err := UnmarshalJSON(invalidJSON, &result)
	if err == nil {
		t.Error("Expected error when unmarshaling invalid JSON, got nil")
	}

	// 测试类型不匹配
	validJSON := []byte(`{"number": "not a number"}`)
	var intResult struct {
		Number int `json:"number"`
	}

	err = UnmarshalJSON(validJSON, &intResult)
	if err == nil {
		t.Error("Expected error when unmarshaling string to int, got nil")
	}
}

// 添加性能基准测试
func BenchmarkMarshalJSON(b *testing.B) {
	obj := &testMessage{
		Field1: "benchmark",
		Field2: "marshal",
		Field3: "json",
		Embed: &testEmbed{
			Level1a: 1,
			Level1b: 2,
			Level1c: 3,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := MarshalJSON(obj)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	data := []byte(`{"a":"benchmark","b":"unmarshal","c":"json","embed":{"a":1,"b":2,"c":3}}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var obj testMessage
		err := UnmarshalJSON(data, &obj)
		if err != nil {
			b.Fatal(err)
		}
	}
}
