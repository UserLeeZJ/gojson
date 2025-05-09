// Package patch 提供gojson库的JSON Patch功能
package patch

import (
	"encoding/json"
	"fmt"
	"strings"

	jsonerrors "github.com/UserLeeZJ/gojson/errors"
	"github.com/UserLeeZJ/gojson/jsonpath"
	"github.com/UserLeeZJ/gojson/parser"
	"github.com/UserLeeZJ/gojson/types"
)

// PatchOperation 表示JSON Patch操作
type PatchOperation struct {
	Op    string          `json:"op"`
	Path  string          `json:"path"`
	Value json.RawMessage `json:"value,omitempty"`
	From  string          `json:"from,omitempty"`
}

// PatchError 表示JSON Patch操作中的错误
type PatchError struct {
	Operation PatchOperation
	Message   string
}

// Error 实现error接口
func (e *PatchError) Error() string {
	return fmt.Sprintf("JSON Patch 错误: %s, 操作: %+v", e.Message, e.Operation)
}

// ApplyPatch 将JSON Patch应用到JSON值
func ApplyPatch(value types.JSONValue, patchJSON string) (types.JSONValue, error) {
	// 解析补丁
	var patchOps []PatchOperation
	err := json.Unmarshal([]byte(patchJSON), &patchOps)
	if err != nil {
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrInvalidPatch, "无效的JSON Patch").WithCause(err)
	}

	// 克隆原始值
	var result types.JSONValue
	switch value.Type() {
	case "object":
		obj, _ := value.AsObject()
		result = obj.Clone()
	case "array":
		arr, _ := value.AsArray()
		// 创建新数组并复制元素
		newArr := types.NewJSONArray()
		for i := 0; i < arr.Size(); i++ {
			newArr.Add(arr.Get(i))
		}
		result = newArr
	default:
		// 对于基本类型，直接使用原值
		result = value
	}

	// 应用每个操作
	for _, op := range patchOps {
		var err error
		result, err = applyOperation(result, op)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// 应用单个补丁操作
func applyOperation(value types.JSONValue, op PatchOperation) (types.JSONValue, error) {
	// 标准化路径
	path := normalizePath(op.Path)
	from := normalizePath(op.From)

	switch op.Op {
	case "add":
		return applyAddOperation(value, path, op.Value)
	case "remove":
		return applyRemoveOperation(value, path)
	case "replace":
		return applyReplaceOperation(value, path, op.Value)
	case "move":
		return applyMoveOperation(value, from, path)
	case "copy":
		return applyCopyOperation(value, from, path)
	case "test":
		return applyTestOperation(value, path, op.Value)
	default:
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrInvalidPatch, fmt.Sprintf("未知的操作类型: %s", op.Op))
	}
}

// 标准化JSON Patch路径
func normalizePath(path string) string {
	if path == "" {
		return ""
	}
	// 将JSON Patch路径转换为JSON Path格式
	// 例如: /foo/bar -> $.foo.bar
	parts := strings.Split(path, "/")
	if len(parts) > 0 && parts[0] == "" {
		parts = parts[1:]
	}

	result := "$"
	for _, part := range parts {
		// 处理转义字符
		part = strings.ReplaceAll(part, "~1", "/")
		part = strings.ReplaceAll(part, "~0", "~")

		// 检查是否为数组索引
		if isArrayIndex(part) {
			result += "[" + part + "]"
		} else {
			result += "." + part
		}
	}
	return result
}

// 检查字符串是否为数组索引
func isArrayIndex(s string) bool {
	// 检查是否为非负整数
	if s == "0" {
		return true
	}
	if len(s) > 0 && s[0] != '0' {
		for _, c := range s {
			if c < '0' || c > '9' {
				return false
			}
		}
		return true
	}
	return false
}

// 应用add操作
func applyAddOperation(value types.JSONValue, path string, rawValue json.RawMessage) (types.JSONValue, error) {
	// 解析要添加的值
	newValue, err := parser.ParseBytesToValue(rawValue)
	if err != nil {
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrInvalidPatch, "无效的值").WithCause(err)
	}

	// 处理根路径
	if path == "$" {
		return newValue, nil
	}

	// 获取父路径和最后一个段
	parentPath, lastSegment := splitPath(path)

	// 获取父对象或数组
	results, err := jsonpath.QueryJSONPath(value, parentPath)
	if err != nil || len(results) == 0 {
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrPathNotFound, "父路径不存在").WithPath(parentPath)
	}
	parent := results[0]

	// 根据父对象类型处理
	if parent.IsObject() {
		obj, _ := parent.AsObject()
		obj.Put(lastSegment, newValue)
	} else if parent.IsArray() {
		arr, _ := parent.AsArray()
		index, err := parseArrayIndex(lastSegment, arr.Size())
		if err != nil {
			return nil, err
		}
		// 在数组中插入元素
		if index == arr.Size() {
			arr.Add(newValue)
		} else {
			// 创建新数组并重新排列元素
			newArr := types.NewJSONArray()
			for i := 0; i < arr.Size()+1; i++ {
				if i < index {
					newArr.Add(arr.Get(i))
				} else if i == index {
					newArr.Add(newValue)
				} else {
					newArr.Add(arr.Get(i - 1))
				}
			}
			// 替换原数组
			*arr = *newArr
		}
	} else {
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrInvalidType, "父路径必须是对象或数组").WithPath(parentPath)
	}

	return value, nil
}

// 应用remove操作
func applyRemoveOperation(value types.JSONValue, path string) (types.JSONValue, error) {
	// 处理根路径
	if path == "$" {
		return types.NewJSONNull(), nil
	}

	// 获取父路径和最后一个段
	parentPath, lastSegment := splitPath(path)

	// 获取父对象或数组
	results, err := jsonpath.QueryJSONPath(value, parentPath)
	if err != nil || len(results) == 0 {
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrPathNotFound, "父路径不存在").WithPath(parentPath)
	}
	parent := results[0]

	// 根据父对象类型处理
	if parent.IsObject() {
		obj, _ := parent.AsObject()
		if !obj.Has(lastSegment) {
			return nil, jsonerrors.NewJSONError(jsonerrors.ErrPathNotFound, "路径不存在").WithPath(path)
		}
		obj.Remove(lastSegment)
	} else if parent.IsArray() {
		arr, _ := parent.AsArray()
		index, err := parseArrayIndex(lastSegment, arr.Size())
		if err != nil {
			return nil, err
		}
		if index >= arr.Size() {
			return nil, jsonerrors.NewJSONError(jsonerrors.ErrIndexOutOfRange, "索引超出范围").WithPath(path)
		}
		arr.Remove(index)
	} else {
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrInvalidType, "父路径必须是对象或数组").WithPath(parentPath)
	}

	return value, nil
}

// 应用replace操作
func applyReplaceOperation(value types.JSONValue, path string, rawValue json.RawMessage) (types.JSONValue, error) {
	// 解析要替换的值
	newValue, err := parser.ParseBytesToValue(rawValue)
	if err != nil {
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrInvalidPatch, "无效的值").WithCause(err)
	}

	// 处理根路径
	if path == "$" {
		return newValue, nil
	}

	// 获取父路径和最后一个段
	parentPath, lastSegment := splitPath(path)

	// 获取父对象或数组
	results, err := jsonpath.QueryJSONPath(value, parentPath)
	if err != nil || len(results) == 0 {
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrPathNotFound, "父路径不存在").WithPath(parentPath)
	}
	parent := results[0]

	// 根据父对象类型处理
	if parent.IsObject() {
		obj, _ := parent.AsObject()
		if !obj.Has(lastSegment) {
			return nil, jsonerrors.NewJSONError(jsonerrors.ErrPathNotFound, "路径不存在").WithPath(path)
		}
		obj.Put(lastSegment, newValue)
	} else if parent.IsArray() {
		arr, _ := parent.AsArray()
		index, err := parseArrayIndex(lastSegment, arr.Size())
		if err != nil {
			return nil, err
		}
		if index >= arr.Size() {
			return nil, jsonerrors.NewJSONError(jsonerrors.ErrIndexOutOfRange, "索引超出范围").WithPath(path)
		}
		arr.Set(index, newValue)
	} else {
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrInvalidType, "父路径必须是对象或数组").WithPath(parentPath)
	}

	return value, nil
}

// 应用move操作
func applyMoveOperation(value types.JSONValue, from, path string) (types.JSONValue, error) {
	// 获取源值
	results, err := jsonpath.QueryJSONPath(value, from)
	if err != nil || len(results) == 0 {
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrPathNotFound, "源路径不存在").WithPath(from)
	}
	sourceValue := results[0]

	// 先移除源值
	value, err = applyRemoveOperation(value, from)
	if err != nil {
		return nil, err
	}

	// 将源值添加到目标路径
	valueBytes, _ := json.Marshal(types.ValueToInterface(sourceValue))
	return applyAddOperation(value, path, valueBytes)
}

// 应用copy操作
func applyCopyOperation(value types.JSONValue, from, path string) (types.JSONValue, error) {
	// 获取源值
	results, err := jsonpath.QueryJSONPath(value, from)
	if err != nil || len(results) == 0 {
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrPathNotFound, "源路径不存在").WithPath(from)
	}
	sourceValue := results[0]

	// 将源值添加到目标路径
	valueBytes, _ := json.Marshal(types.ValueToInterface(sourceValue))
	return applyAddOperation(value, path, valueBytes)
}

// 应用test操作
func applyTestOperation(value types.JSONValue, path string, rawValue json.RawMessage) (types.JSONValue, error) {
	// 解析要测试的值
	testValue, err := parser.ParseBytesToValue(rawValue)
	if err != nil {
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrInvalidPatch, "无效的值").WithCause(err)
	}

	// 获取目标值
	results, err := jsonpath.QueryJSONPath(value, path)
	if err != nil || len(results) == 0 {
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrPathNotFound, "路径不存在").WithPath(path)
	}
	targetValue := results[0]

	// 比较值
	if !compareValues(targetValue, testValue) {
		return nil, jsonerrors.ErrTestFailedWithDetails(path, testValue, targetValue)
	}

	return value, nil
}

// 比较两个JSON值是否相等
func compareValues(a, b types.JSONValue) bool {
	if a.Type() != b.Type() {
		return false
	}

	switch a.Type() {
	case "null":
		return true
	case "boolean":
		aBool, _ := a.AsBoolean()
		bBool, _ := b.AsBoolean()
		return aBool == bBool
	case "number":
		aNum, _ := a.AsNumber()
		bNum, _ := b.AsNumber()
		return aNum == bNum
	case "string":
		aStr, _ := a.AsString()
		bStr, _ := b.AsString()
		return aStr == bStr
	case "array":
		aArr, _ := a.AsArray()
		bArr, _ := b.AsArray()
		if aArr.Size() != bArr.Size() {
			return false
		}
		for i := 0; i < aArr.Size(); i++ {
			if !compareValues(aArr.Get(i), bArr.Get(i)) {
				return false
			}
		}
		return true
	case "object":
		aObj, _ := a.AsObject()
		bObj, _ := b.AsObject()
		if aObj.Size() != bObj.Size() {
			return false
		}
		for _, key := range aObj.Keys() {
			if !bObj.Has(key) || !compareValues(aObj.Get(key), bObj.Get(key)) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// 分割路径为父路径和最后一个段
func splitPath(path string) (string, string) {
	// 处理数组索引
	if strings.HasSuffix(path, "]") {
		lastOpenBracket := strings.LastIndex(path, "[")
		if lastOpenBracket != -1 {
			return path[:lastOpenBracket], path[lastOpenBracket+1 : len(path)-1]
		}
	}

	// 处理对象属性
	lastDot := strings.LastIndex(path, ".")
	if lastDot != -1 {
		return path[:lastDot], path[lastDot+1:]
	}

	return "$", path[1:]
}

// 解析数组索引
func parseArrayIndex(indexStr string, arraySize int) (int, error) {
	// 处理"-"表示数组末尾
	if indexStr == "-" {
		return arraySize, nil
	}

	// 解析数字索引
	var index int
	_, err := fmt.Sscanf(indexStr, "%d", &index)
	if err != nil || index < 0 {
		return 0, jsonerrors.NewJSONError(jsonerrors.ErrInvalidIndex, "无效的数组索引")
	}

	return index, nil
}
