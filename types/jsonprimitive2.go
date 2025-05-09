// Package types 提供gojson库的基本类型定义。
package types

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/UserLeeZJ/gojson/errors"
)

// JSONNumber 表示JSON中的数字值。
type JSONNumber struct {
	value float64
}

// NewJSONNumber 创建一个新的JSONNumber对象。
func NewJSONNumber(value float64) *JSONNumber {
	return &JSONNumber{value: value}
}

// Type 返回JSON值的类型。
func (n *JSONNumber) Type() string {
	return "number"
}

// String 返回JSON值的字符串表示。
func (n *JSONNumber) String() string {
	return strconv.FormatFloat(n.value, 'f', -1, 64)
}

// MarshalJSON 实现json.Marshaler接口。
func (n *JSONNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.value)
}

// IsNull 检查值是否为null。
func (n *JSONNumber) IsNull() bool {
	return false
}

// IsBoolean 检查值是否为布尔值。
func (n *JSONNumber) IsBoolean() bool {
	return false
}

// IsNumber 检查值是否为数字。
func (n *JSONNumber) IsNumber() bool {
	return true
}

// IsString 检查值是否为字符串。
func (n *JSONNumber) IsString() bool {
	return false
}

// IsArray 检查值是否为数组。
func (n *JSONNumber) IsArray() bool {
	return false
}

// IsObject 检查值是否为对象。
func (n *JSONNumber) IsObject() bool {
	return false
}

// AsBoolean 将值转换为布尔值。
func (n *JSONNumber) AsBoolean() (bool, error) {
	return n.value != 0, nil
}

// AsNumber 将值转换为数字。
func (n *JSONNumber) AsNumber() (float64, error) {
	return n.value, nil
}

// AsString 将值转换为字符串。
func (n *JSONNumber) AsString() (string, error) {
	return n.String(), nil
}

// AsArray 将值转换为数组。
func (n *JSONNumber) AsArray() (*JSONArray, error) {
	return nil, errors.ErrInvalidTypeWithDetails("array", "number")
}

// AsObject 将值转换为对象。
func (n *JSONNumber) AsObject() (*JSONObject, error) {
	return nil, errors.ErrInvalidTypeWithDetails("object", "number")
}

// JSONString 表示JSON中的字符串值。
type JSONString struct {
	value string
}

// NewJSONString 创建一个新的JSONString对象。
func NewJSONString(value string) *JSONString {
	return &JSONString{value: value}
}

// Type 返回JSON值的类型。
func (s *JSONString) Type() string {
	return "string"
}

// String 返回JSON值的字符串表示。
func (s *JSONString) String() string {
	bytes, _ := json.Marshal(s.value)
	return string(bytes)
}

// MarshalJSON 实现json.Marshaler接口。
func (s *JSONString) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.value)
}

// IsNull 检查值是否为null。
func (s *JSONString) IsNull() bool {
	return false
}

// IsBoolean 检查值是否为布尔值。
func (s *JSONString) IsBoolean() bool {
	return false
}

// IsNumber 检查值是否为数字。
func (s *JSONString) IsNumber() bool {
	return false
}

// IsString 检查值是否为字符串。
func (s *JSONString) IsString() bool {
	return true
}

// IsArray 检查值是否为数组。
func (s *JSONString) IsArray() bool {
	return false
}

// IsObject 检查值是否为对象。
func (s *JSONString) IsObject() bool {
	return false
}

// AsBoolean 将值转换为布尔值。
func (s *JSONString) AsBoolean() (bool, error) {
	switch s.value {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("无法将字符串 '%s' 转换为布尔值", s.value)
	}
}

// AsNumber 将值转换为数字。
func (s *JSONString) AsNumber() (float64, error) {
	return strconv.ParseFloat(s.value, 64)
}

// AsString 将值转换为字符串。
func (s *JSONString) AsString() (string, error) {
	return s.value, nil
}

// AsArray 将值转换为数组。
func (s *JSONString) AsArray() (*JSONArray, error) {
	return nil, errors.ErrInvalidTypeWithDetails("array", "string")
}

// AsObject 将值转换为对象。
func (s *JSONString) AsObject() (*JSONObject, error) {
	return nil, errors.ErrInvalidTypeWithDetails("object", "string")
}
