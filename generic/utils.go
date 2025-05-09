package generic

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/UserLeeZJ/gojson/errors"
	"github.com/UserLeeZJ/gojson/types"
)

// GetTyped gets a typed value from a JSONObject
// V is the target Go type
// obj is the JSONObject to get the value from
// key is the key to get the value for. If empty, the entire object is used
func GetTyped[V any](obj *types.JSONObject, key string) (V, error) {
	var zero V
	var value types.JSONValue

	// If key is empty, use the entire object
	if key == "" {
		value = obj
	} else {
		value = obj.Get(key)
		if value.IsNull() {
			return zero, errors.ErrPathNotFoundWithDetails(key)
		}
	}

	// Get the target type using reflection
	targetType := reflect.TypeOf(zero)

	// Convert based on target type
	switch targetType.Kind() {
	case reflect.String:
		if !value.IsString() {
			return zero, errors.ErrInvalidTypeWithDetails("string", value.Type())
		}
		str, _ := value.AsString()
		return any(str).(V), nil
	case reflect.Bool:
		if !value.IsBoolean() {
			return zero, errors.ErrInvalidTypeWithDetails("boolean", value.Type())
		}
		b, _ := value.AsBoolean()
		return any(b).(V), nil
	case reflect.Float64, reflect.Float32:
		if !value.IsNumber() {
			return zero, errors.ErrInvalidTypeWithDetails("number", value.Type())
		}
		num, _ := value.AsNumber()
		return any(num).(V), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !value.IsNumber() {
			return zero, errors.ErrInvalidTypeWithDetails("number", value.Type())
		}
		num, _ := value.AsNumber()
		return any(int(num)).(V), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if !value.IsNumber() {
			return zero, errors.ErrInvalidTypeWithDetails("number", value.Type())
		}
		num, _ := value.AsNumber()
		return any(uint(num)).(V), nil
	case reflect.Slice, reflect.Array:
		if !value.IsArray() {
			return zero, errors.ErrInvalidTypeWithDetails("array", value.Type())
		}
		arr, _ := value.AsArray()
		// Use json package for type conversion
		data, err := json.Marshal(arr.ToArray())
		if err != nil {
			return zero, errors.NewJSONError(errors.ErrTypeConversion,
				fmt.Sprintf("failed to marshal array: %v", err)).WithCause(err)
		}
		
		var result V
		if err := json.Unmarshal(data, &result); err != nil {
			return zero, errors.NewJSONError(errors.ErrTypeConversion,
				fmt.Sprintf("cannot convert array to %T", zero)).WithCause(err)
		}
		return result, nil
	case reflect.Map:
		if !value.IsObject() {
			return zero, errors.ErrInvalidTypeWithDetails("object", value.Type())
		}
		objVal, _ := value.AsObject()
		// Use json package for type conversion
		data, err := json.Marshal(types.ValueToInterface(objVal))
		if err != nil {
			return zero, errors.NewJSONError(errors.ErrTypeConversion,
				fmt.Sprintf("failed to marshal object: %v", err)).WithCause(err)
		}
		
		var result V
		if err := json.Unmarshal(data, &result); err != nil {
			return zero, errors.NewJSONError(errors.ErrTypeConversion,
				fmt.Sprintf("cannot convert object to %T", zero)).WithCause(err)
		}
		return result, nil
	case reflect.Struct:
		// For structs, use json package for conversion
		data, err := json.Marshal(types.ValueToInterface(value))
		if err != nil {
			return zero, errors.NewJSONError(errors.ErrTypeConversion,
				fmt.Sprintf("failed to marshal JSON: %v", err)).WithCause(err)
		}
		
		var result V
		if err := json.Unmarshal(data, &result); err != nil {
			return zero, errors.NewJSONError(errors.ErrTypeConversion,
				fmt.Sprintf("cannot convert JSON to %T", zero)).WithCause(err)
		}
		return result, nil
	default:
		return zero, errors.NewJSONError(errors.ErrNotSupported,
			fmt.Sprintf("unsupported type: %T", zero))
	}
}

// ToJSONValue converts a Go value to a JSONValue
func ToJSONValue(v interface{}) (types.JSONValue, error) {
	if v == nil {
		return types.NewJSONNull(), nil
	}

	// Handle primitive types directly
	switch val := v.(type) {
	case string:
		return types.NewJSONString(val), nil
	case float64:
		return types.NewJSONNumber(val), nil
	case float32:
		return types.NewJSONNumber(float64(val)), nil
	case int:
		return types.NewJSONNumber(float64(val)), nil
	case int8:
		return types.NewJSONNumber(float64(val)), nil
	case int16:
		return types.NewJSONNumber(float64(val)), nil
	case int32:
		return types.NewJSONNumber(float64(val)), nil
	case int64:
		return types.NewJSONNumber(float64(val)), nil
	case uint:
		return types.NewJSONNumber(float64(val)), nil
	case uint8:
		return types.NewJSONNumber(float64(val)), nil
	case uint16:
		return types.NewJSONNumber(float64(val)), nil
	case uint32:
		return types.NewJSONNumber(float64(val)), nil
	case uint64:
		return types.NewJSONNumber(float64(val)), nil
	case bool:
		return types.NewJSONBool(val), nil
	}

	// For complex types, use reflection and json marshaling
	data, err := json.Marshal(v)
	if err != nil {
		return nil, errors.NewJSONError(errors.ErrTypeConversion,
			fmt.Sprintf("failed to marshal %T to JSON", v)).WithCause(err)
	}

	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, errors.NewJSONError(errors.ErrTypeConversion,
			fmt.Sprintf("failed to unmarshal JSON to interface{}: %v", err)).WithCause(err)
	}

	return convertToJSONValue(raw)
}

// convertToJSONValue converts a Go native type to JSONValue
func convertToJSONValue(v interface{}) (types.JSONValue, error) {
	if v == nil {
		return types.NewJSONNull(), nil
	}

	switch val := v.(type) {
	case string:
		return types.NewJSONString(val), nil
	case float64:
		return types.NewJSONNumber(val), nil
	case int:
		return types.NewJSONNumber(float64(val)), nil
	case bool:
		return types.NewJSONBool(val), nil
	case map[string]interface{}:
		obj := types.NewJSONObject()
		for k, v := range val {
			jsonVal, err := convertToJSONValue(v)
			if err != nil {
				return nil, err
			}
			obj.Put(k, jsonVal)
		}
		return obj, nil
	case []interface{}:
		arr := types.NewJSONArray()
		for _, v := range val {
			jsonVal, err := convertToJSONValue(v)
			if err != nil {
				return nil, err
			}
			arr.Add(jsonVal)
		}
		return arr, nil
	default:
		return nil, errors.NewJSONError(errors.ErrNotSupported,
			fmt.Sprintf("unsupported type: %T", val))
	}
}
