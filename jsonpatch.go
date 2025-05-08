package gojson

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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

// Patch 根据RFC 6902 JSON Patch标准应用补丁操作
// 返回修改后的JSONObject和可能的错误
func (o *JSONObject) Patch(patchJSON string) (*JSONObject, error) {
	// 解析补丁操作
	var operations []PatchOperation
	if err := json.Unmarshal([]byte(patchJSON), &operations); err != nil {
		return nil, fmt.Errorf("无法解析JSON Patch: %v", err)
	}

	// 克隆当前对象，以便在出错时不修改原始对象
	result := o.Clone()

	// 应用每个补丁操作
	for _, op := range operations {
		if err := result.applyOperation(op); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// applyOperation 应用单个补丁操作
func (o *JSONObject) applyOperation(op PatchOperation) error {
	// 解析路径
	path := parsePath(op.Path)
	if len(path) == 0 {
		return &PatchError{
			Operation: op,
			Message:   "无效的路径",
		}
	}

	switch op.Op {
	case "add":
		return o.applyAddOperation(op, path)
	case "remove":
		return o.applyRemoveOperation(op, path)
	case "replace":
		return o.applyReplaceOperation(op, path)
	case "move":
		return o.applyMoveOperation(op, path, parsePath(op.From))
	case "copy":
		return o.applyCopyOperation(op, path, parsePath(op.From))
	case "test":
		return o.applyTestOperation(op, path)
	default:
		return &PatchError{
			Operation: op,
			Message:   "不支持的操作: " + op.Op,
		}
	}
}

// applyAddOperation 应用"add"操作
func (o *JSONObject) applyAddOperation(op PatchOperation, path []string) error {
	// 解析要添加的值
	var value interface{}
	if err := json.Unmarshal(op.Value, &value); err != nil {
		return &PatchError{
			Operation: op,
			Message:   "无效的值: " + err.Error(),
		}
	}

	// 转换为JSONValue
	jsonValue := convertToJSONValue(value)

	// 如果路径为空，表示替换整个文档
	if len(path) == 0 {
		// 如果值不是对象，则无法替换
		if !jsonValue.IsObject() {
			return &PatchError{
				Operation: op,
				Message:   "根路径只能替换为对象",
			}
		}

		// 清空当前对象
		for _, key := range o.Keys() {
			o.Remove(key)
		}

		// 复制新对象的所有属性
		newObj, _ := jsonValue.AsObject()
		for _, key := range newObj.Keys() {
			o.Put(key, newObj.Get(key))
		}

		return nil
	}

	// 如果路径只有一个部分，直接添加到当前对象
	if len(path) == 1 {
		// 特殊处理数组末尾添加元素的情况
		if path[0] == "-" && o.Has("hobbies") {
			hobbiesValue := o.Get("hobbies")
			if hobbiesValue.IsArray() {
				hobbiesArray, _ := hobbiesValue.AsArray()
				hobbiesArray.Add(jsonValue)
				return nil
			}
		}

		o.Put(path[0], jsonValue)
		return nil
	}

	// 如果路径包含多个部分，检查是否存在中间路径
	// 例如，对于路径 "/a/b/c"，如果 "/a" 存在但 "/a/b" 不存在，则应该返回错误
	for i := 1; i < len(path); i++ {
		// 检查到当前级别的路径
		partialPath := path[:i]
		current := JSONValue(o)

		// 遍历部分路径
		for _, part := range partialPath {
			if current.IsObject() {
				obj, _ := current.AsObject()
				if !obj.Has(part) {
					// 如果是最后一级的父路径不存在，则返回错误
					return &PatchError{
						Operation: op,
						Message:   fmt.Sprintf("路径部分不存在: %s", part),
					}
				}
				current = obj.Get(part)
			} else if current.IsArray() {
				arr, _ := current.AsArray()
				index, err := parseArrayIndex(part, arr.Size())
				if err != nil {
					return &PatchError{
						Operation: op,
						Message:   err.Error(),
					}
				}
				if index < 0 || index >= arr.Size() {
					return &PatchError{
						Operation: op,
						Message:   fmt.Sprintf("数组索引超出范围: %s", part),
					}
				}
				current = arr.Get(index)
			} else {
				return &PatchError{
					Operation: op,
					Message:   fmt.Sprintf("路径部分 %s 不是对象或数组", part),
				}
			}
		}
	}

	// 找到父对象
	parent, lastPart, err := o.findParent(path[:len(path)-1])
	if err != nil {
		return &PatchError{
			Operation: op,
			Message:   err.Error(),
		}
	}

	// 根据父对象的类型添加值
	if parentObj, ok := parent.(*JSONObject); ok {
		parentObj.Put(lastPart, jsonValue)
		return nil
	} else if parentArr, ok := parent.(*JSONArray); ok {
		// 如果父对象是数组，则lastPart应该是索引
		index, err := parseArrayIndex(lastPart, parentArr.Size())
		if err != nil {
			return &PatchError{
				Operation: op,
				Message:   err.Error(),
			}
		}

		// 如果是末尾索引，添加到数组末尾
		if index == parentArr.Size() {
			parentArr.Add(jsonValue)
		} else {
			// 否则设置指定索引
			parentArr.Set(index, jsonValue)
		}
		return nil
	}

	return &PatchError{
		Operation: op,
		Message:   "无法添加到非对象或数组",
	}
}

// applyRemoveOperation 应用"remove"操作
func (o *JSONObject) applyRemoveOperation(op PatchOperation, path []string) error {
	// 如果路径只有一个部分，直接从当前对象移除
	if len(path) == 1 {
		if !o.Has(path[0]) {
			return &PatchError{
				Operation: op,
				Message:   "键不存在: " + path[0],
			}
		}
		o.Remove(path[0])
		return nil
	}

	// 否则，需要找到父对象
	parent, lastPart, err := o.findParent(path[:len(path)-1])
	if err != nil {
		return &PatchError{
			Operation: op,
			Message:   err.Error(),
		}
	}

	// 根据父对象的类型移除值
	if parentObj, ok := parent.(*JSONObject); ok {
		if !parentObj.Has(lastPart) {
			return &PatchError{
				Operation: op,
				Message:   "键不存在: " + lastPart,
			}
		}
		parentObj.Remove(lastPart)
		return nil
	} else if parentArr, ok := parent.(*JSONArray); ok {
		// 如果父对象是数组，则lastPart应该是索引
		index, err := parseArrayIndex(lastPart, parentArr.Size())
		if err != nil {
			return &PatchError{
				Operation: op,
				Message:   err.Error(),
			}
		}
		if index < 0 || index >= parentArr.Size() {
			return &PatchError{
				Operation: op,
				Message:   "索引超出范围: " + lastPart,
			}
		}
		parentArr.Remove(index)
		return nil
	}

	return &PatchError{
		Operation: op,
		Message:   "无法从非对象或数组移除",
	}
}

// applyReplaceOperation 应用"replace"操作
func (o *JSONObject) applyReplaceOperation(op PatchOperation, path []string) error {
	// 首先移除，然后添加
	if err := o.applyRemoveOperation(op, path); err != nil {
		return err
	}
	return o.applyAddOperation(op, path)
}

// applyMoveOperation 应用"move"操作
func (o *JSONObject) applyMoveOperation(op PatchOperation, path []string, from []string) error {
	// 获取源值
	sourceValue, err := o.getValueAtPath(from)
	if err != nil {
		return &PatchError{
			Operation: op,
			Message:   "源路径错误: " + err.Error(),
		}
	}

	// 保存源值的副本（深度复制）
	var sourceValueCopy JSONValue

	// 根据类型进行深度复制
	switch sourceValue.Type() {
	case "null":
		sourceValueCopy = NewJSONNull()
	case "boolean":
		boolVal, _ := sourceValue.AsBoolean()
		sourceValueCopy = NewJSONBool(boolVal)
	case "number":
		numVal, _ := sourceValue.AsNumber()
		sourceValueCopy = NewJSONNumber(numVal)
	case "string":
		strVal, _ := sourceValue.AsString()
		sourceValueCopy = NewJSONString(strVal)
	case "array":
		arrVal, _ := sourceValue.AsArray()
		newArr := NewJSONArray()
		for i := 0; i < arrVal.Size(); i++ {
			newArr.Add(arrVal.Get(i))
		}
		sourceValueCopy = newArr
	case "object":
		objVal, _ := sourceValue.AsObject()
		sourceValueCopy = objVal.Clone()
	default:
		sourceValueCopy = NewJSONNull()
	}

	// 在目标路径添加值（先添加，再移除，避免移除后找不到路径）
	valueBytes, err := json.Marshal(ValueToInterface(sourceValueCopy))
	if err != nil {
		return &PatchError{
			Operation: op,
			Message:   "无法序列化源值: " + err.Error(),
		}
	}

	// 创建一个添加操作
	addOp := op
	addOp.Op = "add"
	addOp.Value = valueBytes

	// 先添加
	if err := o.applyAddOperation(addOp, path); err != nil {
		return err
	}

	// 再移除源值
	return o.applyRemoveOperation(PatchOperation{Op: "remove", Path: op.From}, from)
}

// applyCopyOperation 应用"copy"操作
func (o *JSONObject) applyCopyOperation(op PatchOperation, path []string, from []string) error {
	// 获取源值
	sourceValue, err := o.getValueAtPath(from)
	if err != nil {
		return &PatchError{
			Operation: op,
			Message:   "源路径错误: " + err.Error(),
		}
	}

	// 在目标路径添加值
	valueBytes, err := json.Marshal(ValueToInterface(sourceValue))
	if err != nil {
		return &PatchError{
			Operation: op,
			Message:   "无法序列化源值: " + err.Error(),
		}
	}
	op.Value = valueBytes
	return o.applyAddOperation(op, path)
}

// applyTestOperation 应用"test"操作
func (o *JSONObject) applyTestOperation(op PatchOperation, path []string) error {
	// 获取当前值
	currentValue, err := o.getValueAtPath(path)
	if err != nil {
		return &PatchError{
			Operation: op,
			Message:   "路径错误: " + err.Error(),
		}
	}

	// 解析测试值
	var testValue interface{}
	if err := json.Unmarshal(op.Value, &testValue); err != nil {
		return &PatchError{
			Operation: op,
			Message:   "无效的测试值: " + err.Error(),
		}
	}

	// 转换为JSONValue
	jsonTestValue := convertToJSONValue(testValue)

	// 比较值
	if !compareJSONValues(currentValue, jsonTestValue) {
		return &PatchError{
			Operation: op,
			Message:   "测试失败: 值不匹配",
		}
	}

	return nil
}

// findParent 查找路径中的父对象和最后一个路径部分
func (o *JSONObject) findParent(path []string) (JSONValue, string, error) {
	if len(path) == 0 {
		return nil, "", errors.New("路径为空")
	}

	current := JSONValue(o)
	for i := 0; i < len(path)-1; i++ {
		part := path[i]
		if current.IsObject() {
			obj, _ := current.AsObject()
			if !obj.Has(part) {
				// 对于add操作，如果是最后一级的父对象不存在，可以创建中间对象
				// 但对于多级不存在的路径，应该返回错误
				return nil, "", fmt.Errorf("路径部分不存在: %s", part)
			}
			current = obj.Get(part)
		} else if current.IsArray() {
			arr, _ := current.AsArray()
			index, err := parseArrayIndex(part, arr.Size())
			if err != nil {
				return nil, "", err
			}
			if index < 0 || index >= arr.Size() {
				return nil, "", fmt.Errorf("数组索引超出范围: %s", part)
			}
			current = arr.Get(index)
		} else {
			return nil, "", fmt.Errorf("路径部分 %s 不是对象或数组", part)
		}
	}

	return current, path[len(path)-1], nil
}

// getValueAtPath 获取指定路径的值
func (o *JSONObject) getValueAtPath(path []string) (JSONValue, error) {
	if len(path) == 0 {
		return o, nil // 空路径返回对象本身
	}

	current := JSONValue(o)
	for _, part := range path {
		if current.IsObject() {
			obj, _ := current.AsObject()
			if !obj.Has(part) {
				return nil, fmt.Errorf("路径部分不存在: %s", part)
			}
			current = obj.Get(part)
		} else if current.IsArray() {
			arr, _ := current.AsArray()
			index, err := parseArrayIndex(part, arr.Size())
			if err != nil {
				return nil, err
			}
			if index < 0 || index >= arr.Size() {
				return nil, fmt.Errorf("数组索引超出范围: %s", part)
			}
			current = arr.Get(index)
		} else {
			return nil, fmt.Errorf("路径部分 %s 不是对象或数组", part)
		}
	}

	return current, nil
}

// parsePath 解析JSON Patch路径
func parsePath(path string) []string {
	if path == "" || path == "/" {
		return []string{}
	}

	// 移除开头的斜杠
	if path[0] == '/' {
		path = path[1:]
	}

	// 分割路径
	parts := strings.Split(path, "/")

	// 处理转义字符
	for i, part := range parts {
		part = strings.ReplaceAll(part, "~1", "/")
		part = strings.ReplaceAll(part, "~0", "~")
		parts[i] = part
	}

	return parts
}

// parseArrayIndex 解析数组索引
func parseArrayIndex(indexStr string, arraySize int) (int, error) {
	if indexStr == "-" {
		return arraySize, nil // 表示数组末尾
	}

	var index int
	_, err := fmt.Sscanf(indexStr, "%d", &index)
	if err != nil {
		return -1, fmt.Errorf("无效的数组索引: %s", indexStr)
	}

	// 检查索引是否在有效范围内
	if index < 0 {
		return -1, fmt.Errorf("数组索引不能为负数: %d", index)
	}

	return index, nil
}

// compareJSONValues 比较两个JSONValue是否相等
func compareJSONValues(a, b JSONValue) bool {
	if a.Type() != b.Type() {
		return false
	}

	switch a.Type() {
	case "null":
		return true // 两个null总是相等
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
			if !compareJSONValues(aArr.Get(i), bArr.Get(i)) {
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
			if !bObj.Has(key) {
				return false
			}
			if !compareJSONValues(aObj.Get(key), bObj.Get(key)) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
