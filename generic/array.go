package generic

import (
	"encoding/json"
	"fmt"

	"github.com/UserLeeZJ/gojson/errors"
	"github.com/UserLeeZJ/gojson/types"
)

// JSONArray is a generic version of JSON array
// T is the Go type used to represent the elements in the array
type JSONArray[T any] struct {
	arr *types.JSONArray
}

// NewJSONArray creates a new generic JSON array
func NewJSONArray[T any]() *JSONArray[T] {
	return &JSONArray[T]{
		arr: types.NewJSONArray(),
	}
}

// FromJSONArray creates a generic JSONArray from a non-generic JSONArray
func FromJSONArray[T any](arr *types.JSONArray) *JSONArray[T] {
	return &JSONArray[T]{
		arr: arr,
	}
}

// Value returns the value as Go type
func (a *JSONArray[T]) Value() []T {
	// Convert JSONArray to []interface{}
	arr := a.arr.ToArray()

	// Use json package for type conversion
	data, err := json.Marshal(arr)
	if err != nil {
		return nil
	}
	
	var result []T
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil
	}

	return result
}

// JSONValue returns the non-generic JSONValue
func (a *JSONArray[T]) JSONValue() types.JSONValue {
	return a.arr
}

// Size returns the size of the array
func (a *JSONArray[T]) Size() int {
	return a.arr.Size()
}

// Get gets the value at the specified index
func (a *JSONArray[T]) Get(index int) types.JSONValue {
	return a.arr.Get(index)
}

// GetTyped gets the typed value at the specified index
func (a *JSONArray[T]) GetTyped(index int) (T, error) {
	var zero T
	if index < 0 || index >= a.arr.Size() {
		return zero, errors.ErrIndexOutOfRangeWithDetails(index, a.arr.Size())
	}

	value := a.arr.Get(index)

	// Use json package for type conversion
	data, err := json.Marshal(types.ValueToInterface(value))
	if err != nil {
		return zero, errors.NewJSONError(errors.ErrTypeConversion,
			fmt.Sprintf("failed to marshal JSON: %v", err)).WithCause(err)
	}
	
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return zero, errors.NewJSONError(errors.ErrTypeConversion,
			fmt.Sprintf("cannot convert JSON to %T", zero)).WithCause(err)
	}

	return result, nil
}

// Add adds a value to the array
func (a *JSONArray[T]) Add(value types.JSONValue) *JSONArray[T] {
	a.arr.Add(value)
	return a
}

// AddTyped adds a typed value to the array
func (a *JSONArray[T]) AddTyped(value T) (*JSONArray[T], error) {
	jsonValue, err := ToJSONValue(value)
	if err != nil {
		return a, err
	}
	
	a.arr.Add(jsonValue)
	return a, nil
}

// AddString adds a string value to the array
func (a *JSONArray[T]) AddString(value string) *JSONArray[T] {
	a.arr.AddString(value)
	return a
}

// AddNumber adds a number value to the array
func (a *JSONArray[T]) AddNumber(value float64) *JSONArray[T] {
	a.arr.AddNumber(value)
	return a
}

// AddBoolean adds a boolean value to the array
func (a *JSONArray[T]) AddBoolean(value bool) *JSONArray[T] {
	a.arr.AddBoolean(value)
	return a
}

// AddNull adds a null value to the array
func (a *JSONArray[T]) AddNull() *JSONArray[T] {
	a.arr.AddNull()
	return a
}

// Set sets the value at the specified index
func (a *JSONArray[T]) Set(index int, value types.JSONValue) *JSONArray[T] {
	a.arr.Set(index, value)
	return a
}

// SetTyped sets the typed value at the specified index
func (a *JSONArray[T]) SetTyped(index int, value T) (*JSONArray[T], error) {
	jsonValue, err := ToJSONValue(value)
	if err != nil {
		return a, err
	}
	
	a.arr.Set(index, jsonValue)
	return a, nil
}

// Remove removes the value at the specified index
func (a *JSONArray[T]) Remove(index int) *JSONArray[T] {
	a.arr.Remove(index)
	return a
}

// String returns the JSON string representation
func (a *JSONArray[T]) String() string {
	return a.arr.String()
}

// ToArray converts the JSONArray to a Go slice
func (a *JSONArray[T]) ToArray() []interface{} {
	return a.arr.ToArray()
}

// ForEach executes a function for each element in the array
func (a *JSONArray[T]) ForEach(fn func(value types.JSONValue, index int)) {
	a.arr.ForEach(fn)
}

// Map applies a function to each element in the array and returns a new array
func (a *JSONArray[T]) Map(fn func(value types.JSONValue, index int) types.JSONValue) *JSONArray[T] {
	return FromJSONArray[T](a.arr.Map(fn))
}

// Filter filters the array and returns a new array
func (a *JSONArray[T]) Filter(fn func(value types.JSONValue, index int) bool) *JSONArray[T] {
	return FromJSONArray[T](a.arr.Filter(fn))
}

// Slice returns a slice of the array
func (a *JSONArray[T]) Slice(start, end int) *JSONArray[T] {
	return FromJSONArray[T](a.arr.Slice(start, end))
}
