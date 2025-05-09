package benchmarks

import (
	"encoding/json"
	"testing"

	"github.com/UserLeeZJ/gojson/fast"
)

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
		err := fast.Unmarshal(jsonData, &result)
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
				err := fast.Unmarshal(jsonData, bm.value)
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

// BenchmarkLargeJSONUnmarshal 基准测试大型JSON的反序列化性能
func BenchmarkLargeJSONUnmarshal(b *testing.B) {
	// 准备大型测试数据
	testData := map[string]interface{}{
		"items": make([]interface{}, 1000),
	}

	// 填充大数组
	for i := 0; i < 1000; i++ {
		testData["items"].([]interface{})[i] = map[string]interface{}{
			"id":    i,
			"name":  "Item " + string(rune(i)),
			"value": i * 10,
			"tags":  []string{"tag1", "tag2", "tag3"},
		}
	}

	// 序列化为JSON
	jsonData, _ := json.Marshal(testData)

	// 重置计时器
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		var result interface{}
		err := fast.Unmarshal(jsonData, &result)
		if err != nil {
			b.Fatalf("Unmarshal() error = %v", err)
		}
	}
}

// BenchmarkLargeJSONStandardUnmarshal 基准测试标准库处理大型JSON的反序列化性能
func BenchmarkLargeJSONStandardUnmarshal(b *testing.B) {
	// 准备大型测试数据
	testData := map[string]interface{}{
		"items": make([]interface{}, 1000),
	}

	// 填充大数组
	for i := 0; i < 1000; i++ {
		testData["items"].([]interface{})[i] = map[string]interface{}{
			"id":    i,
			"name":  "Item " + string(rune(i)),
			"value": i * 10,
			"tags":  []string{"tag1", "tag2", "tag3"},
		}
	}

	// 序列化为JSON
	jsonData, _ := json.Marshal(testData)

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
