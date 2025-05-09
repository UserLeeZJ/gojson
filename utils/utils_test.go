package utils

import (
	"encoding/json"
	"sort"
	"strings"
	"testing"

	"github.com/UserLeeZJ/gojson/parser"
	"github.com/UserLeeZJ/gojson/types"
)

func TestPrettyPrint(t *testing.T) {
	// 创建测试对象
	obj := types.NewJSONObject()
	obj.PutString("name", "John")
	obj.PutNumber("age", 30)
	
	// 创建嵌套对象
	address := types.NewJSONObject()
	address.PutString("city", "New York")
	address.PutString("country", "USA")
	obj.PutObject("address", address)
	
	// 创建数组
	hobbies := types.NewJSONArray()
	hobbies.AddString("reading").AddString("swimming")
	obj.PutArray("hobbies", hobbies)
	
	// 测试默认选项
	pretty, err := PrettyPrint(obj, DefaultPrettyOptions())
	if err != nil {
		t.Errorf("PrettyPrint失败: %v", err)
	}
	
	// 验证结果包含预期的格式
	if !strings.Contains(pretty, "{\n") {
		t.Errorf("美化后的JSON不包含换行符")
	}
	if !strings.Contains(pretty, "  \"name\"") {
		t.Errorf("美化后的JSON不包含缩进")
	}
	
	// 测试自定义缩进
	options := DefaultPrettyOptions()
	options.Indent = "    " // 四个空格
	pretty, err = PrettyPrint(obj, options)
	if err != nil {
		t.Errorf("PrettyPrint失败: %v", err)
	}
	
	// 验证结果包含自定义缩进
	if !strings.Contains(pretty, "    \"name\"") {
		t.Errorf("美化后的JSON不包含自定义缩进")
	}
	
	// 测试排序键
	options = DefaultPrettyOptions()
	options.SortKeys = true
	pretty, err = PrettyPrint(obj, options)
	if err != nil {
		t.Errorf("PrettyPrint失败: %v", err)
	}
	
	// 验证键是否已排序
	// 注意：这个测试可能不够严格，因为JSON对象的键顺序在标准库中可能已经是排序的
	firstNameIndex := strings.Index(pretty, "\"address\"")
	firstAgeIndex := strings.Index(pretty, "\"age\"")
	firstHobbiesIndex := strings.Index(pretty, "\"hobbies\"")
	firstNameIndex2 := strings.Index(pretty, "\"name\"")
	
	if firstNameIndex > firstAgeIndex || firstAgeIndex > firstHobbiesIndex || firstHobbiesIndex > firstNameIndex2 {
		t.Errorf("键未按字母顺序排序")
	}
}

func TestCompressJSON(t *testing.T) {
	// 创建测试对象
	obj := types.NewJSONObject()
	obj.PutString("name", "John")
	obj.PutNumber("age", 30)
	
	// 创建嵌套对象
	address := types.NewJSONObject()
	address.PutString("city", "New York")
	address.PutString("country", "USA")
	obj.PutObject("address", address)
	
	// 创建数组
	hobbies := types.NewJSONArray()
	hobbies.AddString("reading").AddString("swimming")
	obj.PutArray("hobbies", hobbies)
	
	// 测试压缩
	compressed, err := CompressJSON(obj)
	if err != nil {
		t.Errorf("CompressJSON失败: %v", err)
	}
	
	// 验证结果不包含空白字符
	if strings.Contains(compressed, "\n") || strings.Contains(compressed, "  ") {
		t.Errorf("压缩后的JSON包含不必要的空白字符")
	}
	
	// 验证结果是有效的JSON
	var result interface{}
	if err := json.Unmarshal([]byte(compressed), &result); err != nil {
		t.Errorf("压缩后的JSON无效: %v", err)
	}
}

func TestExtractPaths(t *testing.T) {
	// 创建测试JSON
	jsonStr := `{
		"name": "John",
		"age": 30,
		"address": {
			"city": "New York",
			"country": "USA"
		},
		"hobbies": ["reading", "swimming"],
		"special.key": "value"
	}`
	
	// 解析JSON
	value, err := parser.ParseToValue(jsonStr)
	if err != nil {
		t.Fatalf("解析JSON失败: %v", err)
	}
	
	// 提取路径
	paths := ExtractPaths(value)
	
	// 验证结果
	expectedPaths := []string{
		"$",
		"$.name",
		"$.age",
		"$.address",
		"$.address.city",
		"$.address.country",
		"$.hobbies",
		"$.hobbies[0]",
		"$.hobbies[1]",
		"$['special.key']",
	}
	
	// 排序路径以便比较
	sort.Strings(paths)
	sort.Strings(expectedPaths)
	
	// 检查路径数量
	if len(paths) != len(expectedPaths) {
		t.Errorf("路径数量不匹配: 期望 %d, 实际 %d", len(expectedPaths), len(paths))
	}
	
	// 检查每个路径
	for i, path := range expectedPaths {
		if i >= len(paths) {
			t.Errorf("缺少路径: %s", path)
			continue
		}
		if paths[i] != path {
			t.Errorf("路径不匹配: 期望 %s, 实际 %s", path, paths[i])
		}
	}
}

func TestAnalyzeStructure(t *testing.T) {
	// 创建测试JSON
	jsonStr := `{
		"name": "John",
		"age": 30,
		"address": {
			"city": "New York",
			"country": "USA"
		},
		"hobbies": ["reading", "swimming"],
		"active": true,
		"data": null
	}`
	
	// 解析JSON
	value, err := parser.ParseToValue(jsonStr)
	if err != nil {
		t.Fatalf("解析JSON失败: %v", err)
	}
	
	// 分析结构
	info := AnalyzeStructure(value)
	
	// 验证结果
	if info.Type != "object" {
		t.Errorf("类型不匹配: 期望 object, 实际 %s", info.Type)
	}
	
	if info.Size != 6 {
		t.Errorf("大小不匹配: 期望 6, 实际 %d", info.Size)
	}
	
	if info.Depth != 3 {
		t.Errorf("深度不匹配: 期望 3, 实际 %d", info.Depth)
	}
	
	if info.KeyCount != 6 {
		t.Errorf("键数不匹配: 期望 6, 实际 %d", info.KeyCount)
	}
	
	// 检查子元素类型
	expectedChildTypes := map[string]string{
		"name":    "string",
		"age":     "number",
		"address": "object",
		"hobbies": "array",
		"active":  "boolean",
		"data":    "null",
	}
	
	for key, expectedType := range expectedChildTypes {
		actualType, ok := info.ChildTypes[key]
		if !ok {
			t.Errorf("缺少键: %s", key)
			continue
		}
		if actualType != expectedType {
			t.Errorf("键 %s 的类型不匹配: 期望 %s, 实际 %s", key, expectedType, actualType)
		}
	}
	
	// 检查值类型统计
	expectedValueCounts := map[string]int{
		"string":  3, // name + 2 hobbies
		"number":  1, // age
		"object":  1, // address
		"array":   1, // hobbies
		"boolean": 1, // active
		"null":    1, // data
	}
	
	for typeName, expectedCount := range expectedValueCounts {
		actualCount, ok := info.ValueCounts[typeName]
		if !ok && expectedCount > 0 {
			t.Errorf("缺少类型统计: %s", typeName)
			continue
		}
		if actualCount != expectedCount {
			t.Errorf("类型 %s 的数量不匹配: 期望 %d, 实际 %d", typeName, expectedCount, actualCount)
		}
	}
	
	// 测试String方法
	infoStr := info.String()
	if !strings.Contains(infoStr, "类型: object") {
		t.Errorf("信息字符串不包含类型")
	}
	if !strings.Contains(infoStr, "大小: 6") {
		t.Errorf("信息字符串不包含大小")
	}
	if !strings.Contains(infoStr, "最大深度: 3") {
		t.Errorf("信息字符串不包含深度")
	}
}
