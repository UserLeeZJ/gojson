// 示例：流式处理JSON
package main

import (
	"fmt"
	"strings"

	"github.com/UserLeeZJ/gojson/stream"
)

func main() {
	// 1. 流式解析示例
	fmt.Println("1. 流式解析示例")

	// 准备JSON数据
	jsonData := `{
		"name": "张三",
		"age": 30,
		"address": {
			"city": "北京",
			"district": "海淀区"
		},
		"hobbies": ["读书", "旅行"],
		"active": true,
		"data": null
	}`

	// 创建流式解析器
	tokenizer := stream.NewJSONTokenizer(strings.NewReader(jsonData))

	// 逐个处理令牌
	fmt.Println("解析令牌:")
	for {
		token := tokenizer.Next()

		// 检查是否到达文件末尾
		if token.Type == stream.TokenEOF {
			fmt.Println("  解析完成")
			break
		}

		// 检查是否有错误
		if token.Type == stream.TokenError {
			fmt.Printf("  错误: %v\n", token.Error)
			break
		}

		// 打印令牌信息
		fmt.Printf("  类型: %-15v 深度: %d", token.Type, token.Depth)

		// 打印值（如果有）
		if token.Value != nil {
			fmt.Printf(" 值: %v", token.Value)
		}

		// 打印路径（如果有）
		if token.Path != "" {
			fmt.Printf(" 路径: %s", token.Path)
		}

		fmt.Println()
	}

	// 2. 流式生成示例
	fmt.Println("\n2. 流式生成示例")

	// 创建一个字符串构建器作为输出目标
	var sb strings.Builder

	// 创建流式生成器
	generator := stream.NewJSONGenerator(&sb)

	// 生成JSON
	fmt.Println("生成JSON:")

	// 开始对象
	generator.BeginObject()

	// 添加简单属性
	generator.WriteProperty("name")
	generator.WriteString("李四")

	generator.WriteProperty("age")
	generator.WriteNumber(25)

	generator.WriteProperty("active")
	generator.WriteBoolean(true)

	generator.WriteProperty("data")
	generator.WriteNull()

	// 添加嵌套对象
	generator.WriteProperty("address")
	generator.BeginObject()

	generator.WriteProperty("city")
	generator.WriteString("上海")

	generator.WriteProperty("district")
	generator.WriteString("浦东新区")

	// 结束嵌套对象
	generator.EndObject()

	// 添加数组
	generator.WriteProperty("scores")
	generator.BeginArray()

	generator.WriteNumber(85)
	generator.WriteNumber(90)
	generator.WriteNumber(95)

	// 结束数组
	generator.EndArray()

	// 结束对象
	generator.EndObject()

	// 刷新缓冲区
	generator.Flush()

	// 打印生成的JSON
	fmt.Println(sb.String())

	// 3. 增量解析示例
	fmt.Println("\n3. 增量解析示例")

	// 准备一个大型JSON，分块处理
	largeJSON := `{
		"id": 12345,
		"name": "产品名称",
		"description": "这是一个示例产品",
		"price": 99.99,
		"categories": ["电子", "家电", "智能设备"],
		"specifications": {
			"weight": "1.5kg",
			"dimensions": "10 x 20 x 5 cm",
			"color": "黑色"
		}
	}`

	// 将JSON分成多个块
	chunks := []string{
		largeJSON[:50],    // 第一块
		largeJSON[50:100], // 第二块
		largeJSON[100:],   // 第三块
	}

	// 创建增量解析器
	parser := stream.NewIncrementalParser()

	// 逐块提供数据
	fmt.Println("增量解析:")
	for i, chunk := range chunks {
		fmt.Printf("  提供第%d块数据 (长度: %d)\n", i+1, len(chunk))
		err := parser.Feed([]byte(chunk))
		if err != nil {
			fmt.Printf("  解析错误: %v\n", err)
			break
		}

		fmt.Printf("  解析完成: %v\n", parser.IsComplete())
	}

	// 获取结果
	if parser.IsComplete() {
		result, err := parser.Result()
		if err != nil {
			fmt.Printf("  获取结果错误: %v\n", err)
		} else {
			fmt.Println("  解析结果:")
			fmt.Printf("  %v\n", result)
		}
	}

	// 4. 实际应用：处理大文件
	fmt.Println("\n4. 实际应用：处理大文件")
	fmt.Println("  (这里只是示例，不实际创建大文件)")

	// 模拟处理大文件的代码
	fmt.Println("  处理大型JSON文件的步骤:")
	fmt.Println("  1. 打开文件")
	fmt.Println("  2. 创建流式解析器")
	fmt.Println("  3. 逐个处理令牌")
	fmt.Println("  4. 根据需要提取数据")
	fmt.Println("  5. 关闭文件")

	// 示例代码（不实际执行）
	fmt.Println("\n  示例代码:")

	// 使用字符串拼接，避免格式化指令问题
	code1 := "  file, err := os.Open(\"large_file.json\")"
	code2 := "  if err != nil {"
	code3 := "      fmt.Println(\"打开文件失败:\", err)"
	code4 := "      return"
	code5 := "  }"
	code6 := "  defer file.Close()"
	code7 := ""
	code8 := "  tokenizer := stream.NewJSONTokenizer(file)"
	code9 := ""
	code10 := "  for {"
	code11 := "      token := tokenizer.Next()"
	code12 := ""
	code13 := "      if token.Type == stream.TokenEOF {"
	code14 := "          break"
	code15 := "      }"
	code16 := ""
	code17 := "      if token.Type == stream.TokenError {"
	code18 := "          fmt.Println(\"解析错误:\", token.Error)"
	code19 := "          break"
	code20 := "      }"
	code21 := ""
	code22 := "      // 处理特定路径的数据"
	code23 := "      if token.Type == stream.TokenString && strings.HasPrefix(token.Path, \"$.items[\") && strings.HasSuffix(token.Path, \"].name\") {"
	code24 := "          fmt.Println(\"找到项目名称:\", token.Value)"
	code25 := "      }"
	code26 := "  }"

	fmt.Println(code1)
	fmt.Println(code2)
	fmt.Println(code3)
	fmt.Println(code4)
	fmt.Println(code5)
	fmt.Println(code6)
	fmt.Println(code7)
	fmt.Println(code8)
	fmt.Println(code9)
	fmt.Println(code10)
	fmt.Println(code11)
	fmt.Println(code12)
	fmt.Println(code13)
	fmt.Println(code14)
	fmt.Println(code15)
	fmt.Println(code16)
	fmt.Println(code17)
	fmt.Println(code18)
	fmt.Println(code19)
	fmt.Println(code20)
	fmt.Println(code21)
	fmt.Println(code22)
	fmt.Println(code23)
	fmt.Println(code24)
	fmt.Println(code25)
	fmt.Println(code26)

	fmt.Println("\n示例结束")
}
