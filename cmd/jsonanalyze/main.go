// jsonanalyze 是一个JSON结构分析工具，用于分析JSON的结构
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/UserLeeZJ/gojson/jsonpath"
	"github.com/UserLeeZJ/gojson/parser"
	"github.com/UserLeeZJ/gojson/utils"
)

var (
	inputFile  string
	outputFile string
	path       string
	showPaths  bool
)

func init() {
	flag.StringVar(&inputFile, "i", "", "输入文件路径，如果为空则从标准输入读取")
	flag.StringVar(&outputFile, "o", "", "输出文件路径，如果为空则输出到标准输出")
	flag.StringVar(&path, "p", "$", "JSON Path表达式，用于分析特定路径的结构")
	flag.BoolVar(&showPaths, "paths", false, "显示所有可能的JSON Path")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "jsonanalyze - JSON结构分析工具\n\n")
	fmt.Fprintf(os.Stderr, "用法:\n")
	fmt.Fprintf(os.Stderr, "  jsonanalyze [选项]\n\n")
	fmt.Fprintf(os.Stderr, "选项:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n示例:\n")
	fmt.Fprintf(os.Stderr, "  jsonanalyze -i input.json\n")
	fmt.Fprintf(os.Stderr, "  cat input.json | jsonanalyze\n")
	fmt.Fprintf(os.Stderr, "  jsonanalyze -i input.json -paths\n")
	fmt.Fprintf(os.Stderr, "  jsonanalyze -i input.json -p \"$.store.book\"\n")
}

func main() {
	flag.Parse()

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

	// 如果指定了路径，先执行查询
	if path != "$" {
		results, err := jsonpath.QueryJSONPath(jsonValue, path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "查询失败: %v\n", err)
			os.Exit(1)
		}
		if len(results) == 0 {
			fmt.Fprintf(os.Stderr, "路径 %s 没有匹配的结果\n", path)
			os.Exit(1)
		}
		// 使用第一个结果
		jsonValue = results[0]
	}

	// 准备输出
	var output strings.Builder

	// 显示基本信息
	output.WriteString(fmt.Sprintf("JSON分析结果 (路径: %s)\n", path))
	output.WriteString("====================\n\n")

	// 如果需要显示所有路径
	if showPaths {
		paths := utils.ExtractPaths(jsonValue)
		sort.Strings(paths)
		
		output.WriteString("所有可能的JSON Path:\n")
		for _, p := range paths {
			output.WriteString(fmt.Sprintf("  %s\n", p))
		}
		output.WriteString("\n")
	}

	// 分析结构
	info := utils.AnalyzeStructure(jsonValue)
	output.WriteString("结构分析:\n")
	output.WriteString(info.String())

	// 写入输出
	if outputFile == "" {
		fmt.Print(output.String())
	} else {
		err = os.WriteFile(outputFile, []byte(output.String()), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "写入输出失败: %v\n", err)
			os.Exit(1)
		}
	}
}
