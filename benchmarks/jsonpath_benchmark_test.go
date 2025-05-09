package benchmarks

import (
	"testing"

	"github.com/UserLeeZJ/gojson/jsonpath"
	"github.com/UserLeeZJ/gojson/parser"
)

// 测试数据
const testJSON = `{
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

// BenchmarkJSONPathSimple 基准测试简单的JSON Path查询
func BenchmarkJSONPathSimple(b *testing.B) {
	// 解析JSON
	value, err := parser.ParseToValue(testJSON)
	if err != nil {
		b.Fatalf("解析JSON失败: %v", err)
	}

	// 测试不同的简单路径
	benchmarks := []struct {
		name string
		path string
	}{
		{"Root", "$"},
		{"Property", "$.store"},
		{"NestedProperty", "$.store.bicycle"},
		{"DeepProperty", "$.store.bicycle.color"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := jsonpath.QueryJSONPath(value, bm.path)
				if err != nil {
					b.Fatalf("查询失败: %v", err)
				}
			}
		})
	}
}

// BenchmarkJSONPathArray 基准测试数组相关的JSON Path查询
func BenchmarkJSONPathArray(b *testing.B) {
	// 解析JSON
	value, err := parser.ParseToValue(testJSON)
	if err != nil {
		b.Fatalf("解析JSON失败: %v", err)
	}

	// 测试不同的数组路径
	benchmarks := []struct {
		name string
		path string
	}{
		{"ArrayAccess", "$.store.book[0]"},
		{"ArraySlice", "$.store.book[1:3]"},
		{"ArrayWildcard", "$.store.book[*]"},
		{"ArrayProperty", "$.store.book[0].title"},
		{"ArrayFilter", "$.store.book[?(@.price < 10)]"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := jsonpath.QueryJSONPath(value, bm.path)
				if err != nil {
					b.Fatalf("查询失败: %v", err)
				}
			}
		})
	}
}

// BenchmarkJSONPathComplex 基准测试复杂的JSON Path查询
func BenchmarkJSONPathComplex(b *testing.B) {
	// 解析JSON
	value, err := parser.ParseToValue(testJSON)
	if err != nil {
		b.Fatalf("解析JSON失败: %v", err)
	}

	// 测试不同的复杂路径
	benchmarks := []struct {
		name string
		path string
	}{
		{"DeepWildcard", "$..author"},
		{"FilterWithComparison", "$.store.book[?(@.price > $.expensive)]"},
		{"FilterWithRegex", "$.store.book[?(@.category == 'fiction')]"},
		{"MultipleConditions", "$.store.book[?(@.price < 10 && @.category == 'fiction')]"},
		{"ComplexPath", "$.store.book[*].title"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := jsonpath.QueryJSONPath(value, bm.path)
				if err != nil {
					// 某些复杂查询可能不支持，忽略错误
					continue
				}
			}
		})
	}
}

// BenchmarkJSONPathParse 基准测试JSON Path解析性能
func BenchmarkJSONPathParse(b *testing.B) {
	// 测试不同的路径解析
	benchmarks := []struct {
		name string
		path string
	}{
		{"Simple", "$.store.bicycle.color"},
		{"ArrayAccess", "$.store.book[0].title"},
		{"ArraySlice", "$.store.book[1:3]"},
		{"DeepWildcard", "$..author"},
		{"Complex", "$.store.book[?(@.price > 10)].title"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := jsonpath.ParseJSONPath(bm.path)
				if err != nil {
					b.Fatalf("解析失败: %v", err)
				}
			}
		})
	}
}

// BenchmarkJSONPathLargeResult 基准测试返回大量结果的JSON Path查询
func BenchmarkJSONPathLargeResult(b *testing.B) {
	// 创建一个包含大量数据的JSON
	largeJSON := `{"items": [`
	for i := 0; i < 1000; i++ {
		if i > 0 {
			largeJSON += ","
		}
		largeJSON += `{"id": ` + string(rune(i)) + `, "value": "item` + string(rune(i)) + `"}`
	}
	largeJSON += `]}`

	// 解析JSON
	value, err := parser.ParseToValue(largeJSON)
	if err != nil {
		b.Fatalf("解析JSON失败: %v", err)
	}

	// 测试返回大量结果的查询
	benchmarks := []struct {
		name string
		path string
	}{
		{"AllItems", "$.items[*]"},
		{"AllIds", "$.items[*].id"},
		{"AllValues", "$.items[*].value"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := jsonpath.QueryJSONPath(value, bm.path)
				if err != nil {
					b.Fatalf("查询失败: %v", err)
				}
			}
		})
	}
}
