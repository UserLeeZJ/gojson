package jsonpath

import (
	"testing"

	"github.com/UserLeeZJ/gojson/parser"
)

// contains 检查字符串切片是否包含指定字符串
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func TestJSONPath(t *testing.T) {
	// 测试数据
	jsonStr := `{
		"store": {
			"book": [
				{
					"category": "reference",
					"author": "Nigel Rees",
					"title": "Sayings of the Century",
					"price": 8.95
				},
				{
					"category": "fiction",
					"author": "Evelyn Waugh",
					"title": "Sword of Honour",
					"price": 12.99
				},
				{
					"category": "fiction",
					"author": "Herman Melville",
					"title": "Moby Dick",
					"isbn": "0-553-21311-3",
					"price": 8.99
				},
				{
					"category": "fiction",
					"author": "J. R. R. Tolkien",
					"title": "The Lord of the Rings",
					"isbn": "0-395-19395-8",
					"price": 22.99
				}
			],
			"bicycle": {
				"color": "red",
				"price": 19.95
			}
		},
		"expensive": 10
	}`

	jsonValue, err := parser.ParseToValue(jsonStr)
	if err != nil {
		t.Fatalf("解析JSON失败: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		expected int
		check    func(results []interface{}) bool
	}{
		{
			name:     "根节点",
			path:     "$",
			expected: 1,
			check: func(results []interface{}) bool {
				return len(results) == 1
			},
		},
		{
			name:     "对象属性",
			path:     "$.store",
			expected: 1,
			check: func(results []interface{}) bool {
				return len(results) == 1
			},
		},
		{
			name:     "嵌套对象属性",
			path:     "$.store.bicycle",
			expected: 1,
			check: func(results []interface{}) bool {
				return len(results) == 1
			},
		},
		{
			name:     "数组",
			path:     "$.store.book",
			expected: 1,
			check: func(results []interface{}) bool {
				return len(results) == 1
			},
		},
		{
			name:     "数组索引",
			path:     "$.store.book[0]",
			expected: 1,
			check: func(results []interface{}) bool {
				return len(results) == 1
			},
		},
		{
			name:     "数组通配符",
			path:     "$.store.book[*]",
			expected: 4,
			check: func(results []interface{}) bool {
				return len(results) == 4
			},
		},
		{
			name:     "数组切片-指定范围",
			path:     "$.store.book[1:3]",
			expected: 2,
			check: func(results []interface{}) bool {
				return len(results) == 2
			},
		},
		{
			name:     "数组切片-省略起始索引",
			path:     "$.store.book[:2]",
			expected: 2,
			check: func(results []interface{}) bool {
				return len(results) == 2
			},
		},
		{
			name:     "数组切片-省略结束索引",
			path:     "$.store.book[2:]",
			expected: 2,
			check: func(results []interface{}) bool {
				return len(results) == 2
			},
		},
		{
			name:     "数组切片-负索引",
			path:     "$.store.book[-2:]",
			expected: 2,
			check: func(results []interface{}) bool {
				return len(results) == 2
			},
		},
		{
			name:     "数组切片-全部",
			path:     "$.store.book[:]",
			expected: 4,
			check: func(results []interface{}) bool {
				return len(results) == 4
			},
		},
		{
			name:     "数组属性",
			path:     "$.store.book[*].author",
			expected: 4,
			check: func(results []interface{}) bool {
				return len(results) == 4
			},
		},
		{
			name:     "对象通配符",
			path:     "$.store.*",
			expected: 2,
			check: func(results []interface{}) bool {
				return len(results) == 2
			},
		},
		{
			name:     "数组切片和属性访问",
			path:     "$.store.book[1:3].title",
			expected: 2,
			check: func(results []interface{}) bool {
				// 只检查结果数量
				return len(results) == 2
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := QueryJSONPath(jsonValue, tt.path)
			if err != nil {
				t.Fatalf("查询失败: %v", err)
			}

			if len(results) != tt.expected {
				t.Errorf("结果数量不匹配: 期望 %d, 实际 %d", tt.expected, len(results))
			}

			// 转换为interface{}切片进行检查
			interfaceResults := make([]interface{}, len(results))
			for i, r := range results {
				interfaceResults[i] = r
			}

			if !tt.check(interfaceResults) {
				t.Errorf("结果验证失败")
			}
		})
	}
}

func TestJSONPathErrors(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "空路径",
			path: "",
		},
		{
			name: "无效路径开头",
			path: "store.book",
		},
		{
			name: "无效属性名",
			path: "$.store.123",
		},
		{
			name: "括号不匹配",
			path: "$.store.book[0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseJSONPath(tt.path)
			if err == nil {
				t.Errorf("期望错误，但没有错误")
			}
		})
	}
}

func TestQueryJSONPathString(t *testing.T) {
	jsonStr := `{"name":"John","age":30,"address":{"city":"New York"}}`

	results, err := QueryJSONPathString(jsonStr, "$.name")
	if err != nil {
		t.Fatalf("查询失败: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("结果数量不匹配: 期望 1, 实际 %d", len(results))
	}

	if !results[0].IsString() {
		t.Errorf("结果类型不匹配: 期望 string, 实际 %s", results[0].Type())
	}

	val, _ := results[0].AsString()
	if val != "John" {
		t.Errorf("结果值不匹配: 期望 John, 实际 %s", val)
	}
}
