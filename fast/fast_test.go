package fast

import (
	"encoding/json"
	"reflect"
	"testing"
)

// TestMarshal 测试Marshal函数
func TestMarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name:    "空对象",
			input:   map[string]interface{}{},
			wantErr: false,
		},
		{
			name: "简单对象",
			input: map[string]interface{}{
				"name": "张三",
				"age":  30,
			},
			wantErr: false,
		},
		{
			name: "嵌套对象",
			input: map[string]interface{}{
				"name": "张三",
				"age":  30,
				"address": map[string]interface{}{
					"city":     "北京",
					"district": "海淀区",
				},
			},
			wantErr: false,
		},
		{
			name: "包含数组",
			input: map[string]interface{}{
				"name":    "张三",
				"hobbies": []string{"阅读", "编程", "旅行"},
			},
			wantErr: false,
		},
		{
			name:    "nil值",
			input:   nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 使用Marshal
			got, err := Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 使用标准库进行比较
			expected, err := json.Marshal(tt.input)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}

			// 比较结果（注意：Marshal可能会在末尾添加换行符，需要处理）
			if len(got) > 0 && got[len(got)-1] == '\n' {
				got = got[:len(got)-1]
			}

			// 解析JSON以比较内容而不是字符串表示
			var gotObj interface{}
			var expectedObj interface{}

			if err := json.Unmarshal(got, &gotObj); err != nil {
				t.Errorf("无法解析Marshal结果: %v", err)
				return
			}

			if err := json.Unmarshal(expected, &expectedObj); err != nil {
				t.Errorf("无法解析expected结果: %v", err)
				return
			}

			if !reflect.DeepEqual(gotObj, expectedObj) {
				t.Errorf("Marshal() = %s, want %s", got, expected)
			}
		})
	}
}

// TestUnmarshal 测试Unmarshal函数
func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    interface{}
		wantErr bool
	}{
		{
			name:    "空对象",
			input:   "{}",
			want:    map[string]interface{}{},
			wantErr: false,
		},
		{
			name:  "简单对象",
			input: `{"name":"张三","age":30}`,
			want: map[string]interface{}{
				"name": "张三",
				"age":  json.Number("30"),
			},
			wantErr: false,
		},
		{
			name:  "嵌套对象",
			input: `{"name":"张三","address":{"city":"北京"}}`,
			want: map[string]interface{}{
				"name": "张三",
				"address": map[string]interface{}{
					"city": "北京",
				},
			},
			wantErr: false,
		},
		{
			name:  "包含数组",
			input: `{"name":"张三","hobbies":["阅读","编程","旅行"]}`,
			want: map[string]interface{}{
				"name":    "张三",
				"hobbies": []interface{}{"阅读", "编程", "旅行"},
			},
			wantErr: false,
		},
		{
			name:    "空字符串",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "无效JSON",
			input:   "{invalid}",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got interface{}
			err := Unmarshal([]byte(tt.input), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 由于Unmarshal使用json.Number，需要特殊处理比较
				if !deepEqual(got, tt.want) {
					// 转换为JSON字符串进行比较
					gotJSON, _ := json.Marshal(got)
					wantJSON, _ := json.Marshal(tt.want)

					// 解析回来再比较
					var gotObj, wantObj interface{}
					json.Unmarshal(gotJSON, &gotObj)
					json.Unmarshal(wantJSON, &wantObj)

					if !reflect.DeepEqual(gotObj, wantObj) {
						t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
					}
				}
			}
		})
	}
}

// TestFragmentCache 测试片段缓存功能
func TestFragmentCache(t *testing.T) {
	// 清空缓存，确保测试环境干净
	ClearFragmentCache()

	// 测试缓存设置和获取
	key := "testKey"
	value := map[string]interface{}{
		"name": "张三",
		"age":  30,
	}

	// 设置缓存
	CacheFragment(key, value)

	// 获取缓存
	got, exists := GetCachedFragment(key)
	if !exists {
		t.Errorf("GetCachedFragment() 未找到缓存的键 %s", key)
		return
	}

	// 比较值
	if !reflect.DeepEqual(got, value) {
		t.Errorf("GetCachedFragment() = %v, want %v", got, value)
	}

	// 测试清空缓存
	ClearFragmentCache()
	_, exists = GetCachedFragment(key)
	if exists {
		t.Errorf("ClearFragmentCache() 后仍能找到缓存的键 %s", key)
	}
}

// TestBufferPool 测试缓冲池功能
func TestBufferPool(t *testing.T) {
	// 获取缓冲区
	buf := getBuffer()
	if buf == nil {
		t.Errorf("getBuffer() 返回nil")
		return
	}

	// 写入数据
	testData := "测试数据"
	buf.WriteString(testData)

	// 检查数据
	if buf.String() != testData {
		t.Errorf("缓冲区内容 = %s, want %s", buf.String(), testData)
	}

	// 释放缓冲区
	releaseBuffer(buf)

	// 再次获取缓冲区，应该是空的
	buf = getBuffer()
	if buf.Len() != 0 {
		t.Errorf("重用的缓冲区不为空，长度 = %d", buf.Len())
	}

	// 再次释放
	releaseBuffer(buf)
}

// 辅助函数：深度比较，处理json.Number类型
func deepEqual(a, b interface{}) bool {
	if a == nil || b == nil {
		return a == b
	}

	// 处理map类型
	aMap, aIsMap := a.(map[string]interface{})
	bMap, bIsMap := b.(map[string]interface{})

	if aIsMap && bIsMap {
		if len(aMap) != len(bMap) {
			return false
		}
		for k, aVal := range aMap {
			bVal, ok := bMap[k]
			if !ok || !deepEqual(aVal, bVal) {
				return false
			}
		}
		return true
	}

	// 处理slice类型
	aSlice, aIsSlice := a.([]interface{})
	bSlice, bIsSlice := b.([]interface{})

	if aIsSlice && bIsSlice {
		if len(aSlice) != len(bSlice) {
			return false
		}
		for i := range aSlice {
			if !deepEqual(aSlice[i], bSlice[i]) {
				return false
			}
		}
		return true
	}

	// 处理json.Number类型
	aNum, aIsNum := a.(json.Number)
	if aIsNum {
		switch b := b.(type) {
		case json.Number:
			return aNum.String() == b.String()
		case float64:
			bFloat, err := aNum.Float64()
			return err == nil && bFloat == b
		case int:
			bInt, err := aNum.Int64()
			return err == nil && bInt == int64(b)
		}
	}

	// 其他类型直接比较
	return reflect.DeepEqual(a, b)
}
