# GoJSON 示例

本目录包含 GoJSON 库的示例代码，展示了库的各种功能和用法。

## 示例列表

1. **basic_usage.go** - 基本用法示例，展示如何创建、修改和序列化JSON对象
2. **parsing.go** - 解析JSON字符串的示例
3. **jsonpath.go** - 使用JSON Path查询JSON的示例
4. **diff_patch.go** - 比较JSON差异和应用JSON Patch的示例
5. **streaming.go** - 流式处理JSON的示例
6. **generic.go** - 使用泛型功能的示例
7. **performance.go** - 性能优化相关的示例
8. **utils.go** - 使用工具函数的示例

## 运行示例

要运行示例，请使用以下命令：

```bash
go run examples/basic_usage.go
```

或者

```bash
cd examples
go run basic_usage.go
```

## 示例说明

### 基本用法示例 (basic_usage.go)

展示了如何创建JSON对象和数组，设置和获取值，以及序列化为字符串。

### 解析示例 (parsing.go)

展示了如何解析JSON字符串为Go对象，以及如何处理解析错误。

### JSON Path示例 (jsonpath.go)

展示了如何使用JSON Path查询JSON数据，包括基本查询、数组访问、过滤器等。

### 差异和补丁示例 (diff_patch.go)

展示了如何比较两个JSON对象的差异，以及如何应用JSON Patch。

### 流式处理示例 (streaming.go)

展示了如何使用流式API处理大型JSON数据，减少内存使用。

### 泛型示例 (generic.go)

展示了如何使用泛型功能，提供类型安全的JSON处理。

### 性能示例 (performance.go)

展示了如何使用性能优化功能，如FastMarshal和FastUnmarshal。

### 工具示例 (utils.go)

展示了如何使用各种工具函数，如JSON美化、压缩、路径提取和结构分析。
