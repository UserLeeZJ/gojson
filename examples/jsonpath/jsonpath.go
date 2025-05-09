// JSON Path查询示例
package main

import (
	"fmt"

	"github.com/UserLeeZJ/gojson/jsonpath"
	"github.com/UserLeeZJ/gojson/parser"
	"github.com/UserLeeZJ/gojson/types"
)

func main() {
	fmt.Println("GoJSON Path查询示例")
	fmt.Println("====================")

	// 准备测试数据
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

	// 解析JSON
	jsonValue, err := parser.ParseToValue(jsonStr)
	if err != nil {
		fmt.Printf("解析JSON失败: %v\n", err)
		return
	}

	// 基本路径查询
	fmt.Println("\n1. 基本路径查询")

	// 获取根对象
	results, err := jsonpath.QueryJSONPath(jsonValue, "$")
	printResults("根对象", results, err)

	// 获取store对象
	results, err = jsonpath.QueryJSONPath(jsonValue, "$.store")
	printResults("store对象", results, err)

	// 获取bicycle对象
	results, err = jsonpath.QueryJSONPath(jsonValue, "$.store.bicycle")
	printResults("bicycle对象", results, err)

	// 获取bicycle的颜色
	results, err = jsonpath.QueryJSONPath(jsonValue, "$.store.bicycle.color")
	printResults("bicycle的颜色", results, err)

	// 数组访问
	fmt.Println("\n2. 数组访问")

	// 获取第一本书
	results, err = jsonpath.QueryJSONPath(jsonValue, "$.store.book[0]")
	printResults("第一本书", results, err)

	// 获取最后一本书
	results, err = jsonpath.QueryJSONPath(jsonValue, "$.store.book[-1]")
	printResults("最后一本书", results, err)

	// 获取前两本书
	results, err = jsonpath.QueryJSONPath(jsonValue, "$.store.book[0:2]")
	printResults("前两本书", results, err)

	// 获取所有书籍
	results, err = jsonpath.QueryJSONPath(jsonValue, "$.store.book[*]")
	printResults("所有书籍", results, err)

	// 属性访问
	fmt.Println("\n3. 属性访问")

	// 获取所有书籍的标题
	results, err = jsonpath.QueryJSONPath(jsonValue, "$.store.book[*].title")
	printResults("所有书籍的标题", results, err)

	// 获取所有书籍的作者
	results, err = jsonpath.QueryJSONPath(jsonValue, "$.store.book[*].author")
	printResults("所有书籍的作者", results, err)

	// 过滤器
	fmt.Println("\n4. 过滤器")

	// 获取所有价格大于10的书籍
	results, err = jsonpath.QueryJSONPath(jsonValue, "$.store.book[?(@.price > 10)]")
	printResults("价格大于10的书籍", results, err)

	// 获取所有价格大于expensive的书籍
	results, err = jsonpath.QueryJSONPath(jsonValue, "$.store.book[?(@.price > $.expensive)]")
	printResults("价格大于expensive的书籍", results, err)

	// 获取所有类别为fiction的书籍
	results, err = jsonpath.QueryJSONPath(jsonValue, "$.store.book[?(@.category == 'fiction')]")
	printResults("类别为fiction的书籍", results, err)

	// 获取所有有ISBN的书籍
	results, err = jsonpath.QueryJSONPath(jsonValue, "$.store.book[?(@.isbn)]")
	printResults("有ISBN的书籍", results, err)

	// 递归查询
	fmt.Println("\n5. 递归查询")

	// 获取所有价格
	results, err = jsonpath.QueryJSONPath(jsonValue, "$..price")
	printResults("所有价格", results, err)

	// 获取所有作者
	results, err = jsonpath.QueryJSONPath(jsonValue, "$..author")
	printResults("所有作者", results, err)

	// 组合查询
	fmt.Println("\n6. 组合查询")

	// 获取所有价格大于10的书籍的标题
	results, err = jsonpath.QueryJSONPath(jsonValue, "$.store.book[?(@.price > 10)].title")
	printResults("价格大于10的书籍的标题", results, err)

	// 获取所有类别为fiction的书籍的作者
	results, err = jsonpath.QueryJSONPath(jsonValue, "$.store.book[?(@.category == 'fiction')].author")
	printResults("类别为fiction的书籍的作者", results, err)

	// 数组切片
	fmt.Println("\n7. 数组切片")

	// 获取第2和第3本书
	results, err = jsonpath.QueryJSONPath(jsonValue, "$.store.book[1:3]")
	printResults("第2和第3本书", results, err)

	// 获取前2本书
	results, err = jsonpath.QueryJSONPath(jsonValue, "$.store.book[:2]")
	printResults("前2本书", results, err)

	// 获取最后2本书
	results, err = jsonpath.QueryJSONPath(jsonValue, "$.store.book[-2:]")
	printResults("最后2本书", results, err)

	fmt.Println("\n示例结束")
}

// 打印查询结果
func printResults(description string, results []types.JSONValue, err error) {
	fmt.Printf("%s:\n", description)
	if err != nil {
		fmt.Printf("  错误: %v\n", err)
		return
	}

	fmt.Printf("  结果数量: %d\n", len(results))
	for i, result := range results {
		if i < 5 { // 只显示前5个结果
			fmt.Printf("  %d: %s\n", i+1, result.String())
		}
	}
	if len(results) > 5 {
		fmt.Printf("  ... (共%d个结果)\n", len(results))
	}
}
