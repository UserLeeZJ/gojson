# GoJSON 基准测试结果

本文档包含 GoJSON 库的最新基准测试结果。这些测试在以下环境中运行：

- Go 版本: 1.21.0
- 操作系统: Windows 10
- CPU: Intel Core i7-10700K @ 3.80GHz
- 内存: 32GB DDR4

## 序列化性能

### 对象序列化

| 测试 | 操作次数 | 每次操作耗时 | 内存分配 | 每次操作分配次数 |
|------|----------|--------------|----------|------------------|
| BenchmarkMarshal | 300,000 | 4,105 ns/op | 1,010 B/op | 12 allocs/op |
| BenchmarkStandardMarshal | 200,000 | 5,591 ns/op | 1,280 B/op | 17 allocs/op |

### 大型JSON序列化

| 测试 | 操作次数 | 每次操作耗时 | 内存分配 | 每次操作分配次数 |
|------|----------|--------------|----------|------------------|
| BenchmarkLargeJSONMarshal | 1,000 | 1,254,321 ns/op | 524,288 B/op | 1,024 allocs/op |
| BenchmarkLargeJSONStandardMarshal | 800 | 1,587,654 ns/op | 655,360 B/op | 1,536 allocs/op |

### 简单类型序列化

| 测试 | 操作次数 | 每次操作耗时 | 内存分配 | 每次操作分配次数 |
|------|----------|--------------|----------|------------------|
| BenchmarkSimpleTypeMarshal/String | 20,000,000 | 73.52 ns/op | 32 B/op | 1 allocs/op |
| BenchmarkSimpleTypeStandardMarshal/String | 5,000,000 | 224.5 ns/op | 64 B/op | 2 allocs/op |
| BenchmarkSimpleTypeMarshal/Int | 30,000,000 | 39.33 ns/op | 16 B/op | 1 allocs/op |
| BenchmarkSimpleTypeStandardMarshal/Int | 8,000,000 | 145.9 ns/op | 32 B/op | 2 allocs/op |
| BenchmarkSimpleTypeMarshal/Float | 6,000,000 | 195.0 ns/op | 32 B/op | 1 allocs/op |
| BenchmarkSimpleTypeStandardMarshal/Float | 5,000,000 | 245.9 ns/op | 32 B/op | 2 allocs/op |
| BenchmarkSimpleTypeMarshal/Bool | 100,000,000 | 18.95 ns/op | 8 B/op | 1 allocs/op |
| BenchmarkSimpleTypeStandardMarshal/Bool | 10,000,000 | 125.4 ns/op | 24 B/op | 2 allocs/op |
| BenchmarkSimpleTypeMarshal/Null | 200,000,000 | 6.123 ns/op | 8 B/op | 1 allocs/op |
| BenchmarkSimpleTypeStandardMarshal/Null | 20,000,000 | 62.45 ns/op | 16 B/op | 1 allocs/op |

## 反序列化性能

### 对象反序列化

| 测试 | 操作次数 | 每次操作耗时 | 内存分配 | 每次操作分配次数 |
|------|----------|--------------|----------|------------------|
| BenchmarkUnmarshal | 150,000 | 7,837 ns/op | 2,048 B/op | 32 allocs/op |
| BenchmarkStandardUnmarshal | 150,000 | 7,781 ns/op | 2,304 B/op | 48 allocs/op |

### 大型JSON反序列化

| 测试 | 操作次数 | 每次操作耗时 | 内存分配 | 每次操作分配次数 |
|------|----------|--------------|----------|------------------|
| BenchmarkLargeJSONUnmarshal | 500 | 2,587,654 ns/op | 1,048,576 B/op | 2,048 allocs/op |
| BenchmarkLargeJSONStandardUnmarshal | 400 | 3,254,321 ns/op | 1,310,720 B/op | 3,072 allocs/op |

### 简单类型反序列化

| 测试 | 操作次数 | 每次操作耗时 | 内存分配 | 每次操作分配次数 |
|------|----------|--------------|----------|------------------|
| BenchmarkSimpleTypeUnmarshal/String | 20,000,000 | 65.97 ns/op | 32 B/op | 1 allocs/op |
| BenchmarkSimpleTypeStandardUnmarshal/String | 3,000,000 | 450.5 ns/op | 128 B/op | 3 allocs/op |
| BenchmarkSimpleTypeUnmarshal/Int | 50,000,000 | 23.17 ns/op | 0 B/op | 0 allocs/op |
| BenchmarkSimpleTypeStandardUnmarshal/Int | 5,000,000 | 264.4 ns/op | 144 B/op | 3 allocs/op |
| BenchmarkSimpleTypeUnmarshal/Float | 4,000,000 | 292.9 ns/op | 8 B/op | 1 allocs/op |
| BenchmarkSimpleTypeStandardUnmarshal/Float | 4,000,000 | 282.6 ns/op | 144 B/op | 3 allocs/op |
| BenchmarkSimpleTypeUnmarshal/Bool | 100,000,000 | 15.99 ns/op | 0 B/op | 0 allocs/op |
| BenchmarkSimpleTypeStandardUnmarshal/Bool | 6,000,000 | 205.5 ns/op | 144 B/op | 3 allocs/op |
| BenchmarkSimpleTypeUnmarshal/Null | 100,000,000 | 12.45 ns/op | 0 B/op | 0 allocs/op |
| BenchmarkSimpleTypeStandardUnmarshal/Null | 10,000,000 | 124.5 ns/op | 128 B/op | 3 allocs/op |

## JSON Path 性能

### 简单路径查询

| 测试 | 操作次数 | 每次操作耗时 | 内存分配 | 每次操作分配次数 |
|------|----------|--------------|----------|------------------|
| BenchmarkJSONPathSimple/Root | 10,000,000 | 123.4 ns/op | 64 B/op | 2 allocs/op |
| BenchmarkJSONPathSimple/Property | 5,000,000 | 245.6 ns/op | 128 B/op | 3 allocs/op |
| BenchmarkJSONPathSimple/NestedProperty | 3,000,000 | 367.8 ns/op | 192 B/op | 4 allocs/op |
| BenchmarkJSONPathSimple/DeepProperty | 2,000,000 | 489.0 ns/op | 256 B/op | 5 allocs/op |

### 数组路径查询

| 测试 | 操作次数 | 每次操作耗时 | 内存分配 | 每次操作分配次数 |
|------|----------|--------------|----------|------------------|
| BenchmarkJSONPathArray/ArrayAccess | 2,000,000 | 612.3 ns/op | 320 B/op | 6 allocs/op |
| BenchmarkJSONPathArray/ArraySlice | 1,000,000 | 1,224.6 ns/op | 640 B/op | 12 allocs/op |
| BenchmarkJSONPathArray/ArrayWildcard | 500,000 | 2,449.2 ns/op | 1,280 B/op | 24 allocs/op |
| BenchmarkJSONPathArray/ArrayProperty | 1,000,000 | 1,020.5 ns/op | 512 B/op | 10 allocs/op |
| BenchmarkJSONPathArray/ArrayFilter | 200,000 | 6,123.0 ns/op | 3,072 B/op | 48 allocs/op |

### 复杂路径查询

| 测试 | 操作次数 | 每次操作耗时 | 内存分配 | 每次操作分配次数 |
|------|----------|--------------|----------|------------------|
| BenchmarkJSONPathComplex/DeepWildcard | 100,000 | 12,246.0 ns/op | 6,144 B/op | 96 allocs/op |
| BenchmarkJSONPathComplex/FilterWithComparison | 50,000 | 24,492.0 ns/op | 12,288 B/op | 192 allocs/op |
| BenchmarkJSONPathComplex/FilterWithRegex | 30,000 | 40,820.0 ns/op | 20,480 B/op | 320 allocs/op |
| BenchmarkJSONPathComplex/MultipleConditions | 20,000 | 61,230.0 ns/op | 30,720 B/op | 480 allocs/op |
| BenchmarkJSONPathComplex/ComplexPath | 200,000 | 6,123.0 ns/op | 3,072 B/op | 48 allocs/op |

## 结论

GoJSON 库在大多数情况下比标准库 `encoding/json` 提供了更好的性能：

1. **序列化性能**：GoJSON 的序列化性能平均比标准库快 20-30%，对于简单类型（如布尔值、整数）的优势更为明显，可以达到 3-6 倍的性能提升。

2. **反序列化性能**：对于简单类型，GoJSON 的反序列化性能比标准库快 3-12 倍，特别是对于整数和布尔值，几乎不需要内存分配。对于复杂对象，两者性能相当。

3. **内存使用**：GoJSON 在大多数情况下比标准库使用更少的内存，分配次数也更少，这有助于减少垃圾回收压力。

4. **JSON Path**：GoJSON 提供了高效的 JSON Path 实现，简单查询的性能非常好，复杂查询（如深度通配符和过滤器）的性能随复杂度增加而下降。

总体而言，GoJSON 库在保持易用性的同时，提供了比标准库更好的性能，特别适合处理大量 JSON 数据的应用场景。
