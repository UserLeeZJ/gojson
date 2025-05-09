# GoJSON

[![Go](https://github.com/UserLeeZJ/gojson/actions/workflows/go.yml/badge.svg)](https://github.com/UserLeeZJ/gojson/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/UserLeeZJ/gojson)](https://goreportcard.com/report/github.com/UserLeeZJ/gojson)
[![codecov](https://codecov.io/gh/UserLeeZJ/gojson/branch/main/graph/badge.svg)](https://codecov.io/gh/UserLeeZJ/gojson)
[![GoDoc](https://godoc.org/github.com/UserLeeZJ/gojson?status.svg)](https://godoc.org/github.com/UserLeeZJ/gojson)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

`GoJSON` 是一个Go语言库，提供类似JavaScript JSON接口的功能，使JSON处理更加简单和直观。该库专注于提供面向对象的JSON操作体验，同时保持高性能和可靠性。

## 特性

- 类似JavaScript的JSON API，熟悉的接口设计
- 面向对象的JSON值类型（JSONObject, JSONArray, JSONString等）
- 简单易用的链式调用接口
- 结构化的错误处理系统
- 模块化的代码结构，便于维护和扩展
- 高性能实现，适用于大型JSON处理
- JSON Path查询支持，轻松从复杂JSON中提取数据
- JSON Diff功能，比较JSON结构差异
- JSON Patch (RFC 6902)标准支持
- 流式处理支持，适用于超大型JSON数据
- 泛型支持，提供类型安全的JSON处理
- 命令行工具，方便在终端中处理JSON数据
- 详细的单元测试和文档
- 纯Go实现，无外部依赖
- 支持Go 1.21+

## 安装

### 库

```bash
go get github.com/UserLeeZJ/gojson
```

### 命令行工具

```bash
# 安装所有工具
go install github.com/UserLeeZJ/gojson/cmd/...@latest

# 或者安装单个工具
go install github.com/UserLeeZJ/gojson/cmd/jsonformat@latest
go install github.com/UserLeeZJ/gojson/cmd/jsonpath@latest
go install github.com/UserLeeZJ/gojson/cmd/jsonanalyze@latest
go install github.com/UserLeeZJ/gojson/cmd/jsonstream@latest
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/UserLeeZJ/gojson"
)

func main() {
    // 创建JSON对象
    person := gojson.NewJSONObject()
    person.PutString("name", "张三")
    person.PutNumber("age", 28)

    // 添加嵌套对象
    address := gojson.NewJSONObject()
    address.PutString("city", "北京")
    address.PutString("district", "海淀区")
    person.PutObject("address", address)

    // 添加数组
    hobbies := gojson.NewJSONArray()
    hobbies.AddString("阅读").AddString("编程").AddString("旅行")
    person.PutArray("hobbies", hobbies)

    // 输出格式化的JSON
    fmt.Println(person.String())

    // 使用JSON Path查询
    results, _ := gojson.QueryJSONPath(person, "$.hobbies[1]")
    fmt.Println("第二个爱好:", results[0].String())
}
```

## 使用示例

### 解析JSON字符串

```go
package main

import (
    "fmt"
    "github.com/UserLeeZJ/gojson"
)

func main() {
    jsonStr := `{"name":"John","age":30,"address":{"city":"New York"}}`

    var data interface{}
    err := gojson.Parse(jsonStr, &data)
    if err != nil {
        fmt.Println("解析错误:", err)
        return
    }

    fmt.Println("解析结果:", data)
}
```

### 将Go对象转换为JSON字符串

```go
package main

import (
    "fmt"
    "github.com/UserLeeZJ/gojson"
)

func main() {
    data := map[string]interface{}{
        "name": "John",
        "age":  30,
        "address": map[string]interface{}{
            "city": "New York",
        },
    }

    jsonStr, err := gojson.Stringify(data)
    if err != nil {
        fmt.Println("序列化错误:", err)
        return
    }

    fmt.Println("JSON字符串:", jsonStr)
}
```

### 使用JSONObject

```go
package main

import (
    "fmt"
    "github.com/UserLeeZJ/gojson"
)

func main() {
    // 创建一个新的JSONObject
    obj := gojson.NewJSONObject()

    // 添加各种类型的值
    obj.PutString("name", "John")
    obj.PutNumber("age", 30)
    obj.PutBoolean("active", true)

    // 创建并添加嵌套对象
    address := gojson.NewJSONObject()
    address.PutString("city", "New York")
    address.PutString("country", "USA")
    obj.PutObject("address", address)

    // 创建并添加数组
    hobbies := gojson.NewJSONArray()
    hobbies.AddString("reading").AddString("swimming")
    obj.PutArray("hobbies", hobbies)

    // 转换为JSON字符串
    jsonStr := obj.String()
    fmt.Println(jsonStr)

    // 获取值
    name, _ := obj.GetString("name")
    age, _ := obj.GetNumber("age")
    fmt.Printf("Name: %s, Age: %.0f\n", name, age)
}
```

### 使用JSONArray

```go
package main

import (
    "fmt"
    "github.com/UserLeeZJ/gojson"
)

func main() {
    // 创建一个新的JSONArray
    arr := gojson.NewJSONArray()

    // 添加各种类型的值
    arr.AddString("hello")
    arr.AddNumber(123)
    arr.AddBoolean(true)

    // 添加嵌套对象
    person := gojson.NewJSONObject()
    person.PutString("name", "John")
    person.PutNumber("age", 30)
    arr.Add(person)

    // 转换为JSON字符串
    jsonStr := arr.String()
    fmt.Println(jsonStr)

    // 使用ForEach遍历数组
    arr.ForEach(func(value gojson.JSONValue, index int) {
        fmt.Printf("Index %d: Type %s\n", index, value.Type())
    })
}
```

### 使用Patch方法

```go
package main

import (
    "fmt"
    "github.com/UserLeeZJ/gojson"
)

func main() {
    // 创建一个JSONObject
    obj := gojson.NewJSONObject()
    obj.PutString("name", "John")
    obj.PutNumber("age", 30)

    // 创建一个JSON Patch
    patchJSON := `[
        {"op":"add","path":"/email","value":"john@example.com"},
        {"op":"remove","path":"/age"},
        {"op":"replace","path":"/name","value":"Jane"}
    ]`

    // 应用补丁
    result, err := obj.Patch(patchJSON)
    if err != nil {
        fmt.Println("补丁应用错误:", err)
        return
    }

    // 输出结果
    fmt.Println("补丁应用后:", result.String())
    // 输出: {"email":"john@example.com","name":"Jane"}
}
```

## 主要功能

### JSONObject

- 创建和操作JSON对象
- 添加、获取和删除属性
- 合并对象
- 克隆对象
- 应用JSON Patch (RFC 6902)

### JSONArray

- 创建和操作JSON数组
- 添加、获取和删除元素
- 遍历、映射和过滤数组元素
- 支持链式调用API

### JSON Path

```go
package main

import (
    "fmt"
    "github.com/UserLeeZJ/gojson"
)

func main() {
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
                }
            ],
            "bicycle": {
                "color": "red",
                "price": 19.95
            }
        }
    }`

    // 解析JSON
    jsonValue, _ := gojson.ParseToValue(jsonStr)

    // 使用JSON Path查询
    // 获取所有书籍的标题
    results, _ := gojson.QueryJSONPath(jsonValue, "$.store.book[*].title")

    fmt.Println("所有书籍标题:")
    for _, result := range results {
        fmt.Println("-", result.String())
    }

    // 获取价格大于10的书籍
    expensiveBooks, _ := gojson.QueryJSONPath(jsonValue, "$.store.book[?(@.price > 10)]")

    fmt.Println("\n价格大于10的书籍:")
    for _, book := range expensiveBooks {
        title, _ := book.AsObject().GetString("title")
        price, _ := book.AsObject().GetNumber("price")
        fmt.Printf("- %s (¥%.2f)\n", title, price)
    }
}
```

### JSON Diff

```go
package main

import (
    "fmt"
    "github.com/UserLeeZJ/gojson"
)

func main() {
    // 原始JSON
    oldJSON := `{"name":"张三","age":30,"address":{"city":"北京"}}`

    // 修改后的JSON
    newJSON := `{"name":"张三","age":31,"address":{"city":"上海","district":"浦东"}}`

    // 比较差异
    diffs, _ := gojson.DiffJSONStrings(oldJSON, newJSON, nil)

    fmt.Println("JSON差异:")
    for _, diff := range diffs {
        fmt.Println("-", diff.String())
    }

    // 生成JSON Patch
    patch := gojson.GeneratePatch(diffs)
    fmt.Println("\nJSON Patch:")
    fmt.Println(patch.String())
}
```

### 错误处理

```go
package main

import (
    "fmt"
    "github.com/UserLeeZJ/gojson"
)

func main() {
    jsonStr := `{"name":"John","age":invalid}`

    _, err := gojson.ParseToValue(jsonStr)
    if err != nil {
        // 使用类型断言获取详细错误信息
        if jsonErr, ok := err.(*gojson.JSONError); ok {
            fmt.Printf("错误代码: %s\n", jsonErr.Code)
            fmt.Printf("错误消息: %s\n", jsonErr.Message)
            if jsonErr.Cause != nil {
                fmt.Printf("原始错误: %v\n", jsonErr.Cause)
            }
        } else {
            fmt.Printf("未知错误: %v\n", err)
        }
    }
}
```

### 性能优化

```go
package main

import (
    "fmt"
    "github.com/UserLeeZJ/gojson"
)

func main() {
    // 使用优化的序列化/反序列化函数
    data := map[string]interface{}{
        "items": make([]interface{}, 1000),
    }

    // 填充大数组
    for i := 0; i < 1000; i++ {
        data["items"].([]interface{})[i] = i
    }

    // 使用优化的序列化函数
    jsonBytes, _ := gojson.FastMarshal(data)
    fmt.Printf("序列化后大小: %d 字节\n", len(jsonBytes))

    // 使用优化的反序列化函数
    var result interface{}
    gojson.FastUnmarshal(jsonBytes, &result)
}
```

#### 优化技术

gojson 使用了多种优化技术来提高性能：

1. **快速路径处理**：为简单类型（字符串、数字、布尔值等）提供专门的处理路径，避免通用处理的开销
2. **预计算值**：预计算常用数字的字符串表示，减少运行时转换
3. **缓冲池优化**：使用对象池减少内存分配和垃圾回收压力
4. **零拷贝技术**：在可能的情况下避免内存复制
5. **分片锁**：使用分片锁减少锁竞争，提高并发性能
6. **字符串处理优化**：特别优化了字符串的序列化和反序列化
7. **内存分配优化**：预分配合理大小的缓冲区，减少扩容操作
8. **单次扫描**：尽可能在一次扫描中完成解析，避免多次处理同一数据

这些优化技术参考了流行的第三方库（如 jsoniter、easyjson），但完全使用纯 Go 实现，无外部依赖。

### 其他类型

- JSONBool - 表示JSON中的布尔值
- JSONNull - 表示JSON中的null值
- JSONNumber - 表示JSON中的数字
- JSONString - 表示JSON中的字符串

## 项目结构

GoJSON采用模块化的代码结构，便于维护和扩展：

- **types**: 包含所有JSON值类型的定义（JSONObject, JSONArray, JSONString等）
- **parser**: 提供JSON解析和序列化功能
- **jsonpath**: 实现JSON Path查询功能
- **diff**: 提供JSON差异比较功能
- **patch**: 实现JSON Patch (RFC 6902)功能
- **errors**: 提供结构化的错误处理系统
- **fast**: 提供高性能的JSON序列化和反序列化功能
- **stream**: 提供流式处理JSON的功能
- **generic**: 提供泛型支持，增强类型安全
- **utils**: 提供各种实用工具函数
- **cmd**: 提供命令行工具
- **benchmarks**: 包含基准测试代码
- **examples**: 包含示例代码

主包（gojson）重新导出所有子包的公共API，使用户可以通过单一导入路径访问所有功能。

### 目录结构

```bash
gojson/
├── benchmarks/       # 基准测试代码
├── cmd/              # 命令行工具
│   ├── gojson/       # 主命令行工具
│   ├── jsonformat/   # JSON格式化工具
│   ├── jsonpath/     # JSON Path查询工具
│   ├── jsonanalyze/  # JSON结构分析工具
│   └── jsonstream/   # JSON流式处理工具
├── diff/             # JSON差异比较功能
├── errors/           # 结构化的错误处理系统
├── examples/         # 示例代码
├── fast/             # 高性能JSON序列化和反序列化
├── generic/          # 泛型支持
├── jsonpath/         # JSON Path查询功能
├── parser/           # JSON解析和序列化功能
├── patch/            # JSON Patch功能
├── stream/           # 流式处理JSON功能
├── types/            # JSON值类型定义
├── utils/            # 实用工具函数
├── .github/          # GitHub配置文件
├── go.mod            # Go模块定义
├── go.sum            # Go模块依赖校验
├── gojson.go         # 主包
├── LICENSE           # 许可证
├── Makefile          # 构建脚本
└── README.md         # 项目说明
```

## 性能比较

与标准库 `encoding/json` 相比，gojson 提供了显著的性能提升：

### 序列化性能

| 数据类型 | GoJSON | 标准库 | 性能提升 |
|---------|--------|--------|---------|
| 复杂对象 | 4105 ns/op | 5591 ns/op | 快 27% |
| 字符串 | 73.52 ns/op | 224.5 ns/op | 快 3倍 |
| 整数 | 39.33 ns/op | 145.9 ns/op | 快 3.7倍 |
| 浮点数 | 195.0 ns/op | 245.9 ns/op | 快 1.3倍 |
| 布尔值 | 18.95 ns/op | 125.4 ns/op | 快 6.6倍 |
| 小型Map | 597.8 ns/op | 1559 ns/op | 快 2.6倍 |
| 小型数组 | 472.5 ns/op | 700.5 ns/op | 快 1.5倍 |

### 反序列化性能

| 数据类型 | GoJSON | 标准库 | 性能提升 |
|---------|--------|--------|---------|
| 字符串 | 65.97 ns/op | 450.5 ns/op | 快 6.8倍 |
| 整数 | 23.17 ns/op | 264.4 ns/op | 快 11.4倍 |
| 浮点数 | 292.9 ns/op | 282.6 ns/op | 相当 |
| 布尔值 | 15.99 ns/op | 205.5 ns/op | 快 12.9倍 |
| 复杂对象 | 7837 ns/op | 7781 ns/op | 相当 |

### 内存使用

| 操作 | 数据类型 | GoJSON | 标准库 | 内存减少 |
|------|---------|--------|--------|---------|
| 序列化 | 复杂对象 | 1010 B/op | 1280 B/op | 减少 21% |
| 序列化 | 小型Map | 160 B/op | 368 B/op | 减少 57% |
| 反序列化 | 整数 | 0 B/op | 144 B/op | 减少 100% |
| 反序列化 | 布尔值 | 0 B/op | 144 B/op | 减少 100% |

这些性能数据基于最新的基准测试结果，使用 Go 1.21+ 在标准硬件上测试。

## 命令行工具详解

GoJSON 提供了一系列命令行工具，方便在命令行中处理 JSON 数据：

### jsonformat

JSON 格式化工具，用于美化和压缩 JSON。

```bash
# 美化 JSON
jsonformat -i input.json -o output.json -p

# 压缩 JSON
jsonformat -i input.json -o output.json -c

# 从标准输入读取，输出到标准输出
cat input.json | jsonformat -p > output.json

# 排序键
jsonformat -i input.json -o output.json -p -s
```

### jsonpath

JSON Path 查询工具，用于从 JSON 中提取数据。

```bash
# 使用 JSON Path 查询
jsonpath -i input.json -p "$.store.book[0].title"

# 从标准输入读取
cat input.json | jsonpath -p "$.store.book[*].author"

# 输出为美化格式
jsonpath -i input.json -p "$.store.book[*]" -pretty
```

### jsonanalyze

JSON 结构分析工具，用于分析 JSON 的结构。

```bash
# 分析 JSON 结构
jsonanalyze -i input.json

# 输出所有 JSON Path
jsonanalyze -i input.json -paths

# 分析特定路径的结构
jsonanalyze -i input.json -p "$.store.book"
```

### jsonstream

JSON 流式处理工具，用于处理大型 JSON 文件。

```bash
# 流式处理大型 JSON 文件
jsonstream -i large.json -f "$.items[*].name"

# 从标准输入读取
cat large.json | jsonstream -f "$.items[*]" > output.json

# 限制输出数量
jsonstream -i large.json -f "$.items[*]" -limit 10
```

### 统一入口

所有工具也可以通过 `gojson` 命令统一访问：

```bash
gojson format -i input.json -o output.json -p
gojson path -i input.json -p "$.store.book[0].title"
gojson analyze -i input.json -paths
gojson stream -i large.json -f "$.items[*].name"
```

## 开发

### 构建和测试

GoJSON 使用 Makefile 来简化常见的开发任务。以下是一些常用的命令：

```bash
# 构建项目
make build

# 运行测试
make test

# 运行基准测试
make bench

# 运行示例
make examples

# 生成代码覆盖率报告
make coverage

# 运行代码检查
make lint

# 清理构建产物
make clean

# 显示帮助信息
make help
```

### 持续集成

GoJSON 使用 GitHub Actions 进行持续集成，包括：

- 在多个 Go 版本上构建和测试
- 运行基准测试
- 生成代码覆盖率报告
- 代码质量检查
- 自动发布

## 贡献

欢迎贡献代码、报告问题或提出改进建议。请通过 GitHub Issues 或 Pull Requests 参与项目。

贡献步骤：

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

GNU GENERAL PUBLIC LICENSE
