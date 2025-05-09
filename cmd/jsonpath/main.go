// jsonpath 是一个JSON Path查询工具，用于从JSON中提取数据
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/UserLeeZJ/gojson/jsonpath"
	"github.com/UserLeeZJ/gojson/parser"
	"github.com/UserLeeZJ/gojson/types"
	"github.com/UserLeeZJ/gojson/utils"
)

var (
	inputFile  string
	path       string
	compact    bool
	pretty     bool
	outputFile string
)

func init() {
	flag.StringVar(&inputFile, "i", "", "输入文件路径，如果为空则从标准输入读取")
	flag.StringVar(&path, "p", "$", "JSON Path表达式")
	flag.BoolVar(&compact, "c", false, "输出为紧凑格式")
	flag.BoolVar(&pretty, "pretty", false, "输出为美化格式")
	flag.StringVar(&outputFile, "o", "", "输出文件路径，如果为空则输出到标准输出")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "jsonpath - JSON Path查询工具\n\n")
	fmt.Fprintf(os.Stderr, "用法:\n")
	fmt.Fprintf(os.Stderr, "  jsonpath [选项]\n\n")
	fmt.Fprintf(os.Stderr, "选项:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n示例:\n")
	fmt.Fprintf(os.Stderr, "  jsonpath -i input.json -p \"$.store.book[0].title\"\n")
	fmt.Fprintf(os.Stderr, "  cat input.json | jsonpath -p \"$.store.book[*].author\"\n")
}

func main() {
	flag.Parse()

	// 检查参数
	if compact && pretty {
		fmt.Fprintf(os.Stderr, "错误: 不能同时指定紧凑格式和美化格式\n")
		os.Exit(1)
	}

	// 读取输入
	var input []byte
	var err error
	if inputFile == "" {
		input, err = io.ReadAll(os.Stdin)
	} else {
		input, err = os.ReadFile(inputFile)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取输入失败: %v\n", err)
		os.Exit(1)
	}

	// 解析JSON
	jsonValue, err := parser.ParseToValue(string(input))
	if err != nil {
		fmt.Fprintf(os.Stderr, "解析JSON失败: %v\n", err)
		os.Exit(1)
	}

	// 执行JSON Path查询
	results, err := jsonpath.QueryJSONPath(jsonValue, path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "查询失败: %v\n", err)
		os.Exit(1)
	}

	// 处理结果
	var output string
	if len(results) == 0 {
		output = "[]" // 空结果
	} else if len(results) == 1 {
		// 单个结果
		if compact {
			output, err = utils.CompressJSON(results[0])
		} else if pretty {
			output, err = utils.PrettyPrint(results[0], utils.DefaultPrettyOptions())
		} else {
			output = results[0].String()
		}
	} else {
		// 多个结果，包装为数组
		array := types.NewJSONArray()
		for _, result := range results {
			array.Add(result)
		}
		
		if compact {
			output, err = utils.CompressJSON(array)
		} else if pretty {
			output, err = utils.PrettyPrint(array, utils.DefaultPrettyOptions())
		} else {
			output = array.String()
		}
	}
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "格式化结果失败: %v\n", err)
		os.Exit(1)
	}

	// 写入输出
	if outputFile == "" {
		fmt.Println(output)
	} else {
		err = os.WriteFile(outputFile, []byte(output), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "写入输出失败: %v\n", err)
			os.Exit(1)
		}
	}
}
