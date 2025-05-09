// Package types 提供gojson库的基本类型定义。
package types

import (
	"encoding/json"
	"strconv"

	"github.com/UserLeeZJ/gojson/errors"
)

// JSONNull 表示JSON中的null值。
type JSONNull struct{}

// NewJSONNull 创建一个新的JSONNull对象。
func NewJSONNull() *JSONNull {
	return &JSONNull{}
}

// Type 返回JSON值的类型。
func (n *JSONNull) Type() string {
	return "null"
}

// String 返回JSON值的字符串表示。
func (n *JSONNull) String() string {
	return "null"
}

// MarshalJSON 实现json.Marshaler接口。
func (n *JSONNull) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

// IsNull 检查值是否为null。
func (n *JSONNull) IsNull() bool {
	return true
}

// IsBoolean 检查值是否为布尔值。
func (n *JSONNull) IsBoolean() bool {
	return false
}

// IsNumber 检查值是否为数字。
func (n *JSONNull) IsNumber() bool {
	return false
}

// IsString 检查值是否为字符串。
func (n *JSONNull) IsString() bool {
	return false
}

// IsArray 检查值是否为数组。
func (n *JSONNull) IsArray() bool {
	return false
}

// IsObject 检查值是否为对象。
func (n *JSONNull) IsObject() bool {
	return false
}

// AsBoolean 将值转换为布尔值。
func (n *JSONNull) AsBoolean() (bool, error) {
	return false, errors.ErrInvalidTypeWithDetails("boolean", "null")
}

// AsNumber 将值转换为数字。
func (n *JSONNull) AsNumber() (float64, error) {
	return 0, errors.ErrInvalidTypeWithDetails("number", "null")
}

// AsString 将值转换为字符串。
func (n *JSONNull) AsString() (string, error) {
	return "", errors.ErrInvalidTypeWithDetails("string", "null")
}

// AsArray 将值转换为数组。
func (n *JSONNull) AsArray() (*JSONArray, error) {
	return nil, errors.ErrInvalidTypeWithDetails("array", "null")
}

// AsObject 将值转换为对象。
func (n *JSONNull) AsObject() (*JSONObject, error) {
	return nil, errors.ErrInvalidTypeWithDetails("object", "null")
}

// JSONBool 表示JSON中的布尔值。
type JSONBool struct {
	value bool
}

// NewJSONBool 创建一个新的JSONBool对象。
func NewJSONBool(value bool) *JSONBool {
	return &JSONBool{value: value}
}

// Type 返回JSON值的类型。
func (b *JSONBool) Type() string {
	return "boolean"
}

// String 返回JSON值的字符串表示。
func (b *JSONBool) String() string {
	return strconv.FormatBool(b.value)
}

// MarshalJSON 实现json.Marshaler接口。
func (b *JSONBool) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.value)
}

// IsNull 检查值是否为null。
func (b *JSONBool) IsNull() bool {
	return false
}

// IsBoolean 检查值是否为布尔值。
func (b *JSONBool) IsBoolean() bool {
	return true
}

// IsNumber 检查值是否为数字。
func (b *JSONBool) IsNumber() bool {
	return false
}

// IsString 检查值是否为字符串。
func (b *JSONBool) IsString() bool {
	return false
}

// IsArray 检查值是否为数组。
func (b *JSONBool) IsArray() bool {
	return false
}

// IsObject 检查值是否为对象。
func (b *JSONBool) IsObject() bool {
	return false
}

// AsBoolean 将值转换为布尔值。
func (b *JSONBool) AsBoolean() (bool, error) {
	return b.value, nil
}

// AsNumber 将值转换为数字。
func (b *JSONBool) AsNumber() (float64, error) {
	if b.value {
		return 1, nil
	}
	return 0, nil
}

// AsString 将值转换为字符串。
func (b *JSONBool) AsString() (string, error) {
	return b.String(), nil
}

// AsArray 将值转换为数组。
func (b *JSONBool) AsArray() (*JSONArray, error) {
	return nil, errors.ErrInvalidTypeWithDetails("array", "boolean")
}

// AsObject 将值转换为对象。
func (b *JSONBool) AsObject() (*JSONObject, error) {
	return nil, errors.ErrInvalidTypeWithDetails("object", "boolean")
}
