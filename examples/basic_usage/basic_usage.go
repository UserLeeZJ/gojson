// 基本用法示例
package main

import (
	"fmt"

	"github.com/UserLeeZJ/gojson/types"
)

func main() {
	fmt.Println("GoJSON 基本用法示例")
	fmt.Println("====================")

	// 创建JSON对象
	fmt.Println("\n1. 创建JSON对象")
	person := types.NewJSONObject()
	person.PutString("name", "张三")
	person.PutNumber("age", 28)
	person.PutBoolean("active", true)
	person.PutNull("data")

	// 输出JSON对象
	fmt.Println("创建的JSON对象:")
	fmt.Println(person.String())

	// 添加嵌套对象
	fmt.Println("\n2. 添加嵌套对象")
	address := types.NewJSONObject()
	address.PutString("city", "北京")
	address.PutString("district", "海淀区")
	address.PutString("zipcode", "100000")
	person.PutObject("address", address)

	// 输出更新后的JSON对象
	fmt.Println("添加嵌套对象后:")
	fmt.Println(person.String())

	// 添加数组
	fmt.Println("\n3. 添加数组")
	hobbies := types.NewJSONArray()
	hobbies.AddString("阅读").AddString("编程").AddString("旅行")
	person.PutArray("hobbies", hobbies)

	// 输出更新后的JSON对象
	fmt.Println("添加数组后:")
	fmt.Println(person.String())

	// 获取值
	fmt.Println("\n4. 获取值")
	name, _ := person.GetString("name")
	fmt.Println("姓名:", name)

	age, _ := person.GetNumber("age")
	fmt.Println("年龄:", age)

	active, _ := person.GetBoolean("active")
	fmt.Println("是否活跃:", active)

	// 获取嵌套对象的值
	fmt.Println("\n5. 获取嵌套对象的值")
	addressObj, _ := person.GetObject("address")
	city, _ := addressObj.GetString("city")
	fmt.Println("城市:", city)

	// 获取数组的值
	fmt.Println("\n6. 获取数组的值")
	hobbiesArr, _ := person.GetArray("hobbies")
	fmt.Println("爱好数量:", hobbiesArr.Size())
	firstHobby, _ := hobbiesArr.GetString(0)
	fmt.Println("第一个爱好:", firstHobby)

	// 遍历数组
	fmt.Println("\n7. 遍历数组")
	fmt.Println("所有爱好:")
	for i := 0; i < hobbiesArr.Size(); i++ {
		hobby, _ := hobbiesArr.GetString(i)
		fmt.Printf("  %d. %s\n", i+1, hobby)
	}

	// 修改值
	fmt.Println("\n8. 修改值")
	person.PutNumber("age", 29)
	addressObj.PutString("district", "朝阳区")
	hobbiesArr.Set(2, types.NewJSONString("摄影"))

	// 输出修改后的JSON对象
	fmt.Println("修改后:")
	fmt.Println(person.String())

	// 删除值
	fmt.Println("\n9. 删除值")
	person.Remove("data")
	addressObj.Remove("zipcode")

	// 输出删除后的JSON对象
	fmt.Println("删除后:")
	fmt.Println(person.String())

	// 检查键是否存在
	fmt.Println("\n10. 检查键是否存在")
	fmt.Println("'name'键是否存在:", person.Has("name"))
	fmt.Println("'data'键是否存在:", person.Has("data"))

	// 获取所有键
	fmt.Println("\n11. 获取所有键")
	keys := person.Keys()
	fmt.Println("所有键:", keys)

	// 创建复杂的嵌套结构
	fmt.Println("\n12. 创建复杂的嵌套结构")
	company := types.NewJSONObject()
	company.PutString("name", "示例公司")
	company.PutNumber("founded", 2010)

	// 添加员工数组
	employees := types.NewJSONArray()

	employee1 := types.NewJSONObject()
	employee1.PutString("name", "李四")
	employee1.PutNumber("age", 30)
	employees.Add(employee1)

	employee2 := types.NewJSONObject()
	employee2.PutString("name", "王五")
	employee2.PutNumber("age", 25)
	employees.Add(employee2)

	company.PutArray("employees", employees)

	// 添加部门对象
	departments := types.NewJSONObject()

	tech := types.NewJSONArray()
	tech.AddString("开发").AddString("测试").AddString("运维")
	departments.PutArray("技术部", tech)

	sales := types.NewJSONArray()
	sales.AddString("国内销售").AddString("海外销售")
	departments.PutArray("销售部", sales)

	company.PutObject("departments", departments)

	// 输出复杂结构
	fmt.Println("复杂的嵌套结构:")
	fmt.Println(company.String())

	fmt.Println("\n示例结束")
}
