package gojson

import (
	"reflect"
	"testing"
)

// 测试ParseToValue函数
func TestParseToValue(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		want    string // 期望的类型
		wantErr bool
	}{
		{
			name:    "解析空JSON字符串",
			jsonStr: "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "解析null",
			jsonStr: "null",
			want:    "null",
			wantErr: false,
		},
		{
			name:    "解析布尔值",
			jsonStr: "true",
			want:    "boolean",
			wantErr: false,
		},
		{
			name:    "解析数字",
			jsonStr: "123",
			want:    "number",
			wantErr: false,
		},
		{
			name:    "解析字符串",
			jsonStr: `"hello"`,
			want:    "string",
			wantErr: false,
		},
		{
			name:    "解析数组",
			jsonStr: "[1,2,3]",
			want:    "array",
			wantErr: false,
		},
		{
			name:    "解析对象",
			jsonStr: `{"name":"John","age":30}`,
			want:    "object",
			wantErr: false,
		},
		{
			name:    "解析无效JSON",
			jsonStr: "{invalid}",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseToValue(tt.jsonStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Type() != tt.want {
					t.Errorf("ParseToValue() got type = %v, want %v", got.Type(), tt.want)
				}
			}
		})
	}
}

// 测试ValueToInterface函数
func TestValueToInterface(t *testing.T) {
	tests := []struct {
		name  string
		value JSONValue
		want  interface{}
	}{
		{
			name:  "null值",
			value: NewJSONNull(),
			want:  nil,
		},
		{
			name:  "布尔值",
			value: NewJSONBool(true),
			want:  true,
		},
		{
			name:  "数字",
			value: NewJSONNumber(123),
			want:  float64(123),
		},
		{
			name:  "字符串",
			value: NewJSONString("hello"),
			want:  "hello",
		},
		{
			name: "数组",
			value: func() *JSONArray {
				arr := NewJSONArray()
				arr.AddNumber(1).AddNumber(2).AddNumber(3)
				return arr
			}(),
			want: []interface{}{float64(1), float64(2), float64(3)},
		},
		{
			name: "对象",
			value: func() *JSONObject {
				obj := NewJSONObject()
				obj.PutString("name", "John").PutNumber("age", 30)
				return obj
			}(),
			want: map[string]interface{}{"name": "John", "age": float64(30)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValueToInterface(tt.value)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValueToInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试convertToJSONValue函数
func TestConvertToJSONValue(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  string // 期望的类型
	}{
		{
			name:  "nil",
			value: nil,
			want:  "null",
		},
		{
			name:  "布尔值",
			value: true,
			want:  "boolean",
		},
		{
			name:  "数字",
			value: 123,
			want:  "number",
		},
		{
			name:  "字符串",
			value: "hello",
			want:  "string",
		},
		{
			name:  "数组",
			value: []interface{}{1, 2, 3},
			want:  "array",
		},
		{
			name:  "对象",
			value: map[string]interface{}{"name": "John", "age": 30},
			want:  "object",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertToJSONValue(tt.value)
			if got.Type() != tt.want {
				t.Errorf("convertToJSONValue() got type = %v, want %v", got.Type(), tt.want)
			}
		})
	}
}
