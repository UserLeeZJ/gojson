package stream

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestJSONTokenizer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []JSONTokenType
	}{
		{
			name:     "空对象",
			input:    "{}",
			expected: []JSONTokenType{TokenObjectStart, TokenObjectEnd, TokenEOF},
		},
		{
			name:     "简单对象",
			input:    `{"name":"John","age":30}`,
			expected: []JSONTokenType{TokenObjectStart, TokenPropertyName, TokenString, TokenPropertyName, TokenNumber, TokenObjectEnd, TokenEOF},
		},
		{
			name:     "嵌套对象",
			input:    `{"person":{"name":"John","age":30}}`,
			expected: []JSONTokenType{TokenObjectStart, TokenPropertyName, TokenObjectStart, TokenPropertyName, TokenString, TokenPropertyName, TokenNumber, TokenObjectEnd, TokenObjectEnd, TokenEOF},
		},
		{
			name:     "空数组",
			input:    "[]",
			expected: []JSONTokenType{TokenArrayStart, TokenArrayEnd, TokenEOF},
		},
		{
			name:     "简单数组",
			input:    `[1,2,3]`,
			expected: []JSONTokenType{TokenArrayStart, TokenNumber, TokenNumber, TokenNumber, TokenArrayEnd, TokenEOF},
		},
		{
			name:     "对象数组",
			input:    `[{"name":"John"},{"name":"Jane"}]`,
			expected: []JSONTokenType{TokenArrayStart, TokenObjectStart, TokenPropertyName, TokenString, TokenObjectEnd, TokenObjectStart, TokenPropertyName, TokenString, TokenObjectEnd, TokenArrayEnd, TokenEOF},
		},
		{
			name:     "布尔值和null",
			input:    `{"active":true,"deleted":false,"data":null}`,
			expected: []JSONTokenType{TokenObjectStart, TokenPropertyName, TokenBoolean, TokenPropertyName, TokenBoolean, TokenPropertyName, TokenNull, TokenObjectEnd, TokenEOF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewJSONTokenizer(strings.NewReader(tt.input))
			
			for i, expected := range tt.expected {
				token := tokenizer.Next()
				
				if token.Type != expected {
					t.Errorf("令牌 %d: 期望类型 %v, 实际类型 %v", i, expected, token.Type)
				}
				
				if token.Type == TokenError && token.Error != nil {
					t.Errorf("令牌 %d: 错误 %v", i, token.Error)
				}
			}
		})
	}
}

func TestJSONGenerator(t *testing.T) {
	tests := []struct {
		name     string
		generate func(*JSONGenerator) error
		expected string
	}{
		{
			name: "空对象",
			generate: func(g *JSONGenerator) error {
				if err := g.BeginObject(); err != nil {
					return err
				}
				return g.EndObject()
			},
			expected: "{}",
		},
		{
			name: "简单对象",
			generate: func(g *JSONGenerator) error {
				if err := g.BeginObject(); err != nil {
					return err
				}
				if err := g.WriteProperty("name"); err != nil {
					return err
				}
				if err := g.WriteString("John"); err != nil {
					return err
				}
				if err := g.WriteProperty("age"); err != nil {
					return err
				}
				if err := g.WriteNumber(30); err != nil {
					return err
				}
				return g.EndObject()
			},
			expected: `{"name":"John","age":30}`,
		},
		{
			name: "嵌套对象",
			generate: func(g *JSONGenerator) error {
				if err := g.BeginObject(); err != nil {
					return err
				}
				if err := g.WriteProperty("person"); err != nil {
					return err
				}
				if err := g.BeginObject(); err != nil {
					return err
				}
				if err := g.WriteProperty("name"); err != nil {
					return err
				}
				if err := g.WriteString("John"); err != nil {
					return err
				}
				if err := g.WriteProperty("age"); err != nil {
					return err
				}
				if err := g.WriteNumber(30); err != nil {
					return err
				}
				if err := g.EndObject(); err != nil {
					return err
				}
				return g.EndObject()
			},
			expected: `{"person":{"name":"John","age":30}}`,
		},
		{
			name: "空数组",
			generate: func(g *JSONGenerator) error {
				if err := g.BeginArray(); err != nil {
					return err
				}
				return g.EndArray()
			},
			expected: "[]",
		},
		{
			name: "简单数组",
			generate: func(g *JSONGenerator) error {
				if err := g.BeginArray(); err != nil {
					return err
				}
				if err := g.WriteNumber(1); err != nil {
					return err
				}
				if err := g.WriteNumber(2); err != nil {
					return err
				}
				if err := g.WriteNumber(3); err != nil {
					return err
				}
				return g.EndArray()
			},
			expected: "[1,2,3]",
		},
		{
			name: "布尔值和null",
			generate: func(g *JSONGenerator) error {
				if err := g.BeginObject(); err != nil {
					return err
				}
				if err := g.WriteProperty("active"); err != nil {
					return err
				}
				if err := g.WriteBoolean(true); err != nil {
					return err
				}
				if err := g.WriteProperty("deleted"); err != nil {
					return err
				}
				if err := g.WriteBoolean(false); err != nil {
					return err
				}
				if err := g.WriteProperty("data"); err != nil {
					return err
				}
				if err := g.WriteNull(); err != nil {
					return err
				}
				return g.EndObject()
			},
			expected: `{"active":true,"deleted":false,"data":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			generator := NewJSONGenerator(&buf)
			
			err := tt.generate(generator)
			if err != nil {
				t.Fatalf("生成JSON失败: %v", err)
			}
			
			if err := generator.Flush(); err != nil {
				t.Fatalf("刷新缓冲区失败: %v", err)
			}
			
			result := buf.String()
			
			// 解析并比较JSON
			var expected, actual interface{}
			if err := json.Unmarshal([]byte(tt.expected), &expected); err != nil {
				t.Fatalf("解析期望的JSON失败: %v", err)
			}
			if err := json.Unmarshal([]byte(result), &actual); err != nil {
				t.Fatalf("解析生成的JSON失败: %v", err)
			}
			
			// 转换为JSON字符串进行比较
			expectedJSON, _ := json.Marshal(expected)
			actualJSON, _ := json.Marshal(actual)
			
			if string(expectedJSON) != string(actualJSON) {
				t.Errorf("期望: %s, 实际: %s", tt.expected, result)
			}
		})
	}
}
