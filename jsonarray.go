package gojson

import (
	"encoding/json"
	"errors"
	"strings"
)

// JSONArray 表示JSON中的数组
type JSONArray struct {
	elements []JSONValue
}

// NewJSONArray 创建一个新的空JSONArray
func NewJSONArray() *JSONArray {
	return &JSONArray{
		elements: make([]JSONValue, 0),
	}
}

// NewJSONArrayFromValues 从JSONValue切片创建一个新的JSONArray
func NewJSONArrayFromValues(values []JSONValue) *JSONArray {
	return &JSONArray{
		elements: values,
	}
}

// Type 返回JSON值的类型
func (a *JSONArray) Type() string {
	return "array"
}

// String 返回JSON值的字符串表示
func (a *JSONArray) String() string {
	bytes, err := a.MarshalJSON()
	if err != nil {
		return "[]"
	}
	return string(bytes)
}

// MarshalJSON 实现json.Marshaler接口
func (a *JSONArray) MarshalJSON() ([]byte, error) {
	values := make([]interface{}, len(a.elements))
	for i, v := range a.elements {
		if v == nil {
			values[i] = nil
		} else {
			values[i] = ValueToInterface(v)
		}
	}
	return json.Marshal(values)
}

// IsNull 检查值是否为null
func (a *JSONArray) IsNull() bool {
	return false
}

// IsBoolean 检查值是否为布尔值
func (a *JSONArray) IsBoolean() bool {
	return false
}

// IsNumber 检查值是否为数字
func (a *JSONArray) IsNumber() bool {
	return false
}

// IsString 检查值是否为字符串
func (a *JSONArray) IsString() bool {
	return false
}

// IsArray 检查值是否为数组
func (a *JSONArray) IsArray() bool {
	return true
}

// IsObject 检查值是否为对象
func (a *JSONArray) IsObject() bool {
	return false
}

// AsBoolean 将值转换为布尔值
func (a *JSONArray) AsBoolean() (bool, error) {
	return false, errors.New("无法将数组转换为布尔值")
}

// AsNumber 将值转换为数字
func (a *JSONArray) AsNumber() (float64, error) {
	return 0, errors.New("无法将数组转换为数字")
}

// AsString 将值转换为字符串
func (a *JSONArray) AsString() (string, error) {
	return a.String(), nil
}

// AsArray 将值转换为数组
func (a *JSONArray) AsArray() (*JSONArray, error) {
	return a, nil
}

// AsObject 将值转换为对象
func (a *JSONArray) AsObject() (*JSONObject, error) {
	return nil, errors.New("无法将数组转换为对象")
}

// Size 返回数组的大小
func (a *JSONArray) Size() int {
	return len(a.elements)
}

// Get 获取指定索引的元素
func (a *JSONArray) Get(index int) JSONValue {
	if index < 0 || index >= len(a.elements) {
		return NewJSONNull()
	}
	return a.elements[index]
}

// GetBoolean 获取指定索引的布尔值
func (a *JSONArray) GetBoolean(index int) (bool, error) {
	value := a.Get(index)
	if value == nil {
		return false, errors.New("索引超出范围")
	}
	return value.AsBoolean()
}

// GetNumber 获取指定索引的数字
func (a *JSONArray) GetNumber(index int) (float64, error) {
	value := a.Get(index)
	if value == nil {
		return 0, errors.New("索引超出范围")
	}
	return value.AsNumber()
}

// GetString 获取指定索引的字符串
func (a *JSONArray) GetString(index int) (string, error) {
	value := a.Get(index)
	if value == nil {
		return "", errors.New("索引超出范围")
	}
	return value.AsString()
}

// GetArray 获取指定索引的数组
func (a *JSONArray) GetArray(index int) (*JSONArray, error) {
	value := a.Get(index)
	if value == nil {
		return nil, errors.New("索引超出范围")
	}
	return value.AsArray()
}

// GetObject 获取指定索引的对象
func (a *JSONArray) GetObject(index int) (*JSONObject, error) {
	value := a.Get(index)
	if value == nil {
		return nil, errors.New("索引超出范围")
	}
	return value.AsObject()
}

// Add 添加一个元素到数组末尾
func (a *JSONArray) Add(value JSONValue) *JSONArray {
	a.elements = append(a.elements, value)
	return a
}

// AddBoolean 添加一个布尔值到数组末尾
func (a *JSONArray) AddBoolean(value bool) *JSONArray {
	return a.Add(NewJSONBool(value))
}

// AddNumber 添加一个数字到数组末尾
func (a *JSONArray) AddNumber(value float64) *JSONArray {
	return a.Add(NewJSONNumber(value))
}

// AddString 添加一个字符串到数组末尾
func (a *JSONArray) AddString(value string) *JSONArray {
	return a.Add(NewJSONString(value))
}

// AddNull 添加一个null值到数组末尾
func (a *JSONArray) AddNull() *JSONArray {
	return a.Add(NewJSONNull())
}

// Set 设置指定索引的元素
func (a *JSONArray) Set(index int, value JSONValue) *JSONArray {
	// 如果索引超出范围，自动扩展数组
	for len(a.elements) <= index {
		a.elements = append(a.elements, NewJSONNull())
	}
	a.elements[index] = value
	return a
}

// SetBoolean 设置指定索引的布尔值
func (a *JSONArray) SetBoolean(index int, value bool) *JSONArray {
	return a.Set(index, NewJSONBool(value))
}

// SetNumber 设置指定索引的数字
func (a *JSONArray) SetNumber(index int, value float64) *JSONArray {
	return a.Set(index, NewJSONNumber(value))
}

// SetString 设置指定索引的字符串
func (a *JSONArray) SetString(index int, value string) *JSONArray {
	return a.Set(index, NewJSONString(value))
}

// SetNull 设置指定索引的null值
func (a *JSONArray) SetNull(index int) *JSONArray {
	return a.Set(index, NewJSONNull())
}

// Remove 移除指定索引的元素
func (a *JSONArray) Remove(index int) *JSONArray {
	if index < 0 || index >= len(a.elements) {
		return a
	}
	a.elements = append(a.elements[:index], a.elements[index+1:]...)
	return a
}

// ToArray 将JSONArray转换为Go切片
func (a *JSONArray) ToArray() []interface{} {
	result := make([]interface{}, len(a.elements))
	for i, v := range a.elements {
		result[i] = ValueToInterface(v)
	}
	return result
}

// Join 将数组元素连接为字符串
func (a *JSONArray) Join(separator string) string {
	strs := make([]string, len(a.elements))
	for i, v := range a.elements {
		if v == nil {
			strs[i] = "null"
		} else if v.IsString() {
			str, _ := v.AsString()
			strs[i] = str
		} else {
			strs[i] = v.String()
		}
	}
	return strings.Join(strs, separator)
}

// ForEach 对数组中的每个元素执行函数
func (a *JSONArray) ForEach(fn func(value JSONValue, index int)) {
	for i, v := range a.elements {
		fn(v, i)
	}
}

// Map 对数组中的每个元素应用函数，并返回新数组
func (a *JSONArray) Map(fn func(value JSONValue, index int) JSONValue) *JSONArray {
	result := NewJSONArray()
	for i, v := range a.elements {
		result.Add(fn(v, i))
	}
	return result
}

// Filter 过滤数组中的元素，返回新数组
func (a *JSONArray) Filter(fn func(value JSONValue, index int) bool) *JSONArray {
	result := NewJSONArray()
	for i, v := range a.elements {
		if fn(v, i) {
			result.Add(v)
		}
	}
	return result
}
