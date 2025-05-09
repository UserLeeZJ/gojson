# GoJSON 命令行工具

本目录包含 GoJSON 库的命令行工具，提供了一系列实用功能，方便用户在命令行中处理 JSON 数据。

## 工具列表

1. **jsonformat** - JSON 格式化工具
2. **jsonpath** - JSON Path 查询工具
3. **jsonanalyze** - JSON 结构分析工具
4. **jsonstream** - JSON 流式处理工具

## 安装

要安装所有工具，请使用以下命令：

```bash
go install github.com/UserLeeZJ/gojson/cmd/...@latest
```

或者安装单个工具：

```bash
go install github.com/UserLeeZJ/gojson/cmd/jsonformat@latest
```

## 使用说明

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

# 设置缩进
jsonformat -i input.json -o output.json -p -indent "    "
```

### jsonpath

JSON Path 查询工具，用于从 JSON 中提取数据。

```bash
# 使用 JSON Path 查询
jsonpath -i input.json -p "$.store.book[0].title"

# 从标准输入读取
cat input.json | jsonpath -p "$.store.book[*].author"

# 输出为紧凑格式
jsonpath -i input.json -p "$.store.book[?(@.price < 10)]" -c

# 输出为美化格式
jsonpath -i input.json -p "$.store.book[*]" -pretty
```

### jsonanalyze

JSON 结构分析工具，用于分析 JSON 的结构。

```bash
# 分析 JSON 结构
jsonanalyze -i input.json

# 从标准输入读取
cat input.json | jsonanalyze

# 输出所有 JSON Path
jsonanalyze -i input.json -paths

# 分析特定路径的结构
jsonanalyze -i input.json -p "$.store.book"
```

### jsonstream

JSON 流式处理工具，用于处理大型 JSON 文件。

```bash
# 流式处理大型 JSON 文件
jsonstream -i large.json -o output.json -f "$.items[*].name"

# 从标准输入读取
cat large.json | jsonstream -f "$.items[*]" > output.json

# 限制输出数量
jsonstream -i large.json -f "$.items[*]" -limit 10

# 过滤数据
jsonstream -i large.json -f "$.items[*]" -filter "price > 100"
```

## 示例

### 格式化 JSON

```bash
echo '{"name":"John","age":30,"city":"New York"}' | jsonformat -p
```

输出：

```json
{
  "name": "John",
  "age": 30,
  "city": "New York"
}
```

### 查询 JSON

```bash
echo '{"store":{"book":[{"title":"Book 1","price":10},{"title":"Book 2","price":20}]}}' | jsonpath -p "$.store.book[*].title"
```

输出：

```json
["Book 1", "Book 2"]
```

### 分析 JSON

```bash
echo '{"name":"John","age":30,"address":{"city":"New York","country":"USA"},"hobbies":["reading","swimming"]}' | jsonanalyze
```

输出：

```text
类型: object
大小: 4
最大深度: 3
键数: 4
子元素类型:
  name: string
  age: number
  address: object
  hobbies: array
值类型统计:
  string: 3
  number: 1
  object: 1
  array: 1
```

### 流式处理 JSON

```bash
cat large.json | jsonstream -f "$.items[*].name" | head -5
```

输出前 5 个商品名称。
