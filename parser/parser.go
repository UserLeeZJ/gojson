// Package parser 提供gojson库的JSON解析功能。
package parser

import (
	"encoding/json"
	"strconv"

	jsonerrors "github.com/UserLeeZJ/gojson/errors"
	"github.com/UserLeeZJ/gojson/fast"
	"github.com/UserLeeZJ/gojson/types"
)

// ParseToValue 将JSON字符串解析为JSONValue。
func ParseToValue(jsonStr string) (types.JSONValue, error) {
	if jsonStr == "" {
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrEmptyInput, "输入的JSON字符串为空")
	}

	var raw interface{}
	err := fast.Unmarshal([]byte(jsonStr), &raw)
	if err != nil {
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrInvalidJSON, "解析JSON失败").WithCause(err)
	}

	return convertToJSONValue(raw), nil
}

// ParseBytesToValue 将JSON字节数组解析为JSONValue。
func ParseBytesToValue(jsonBytes []byte) (types.JSONValue, error) {
	if len(jsonBytes) == 0 {
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrEmptyInput, "输入的JSON字节数组为空")
	}

	var raw interface{}
	err := fast.Unmarshal(jsonBytes, &raw)
	if err != nil {
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrInvalidJSON, "解析JSON失败").WithCause(err)
	}

	return convertToJSONValue(raw), nil
}

// Parse 将JSON字符串解析为Go对象。
func Parse(jsonStr string, v interface{}) error {
	if jsonStr == "" {
		return jsonerrors.NewJSONError(jsonerrors.ErrEmptyInput, "输入的JSON字符串为空")
	}

	err := fast.Unmarshal([]byte(jsonStr), v)
	if err != nil {
		return jsonerrors.NewJSONError(jsonerrors.ErrInvalidJSON, "解析JSON失败").WithCause(err)
	}

	return nil
}

// ParseBytes 将JSON字节数组解析为Go对象。
func ParseBytes(jsonBytes []byte, v interface{}) error {
	if len(jsonBytes) == 0 {
		return jsonerrors.NewJSONError(jsonerrors.ErrEmptyInput, "输入的JSON字节数组为空")
	}

	err := fast.Unmarshal(jsonBytes, v)
	if err != nil {
		return jsonerrors.NewJSONError(jsonerrors.ErrInvalidJSON, "解析JSON失败").WithCause(err)
	}

	return nil
}

// Stringify 将Go对象转换为JSON字符串。
func Stringify(v interface{}) (string, error) {
	if v == nil {
		return "null", nil
	}

	jsonBytes, err := fast.Marshal(v)
	if err != nil {
		return "", jsonerrors.NewJSONError(jsonerrors.ErrInvalidJSON, "序列化JSON失败").WithCause(err)
	}

	return string(jsonBytes), nil
}

// StringifyBytes 将Go对象转换为JSON字节数组。
func StringifyBytes(v interface{}) ([]byte, error) {
	if v == nil {
		return []byte("null"), nil
	}

	jsonBytes, err := fast.Marshal(v)
	if err != nil {
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrInvalidJSON, "序列化JSON失败").WithCause(err)
	}

	return jsonBytes, nil
}

// StringifyIndent 将Go对象转换为格式化的JSON字符串。
func StringifyIndent(v interface{}, prefix, indent string) (string, error) {
	if v == nil {
		return "null", nil
	}

	// 由于fast.Marshal不支持缩进，这里仍使用标准库
	jsonBytes, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		return "", jsonerrors.NewJSONError(jsonerrors.ErrInvalidJSON, "序列化JSON失败").WithCause(err)
	}

	return string(jsonBytes), nil
}

// convertToJSONValue 将Go原生类型转换为JSONValue。
func convertToJSONValue(v interface{}) types.JSONValue {
	if v == nil {
		return types.NewJSONNull()
	}

	switch val := v.(type) {
	case bool:
		return types.NewJSONBool(val)
	case float64:
		return types.NewJSONNumber(val)
	case json.Number:
		// 尝试转换为float64
		f, err := val.Float64()
		if err == nil {
			return types.NewJSONNumber(f)
		}
		// 如果转换失败，尝试转换为字符串
		s := val.String()
		num, _ := strconv.ParseFloat(s, 64)
		return types.NewJSONNumber(num)
	case string:
		return types.NewJSONString(val)
	case []interface{}:
		arr := types.NewJSONArray()
		for _, item := range val {
			arr.Add(convertToJSONValue(item))
		}
		return arr
	case map[string]interface{}:
		obj := types.NewJSONObject()
		for k, v := range val {
			obj.Put(k, convertToJSONValue(v))
		}
		return obj
	default:
		// 尝试将其他类型转换为JSON
		data, err := json.Marshal(val)
		if err != nil {
			return types.NewJSONNull()
		}

		var raw interface{}
		err = json.Unmarshal(data, &raw)
		if err != nil {
			return types.NewJSONNull()
		}

		return convertToJSONValue(raw)
	}
}
