package gojson

import (
	"encoding/json"
	"errors"
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

// ParseToValue 将JSON字符串解析为JSONValue
func ParseToValue(jsonStr string) (JSONValue, error) {
	if jsonStr == "" {
		return nil, errors.New("输入的JSON字符串为空")
	}

	var raw interface{}
	err := json.Unmarshal([]byte(jsonStr), &raw)
	if err != nil {
		return nil, err
	}

	return convertToJSONValue(raw), nil
}

// 将Go原生类型转换为JSONValue
func convertToJSONValue(v interface{}) JSONValue {
	if v == nil {
		return NewJSONNull()
	}

	switch val := v.(type) {
	case bool:
		return NewJSONBool(val)
	case float64:
		return NewJSONNumber(val)
	case float32:
		return NewJSONNumber(float64(val))
	case int:
		return NewJSONNumber(float64(val))
	case int8:
		return NewJSONNumber(float64(val))
	case int16:
		return NewJSONNumber(float64(val))
	case int32:
		return NewJSONNumber(float64(val))
	case int64:
		return NewJSONNumber(float64(val))
	case uint:
		return NewJSONNumber(float64(val))
	case uint8:
		return NewJSONNumber(float64(val))
	case uint16:
		return NewJSONNumber(float64(val))
	case uint32:
		return NewJSONNumber(float64(val))
	case uint64:
		return NewJSONNumber(float64(val))
	case string:
		return NewJSONString(val)
	case []interface{}:
		arr := NewJSONArray()
		for _, item := range val {
			arr.Add(convertToJSONValue(item))
		}
		return arr
	case map[string]interface{}:
		obj := NewJSONObject()
		for k, v := range val {
			obj.Put(k, convertToJSONValue(v))
		}
		return obj
	default:
		// 尝试将其他类型转换为JSON
		data, err := json.Marshal(val)
		if err != nil {
			return NewJSONNull()
		}

		var raw interface{}
		err = json.Unmarshal(data, &raw)
		if err != nil {
			return NewJSONNull()
		}

		return convertToJSONValue(raw)
	}
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
