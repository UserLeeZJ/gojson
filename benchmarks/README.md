# GoJSON 基准测试

本目录包含 GoJSON 库的基准测试代码。这些测试用于评估库的性能，并与其他 JSON 库进行比较。

## 运行基准测试

要运行所有基准测试，请使用以下命令：

```bash
go test -bench=. ./benchmarks
```

要运行特定的基准测试，请使用以下命令：

```bash
go test -bench=BenchmarkParse ./benchmarks
```

要包含内存分配统计信息，请添加 `-benchmem` 标志：

```bash
go test -bench=. -benchmem ./benchmarks
```

## 基准测试类别

基准测试分为以下几类：

1. **解析测试**：测试 JSON 字符串解析为 Go 对象的性能
2. **序列化测试**：测试 Go 对象序列化为 JSON 字符串的性能
3. **查询测试**：测试 JSON Path 查询的性能
4. **修改测试**：测试修改 JSON 对象的性能
5. **比较测试**：与其他流行的 JSON 库进行性能比较

## 测试数据

基准测试使用以下数据集：

1. **小型 JSON**：少于 1KB 的简单 JSON 对象
2. **中型 JSON**：1KB 到 10KB 的 JSON 对象，包含嵌套结构
3. **大型 JSON**：大于 10KB 的复杂 JSON 对象
4. **数组 JSON**：包含大量数组元素的 JSON
5. **深度嵌套 JSON**：具有多层嵌套的 JSON 对象

## 比较结果

最新的基准测试结果可以在 [RESULTS.md](RESULTS.md) 文件中找到。
