package fast

import (
	"encoding/json"
	"fmt"
	"testing"
)

// BenchmarkMarshal 基准测试Marshal函数
func BenchmarkMarshal(b *testing.B) {
	// 准备测试数据
	testData := map[string]interface{}{
		"name": "张三",
		"age":  30,
		"address": map[string]interface{}{
			"city":     "北京",
			"district": "海淀区",
			"zipcode":  "100000",
		},
		"hobbies": []string{"阅读", "编程", "旅行", "摄影", "音乐"},
		"skills": []map[string]interface{}{
			{"name": "Go", "level": 9},
			{"name": "Java", "level": 8},
			{"name": "Python", "level": 7},
		},
	}

	// 重置计时器
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		_, err := Marshal(testData)
		if err != nil {
			b.Fatalf("Marshal() error = %v", err)
		}
	}
}

// BenchmarkStandardMarshal 基准测试标准库的Marshal函数
func BenchmarkStandardMarshal(b *testing.B) {
	// 准备测试数据
	testData := map[string]interface{}{
		"name": "张三",
		"age":  30,
		"address": map[string]interface{}{
			"city":     "北京",
			"district": "海淀区",
			"zipcode":  "100000",
		},
		"hobbies": []string{"阅读", "编程", "旅行", "摄影", "音乐"},
		"skills": []map[string]interface{}{
			{"name": "Go", "level": 9},
			{"name": "Java", "level": 8},
			{"name": "Python", "level": 7},
		},
	}

	// 重置计时器
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(testData)
		if err != nil {
			b.Fatalf("json.Marshal() error = %v", err)
		}
	}
}

// BenchmarkUnmarshal 基准测试Unmarshal函数
func BenchmarkUnmarshal(b *testing.B) {
	// 准备测试数据
	jsonData := []byte(`{
		"name": "张三",
		"age": 30,
		"address": {
			"city": "北京",
			"district": "海淀区",
			"zipcode": "100000"
		},
		"hobbies": ["阅读", "编程", "旅行", "摄影", "音乐"],
		"skills": [
			{"name": "Go", "level": 9},
			{"name": "Java", "level": 8},
			{"name": "Python", "level": 7}
		]
	}`)

	// 重置计时器
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		var result interface{}
		err := Unmarshal(jsonData, &result)
		if err != nil {
			b.Fatalf("Unmarshal() error = %v", err)
		}
	}
}

// BenchmarkStandardUnmarshal 基准测试标准库的Unmarshal函数
func BenchmarkStandardUnmarshal(b *testing.B) {
	// 准备测试数据
	jsonData := []byte(`{
		"name": "张三",
		"age": 30,
		"address": {
			"city": "北京",
			"district": "海淀区",
			"zipcode": "100000"
		},
		"hobbies": ["阅读", "编程", "旅行", "摄影", "音乐"],
		"skills": [
			{"name": "Go", "level": 9},
			{"name": "Java", "level": 8},
			{"name": "Python", "level": 7}
		]
	}`)

	// 重置计时器
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		var result interface{}
		err := json.Unmarshal(jsonData, &result)
		if err != nil {
			b.Fatalf("json.Unmarshal() error = %v", err)
		}
	}
}

// BenchmarkLargeJSONMarshal 基准测试大型JSON的序列化性能
func BenchmarkLargeJSONMarshal(b *testing.B) {
	// 准备大型测试数据
	testData := map[string]interface{}{
		"items": make([]interface{}, 1000),
	}

	// 填充大数组
	for i := 0; i < 1000; i++ {
		testData["items"].([]interface{})[i] = map[string]interface{}{
			"id":    i,
			"name":  fmt.Sprintf("Item %d", i),
			"value": i * 10,
			"tags":  []string{"tag1", "tag2", "tag3"},
		}
	}

	// 重置计时器
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		_, err := Marshal(testData)
		if err != nil {
			b.Fatalf("Marshal() error = %v", err)
		}
	}
}

// BenchmarkLargeJSONStandardMarshal 基准测试标准库处理大型JSON的序列化性能
func BenchmarkLargeJSONStandardMarshal(b *testing.B) {
	// 准备大型测试数据
	testData := map[string]interface{}{
		"items": make([]interface{}, 1000),
	}

	// 填充大数组
	for i := 0; i < 1000; i++ {
		testData["items"].([]interface{})[i] = map[string]interface{}{
			"id":    i,
			"name":  fmt.Sprintf("Item %d", i),
			"value": i * 10,
			"tags":  []string{"tag1", "tag2", "tag3"},
		}
	}

	// 重置计时器
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(testData)
		if err != nil {
			b.Fatalf("json.Marshal() error = %v", err)
		}
	}
}

// BenchmarkSimpleTypeMarshal 基准测试简单类型的序列化性能
func BenchmarkSimpleTypeMarshal(b *testing.B) {
	// 测试不同的简单类型
	benchmarks := []struct {
		name  string
		value interface{}
	}{
		{"String", "这是一个测试字符串"},
		{"Int", 12345},
		{"Float", 123.45},
		{"Bool", true},
		{"Null", nil},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := Marshal(bm.value)
				if err != nil {
					b.Fatalf("Marshal() error = %v", err)
				}
			}
		})
	}
}

// BenchmarkSimpleTypeStandardMarshal 基准测试标准库处理简单类型的序列化性能
func BenchmarkSimpleTypeStandardMarshal(b *testing.B) {
	// 测试不同的简单类型
	benchmarks := []struct {
		name  string
		value interface{}
	}{
		{"String", "这是一个测试字符串"},
		{"Int", 12345},
		{"Float", 123.45},
		{"Bool", true},
		{"Null", nil},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := json.Marshal(bm.value)
				if err != nil {
					b.Fatalf("json.Marshal() error = %v", err)
				}
			}
		})
	}
}

// BenchmarkSimpleTypeUnmarshal 基准测试简单类型的反序列化性能
func BenchmarkSimpleTypeUnmarshal(b *testing.B) {
	// 测试不同的简单类型
	benchmarks := []struct {
		name  string
		json  string
		value interface{}
	}{
		{"String", `"这是一个测试字符串"`, new(string)},
		{"Int", `12345`, new(int)},
		{"Float", `123.45`, new(float64)},
		{"Bool", `true`, new(bool)},
		{"Null", `null`, new(interface{})},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			jsonData := []byte(bm.json)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err := Unmarshal(jsonData, bm.value)
				if err != nil {
					b.Fatalf("Unmarshal() error = %v", err)
				}
			}
		})
	}
}

// BenchmarkSimpleTypeStandardUnmarshal 基准测试标准库处理简单类型的反序列化性能
func BenchmarkSimpleTypeStandardUnmarshal(b *testing.B) {
	// 测试不同的简单类型
	benchmarks := []struct {
		name  string
		json  string
		value interface{}
	}{
		{"String", `"这是一个测试字符串"`, new(string)},
		{"Int", `12345`, new(int)},
		{"Float", `123.45`, new(float64)},
		{"Bool", `true`, new(bool)},
		{"Null", `null`, new(interface{})},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			jsonData := []byte(bm.json)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err := json.Unmarshal(jsonData, bm.value)
				if err != nil {
					b.Fatalf("json.Unmarshal() error = %v", err)
				}
			}
		})
	}
}
