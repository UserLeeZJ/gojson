package parser

import (
	"reflect"
	"testing"
)

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

func TestParseBytesToValue(t *testing.T) {
	tests := []struct {
		name      string
		jsonBytes []byte
		want      string // 期望的类型
		wantErr   bool
	}{
		{
			name:      "解析空JSON字节数组",
			jsonBytes: []byte{},
			want:      "",
			wantErr:   true,
		},
		{
			name:      "解析null",
			jsonBytes: []byte("null"),
			want:      "null",
			wantErr:   false,
		},
		{
			name:      "解析布尔值",
			jsonBytes: []byte("true"),
			want:      "boolean",
			wantErr:   false,
		},
		{
			name:      "解析数字",
			jsonBytes: []byte("123"),
			want:      "number",
			wantErr:   false,
		},
		{
			name:      "解析字符串",
			jsonBytes: []byte(`"hello"`),
			want:      "string",
			wantErr:   false,
		},
		{
			name:      "解析数组",
			jsonBytes: []byte("[1,2,3]"),
			want:      "array",
			wantErr:   false,
		},
		{
			name:      "解析对象",
			jsonBytes: []byte(`{"name":"John","age":30}`),
			want:      "object",
			wantErr:   false,
		},
		{
			name:      "解析无效JSON",
			jsonBytes: []byte("{invalid}"),
			want:      "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBytesToValue(tt.jsonBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBytesToValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Type() != tt.want {
					t.Errorf("ParseBytesToValue() got type = %v, want %v", got.Type(), tt.want)
				}
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		wantErr bool
	}{
		{
			name:    "解析空JSON字符串",
			jsonStr: "",
			wantErr: true,
		},
		{
			name:    "解析有效JSON",
			jsonStr: `{"name":"John","age":30}`,
			wantErr: false,
		},
		{
			name:    "解析无效JSON",
			jsonStr: "{invalid}",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data interface{}
			err := Parse(tt.jsonStr, &data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStringify(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		want    string
		wantErr bool
	}{
		{
			name:    "序列化null",
			value:   nil,
			want:    "null",
			wantErr: false,
		},
		{
			name:    "序列化布尔值",
			value:   true,
			want:    "true",
			wantErr: false,
		},
		{
			name:    "序列化数字",
			value:   123,
			want:    "123",
			wantErr: false,
		},
		{
			name:    "序列化字符串",
			value:   "hello",
			want:    `"hello"`,
			wantErr: false,
		},
		{
			name:    "序列化数组",
			value:   []int{1, 2, 3},
			want:    "[1,2,3]",
			wantErr: false,
		},
		{
			name:    "序列化对象",
			value:   map[string]interface{}{"name": "John", "age": 30},
			want:    `{"age":30,"name":"John"}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Stringify(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stringify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// 对于对象和数组，我们需要比较JSON内容而不是字符串表示
				if tt.name == "序列化对象" || tt.name == "序列化数组" {
					var gotObj, wantObj interface{}
					if err := Parse(got, &gotObj); err != nil {
						t.Errorf("无法解析Stringify结果: %v", err)
						return
					}
					if err := Parse(tt.want, &wantObj); err != nil {
						t.Errorf("无法解析期望结果: %v", err)
						return
					}

					// 使用深度比较
					if !reflect.DeepEqual(gotObj, wantObj) {
						t.Errorf("Stringify() = %v, want %v", got, tt.want)
					}
				} else if got != tt.want {
					t.Errorf("Stringify() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestStringifyIndent(t *testing.T) {
	obj := map[string]interface{}{
		"name": "John",
		"age":  30,
	}

	got, err := StringifyIndent(obj, "", "  ")
	if err != nil {
		t.Errorf("StringifyIndent() error = %v", err)
		return
	}

	// 解析JSON字符串进行比较，而不是比较字符串表示
	var gotObj, expectedObj map[string]interface{}

	if err := Parse(got, &gotObj); err != nil {
		t.Errorf("无法解析StringifyIndent结果: %v", err)
		return
	}

	// 创建期望的对象
	expectedObj = map[string]interface{}{
		"name": "John",
		"age":  float64(30), // JSON解析会将数字转为float64
	}

	// 使用深度比较
	if !reflect.DeepEqual(gotObj, expectedObj) {
		t.Errorf("StringifyIndent() = %v, 内容不匹配", got)
	}
}

func TestConvertToJSONValue(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  string // 期望的类型
	}{
		{
			name:  "转换null",
			value: nil,
			want:  "null",
		},
		{
			name:  "转换布尔值",
			value: true,
			want:  "boolean",
		},
		{
			name:  "转换数字",
			value: 123,
			want:  "number",
		},
		{
			name:  "转换字符串",
			value: "hello",
			want:  "string",
		},
		{
			name:  "转换数组",
			value: []interface{}{1, 2, 3},
			want:  "array",
		},
		{
			name:  "转换对象",
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
