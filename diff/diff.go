// Package diff 提供gojson库的JSON差异比较功能
package diff

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/UserLeeZJ/gojson/parser"
	"github.com/UserLeeZJ/gojson/types"
)

// DiffType 表示差异的类型
type DiffType string

const (
	DiffAdded       DiffType = "added"        // 添加了新值
	DiffRemoved     DiffType = "removed"      // 移除了值
	DiffModified    DiffType = "modified"     // 修改了值
	DiffSame        DiffType = "same"         // 值相同
	DiffMoved       DiffType = "moved"        // 值被移动（数组中）
	DiffTypeChanged DiffType = "type_changed" // 类型改变
)

// DiffOptions 表示比较选项
type DiffOptions struct {
	IgnoreCase       bool // 忽略字符串大小写
	IgnoreWhitespace bool // 忽略空白字符
	IgnoreOrder      bool // 忽略数组顺序
	IncludeSame      bool // 包含相同的值
	MaxDepth         int  // 最大递归深度，0表示无限制
}

// DefaultDiffOptions 返回默认的比较选项
func DefaultDiffOptions() *DiffOptions {
	return &DiffOptions{
		IgnoreCase:       false,
		IgnoreWhitespace: false,
		IgnoreOrder:      false,
		IncludeSame:      false,
		MaxDepth:         0,
	}
}

// Diff 表示JSON值之间的差异
type Diff struct {
	Type     DiffType        // 差异类型
	Path     string          // 差异路径
	OldValue types.JSONValue // 旧值
	NewValue types.JSONValue // 新值
}

// String 返回差异的字符串表示
func (d *Diff) String() string {
	switch d.Type {
	case DiffAdded:
		return fmt.Sprintf("添加: %s = %s", d.Path, d.NewValue.String())
	case DiffRemoved:
		return fmt.Sprintf("移除: %s = %s", d.Path, d.OldValue.String())
	case DiffModified:
		return fmt.Sprintf("修改: %s = %s -> %s", d.Path, d.OldValue.String(), d.NewValue.String())
	case DiffSame:
		return fmt.Sprintf("相同: %s = %s", d.Path, d.OldValue.String())
	case DiffMoved:
		return fmt.Sprintf("移动: %s", d.Path)
	case DiffTypeChanged:
		return fmt.Sprintf("类型改变: %s = %s -> %s", d.Path, d.OldValue.Type(), d.NewValue.Type())
	default:
		return fmt.Sprintf("未知差异: %s", d.Path)
	}
}

// DiffJSON 比较两个JSON值的差异
func DiffJSON(oldValue, newValue types.JSONValue, options *DiffOptions) ([]*Diff, error) {
	if options == nil {
		options = DefaultDiffOptions()
	}

	diffs := make([]*Diff, 0)
	diffValues("$", oldValue, newValue, options, &diffs, 0)
	return diffs, nil
}

// DiffJSONStrings 比较两个JSON字符串的差异
func DiffJSONStrings(oldJSON, newJSON string, options *DiffOptions) ([]*Diff, error) {
	oldValue, err := parser.ParseToValue(oldJSON)
	if err != nil {
		return nil, err
	}

	newValue, err := parser.ParseToValue(newJSON)
	if err != nil {
		return nil, err
	}

	return DiffJSON(oldValue, newValue, options)
}

// 递归比较两个JSON值的差异
func diffValues(path string, oldValue, newValue types.JSONValue, options *DiffOptions, diffs *[]*Diff, depth int) {
	// 检查最大递归深度
	if options.MaxDepth > 0 && depth > options.MaxDepth {
		return
	}

	// 处理null值
	if oldValue.IsNull() && newValue.IsNull() {
		if options.IncludeSame {
			*diffs = append(*diffs, &Diff{
				Type:     DiffSame,
				Path:     path,
				OldValue: oldValue,
				NewValue: newValue,
			})
		}
		return
	}

	if oldValue.IsNull() {
		*diffs = append(*diffs, &Diff{
			Type:     DiffAdded,
			Path:     path,
			OldValue: oldValue,
			NewValue: newValue,
		})
		return
	}

	if newValue.IsNull() {
		*diffs = append(*diffs, &Diff{
			Type:     DiffRemoved,
			Path:     path,
			OldValue: oldValue,
			NewValue: newValue,
		})
		return
	}

	// 处理类型不同的情况
	if oldValue.Type() != newValue.Type() {
		*diffs = append(*diffs, &Diff{
			Type:     DiffTypeChanged,
			Path:     path,
			OldValue: oldValue,
			NewValue: newValue,
		})
		return
	}

	// 根据类型进行比较
	switch oldValue.Type() {
	case "boolean":
		diffBooleans(path, oldValue, newValue, options, diffs)
	case "number":
		diffNumbers(path, oldValue, newValue, options, diffs)
	case "string":
		diffStrings(path, oldValue, newValue, options, diffs)
	case "array":
		diffArrays(path, oldValue, newValue, options, diffs, depth)
	case "object":
		diffObjects(path, oldValue, newValue, options, diffs, depth)
	}
}

// 比较布尔值
func diffBooleans(path string, oldValue, newValue types.JSONValue, options *DiffOptions, diffs *[]*Diff) {
	oldBool, _ := oldValue.AsBoolean()
	newBool, _ := newValue.AsBoolean()

	if oldBool == newBool {
		if options.IncludeSame {
			*diffs = append(*diffs, &Diff{
				Type:     DiffSame,
				Path:     path,
				OldValue: oldValue,
				NewValue: newValue,
			})
		}
	} else {
		*diffs = append(*diffs, &Diff{
			Type:     DiffModified,
			Path:     path,
			OldValue: oldValue,
			NewValue: newValue,
		})
	}
}

// 比较数字
func diffNumbers(path string, oldValue, newValue types.JSONValue, options *DiffOptions, diffs *[]*Diff) {
	oldNum, _ := oldValue.AsNumber()
	newNum, _ := newValue.AsNumber()

	if oldNum == newNum {
		if options.IncludeSame {
			*diffs = append(*diffs, &Diff{
				Type:     DiffSame,
				Path:     path,
				OldValue: oldValue,
				NewValue: newValue,
			})
		}
	} else {
		*diffs = append(*diffs, &Diff{
			Type:     DiffModified,
			Path:     path,
			OldValue: oldValue,
			NewValue: newValue,
		})
	}
}

// 比较字符串
func diffStrings(path string, oldValue, newValue types.JSONValue, options *DiffOptions, diffs *[]*Diff) {
	oldStr, _ := oldValue.AsString()
	newStr, _ := newValue.AsString()

	// 应用选项
	if options.IgnoreCase {
		oldStr = strings.ToLower(oldStr)
		newStr = strings.ToLower(newStr)
	}

	if options.IgnoreWhitespace {
		oldStr = removeWhitespace(oldStr)
		newStr = removeWhitespace(newStr)
	}

	if oldStr == newStr {
		if options.IncludeSame {
			*diffs = append(*diffs, &Diff{
				Type:     DiffSame,
				Path:     path,
				OldValue: oldValue,
				NewValue: newValue,
			})
		}
	} else {
		*diffs = append(*diffs, &Diff{
			Type:     DiffModified,
			Path:     path,
			OldValue: oldValue,
			NewValue: newValue,
		})
	}
}

// 移除字符串中的空白字符
func removeWhitespace(s string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(s, "")
}

// 比较数组
func diffArrays(path string, oldValue, newValue types.JSONValue, options *DiffOptions, diffs *[]*Diff, depth int) {
	oldArr, _ := oldValue.AsArray()
	newArr, _ := newValue.AsArray()

	if options.IgnoreOrder {
		// 忽略顺序时，将数组视为集合进行比较
		diffArraysAsSet(path, oldArr, newArr, options, diffs, depth)
	} else {
		// 保持顺序时，按索引比较
		diffArraysInOrder(path, oldArr, newArr, options, diffs, depth)
	}
}

// 按顺序比较数组
func diffArraysInOrder(path string, oldArr, newArr *types.JSONArray, options *DiffOptions, diffs *[]*Diff, depth int) {
	maxLen := oldArr.Size()
	if newArr.Size() > maxLen {
		maxLen = newArr.Size()
	}

	for i := 0; i < maxLen; i++ {
		itemPath := fmt.Sprintf("%s[%d]", path, i)

		if i >= oldArr.Size() {
			// 新数组中添加的元素
			*diffs = append(*diffs, &Diff{
				Type:     DiffAdded,
				Path:     itemPath,
				OldValue: types.NewJSONNull(),
				NewValue: newArr.Get(i),
			})
		} else if i >= newArr.Size() {
			// 旧数组中移除的元素
			*diffs = append(*diffs, &Diff{
				Type:     DiffRemoved,
				Path:     itemPath,
				OldValue: oldArr.Get(i),
				NewValue: types.NewJSONNull(),
			})
		} else {
			// 比较相同位置的元素
			diffValues(itemPath, oldArr.Get(i), newArr.Get(i), options, diffs, depth+1)
		}
	}
}

// 将数组视为集合进行比较
func diffArraysAsSet(path string, oldArr, newArr *types.JSONArray, options *DiffOptions, diffs *[]*Diff, depth int) {
	// TODO: 实现将数组视为集合的比较逻辑
	// 这需要一个复杂的算法来匹配最相似的元素
	// 简化起见，这里仍然使用按顺序比较
	diffArraysInOrder(path, oldArr, newArr, options, diffs, depth)
}

// 比较对象
func diffObjects(path string, oldValue, newValue types.JSONValue, options *DiffOptions, diffs *[]*Diff, depth int) {
	oldObj, _ := oldValue.AsObject()
	newObj, _ := newValue.AsObject()

	// 获取所有键
	oldKeys := oldObj.Keys()
	newKeys := newObj.Keys()
	allKeys := mergeKeys(oldKeys, newKeys)

	// 比较每个键
	for _, key := range allKeys {
		propPath := path
		if path == "$" {
			propPath = "$." + key
		} else {
			// 处理键中的特殊字符
			if isValidIdentifier(key) {
				propPath = path + "." + key
			} else {
				propPath = path + "['" + key + "']"
			}
		}

		oldHas := oldObj.Has(key)
		newHas := newObj.Has(key)

		if oldHas && newHas {
			// 两个对象都有该键，比较值
			diffValues(propPath, oldObj.Get(key), newObj.Get(key), options, diffs, depth+1)
		} else if oldHas {
			// 只有旧对象有该键，表示移除
			*diffs = append(*diffs, &Diff{
				Type:     DiffRemoved,
				Path:     propPath,
				OldValue: oldObj.Get(key),
				NewValue: types.NewJSONNull(),
			})
		} else {
			// 只有新对象有该键，表示添加
			*diffs = append(*diffs, &Diff{
				Type:     DiffAdded,
				Path:     propPath,
				OldValue: types.NewJSONNull(),
				NewValue: newObj.Get(key),
			})
		}
	}
}

// 合并两个字符串切片，去除重复项
func mergeKeys(a, b []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(a)+len(b))

	for _, key := range a {
		if !seen[key] {
			seen[key] = true
			result = append(result, key)
		}
	}

	for _, key := range b {
		if !seen[key] {
			seen[key] = true
			result = append(result, key)
		}
	}

	sort.Strings(result)
	return result
}

// 检查字符串是否为有效的标识符
func isValidIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}
	if !isLetter(s[0]) && s[0] != '_' {
		return false
	}
	for i := 1; i < len(s); i++ {
		if !isLetter(s[i]) && !isDigit(s[i]) && s[i] != '_' {
			return false
		}
	}
	return true
}

// 检查字符是否为字母
func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// 检查字符是否为数字
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// GeneratePatch 从差异生成JSON Patch
func GeneratePatch(diffs []*Diff) *types.JSONArray {
	patch := types.NewJSONArray()

	for _, d := range diffs {
		switch d.Type {
		case DiffAdded:
			op := types.NewJSONObject()
			op.PutString("op", "add")
			op.PutString("path", jsonPathToPatchPath(d.Path))
			op.Put("value", d.NewValue)
			patch.Add(op)
		case DiffRemoved:
			op := types.NewJSONObject()
			op.PutString("op", "remove")
			op.PutString("path", jsonPathToPatchPath(d.Path))
			patch.Add(op)
		case DiffModified:
			op := types.NewJSONObject()
			op.PutString("op", "replace")
			op.PutString("path", jsonPathToPatchPath(d.Path))
			op.Put("value", d.NewValue)
			patch.Add(op)
		}
	}

	return patch
}

// 将JSON Path转换为JSON Patch路径
func jsonPathToPatchPath(path string) string {
	if path == "$" {
		return ""
	}

	// 移除开头的$
	path = path[1:]

	// 替换.为/
	result := strings.ReplaceAll(path, ".", "/")

	// 处理数组索引
	result = strings.ReplaceAll(result, "[", "/")
	result = strings.ReplaceAll(result, "]", "")

	// 处理转义字符
	result = strings.ReplaceAll(result, "~", "~0")
	result = strings.ReplaceAll(result, "/", "~1")

	return "/" + result
}
