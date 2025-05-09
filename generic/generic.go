// Package generic provides generic support for the gojson library
package generic

import (
	"encoding/json"

	"github.com/UserLeeZJ/gojson/types"
)

// JSONValue is a generic version of JSON value
// T is the Go type used to represent the JSON value
type JSONValue[T any] interface {
	// Value returns the value as Go type
	Value() T
	// JSONValue returns the non-generic JSONValue
	JSONValue() types.JSONValue
}

// JSONObject is a generic version of JSON object
// T is the Go type used to represent the JSON object
type JSONObject[T any] struct {
	obj *types.JSONObject
}

// NewJSONObject creates a new generic JSON object
func NewJSONObject[T any]() *JSONObject[T] {
	return &JSONObject[T]{
		obj: types.NewJSONObject(),
	}
}

// FromJSONObject creates a generic JSONObject from a non-generic JSONObject
func FromJSONObject[T any](obj *types.JSONObject) *JSONObject[T] {
	return &JSONObject[T]{
		obj: obj,
	}
}

// Value returns the value as Go type
func (o *JSONObject[T]) Value() T {
	var result T
	// Convert JSONObject to map[string]interface{}
	m := make(map[string]interface{})
	for _, key := range o.obj.Keys() {
		m[key] = types.ValueToInterface(o.obj.Get(key))
	}

	// Use json package for type conversion
	data, err := json.Marshal(m)
	if err != nil {
		return result
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return result
	}

	return result
}

// JSONValue returns the non-generic JSONValue
func (o *JSONObject[T]) JSONValue() types.JSONValue {
	return o.obj
}

// Get gets the value for the specified key
func (o *JSONObject[T]) Get(key string) types.JSONValue {
	return o.obj.Get(key)
}

// GetTyped gets a typed value for the specified key
func (o *JSONObject[T]) GetTyped(key string) (interface{}, error) {
	value := o.obj.Get(key)
	if value.IsNull() {
		return nil, nil
	}
	return types.ValueToInterface(value), nil
}

// GetString gets the string value for the specified key
func (o *JSONObject[T]) GetString(key string) (string, error) {
	return o.obj.GetString(key)
}

// GetNumber gets the number value for the specified key
func (o *JSONObject[T]) GetNumber(key string) (float64, error) {
	return o.obj.GetNumber(key)
}

// GetBoolean gets the boolean value for the specified key
func (o *JSONObject[T]) GetBoolean(key string) (bool, error) {
	return o.obj.GetBoolean(key)
}

// GetObject gets the object value for the specified key
func (o *JSONObject[T]) GetObject(key string) (*types.JSONObject, error) {
	return o.obj.GetObject(key)
}

// GetArray gets the array value for the specified key
func (o *JSONObject[T]) GetArray(key string) (*types.JSONArray, error) {
	return o.obj.GetArray(key)
}

// Put sets the value for the specified key
func (o *JSONObject[T]) Put(key string, value types.JSONValue) *JSONObject[T] {
	o.obj.Put(key, value)
	return o
}

// PutTyped sets a typed value for the specified key
func (o *JSONObject[T]) PutTyped(key string, value interface{}) (*JSONObject[T], error) {
	jsonValue, err := ToJSONValue(value)
	if err != nil {
		return o, err
	}

	o.obj.Put(key, jsonValue)
	return o, nil
}

// PutString sets the string value for the specified key
func (o *JSONObject[T]) PutString(key string, value string) *JSONObject[T] {
	o.obj.PutString(key, value)
	return o
}

// PutNumber sets the number value for the specified key
func (o *JSONObject[T]) PutNumber(key string, value float64) *JSONObject[T] {
	o.obj.PutNumber(key, value)
	return o
}

// PutBoolean sets the boolean value for the specified key
func (o *JSONObject[T]) PutBoolean(key string, value bool) *JSONObject[T] {
	o.obj.PutBoolean(key, value)
	return o
}

// PutObject sets the object value for the specified key
func (o *JSONObject[T]) PutObject(key string, value *types.JSONObject) *JSONObject[T] {
	o.obj.PutObject(key, value)
	return o
}

// PutArray sets the array value for the specified key
func (o *JSONObject[T]) PutArray(key string, value *types.JSONArray) *JSONObject[T] {
	o.obj.PutArray(key, value)
	return o
}

// String returns the JSON string representation
func (o *JSONObject[T]) String() string {
	return o.obj.String()
}

// Keys returns all keys in the object
func (o *JSONObject[T]) Keys() []string {
	return o.obj.Keys()
}

// Has checks if the object has the specified key
func (o *JSONObject[T]) Has(key string) bool {
	return o.obj.Has(key)
}

// Size returns the size of the object
func (o *JSONObject[T]) Size() int {
	return o.obj.Size()
}

// Remove removes the specified key
func (o *JSONObject[T]) Remove(key string) *JSONObject[T] {
	o.obj.Remove(key)
	return o
}

// ToMap converts the JSONObject to a Go map
func (o *JSONObject[T]) ToMap() map[string]interface{} {
	return o.obj.ToMap()
}

// ForEach executes a function for each property in the object
func (o *JSONObject[T]) ForEach(fn func(key string, value types.JSONValue)) {
	o.obj.ForEach(fn)
}

// Merge merges another JSONObject into this object
func (o *JSONObject[T]) Merge(other *JSONObject[T]) *JSONObject[T] {
	o.obj.Merge(other.obj)
	return o
}

// Clone clones the JSONObject
func (o *JSONObject[T]) Clone() *JSONObject[T] {
	return FromJSONObject[T](o.obj.Clone())
}
