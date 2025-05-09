// jsonstream 是一个JSON流式处理工具，用于处理大型JSON文件
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/UserLeeZJ/gojson/stream"
	"github.com/UserLeeZJ/gojson/types"
	"github.com/UserLeeZJ/gojson/utils"
)

var (
	inputFile  string
	outputFile string
	filter     string
	limit      int
	pretty     bool
	compact    bool
)

func init() {
	flag.StringVar(&inputFile, "i", "", "输入文件路径，如果为空则从标准输入读取")
	flag.StringVar(&outputFile, "o", "", "输出文件路径，如果为空则输出到标准输出")
	flag.StringVar(&filter, "f", "$", "JSON Path过滤器，用于选择要处理的元素")
	flag.IntVar(&limit, "limit", 0, "限制输出的元素数量，0表示不限制")
	flag.BoolVar(&pretty, "pretty", false, "输出为美化格式")
	flag.BoolVar(&compact, "c", false, "输出为紧凑格式")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "jsonstream - JSON流式处理工具\n\n")
	fmt.Fprintf(os.Stderr, "用法:\n")
	fmt.Fprintf(os.Stderr, "  jsonstream [选项]\n\n")
	fmt.Fprintf(os.Stderr, "选项:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n示例:\n")
	fmt.Fprintf(os.Stderr, "  jsonstream -i large.json -o output.json -f \"$.items[*].name\"\n")
	fmt.Fprintf(os.Stderr, "  cat large.json | jsonstream -f \"$.items[*]\" > output.json\n")
}

func main() {
	flag.Parse()

	// 检查参数
	if pretty && compact {
		fmt.Fprintf(os.Stderr, "错误: 不能同时指定美化格式和紧凑格式\n")
		os.Exit(1)
	}

	// 打开输入
	var input io.Reader

	if inputFile == "" {
		input = os.Stdin
	} else {
		file, err := os.Open(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "打开输入文件失败: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		input = file
	}

	// 打开输出
	var output io.Writer
	if outputFile == "" {
		output = os.Stdout
	} else {
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "创建输出文件失败: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		output = file
	}

	// 创建流式解析器
	tokenizer := stream.NewJSONTokenizer(input)

	// 创建输出缓冲区
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	// 处理流
	processStream(tokenizer, writer)
}

func processStream(tokenizer *stream.JSONTokenizer, writer *bufio.Writer) {
	// 解析过滤器
	segments := parseFilter(filter)

	// 当前路径
	currentPath := make([]string, 0)

	// 当前深度
	depth := 0

	// 是否匹配
	matching := false

	// 匹配深度
	matchDepth := 0

	// 当前匹配的值
	var currentValue types.JSONValue

	// 计数器
	count := 0

	// 是否是第一个输出
	first := true

	// 写入数组开始
	writer.WriteString("[\n")

	// 处理令牌
	for {
		token := tokenizer.Next()

		// 检查是否结束
		if token.Type == stream.TokenEOF {
			break
		}

		// 检查是否有错误
		if token.Type == stream.TokenError {
			fmt.Fprintf(os.Stderr, "解析错误: %v\n", token.Error)
			break
		}

		// 更新路径
		switch token.Type {
		case stream.TokenObjectStart:
			depth++
			currentPath = append(currentPath, "{")
		case stream.TokenObjectEnd:
			if len(currentPath) > 0 {
				currentPath = currentPath[:len(currentPath)-1]
			}
			depth--
		case stream.TokenArrayStart:
			depth++
			currentPath = append(currentPath, "[")
		case stream.TokenArrayEnd:
			if len(currentPath) > 0 {
				currentPath = currentPath[:len(currentPath)-1]
			}
			depth--
		case stream.TokenPropertyName:
			if len(currentPath) > 0 && currentPath[len(currentPath)-1] == "{" {
				currentPath[len(currentPath)-1] = token.Value.(string)
			} else {
				currentPath = append(currentPath, token.Value.(string))
			}
		}

		// 检查是否匹配过滤器
		pathStr := buildPath(currentPath)
		if matchesFilter(pathStr, segments) && !matching {
			matching = true
			matchDepth = depth
			currentValue = nil
		}

		// 如果匹配，收集值
		if matching {
			// 如果深度小于匹配深度，结束匹配
			if depth < matchDepth {
				matching = false

				// 如果有值，输出
				if currentValue != nil {
					// 检查是否达到限制
					if limit > 0 && count >= limit {
						break
					}

					// 输出分隔符
					if !first {
						writer.WriteString(",\n")
					} else {
						first = false
					}

					// 格式化输出
					var output string
					if pretty {
						output, _ = utils.PrettyPrint(currentValue, utils.DefaultPrettyOptions())
					} else if compact {
						output, _ = utils.CompressJSON(currentValue)
					} else {
						output = currentValue.String()
					}

					writer.WriteString(output)
					count++
				}
			}
		}
	}

	// 写入数组结束
	writer.WriteString("\n]")
}

// 解析过滤器
func parseFilter(filter string) []string {
	// 简单实现，只支持基本路径
	segments := strings.Split(filter, ".")
	if segments[0] == "$" {
		segments = segments[1:]
	}
	return segments
}

// 构建路径
func buildPath(path []string) string {
	return "$." + strings.Join(path, ".")
}

// 检查路径是否匹配过滤器
func matchesFilter(path string, segments []string) bool {
	// 简单实现，只支持基本匹配
	return strings.Contains(path, strings.Join(segments, "."))
}
