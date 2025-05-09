// jsonformat 是一个JSON格式化工具，用于美化和压缩JSON
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/UserLeeZJ/gojson/parser"
	"github.com/UserLeeZJ/gojson/utils"
)

var (
	inputFile  string
	outputFile string
	pretty     bool
	compress   bool
	sortKeys   bool
	indent     string
	escapeHTML bool
)

func init() {
	flag.StringVar(&inputFile, "i", "", "输入文件路径，如果为空则从标准输入读取")
	flag.StringVar(&outputFile, "o", "", "输出文件路径，如果为空则输出到标准输出")
	flag.BoolVar(&pretty, "p", false, "美化JSON")
	flag.BoolVar(&compress, "c", false, "压缩JSON")
	flag.BoolVar(&sortKeys, "s", false, "排序键")
	flag.StringVar(&indent, "indent", "  ", "缩进字符串")
	flag.BoolVar(&escapeHTML, "escape-html", false, "转义HTML字符")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "jsonformat - JSON格式化工具\n\n")
	fmt.Fprintf(os.Stderr, "用法:\n")
	fmt.Fprintf(os.Stderr, "  jsonformat [选项]\n\n")
	fmt.Fprintf(os.Stderr, "选项:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n示例:\n")
	fmt.Fprintf(os.Stderr, "  jsonformat -i input.json -o output.json -p\n")
	fmt.Fprintf(os.Stderr, "  cat input.json | jsonformat -p > output.json\n")
}

func main() {
	flag.Parse()

	// 检查参数
	if !pretty && !compress {
		pretty = true // 默认美化
	}
	if pretty && compress {
		fmt.Fprintf(os.Stderr, "错误: 不能同时指定美化和压缩\n")
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

	// 格式化JSON
	var output string
	if pretty {
		options := utils.PrettyOptions{
			Indent:     indent,
			SortKeys:   sortKeys,
			EscapeHTML: escapeHTML,
		}
		output, err = utils.PrettyPrint(jsonValue, options)
	} else {
		output, err = utils.CompressJSON(jsonValue)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "格式化JSON失败: %v\n", err)
		os.Exit(1)
	}

	// 写入输出
	if outputFile == "" {
		fmt.Print(output)
	} else {
		err = os.WriteFile(outputFile, []byte(output), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "写入输出失败: %v\n", err)
			os.Exit(1)
		}
	}
}
