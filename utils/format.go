package utils

import (
	"encoding/json"

	jsonerrors "github.com/UserLeeZJ/gojson/errors"
	"github.com/UserLeeZJ/gojson/parser"
	"github.com/UserLeeZJ/gojson/types"
)

// FormatJSON 格式化JSON字符串。
func FormatJSON(jsonStr string, indent string, sortKeys bool) (string, error) {
	// 解析JSON
	jsonValue, err := parser.ParseToValue(jsonStr)
	if err != nil {
		return "", err
	}

	// 使用PrettyPrint格式化
	options := PrettyOptions{
		Indent:     indent,
		SortKeys:   sortKeys,
		EscapeHTML: false,
	}
	return PrettyPrint(jsonValue, options)
}

// CompactJSON 压缩JSON字符串。
func CompactJSON(jsonStr string) (string, error) {
	// 解析JSON
	jsonValue, err := parser.ParseToValue(jsonStr)
	if err != nil {
		return "", err
	}

	// 使用CompressJSON压缩
	return CompressJSON(jsonValue)
}

// SortJSONKeys 对JSON对象的键进行排序。
func SortJSONKeys(jsonStr string) (string, error) {
	// 解析JSON
	jsonValue, err := parser.ParseToValue(jsonStr)
	if err != nil {
		return "", err
	}

	// 如果不是对象，直接返回
	if !jsonValue.IsObject() {
		return jsonValue.String(), nil
	}

	// 使用sortMapKeys排序键
	native := types.ValueToInterface(jsonValue)
	sorted := sortMapKeys(native)

	// 转换回字符串
	bytes, err := json.Marshal(sorted)
	if err != nil {
		return "", jsonerrors.NewJSONError(jsonerrors.ErrOperationFailed, "排序JSON键失败").WithCause(err)
	}

	return string(bytes), nil
}

// ValidateJSON 验证JSON字符串是否有效。
func ValidateJSON(jsonStr string) error {
	_, err := parser.ParseToValue(jsonStr)
	return err
}

// MergeJSON 合并两个JSON对象。
func MergeJSON(target, source string) (string, error) {
	// 解析目标JSON
	targetValue, err := parser.ParseToValue(target)
	if err != nil {
		return "", err
	}

	// 解析源JSON
	sourceValue, err := parser.ParseToValue(source)
	if err != nil {
		return "", err
	}

	// 如果目标不是对象，返回错误
	if !targetValue.IsObject() {
		return "", jsonerrors.NewJSONError(jsonerrors.ErrInvalidType, "目标JSON不是对象")
	}

	// 如果源不是对象，返回错误
	if !sourceValue.IsObject() {
		return "", jsonerrors.NewJSONError(jsonerrors.ErrInvalidType, "源JSON不是对象")
	}

	// 获取对象
	targetObj, _ := targetValue.AsObject()
	sourceObj, _ := sourceValue.AsObject()

	// 合并对象
	result := mergeObjects(targetObj, sourceObj)

	// 转换为字符串
	return result.String(), nil
}

// mergeObjects 合并两个JSONObject
func mergeObjects(target, source *types.JSONObject) *types.JSONObject {
	// 创建结果对象
	result := types.NewJSONObject()

	// 复制目标对象的所有属性
	for _, key := range target.Keys() {
		result.Put(key, target.Get(key))
	}

	// 添加或覆盖源对象的属性
	for _, key := range source.Keys() {
		sourceValue := source.Get(key)

		// 如果目标也有这个键，并且两者都是对象，则递归合并
		if target.Has(key) {
			targetValue := target.Get(key)
			if targetValue.IsObject() && sourceValue.IsObject() {
				targetObj, _ := targetValue.AsObject()
				sourceObj, _ := sourceValue.AsObject()
				result.Put(key, mergeObjects(targetObj, sourceObj))
				continue
			}
		}

		// 否则直接使用源对象的值
		result.Put(key, sourceValue)
	}

	return result
}

// DeepCopy 深度复制JSON值。
func DeepCopy(value types.JSONValue) types.JSONValue {
	if value == nil || value.IsNull() {
		return types.NewJSONNull()
	}

	switch {
	case value.IsObject():
		obj, _ := value.AsObject()
		result := types.NewJSONObject()
		for _, key := range obj.Keys() {
			result.Put(key, DeepCopy(obj.Get(key)))
		}
		return result
	case value.IsArray():
		arr, _ := value.AsArray()
		result := types.NewJSONArray()
		for i := 0; i < arr.Size(); i++ {
			result.Add(DeepCopy(arr.Get(i)))
		}
		return result
	case value.IsString():
		str, _ := value.AsString()
		return types.NewJSONString(str)
	case value.IsNumber():
		num, _ := value.AsNumber()
		return types.NewJSONNumber(num)
	case value.IsBoolean():
		b, _ := value.AsBoolean()
		return types.NewJSONBool(b)
	default:
		return types.NewJSONNull()
	}
}
