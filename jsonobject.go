package gojson

import (
	"encoding/json"
	"errors"
	"sort"
)

// JSONObject 表示JSON中的对象
type JSONObject struct {
	properties map[string]JSONValue
	keys       []string // 保持键的顺序
}

// NewJSONObject 创建一个新的空JSONObject
func NewJSONObject() *JSONObject {
	return &JSONObject{
		properties: make(map[string]JSONValue),
		keys:       make([]string, 0),
	}
}

// Type 返回JSON值的类型
func (o *JSONObject) Type() string {
	return "object"
}

// String 返回JSON值的字符串表示
func (o *JSONObject) String() string {
	bytes, err := o.MarshalJSON()
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

// MarshalJSON 实现json.Marshaler接口
func (o *JSONObject) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	for k, v := range o.properties {
		if v == nil {
			m[k] = nil
		} else {
			m[k] = ValueToInterface(v)
		}
	}
	return json.Marshal(m)
}

// IsNull 检查值是否为null
func (o *JSONObject) IsNull() bool {
	return false
}

// IsBoolean 检查值是否为布尔值
func (o *JSONObject) IsBoolean() bool {
	return false
}

// IsNumber 检查值是否为数字
func (o *JSONObject) IsNumber() bool {
	return false
}

// IsString 检查值是否为字符串
func (o *JSONObject) IsString() bool {
	return false
}

// IsArray 检查值是否为数组
func (o *JSONObject) IsArray() bool {
	return false
}

// IsObject 检查值是否为对象
func (o *JSONObject) IsObject() bool {
	return true
}

// AsBoolean 将值转换为布尔值
func (o *JSONObject) AsBoolean() (bool, error) {
	return false, errors.New("无法将对象转换为布尔值")
}

// AsNumber 将值转换为数字
func (o *JSONObject) AsNumber() (float64, error) {
	return 0, errors.New("无法将对象转换为数字")
}

// AsString 将值转换为字符串
func (o *JSONObject) AsString() (string, error) {
	return o.String(), nil
}

// AsArray 将值转换为数组
func (o *JSONObject) AsArray() (*JSONArray, error) {
	return nil, errors.New("无法将对象转换为数组")
}

// AsObject 将值转换为对象
func (o *JSONObject) AsObject() (*JSONObject, error) {
	return o, nil
}

// Size 返回对象的属性数量
func (o *JSONObject) Size() int {
	return len(o.properties)
}

// Has 检查对象是否包含指定键
func (o *JSONObject) Has(key string) bool {
	_, ok := o.properties[key]
	return ok
}

// Get 获取指定键的值
func (o *JSONObject) Get(key string) JSONValue {
	if value, ok := o.properties[key]; ok {
		return value
	}
	return NewJSONNull()
}

// GetBoolean 获取指定键的布尔值
func (o *JSONObject) GetBoolean(key string) (bool, error) {
	value := o.Get(key)
	if value.IsNull() {
		return false, errors.New("键不存在或值为null")
	}
	return value.AsBoolean()
}

// GetNumber 获取指定键的数字
func (o *JSONObject) GetNumber(key string) (float64, error) {
	value := o.Get(key)
	if value.IsNull() {
		return 0, errors.New("键不存在或值为null")
	}
	return value.AsNumber()
}

// GetString 获取指定键的字符串
func (o *JSONObject) GetString(key string) (string, error) {
	value := o.Get(key)
	if value.IsNull() {
		return "", errors.New("键不存在或值为null")
	}
	return value.AsString()
}

// GetArray 获取指定键的数组
func (o *JSONObject) GetArray(key string) (*JSONArray, error) {
	value := o.Get(key)
	if value.IsNull() {
		return nil, errors.New("键不存在或值为null")
	}
	return value.AsArray()
}

// GetObject 获取指定键的对象
func (o *JSONObject) GetObject(key string) (*JSONObject, error) {
	value := o.Get(key)
	if value.IsNull() {
		return nil, errors.New("键不存在或值为null")
	}
	return value.AsObject()
}

// Put 设置指定键的值
func (o *JSONObject) Put(key string, value JSONValue) *JSONObject {
	if !o.Has(key) {
		o.keys = append(o.keys, key)
	}
	o.properties[key] = value
	return o
}

// PutBoolean 设置指定键的布尔值
func (o *JSONObject) PutBoolean(key string, value bool) *JSONObject {
	return o.Put(key, NewJSONBool(value))
}

// PutNumber 设置指定键的数字
func (o *JSONObject) PutNumber(key string, value float64) *JSONObject {
	return o.Put(key, NewJSONNumber(value))
}

// PutString 设置指定键的字符串
func (o *JSONObject) PutString(key string, value string) *JSONObject {
	return o.Put(key, NewJSONString(value))
}

// PutNull 设置指定键的null值
func (o *JSONObject) PutNull(key string) *JSONObject {
	return o.Put(key, NewJSONNull())
}

// PutArray 设置指定键的数组
func (o *JSONObject) PutArray(key string, value *JSONArray) *JSONObject {
	return o.Put(key, value)
}

// PutObject 设置指定键的对象
func (o *JSONObject) PutObject(key string, value *JSONObject) *JSONObject {
	return o.Put(key, value)
}

// Remove 移除指定键的属性
func (o *JSONObject) Remove(key string) *JSONObject {
	if o.Has(key) {
		delete(o.properties, key)

		// 从keys中移除
		for i, k := range o.keys {
			if k == key {
				o.keys = append(o.keys[:i], o.keys[i+1:]...)
				break
			}
		}
	}
	return o
}

// Keys 返回对象的所有键
func (o *JSONObject) Keys() []string {
	return o.keys
}

// SortedKeys 返回对象的所有键（按字母顺序排序）
func (o *JSONObject) SortedKeys() []string {
	keys := make([]string, len(o.keys))
	copy(keys, o.keys)
	sort.Strings(keys)
	return keys
}

// ToMap 将JSONObject转换为Go map
func (o *JSONObject) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range o.properties {
		result[k] = ValueToInterface(v)
	}
	return result
}

// ForEach 对对象中的每个属性执行函数
func (o *JSONObject) ForEach(fn func(key string, value JSONValue)) {
	for _, key := range o.keys {
		fn(key, o.properties[key])
	}
}

// Merge 合并另一个JSONObject到当前对象
func (o *JSONObject) Merge(other *JSONObject) *JSONObject {
	if other == nil {
		return o
	}

	other.ForEach(func(key string, value JSONValue) {
		o.Put(key, value)
	})

	return o
}

// Clone 克隆当前JSONObject
func (o *JSONObject) Clone() *JSONObject {
	clone := NewJSONObject()
	o.ForEach(func(key string, value JSONValue) {
		clone.Put(key, value)
	})
	return clone
}
