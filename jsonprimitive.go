package gojson

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// JSONNull 表示JSON中的null值
type JSONNull struct{}

// NewJSONNull 创建一个新的JSONNull
func NewJSONNull() *JSONNull {
	return &JSONNull{}
}

// Type 返回JSON值的类型
func (n *JSONNull) Type() string {
	return "null"
}

// String 返回JSON值的字符串表示
func (n *JSONNull) String() string {
	return "null"
}

// MarshalJSON 实现json.Marshaler接口
func (n *JSONNull) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

// IsNull 检查值是否为null
func (n *JSONNull) IsNull() bool {
	return true
}

// IsBoolean 检查值是否为布尔值
func (n *JSONNull) IsBoolean() bool {
	return false
}

// IsNumber 检查值是否为数字
func (n *JSONNull) IsNumber() bool {
	return false
}

// IsString 检查值是否为字符串
func (n *JSONNull) IsString() bool {
	return false
}

// IsArray 检查值是否为数组
func (n *JSONNull) IsArray() bool {
	return false
}

// IsObject 检查值是否为对象
func (n *JSONNull) IsObject() bool {
	return false
}

// AsBoolean 将值转换为布尔值
func (n *JSONNull) AsBoolean() (bool, error) {
	return false, errors.New("无法将null转换为布尔值")
}

// AsNumber 将值转换为数字
func (n *JSONNull) AsNumber() (float64, error) {
	return 0, errors.New("无法将null转换为数字")
}

// AsString 将值转换为字符串
func (n *JSONNull) AsString() (string, error) {
	return "", errors.New("无法将null转换为字符串")
}

// AsArray 将值转换为数组
func (n *JSONNull) AsArray() (*JSONArray, error) {
	return nil, errors.New("无法将null转换为数组")
}

// AsObject 将值转换为对象
func (n *JSONNull) AsObject() (*JSONObject, error) {
	return nil, errors.New("无法将null转换为对象")
}

// JSONBool 表示JSON中的布尔值
type JSONBool struct {
	value bool
}

// NewJSONBool 创建一个新的JSONBool
func NewJSONBool(value bool) *JSONBool {
	return &JSONBool{value: value}
}

// Type 返回JSON值的类型
func (b *JSONBool) Type() string {
	return "boolean"
}

// String 返回JSON值的字符串表示
func (b *JSONBool) String() string {
	return strconv.FormatBool(b.value)
}

// MarshalJSON 实现json.Marshaler接口
func (b *JSONBool) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.value)
}

// IsNull 检查值是否为null
func (b *JSONBool) IsNull() bool {
	return false
}

// IsBoolean 检查值是否为布尔值
func (b *JSONBool) IsBoolean() bool {
	return true
}

// IsNumber 检查值是否为数字
func (b *JSONBool) IsNumber() bool {
	return false
}

// IsString 检查值是否为字符串
func (b *JSONBool) IsString() bool {
	return false
}

// IsArray 检查值是否为数组
func (b *JSONBool) IsArray() bool {
	return false
}

// IsObject 检查值是否为对象
func (b *JSONBool) IsObject() bool {
	return false
}

// AsBoolean 将值转换为布尔值
func (b *JSONBool) AsBoolean() (bool, error) {
	return b.value, nil
}

// AsNumber 将值转换为数字
func (b *JSONBool) AsNumber() (float64, error) {
	if b.value {
		return 1, nil
	}
	return 0, nil
}

// AsString 将值转换为字符串
func (b *JSONBool) AsString() (string, error) {
	return b.String(), nil
}

// AsArray 将值转换为数组
func (b *JSONBool) AsArray() (*JSONArray, error) {
	return nil, errors.New("无法将布尔值转换为数组")
}

// AsObject 将值转换为对象
func (b *JSONBool) AsObject() (*JSONObject, error) {
	return nil, errors.New("无法将布尔值转换为对象")
}

// GetValue 返回布尔值
func (b *JSONBool) GetValue() bool {
	return b.value
}

// JSONNumber 表示JSON中的数字
type JSONNumber struct {
	value float64
}

// NewJSONNumber 创建一个新的JSONNumber
func NewJSONNumber(value float64) *JSONNumber {
	return &JSONNumber{value: value}
}

// Type 返回JSON值的类型
func (n *JSONNumber) Type() string {
	return "number"
}

// String 返回JSON值的字符串表示
func (n *JSONNumber) String() string {
	return strconv.FormatFloat(n.value, 'f', -1, 64)
}

// MarshalJSON 实现json.Marshaler接口
func (n *JSONNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.value)
}

// IsNull 检查值是否为null
func (n *JSONNumber) IsNull() bool {
	return false
}

// IsBoolean 检查值是否为布尔值
func (n *JSONNumber) IsBoolean() bool {
	return false
}

// IsNumber 检查值是否为数字
func (n *JSONNumber) IsNumber() bool {
	return true
}

// IsString 检查值是否为字符串
func (n *JSONNumber) IsString() bool {
	return false
}

// IsArray 检查值是否为数组
func (n *JSONNumber) IsArray() bool {
	return false
}

// IsObject 检查值是否为对象
func (n *JSONNumber) IsObject() bool {
	return false
}

// AsBoolean 将值转换为布尔值
func (n *JSONNumber) AsBoolean() (bool, error) {
	return n.value != 0, nil
}

// AsNumber 将值转换为数字
func (n *JSONNumber) AsNumber() (float64, error) {
	return n.value, nil
}

// AsString 将值转换为字符串
func (n *JSONNumber) AsString() (string, error) {
	return n.String(), nil
}

// AsArray 将值转换为数组
func (n *JSONNumber) AsArray() (*JSONArray, error) {
	return nil, errors.New("无法将数字转换为数组")
}

// AsObject 将值转换为对象
func (n *JSONNumber) AsObject() (*JSONObject, error) {
	return nil, errors.New("无法将数字转换为对象")
}

// GetValue 返回数字值
func (n *JSONNumber) GetValue() float64 {
	return n.value
}

// JSONString 表示JSON中的字符串
type JSONString struct {
	value string
}

// NewJSONString 创建一个新的JSONString
func NewJSONString(value string) *JSONString {
	return &JSONString{value: value}
}

// Type 返回JSON值的类型
func (s *JSONString) Type() string {
	return "string"
}

// String 返回JSON值的字符串表示
func (s *JSONString) String() string {
	bytes, _ := json.Marshal(s.value)
	return string(bytes)
}

// MarshalJSON 实现json.Marshaler接口
func (s *JSONString) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.value)
}

// IsNull 检查值是否为null
func (s *JSONString) IsNull() bool {
	return false
}

// IsBoolean 检查值是否为布尔值
func (s *JSONString) IsBoolean() bool {
	return false
}

// IsNumber 检查值是否为数字
func (s *JSONString) IsNumber() bool {
	return false
}

// IsString 检查值是否为字符串
func (s *JSONString) IsString() bool {
	return true
}

// IsArray 检查值是否为数组
func (s *JSONString) IsArray() bool {
	return false
}

// IsObject 检查值是否为对象
func (s *JSONString) IsObject() bool {
	return false
}

// AsBoolean 将值转换为布尔值
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

// AsNumber 将值转换为数字
func (s *JSONString) AsNumber() (float64, error) {
	return strconv.ParseFloat(s.value, 64)
}

// AsString 将值转换为字符串
func (s *JSONString) AsString() (string, error) {
	return s.value, nil
}

// AsArray 将值转换为数组
func (s *JSONString) AsArray() (*JSONArray, error) {
	return nil, errors.New("无法将字符串转换为数组")
}

// AsObject 将值转换为对象
func (s *JSONString) AsObject() (*JSONObject, error) {
	return nil, errors.New("无法将字符串转换为对象")
}

// GetValue 返回字符串值
func (s *JSONString) GetValue() string {
	return s.value
}
