package gojson

import (
	"reflect"
	"testing"
)

// 测试Parse函数
func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		want    interface{}
		wantErr bool
	}{
		{
			name:    "解析空JSON字符串",
			jsonStr: "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "解析null",
			jsonStr: "null",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "解析布尔值",
			jsonStr: "true",
			want:    true,
			wantErr: false,
		},
		{
			name:    "解析数字",
			jsonStr: "123",
			want:    float64(123),
			wantErr: false,
		},
		{
			name:    "解析字符串",
			jsonStr: `"hello"`,
			want:    "hello",
			wantErr: false,
		},
		{
			name:    "解析数组",
			jsonStr: "[1,2,3]",
			want:    []interface{}{float64(1), float64(2), float64(3)},
			wantErr: false,
		},
		{
			name:    "解析对象",
			jsonStr: `{"name":"John","age":30}`,
			want:    map[string]interface{}{"name": "John", "age": float64(30)},
			wantErr: false,
		},
		{
			name:    "解析无效JSON",
			jsonStr: "{invalid}",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got interface{}
			err := Parse(tt.jsonStr, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试ParseBytes函数
func TestParseBytes(t *testing.T) {
	tests := []struct {
		name      string
		jsonBytes []byte
		want      interface{}
		wantErr   bool
	}{
		{
			name:      "解析空JSON字节数组",
			jsonBytes: []byte{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "解析null",
			jsonBytes: []byte("null"),
			want:      nil,
			wantErr:   false,
		},
		{
			name:      "解析对象",
			jsonBytes: []byte(`{"name":"John","age":30}`),
			want:      map[string]interface{}{"name": "John", "age": float64(30)},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got interface{}
			err := ParseBytes(tt.jsonBytes, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试Stringify函数
func TestStringify(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    string
		wantErr bool
	}{
		{
			name:    "序列化nil",
			input:   nil,
			want:    "null",
			wantErr: false,
		},
		{
			name:    "序列化布尔值",
			input:   true,
			want:    "true",
			wantErr: false,
		},
		{
			name:    "序列化数字",
			input:   123,
			want:    "123",
			wantErr: false,
		},
		{
			name:    "序列化字符串",
			input:   "hello",
			want:    `"hello"`,
			wantErr: false,
		},
		{
			name:    "序列化数组",
			input:   []int{1, 2, 3},
			want:    "[1,2,3]",
			wantErr: false,
		},
		{
			name:    "序列化对象",
			input:   map[string]interface{}{"name": "John", "age": 30},
			want:    `{"age":30,"name":"John"}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Stringify(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stringify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Stringify() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试StringifyBytes函数
func TestStringifyBytes(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    []byte
		wantErr bool
	}{
		{
			name:    "序列化nil",
			input:   nil,
			want:    []byte("null"),
			wantErr: false,
		},
		{
			name:    "序列化对象",
			input:   map[string]interface{}{"name": "John", "age": 30},
			want:    []byte(`{"age":30,"name":"John"}`),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringifyBytes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringifyBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringifyBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试StringifyIndent函数
func TestStringifyIndent(t *testing.T) {
	obj := map[string]interface{}{
		"name": "John",
		"age":  30,
		"address": map[string]interface{}{
			"city":  "New York",
			"state": "NY",
		},
	}

	expected := `{
  "address": {
    "city": "New York",
    "state": "NY"
  },
  "age": 30,
  "name": "John"
}`

	got, err := StringifyIndent(obj, "", "  ")
	if err != nil {
		t.Errorf("StringifyIndent() error = %v", err)
		return
	}
	if got != expected {
		t.Errorf("StringifyIndent() got = %v, want %v", got, expected)
	}
}

// 测试IsValidJSON函数
func TestIsValidJSON(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		want    bool
	}{
		{
			name:    "空字符串",
			jsonStr: "",
			want:    false,
		},
		{
			name:    "有效JSON - null",
			jsonStr: "null",
			want:    true,
		},
		{
			name:    "有效JSON - 布尔值",
			jsonStr: "true",
			want:    true,
		},
		{
			name:    "有效JSON - 数字",
			jsonStr: "123",
			want:    true,
		},
		{
			name:    "有效JSON - 字符串",
			jsonStr: `"hello"`,
			want:    true,
		},
		{
			name:    "有效JSON - 数组",
			jsonStr: "[1,2,3]",
			want:    true,
		},
		{
			name:    "有效JSON - 对象",
			jsonStr: `{"name":"John","age":30}`,
			want:    true,
		},
		{
			name:    "无效JSON",
			jsonStr: "{invalid}",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidJSON(tt.jsonStr); got != tt.want {
				t.Errorf("IsValidJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试Prettify函数
func TestPrettify(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		want    string
		wantErr bool
	}{
		{
			name:    "空字符串",
			jsonStr: "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "格式化对象",
			jsonStr: `{"name":"John","age":30}`,
			want: `{
  "name": "John",
  "age": 30
}`,
			wantErr: false,
		},
		{
			name:    "无效JSON",
			jsonStr: "{invalid}",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Prettify(tt.jsonStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Prettify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Prettify() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试Minify函数
func TestMinify(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		want    string
		wantErr bool
	}{
		{
			name:    "空字符串",
			jsonStr: "",
			want:    "",
			wantErr: true,
		},
		{
			name: "最小化格式化的JSON",
			jsonStr: `{
  "name": "John",
  "age": 30
}`,
			want:    `{"age":30,"name":"John"}`,
			wantErr: false,
		},
		{
			name:    "无效JSON",
			jsonStr: "{invalid}",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Minify(tt.jsonStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Minify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Minify() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试GetType函数
func TestGetType(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		want    string
		wantErr bool
	}{
		{
			name:    "空字符串",
			jsonStr: "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "null类型",
			jsonStr: "null",
			want:    "null",
			wantErr: false,
		},
		{
			name:    "布尔类型",
			jsonStr: "true",
			want:    "boolean",
			wantErr: false,
		},
		{
			name:    "数字类型",
			jsonStr: "123",
			want:    "number",
			wantErr: false,
		},
		{
			name:    "字符串类型",
			jsonStr: `"hello"`,
			want:    "string",
			wantErr: false,
		},
		{
			name:    "数组类型",
			jsonStr: "[1,2,3]",
			want:    "array",
			wantErr: false,
		},
		{
			name:    "对象类型",
			jsonStr: `{"name":"John"}`,
			want:    "object",
			wantErr: false,
		},
		{
			name:    "无效JSON",
			jsonStr: "{invalid}",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetType(tt.jsonStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试GetValue函数
func TestGetValue(t *testing.T) {
	jsonStr := `{
		"name": "John",
		"age": 30,
		"address": {
			"city": "New York",
			"state": "NY",
			"zip": 10001
		},
		"hobbies": ["reading", "swimming", "coding"]
	}`

	tests := []struct {
		name    string
		jsonStr string
		path    string
		want    interface{}
		wantErr bool
	}{
		{
			name:    "空字符串",
			jsonStr: "",
			path:    "name",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "获取顶级属性",
			jsonStr: jsonStr,
			path:    "name",
			want:    "John",
			wantErr: false,
		},
		{
			name:    "获取嵌套属性",
			jsonStr: jsonStr,
			path:    "address.city",
			want:    "New York",
			wantErr: false,
		},
		{
			name:    "获取数组元素",
			jsonStr: jsonStr,
			path:    "hobbies.1",
			want:    "swimming",
			wantErr: false,
		},
		{
			name:    "路径不存在",
			jsonStr: jsonStr,
			path:    "nonexistent",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "嵌套路径不存在",
			jsonStr: jsonStr,
			path:    "address.nonexistent",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "无效数组索引",
			jsonStr: jsonStr,
			path:    "hobbies.10",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetValue(tt.jsonStr, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}
