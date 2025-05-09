/*
Package gojson 提供类似JavaScript JSON接口的Go JSON工具库。

该库专注于提供面向对象的JSON操作体验，同时保持高性能和可靠性。
它提供了一组类型和函数，使JSON处理更加直观和简单。

基本类型

gojson定义了以下基本类型：

- JSONValue: 所有JSON值类型的通用接口
- JSONObject: 表示JSON中的对象
- JSONArray: 表示JSON中的数组
- JSONString: 表示JSON中的字符串
- JSONNumber: 表示JSON中的数字
- JSONBool: 表示JSON中的布尔值
- JSONNull: 表示JSON中的null值

基本用法

创建和操作JSON对象：

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

解析JSON字符串：

	jsonStr := `{"name":"张三","age":28,"address":{"city":"北京"}}`
	value, err := gojson.ParseToValue(jsonStr)
	if err != nil {
		// 处理错误
	}
	
	// 转换为对象
	obj, err := value.AsObject()
	if err != nil {
		// 处理错误
	}
	
	// 获取属性
	name, err := obj.GetString("name")
	age, err := obj.GetNumber("age")

JSON Path

gojson支持JSON Path查询，可以从复杂的JSON结构中提取数据：

	// 使用JSON Path查询
	results, err := gojson.QueryJSONPath(jsonValue, "$.store.book[*].author")
	if err != nil {
		// 处理错误
	}
	
	// 处理结果
	for _, result := range results {
		author, _ := result.AsString()
		fmt.Println(author)
	}

支持的JSON Path语法：

- $: 根节点
- .property: 属性访问
- ['property']: 属性访问（替代语法）
- [index]: 数组索引访问
- [start:end]: 数组切片
- [*]: 通配符，匹配所有元素
- [?(@.property == value)]: 过滤器表达式

JSON Diff

gojson提供了比较两个JSON结构差异的功能：

	oldJSON := `{"name":"张三","age":30}`
	newJSON := `{"name":"张三","age":31,"email":"zhangsan@example.com"}`
	
	// 比较差异
	diffs, err := gojson.DiffJSONStrings(oldJSON, newJSON, nil)
	if err != nil {
		// 处理错误
	}
	
	// 处理差异
	for _, diff := range diffs {
		fmt.Printf("路径: %s, 类型: %s\n", diff.Path, diff.Type)
		if diff.Type == gojson.DiffModified {
			fmt.Printf("旧值: %s, 新值: %s\n", diff.OldValue, diff.NewValue)
		}
	}

差异类型：

- DiffAdded: 添加了新值
- DiffRemoved: 移除了值
- DiffModified: 修改了值
- DiffSame: 值相同
- DiffMoved: 值被移动（数组中）
- DiffTypeChanged: 类型改变

JSON Patch

gojson支持根据差异生成JSON Patch（RFC 6902）：

	// 生成JSON Patch
	patch := gojson.GeneratePatch(diffs)
	
	// 输出Patch
	fmt.Println(patch.String())

性能优化

gojson提供了优化的序列化和反序列化函数：

	// 使用优化的序列化函数
	jsonBytes, err := gojson.FastMarshal(data)
	if err != nil {
		// 处理错误
	}
	
	// 使用优化的反序列化函数
	var result interface{}
	err = gojson.FastUnmarshal(jsonBytes, &result)
	if err != nil {
		// 处理错误
	}

错误处理

gojson使用结构化的错误处理系统：

	_, err := gojson.ParseToValue(jsonStr)
	if err != nil {
		// 使用类型断言获取详细错误信息
		if jsonErr, ok := err.(*gojson.JSONError); ok {
			fmt.Printf("错误代码: %s\n", jsonErr.Code)
			fmt.Printf("错误消息: %s\n", jsonErr.Message)
			if jsonErr.Cause != nil {
				fmt.Printf("原始错误: %v\n", jsonErr.Cause)
			}
		}
	}
*/
package gojson
