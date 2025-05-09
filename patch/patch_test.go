package patch

import (
	"strings"
	"testing"

	"github.com/UserLeeZJ/gojson/types"
)

func TestApplyPatch(t *testing.T) {
	// 创建测试对象
	createTestObject := func() types.JSONValue {
		obj := types.NewJSONObject()
		obj.PutString("name", "John")
		obj.PutNumber("age", 30)
		address := types.NewJSONObject()
		address.PutString("city", "New York")
		obj.PutObject("address", address)
		return obj
	}

	tests := []struct {
		name      string
		patchJSON string
		wantErr   bool
		debug     bool
		check     func(o *types.JSONObject) bool
	}{
		{
			name:      "添加属性",
			patchJSON: `[{"op":"add","path":"/email","value":"john@example.com"}]`,
			wantErr:   false,
			check: func(o *types.JSONObject) bool {
				email, err := o.GetString("email")
				return err == nil && email == "john@example.com"
			},
		},
		{
			name:      "移除属性",
			patchJSON: `[{"op":"remove","path":"/age"}]`,
			wantErr:   false,
			check: func(o *types.JSONObject) bool {
				return !o.Has("age")
			},
		},
		{
			name:      "替换属性",
			patchJSON: `[{"op":"replace","path":"/name","value":"Jane"}]`,
			wantErr:   false,
			check: func(o *types.JSONObject) bool {
				name, err := o.GetString("name")
				return err == nil && name == "Jane"
			},
		},
		{
			name:      "移动属性",
			patchJSON: `[{"op":"move","from":"/name","path":"/fullName"}]`,
			wantErr:   false,
			check: func(o *types.JSONObject) bool {
				return !o.Has("name") && o.Has("fullName")
			},
		},
		{
			name:      "复制属性",
			patchJSON: `[{"op":"copy","from":"/name","path":"/fullName"}]`,
			wantErr:   false,
			check: func(o *types.JSONObject) bool {
				name, _ := o.GetString("name")
				fullName, _ := o.GetString("fullName")
				return name == fullName
			},
		},
		{
			name:      "测试属性",
			patchJSON: `[{"op":"test","path":"/name","value":"John"}]`,
			wantErr:   false,
			check: func(o *types.JSONObject) bool {
				return true // 如果测试通过，对象不会改变
			},
		},
		{
			name:      "测试失败",
			patchJSON: `[{"op":"test","path":"/name","value":"Jane"}]`,
			wantErr:   true,
			check: func(o *types.JSONObject) bool {
				return true // 不会执行到这里，因为会返回错误
			},
		},
		{
			name:      "组合操作",
			patchJSON: `[{"op":"add","path":"/email","value":"john@example.com"},{"op":"remove","path":"/age"},{"op":"replace","path":"/name","value":"Jane"}]`,
			wantErr:   false,
			check: func(o *types.JSONObject) bool {
				name, _ := o.GetString("name")
				email, _ := o.GetString("email")
				return name == "Jane" && email == "john@example.com" && !o.Has("age")
			},
		},
		{
			name:      "嵌套属性",
			patchJSON: `[{"op":"add","path":"/address/zipcode","value":"10001"}]`,
			wantErr:   false,
			check: func(o *types.JSONObject) bool {
				address, _ := o.GetObject("address")
				zipcode, err := address.GetString("zipcode")
				return err == nil && zipcode == "10001"
			},
		},
		{
			name:      "无效的路径",
			patchJSON: `[{"op":"add","path":"/nonexistent/field","value":"test"}]`,
			wantErr:   true,
			check: func(o *types.JSONObject) bool {
				// 检查结果中是否不包含"nonexistent"字符串
				jsonStr := o.String()
				return !strings.Contains(jsonStr, "nonexistent/field")
			},
		},
		{
			name:      "无效的JSON",
			patchJSON: `invalid json`,
			wantErr:   true,
			check: func(o *types.JSONObject) bool {
				return true // 不会执行到这里，因为会返回错误
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 为每个测试用例创建新的对象，而不是克隆
			testObj := createTestObject()

			// 如果是调试模式，打印原始对象
			if tt.debug {
				t.Logf("Test: %s", tt.name)
				t.Logf("Original Object: %s", testObj.String())
				t.Logf("Patch JSON: %s", tt.patchJSON)
			}

			// 应用补丁
			result, err := ApplyPatch(testObj, tt.patchJSON)

			// 如果是调试模式，打印更多信息
			if tt.debug {
				t.Logf("Error: %v", err)
				if result != nil {
					t.Logf("Result: %s", result.String())
				} else {
					t.Logf("Result is nil")
				}
			}

			// 检查错误
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyPatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 如果期望成功，检查结果
			if !tt.wantErr {
				obj, err := result.AsObject()
				if err != nil {
					t.Errorf("Result is not an object: %v", err)
					return
				}

				if !tt.check(obj) {
					t.Errorf("Result check failed: %s", obj.String())
				}
			}
		})
	}
}
