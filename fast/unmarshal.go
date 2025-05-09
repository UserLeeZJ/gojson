// Package fast 提供高性能的JSON序列化和反序列化功能。
package fast

import (
	"bytes"
	"encoding/json"
	"strconv"

	jsonerrors "github.com/UserLeeZJ/gojson/errors"
)

// Unmarshal 是一个优化的JSON反序列化函数。
func Unmarshal(data []byte, v interface{}) error {
	// 快速检查空输入。
	if len(data) == 0 {
		return jsonerrors.NewJSONError(ErrEmptyInput, "输入的JSON字节数组为空")
	}

	// 快速检查无效JSON。
	if !isValidJSON(data) {
		return jsonerrors.NewJSONError(ErrInvalidJSON, "无效的JSON格式")
	}

	// 优化：检查目标类型。
	switch target := v.(type) {
	case *string:
		// 字符串类型的快速路径。
		return unmarshalString(data, target)
	case *int:
		// 整数类型的快速路径。
		return unmarshalInt(data, target)
	case *bool:
		// 布尔类型的快速路径。
		return unmarshalBool(data, target)
	case *map[string]interface{}:
		// 空map的快速路径。
		if len(*target) == 0 && isEmptyObject(data) {
			return nil // 空对象，不需要处理。
		}
	case *[]interface{}:
		// 空数组的快速路径。
		if len(*target) == 0 && isEmptyArray(data) {
			return nil // 空数组，不需要处理。
		}
	}

	// 尝试直接使用标准库的Unmarshal。
	err := json.Unmarshal(data, v)
	if err != nil {
		// 如果失败，尝试使用Decoder，它可能处理某些特殊情况更好。
		r := bytes.NewReader(data)
		dec := json.NewDecoder(r)
		dec.UseNumber() // 使用Number类型保持数字精度。

		if decErr := dec.Decode(v); decErr != nil {
			// 返回原始错误，它通常更有信息。
			return jsonerrors.NewJSONError(ErrInvalidJSON, "反序列化失败").WithCause(err)
		}
	}

	return nil
}

// unmarshalString 快速解析字符串值。
func unmarshalString(data []byte, target *string) error {
	// 跳过空白。
	start := 0
	for start < len(data) && isWhitespace(data[start]) {
		start++
	}

	// 检查是否为字符串。
	if start >= len(data) || data[start] != '"' {
		return jsonerrors.NewJSONError(ErrInvalidJSON, "期望字符串值")
	}

	// 查找结束引号。
	end := len(data) - 1
	for end > start && isWhitespace(data[end]) {
		end--
	}

	if end <= start || data[end] != '"' {
		return jsonerrors.NewJSONError(ErrInvalidJSON, "字符串未正确结束")
	}

	// 检查是否需要处理转义字符。
	needUnescape := false
	for i := start + 1; i < end; i++ {
		if data[i] == '\\' {
			needUnescape = true
			break
		}
	}

	if !needUnescape {
		// 不需要处理转义，直接提取字符串。
		*target = string(data[start+1 : end])
		return nil
	}

	// 需要处理转义，使用标准库。
	return json.Unmarshal(data, target)
}

// unmarshalInt 快速解析整数值。
func unmarshalInt(data []byte, target *int) error {
	// 跳过空白。
	start := 0
	for start < len(data) && isWhitespace(data[start]) {
		start++
	}

	// 检查是否为数字。
	if start >= len(data) || (!isDigit(data[start]) && data[start] != '-') {
		return jsonerrors.NewJSONError(ErrInvalidJSON, "期望数字值")
	}

	// 查找结束位置。
	end := len(data) - 1
	for end > start && isWhitespace(data[end]) {
		end--
	}

	// 解析数字。
	val, err := strconv.Atoi(string(data[start : end+1]))
	if err != nil {
		return jsonerrors.NewJSONError(ErrInvalidJSON, "无效的整数值").WithCause(err)
	}

	*target = val
	return nil
}

// unmarshalBool 快速解析布尔值。
func unmarshalBool(data []byte, target *bool) error {
	// 跳过空白。
	start := 0
	for start < len(data) && isWhitespace(data[start]) {
		start++
	}

	// 检查是否为布尔值。
	if start >= len(data) {
		return jsonerrors.NewJSONError(ErrInvalidJSON, "期望布尔值")
	}

	// 查找结束位置。
	end := len(data) - 1
	for end > start && isWhitespace(data[end]) {
		end--
	}

	// 解析布尔值。
	s := string(data[start : end+1])
	if s == "true" {
		*target = true
		return nil
	} else if s == "false" {
		*target = false
		return nil
	}

	return jsonerrors.NewJSONError(ErrInvalidJSON, "无效的布尔值")
}

// isEmptyObject 检查是否为空对象 {}。
func isEmptyObject(data []byte) bool {
	start := 0
	for start < len(data) && isWhitespace(data[start]) {
		start++
	}

	if start >= len(data) || data[start] != '{' {
		return false
	}

	end := len(data) - 1
	for end > start && isWhitespace(data[end]) {
		end--
	}

	if end <= start || data[end] != '}' {
		return false
	}

	// 检查中间是否只有空白。
	for i := start + 1; i < end; i++ {
		if !isWhitespace(data[i]) {
			return false
		}
	}

	return true
}

// isEmptyArray 检查是否为空数组 []。
func isEmptyArray(data []byte) bool {
	start := 0
	for start < len(data) && isWhitespace(data[start]) {
		start++
	}

	if start >= len(data) || data[start] != '[' {
		return false
	}

	end := len(data) - 1
	for end > start && isWhitespace(data[end]) {
		end--
	}

	if end <= start || data[end] != ']' {
		return false
	}

	// 检查中间是否只有空白。
	for i := start + 1; i < end; i++ {
		if !isWhitespace(data[i]) {
			return false
		}
	}

	return true
}

// isValidJSON 快速检查JSON是否有效。
// 这是一个简单的检查，只检查基本的JSON语法。
func isValidJSON(data []byte) bool {
	// 跳过空白。
	i := 0
	for i < len(data) && isWhitespace(data[i]) {
		i++
	}

	// 空JSON不是有效的JSON。
	if i >= len(data) {
		return false
	}

	// 检查第一个非空白字符。
	firstChar := data[i]

	// JSON必须以{、[、"、数字、true、false或null开头。
	return firstChar == '{' || firstChar == '[' || firstChar == '"' ||
		isDigit(firstChar) || firstChar == '-' || firstChar == 't' ||
		firstChar == 'f' || firstChar == 'n'
}

// isWhitespace 检查字符是否为空白。
func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// isDigit 检查字符是否为数字。
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
