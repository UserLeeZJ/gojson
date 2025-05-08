# gojson

`gojson` 是一个Go语言库，提供类似JavaScript JSON接口的功能，使JSON处理更加简单和直观。

## 特性

- 类似JavaScript的JSON API
- 面向对象的JSON值类型（JSONObject, JSONArray, JSONString等）
- 简单易用的接口
- 全面的错误处理
- 详细的单元测试
- 纯Go实现，无外部依赖

## 安装

```bash
go get github.com/yourusername/gojson
```

## 使用示例

### 解析JSON字符串

```go
package main

import (
    "fmt"
    "github.com/yourusername/gojson"
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
    "github.com/yourusername/gojson"
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
    "github.com/yourusername/gojson"
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
    "github.com/yourusername/gojson"
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
    "github.com/yourusername/gojson"
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

### 其他类型

- JSONBool - 表示JSON中的布尔值
- JSONNull - 表示JSON中的null值
- JSONNumber - 表示JSON中的数字
- JSONString - 表示JSON中的字符串

## 许可证

MIT
