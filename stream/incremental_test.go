package stream

import (
	"encoding/json"
	"testing"
)

func TestIncrementalParser(t *testing.T) {
	tests := []struct {
		name   string
		chunks []string
		valid  bool
	}{
		{
			name:   "一次性解析",
			chunks: []string{`{"name":"John","age":30}`},
			valid:  true,
		},
		{
			name:   "分块解析-对象",
			chunks: []string{`{"name":"`, `John","age":`, `30}`},
			valid:  true,
		},
		{
			name:   "分块解析-数组",
			chunks: []string{`[1,2,`, `3,4,`, `5]`},
			valid:  true,
		},
		{
			name:   "分块解析-嵌套",
			chunks: []string{`{"person":{"name":"`, `John","age":30},"active":`, `true}`},
			valid:  true,
		},
		{
			name:   "无效JSON",
			chunks: []string{`{"name":"John"`, `,"age":}`},
			valid:  false,
		},
		{
			name:   "不完整JSON",
			chunks: []string{`{"name":"John"`},
			valid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewIncrementalParser()

			// 逐块提供数据
			for i, chunk := range tt.chunks {
				err := parser.Feed([]byte(chunk))
				if err != nil {
					if tt.valid {
						t.Fatalf("块 %d: 解析失败: %v", i, err)
					} else {
						// 预期失败，测试通过
						return
					}
				}
			}

			// 检查是否完成
			if !parser.IsComplete() && tt.valid {
				t.Errorf("解析未完成")
			}

			// 获取结果
			result, err := parser.Result()
			if err != nil {
				if tt.valid {
					t.Fatalf("获取结果失败: %v", err)
				} else {
					// 预期失败，测试通过
					return
				}
			}

			if !tt.valid {
				t.Errorf("预期解析失败，但成功了")
				return
			}

			// 验证结果
			// 将所有块连接起来，然后使用标准库解析
			var fullJSON string
			for _, chunk := range tt.chunks {
				fullJSON += chunk
			}

			var expected interface{}
			if err := json.Unmarshal([]byte(fullJSON), &expected); err != nil {
				t.Fatalf("解析完整JSON失败: %v", err)
			}

			// 转换为JSON字符串进行比较
			expectedJSON, _ := json.Marshal(expected)
			actualJSON, _ := json.Marshal(result)

			if string(expectedJSON) != string(actualJSON) {
				t.Errorf("结果不匹配:\n期望: %s\n实际: %s", string(expectedJSON), string(actualJSON))
			}
		})
	}
}

func TestIncrementalParserReset(t *testing.T) {
	parser := NewIncrementalParser()

	// 第一次解析
	err := parser.Feed([]byte(`{"name":"John"}`))
	if err != nil {
		t.Fatalf("第一次解析失败: %v", err)
	}

	if !parser.IsComplete() {
		t.Errorf("第一次解析未完成")
	}

	// 重置解析器
	parser.Reset()

	if parser.IsComplete() {
		t.Errorf("重置后解析器仍然标记为完成")
	}

	// 第二次解析
	err = parser.Feed([]byte(`[1,2,3]`))
	if err != nil {
		t.Fatalf("第二次解析失败: %v", err)
	}

	if !parser.IsComplete() {
		t.Errorf("第二次解析未完成")
	}

	// 获取结果
	result, err := parser.Result()
	if err != nil {
		t.Fatalf("获取结果失败: %v", err)
	}

	// 验证结果
	arr, ok := result.([]interface{})
	if !ok {
		t.Fatalf("结果类型不是数组")
	}

	if len(arr) != 3 {
		t.Errorf("数组长度不匹配: 期望 3, 实际 %d", len(arr))
	}
}

func TestResultValue(t *testing.T) {
	parser := NewIncrementalParser()

	// 解析对象
	err := parser.Feed([]byte(`{"name":"John","age":30}`))
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	// 获取JSONValue结果
	result, err := parser.ResultValue()
	if err != nil {
		t.Fatalf("获取JSONValue结果失败: %v", err)
	}

	// 验证结果
	if !result.IsObject() {
		t.Fatalf("结果类型不是对象")
	}

	obj, _ := result.AsObject()
	
	name, err := obj.GetString("name")
	if err != nil {
		t.Fatalf("获取name属性失败: %v", err)
	}
	if name != "John" {
		t.Errorf("name属性不匹配: 期望 John, 实际 %s", name)
	}

	age, err := obj.GetNumber("age")
	if err != nil {
		t.Fatalf("获取age属性失败: %v", err)
	}
	if age != 30 {
		t.Errorf("age属性不匹配: 期望 30, 实际 %f", age)
	}
}
