// Package types 提供gojson库的基本类型定义
package types

import (
	"encoding/json"
)

// JSONValue 是所有JSON值类型的通用接口
// 类似于JavaScript中的JSON值
type JSONValue interface {
	// Type 返回JSON值的类型
	Type() string

	// String 返回JSON值的字符串表示
	String() string

	// MarshalJSON 实现json.Marshaler接口
	MarshalJSON() ([]byte, error)

	// IsNull 检查值是否为null
	IsNull() bool

	// IsBoolean 检查值是否为布尔值
	IsBoolean() bool

	// IsNumber 检查值是否为数字
	IsNumber() bool

	// IsString 检查值是否为字符串
	IsString() bool

	// IsArray 检查值是否为数组
	IsArray() bool

	// IsObject 检查值是否为对象
	IsObject() bool

	// AsBoolean 将值转换为布尔值
	AsBoolean() (bool, error)

	// AsNumber 将值转换为数字
	AsNumber() (float64, error)

	// AsString 将值转换为字符串
	AsString() (string, error)

	// AsArray 将值转换为数组
	AsArray() (*JSONArray, error)

	// AsObject 将值转换为对象
	AsObject() (*JSONObject, error)
}

// ValueToInterface 将JSONValue转换为Go原生类型
func ValueToInterface(v JSONValue) interface{} {
	if v == nil || v.IsNull() {
		return nil
	}

	switch v.Type() {
	case "boolean":
		val, _ := v.AsBoolean()
		return val
	case "number":
		val, _ := v.AsNumber()
		return val
	case "string":
		val, _ := v.AsString()
		return val
	case "array":
		arr, _ := v.AsArray()
		result := make([]interface{}, arr.Size())
		for i := 0; i < arr.Size(); i++ {
			result[i] = ValueToInterface(arr.Get(i))
		}
		return result
	case "object":
		obj, _ := v.AsObject()
		result := make(map[string]interface{})
		for _, key := range obj.Keys() {
			result[key] = ValueToInterface(obj.Get(key))
		}
		return result
	default:
		return nil
	}
}

// FromGoValue 将Go原生类型转换为JSONValue
func FromGoValue(v interface{}) (JSONValue, error) {
	if v == nil {
		return NewJSONNull(), nil
	}

	switch val := v.(type) {
	case bool:
		return NewJSONBool(val), nil
	case float64:
		return NewJSONNumber(val), nil
	case float32:
		return NewJSONNumber(float64(val)), nil
	case int:
		return NewJSONNumber(float64(val)), nil
	case int8:
		return NewJSONNumber(float64(val)), nil
	case int16:
		return NewJSONNumber(float64(val)), nil
	case int32:
		return NewJSONNumber(float64(val)), nil
	case int64:
		return NewJSONNumber(float64(val)), nil
	case uint:
		return NewJSONNumber(float64(val)), nil
	case uint8:
		return NewJSONNumber(float64(val)), nil
	case uint16:
		return NewJSONNumber(float64(val)), nil
	case uint32:
		return NewJSONNumber(float64(val)), nil
	case uint64:
		return NewJSONNumber(float64(val)), nil
	case string:
		return NewJSONString(val), nil
	case []interface{}:
		arr := NewJSONArray()
		for _, item := range val {
			itemValue, err := FromGoValue(item)
			if err != nil {
				return nil, err
			}
			arr.Add(itemValue)
		}
		return arr, nil
	case map[string]interface{}:
		obj := NewJSONObject()
		for key, item := range val {
			itemValue, err := FromGoValue(item)
			if err != nil {
				return nil, err
			}
			obj.Put(key, itemValue)
		}
		return obj, nil
	default:
		// 尝试使用json.Marshal和json.Unmarshal进行转换
		data, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}
		var result interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, err
		}
		return FromGoValue(result)
	}
}
