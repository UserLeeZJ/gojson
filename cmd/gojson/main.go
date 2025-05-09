// gojson 是一个集成了多个JSON工具的命令行程序
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	version = "1.0.0" // 版本号
)

func main() {
	// 检查命令行参数
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// 获取子命令
	subcommand := os.Args[1]

	// 处理版本和帮助命令
	if subcommand == "-v" || subcommand == "--version" {
		fmt.Printf("gojson version %s\n", version)
		os.Exit(0)
	}
	if subcommand == "-h" || subcommand == "--help" {
		printUsage()
		os.Exit(0)
	}

	// 获取可执行文件路径
	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取可执行文件路径失败: %v\n", err)
		os.Exit(1)
	}
	exeDir := filepath.Dir(exePath)

	// 构建子命令路径
	var cmdPath string
	switch subcommand {
	case "format":
		cmdPath = filepath.Join(exeDir, "jsonformat")
	case "path":
		cmdPath = filepath.Join(exeDir, "jsonpath")
	case "analyze":
		cmdPath = filepath.Join(exeDir, "jsonanalyze")
	case "stream":
		cmdPath = filepath.Join(exeDir, "jsonstream")
	default:
		fmt.Fprintf(os.Stderr, "未知的子命令: %s\n", subcommand)
		printUsage()
		os.Exit(1)
	}

	// 检查子命令是否存在
	if _, err := os.Stat(cmdPath); os.IsNotExist(err) {
		// 尝试在PATH中查找
		cmdPath = "json" + subcommand
	}

	// 执行子命令
	cmd := exec.Command(cmdPath, os.Args[2:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		} else {
			fmt.Fprintf(os.Stderr, "执行子命令失败: %v\n", err)
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "gojson - JSON工具集\n\n")
	fmt.Fprintf(os.Stderr, "用法:\n")
	fmt.Fprintf(os.Stderr, "  gojson <子命令> [选项]\n\n")
	fmt.Fprintf(os.Stderr, "子命令:\n")
	fmt.Fprintf(os.Stderr, "  format   格式化JSON (美化或压缩)\n")
	fmt.Fprintf(os.Stderr, "  path     使用JSON Path查询JSON\n")
	fmt.Fprintf(os.Stderr, "  analyze  分析JSON结构\n")
	fmt.Fprintf(os.Stderr, "  stream   流式处理大型JSON文件\n\n")
	fmt.Fprintf(os.Stderr, "全局选项:\n")
	fmt.Fprintf(os.Stderr, "  -v, --version  显示版本信息\n")
	fmt.Fprintf(os.Stderr, "  -h, --help     显示帮助信息\n\n")
	fmt.Fprintf(os.Stderr, "示例:\n")
	fmt.Fprintf(os.Stderr, "  gojson format -i input.json -o output.json -p\n")
	fmt.Fprintf(os.Stderr, "  gojson path -i input.json -p \"$.store.book[0].title\"\n")
	fmt.Fprintf(os.Stderr, "  gojson analyze -i input.json -paths\n")
	fmt.Fprintf(os.Stderr, "  gojson stream -i large.json -f \"$.items[*].name\"\n\n")
	fmt.Fprintf(os.Stderr, "使用 'gojson <子命令> --help' 获取子命令的详细帮助信息\n")
}
