package diff

import (
	"testing"

	"github.com/UserLeeZJ/gojson/parser"
	"github.com/UserLeeZJ/gojson/types"
)

func TestDiffJSON(t *testing.T) {
	// 测试数据
	oldJSON := `{
		"name": "张三",
		"age": 30,
		"address": {
			"city": "北京",
			"district": "海淀区"
		},
		"hobbies": ["阅读", "编程", "旅行"]
	}`

	newJSON := `{
		"name": "张三",
		"age": 31,
		"address": {
			"city": "上海",
			"district": "浦东新区"
		},
		"hobbies": ["阅读", "编程", "摄影"],
		"email": "zhangsan@example.com"
	}`

	oldValue, err := parser.ParseToValue(oldJSON)
	if err != nil {
		t.Fatalf("解析旧JSON失败: %v", err)
	}

	newValue, err := parser.ParseToValue(newJSON)
	if err != nil {
		t.Fatalf("解析新JSON失败: %v", err)
	}

	// 测试默认选项
	diffs, err := DiffJSON(oldValue, newValue, nil)
	if err != nil {
		t.Fatalf("比较JSON失败: %v", err)
	}

	// 验证差异数量
	if len(diffs) < 4 {
		t.Errorf("差异数量不足: 期望至少 4, 实际 %d", len(diffs))
	}

	// 测试自定义选项
	options := DefaultDiffOptions()
	options.IncludeSame = true
	diffs, err = DiffJSON(oldValue, newValue, options)
	if err != nil {
		t.Fatalf("使用自定义选项比较JSON失败: %v", err)
	}

	// 验证包含相同值的差异数量
	if len(diffs) <= len(oldValue.(*types.JSONObject).Keys()) {
		t.Errorf("包含相同值的差异数量不足: 期望 > %d, 实际 %d", len(oldValue.(*types.JSONObject).Keys()), len(diffs))
	}
}

func TestDiffJSONStrings(t *testing.T) {
	oldJSON := `{"name":"张三","age":30}`
	newJSON := `{"name":"张三","age":31,"email":"zhangsan@example.com"}`

	diffs, err := DiffJSONStrings(oldJSON, newJSON, nil)
	if err != nil {
		t.Fatalf("比较JSON字符串失败: %v", err)
	}

	// 验证差异数量
	if len(diffs) != 2 { // age修改和email添加
		t.Errorf("差异数量不匹配: 期望 2, 实际 %d", len(diffs))
	}

	// 验证差异类型
	foundModified := false
	foundAdded := false
	for _, diff := range diffs {
		if diff.Type == DiffModified && diff.Path == "$.age" {
			foundModified = true
		}
		if diff.Type == DiffAdded && diff.Path == "$.email" {
			foundAdded = true
		}
	}

	if !foundModified {
		t.Errorf("未找到age的修改差异")
	}
	if !foundAdded {
		t.Errorf("未找到email的添加差异")
	}
}

func TestGeneratePatch(t *testing.T) {
	oldJSON := `{"name":"张三","age":30}`
	newJSON := `{"name":"张三","age":31,"email":"zhangsan@example.com"}`

	diffs, _ := DiffJSONStrings(oldJSON, newJSON, nil)
	patchArray := GeneratePatch(diffs)

	// 验证patch数组长度
	if patchArray.Size() != 2 {
		t.Fatalf("Patch操作数量不匹配: 期望 2, 实际 %d", patchArray.Size())
	}

	// 验证第一个操作是replace
	op1, _ := patchArray.GetObject(0)
	opType1, _ := op1.GetString("op")

	if opType1 != "replace" {
		t.Errorf("第一个操作类型不匹配: 期望 op=replace, 实际 op=%s", opType1)
	}

	// 验证第二个操作是add
	op2, _ := patchArray.GetObject(1)
	opType2, _ := op2.GetString("op")

	if opType2 != "add" {
		t.Errorf("第二个操作类型不匹配: 期望 op=add, 实际 op=%s", opType2)
	}
}
