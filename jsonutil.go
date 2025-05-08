// Package gojson 提供类似JavaScript JSON接口的Go JSON工具函数
package gojson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Parse 将JSON字符串解析为Go对象
// 类似于JavaScript的JSON.parse()
func Parse(jsonStr string, v interface{}) error {
	if jsonStr == "" {
		return errors.New("输入的JSON字符串为空")
	}

	return json.Unmarshal([]byte(jsonStr), v)
}

// ParseBytes 将JSON字节数组解析为Go对象
func ParseBytes(jsonBytes []byte, v interface{}) error {
	if len(jsonBytes) == 0 {
		return errors.New("输入的JSON字节数组为空")
	}

	return json.Unmarshal(jsonBytes, v)
}

// Stringify 将Go对象转换为JSON字符串
// 类似于JavaScript的JSON.stringify()
func Stringify(v interface{}) (string, error) {
	if v == nil {
		return "null", nil
	}

	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// StringifyBytes 将Go对象转换为JSON字节数组
func StringifyBytes(v interface{}) ([]byte, error) {
	if v == nil {
		return []byte("null"), nil
	}

	return json.Marshal(v)
}

// StringifyIndent 将Go对象转换为格式化的JSON字符串
// prefix 是每行输出的前缀
// indent 是每个嵌套级别的缩进
func StringifyIndent(v interface{}, prefix, indent string) (string, error) {
	if v == nil {
		return "null", nil
	}

	jsonBytes, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// IsValidJSON 检查字符串是否为有效的JSON
func IsValidJSON(jsonStr string) bool {
	if jsonStr == "" {
		return false
	}

	var js interface{}
	return json.Unmarshal([]byte(jsonStr), &js) == nil
}

// IsValidJSONBytes 检查字节数组是否为有效的JSON
func IsValidJSONBytes(jsonBytes []byte) bool {
	if len(jsonBytes) == 0 {
		return false
	}

	var js interface{}
	return json.Unmarshal(jsonBytes, &js) == nil
}

// Prettify 格式化JSON字符串，使其更易读
func Prettify(jsonStr string) (string, error) {
	if jsonStr == "" {
		return "", errors.New("输入的JSON字符串为空")
	}

	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, []byte(jsonStr), "", "  ")
	if err != nil {
		return "", err
	}

	return prettyJSON.String(), nil
}

// Minify 最小化JSON字符串，移除不必要的空白
func Minify(jsonStr string) (string, error) {
	if jsonStr == "" {
		return "", errors.New("输入的JSON字符串为空")
	}

	var js interface{}
	err := json.Unmarshal([]byte(jsonStr), &js)
	if err != nil {
		return "", err
	}

	minified, err := json.Marshal(js)
	if err != nil {
		return "", err
	}

	return string(minified), nil
}

// GetType 获取JSON值的类型
func GetType(jsonStr string) (string, error) {
	if jsonStr == "" {
		return "", errors.New("输入的JSON字符串为空")
	}

	var js interface{}
	err := json.Unmarshal([]byte(jsonStr), &js)
	if err != nil {
		return "", err
	}

	return getTypeOf(js), nil
}

// 内部函数，用于获取接口值的类型
func getTypeOf(v interface{}) string {
	if v == nil {
		return "null"
	}

	switch v.(type) {
	case bool:
		return "boolean"
	case float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "number"
	case string:
		return "string"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return reflect.TypeOf(v).String()
	}
}

// GetValue 从JSON对象中获取指定路径的值
// 路径格式为点分隔的字符串，例如 "user.address.city"
func GetValue(jsonStr string, path string) (interface{}, error) {
	if jsonStr == "" {
		return nil, errors.New("输入的JSON字符串为空")
	}

	var js interface{}
	err := json.Unmarshal([]byte(jsonStr), &js)
	if err != nil {
		return nil, err
	}

	return extractValue(js, path)
}

// 内部函数，用于从接口值中提取指定路径的值
func extractValue(data interface{}, path string) (interface{}, error) {
	if path == "" {
		return data, nil
	}

	// 将路径分割为段
	parts := splitPath(path)
	current := data

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			var ok bool
			current, ok = v[part]
			if !ok {
				return nil, fmt.Errorf("路径 '%s' 中的键 '%s' 不存在", path, part)
			}
		case []interface{}:
			// 尝试将部分解析为数组索引
			var index int
			_, err := fmt.Sscanf(part, "%d", &index)
			if err != nil || index < 0 || index >= len(v) {
				return nil, fmt.Errorf("路径 '%s' 中的索引 '%s' 无效", path, part)
			}
			current = v[index]
		default:
			return nil, fmt.Errorf("路径 '%s' 中的 '%s' 不是对象或数组", path, part)
		}
	}

	return current, nil
}

// 内部函数，用于分割路径字符串
func splitPath(path string) []string {
	// 简单实现，仅支持点分隔的路径
	// 可以扩展以支持更复杂的路径表达式
	return strings.Split(path, ".")
}
