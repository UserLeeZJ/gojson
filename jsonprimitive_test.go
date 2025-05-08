package gojson

import (
	"testing"
)

// 测试JSONNull
func TestJSONNull(t *testing.T) {
	null := NewJSONNull()

	// 测试类型
	if null.Type() != "null" {
		t.Errorf("JSONNull.Type() = %v, want %v", null.Type(), "null")
	}

	// 测试字符串表示
	if null.String() != "null" {
		t.Errorf("JSONNull.String() = %v, want %v", null.String(), "null")
	}

	// 测试类型检查
	if !null.IsNull() {
		t.Errorf("JSONNull.IsNull() = %v, want %v", null.IsNull(), true)
	}
	if null.IsBoolean() {
		t.Errorf("JSONNull.IsBoolean() = %v, want %v", null.IsBoolean(), false)
	}
	if null.IsNumber() {
		t.Errorf("JSONNull.IsNumber() = %v, want %v", null.IsNumber(), false)
	}
	if null.IsString() {
		t.Errorf("JSONNull.IsString() = %v, want %v", null.IsString(), false)
	}
	if null.IsArray() {
		t.Errorf("JSONNull.IsArray() = %v, want %v", null.IsArray(), false)
	}
	if null.IsObject() {
		t.Errorf("JSONNull.IsObject() = %v, want %v", null.IsObject(), false)
	}

	// 测试转换方法
	if _, err := null.AsBoolean(); err == nil {
		t.Errorf("JSONNull.AsBoolean() should return error")
	}
	if _, err := null.AsNumber(); err == nil {
		t.Errorf("JSONNull.AsNumber() should return error")
	}
	if _, err := null.AsString(); err == nil {
		t.Errorf("JSONNull.AsString() should return error")
	}
	if _, err := null.AsArray(); err == nil {
		t.Errorf("JSONNull.AsArray() should return error")
	}
	if _, err := null.AsObject(); err == nil {
		t.Errorf("JSONNull.AsObject() should return error")
	}
}

// 测试JSONBool
func TestJSONBool(t *testing.T) {
	boolTrue := NewJSONBool(true)
	boolFalse := NewJSONBool(false)

	// 测试类型
	if boolTrue.Type() != "boolean" {
		t.Errorf("JSONBool.Type() = %v, want %v", boolTrue.Type(), "boolean")
	}

	// 测试字符串表示
	if boolTrue.String() != "true" {
		t.Errorf("JSONBool(true).String() = %v, want %v", boolTrue.String(), "true")
	}
	if boolFalse.String() != "false" {
		t.Errorf("JSONBool(false).String() = %v, want %v", boolFalse.String(), "false")
	}

	// 测试类型检查
	if boolTrue.IsNull() {
		t.Errorf("JSONBool.IsNull() = %v, want %v", boolTrue.IsNull(), false)
	}
	if !boolTrue.IsBoolean() {
		t.Errorf("JSONBool.IsBoolean() = %v, want %v", boolTrue.IsBoolean(), true)
	}
	if boolTrue.IsNumber() {
		t.Errorf("JSONBool.IsNumber() = %v, want %v", boolTrue.IsNumber(), false)
	}
	if boolTrue.IsString() {
		t.Errorf("JSONBool.IsString() = %v, want %v", boolTrue.IsString(), false)
	}
	if boolTrue.IsArray() {
		t.Errorf("JSONBool.IsArray() = %v, want %v", boolTrue.IsArray(), false)
	}
	if boolTrue.IsObject() {
		t.Errorf("JSONBool.IsObject() = %v, want %v", boolTrue.IsObject(), false)
	}

	// 测试转换方法
	if val, err := boolTrue.AsBoolean(); err != nil || val != true {
		t.Errorf("JSONBool(true).AsBoolean() = %v, %v, want %v, nil", val, err, true)
	}
	if val, err := boolFalse.AsBoolean(); err != nil || val != false {
		t.Errorf("JSONBool(false).AsBoolean() = %v, %v, want %v, nil", val, err, false)
	}

	if val, err := boolTrue.AsNumber(); err != nil || val != 1 {
		t.Errorf("JSONBool(true).AsNumber() = %v, %v, want %v, nil", val, err, 1)
	}
	if val, err := boolFalse.AsNumber(); err != nil || val != 0 {
		t.Errorf("JSONBool(false).AsNumber() = %v, %v, want %v, nil", val, err, 0)
	}

	if val, err := boolTrue.AsString(); err != nil || val != "true" {
		t.Errorf("JSONBool(true).AsString() = %v, %v, want %v, nil", val, err, "true")
	}
	if val, err := boolFalse.AsString(); err != nil || val != "false" {
		t.Errorf("JSONBool(false).AsString() = %v, %v, want %v, nil", val, err, "false")
	}

	if _, err := boolTrue.AsArray(); err == nil {
		t.Errorf("JSONBool.AsArray() should return error")
	}
	if _, err := boolTrue.AsObject(); err == nil {
		t.Errorf("JSONBool.AsObject() should return error")
	}

	// 测试GetValue方法
	if boolTrue.GetValue() != true {
		t.Errorf("JSONBool(true).GetValue() = %v, want %v", boolTrue.GetValue(), true)
	}
	if boolFalse.GetValue() != false {
		t.Errorf("JSONBool(false).GetValue() = %v, want %v", boolFalse.GetValue(), false)
	}
}

// 测试JSONNumber
func TestJSONNumber(t *testing.T) {
	num := NewJSONNumber(123.45)

	// 测试类型
	if num.Type() != "number" {
		t.Errorf("JSONNumber.Type() = %v, want %v", num.Type(), "number")
	}

	// 测试字符串表示
	if num.String() != "123.45" {
		t.Errorf("JSONNumber(123.45).String() = %v, want %v", num.String(), "123.45")
	}

	// 测试类型检查
	if num.IsNull() {
		t.Errorf("JSONNumber.IsNull() = %v, want %v", num.IsNull(), false)
	}
	if num.IsBoolean() {
		t.Errorf("JSONNumber.IsBoolean() = %v, want %v", num.IsBoolean(), false)
	}
	if !num.IsNumber() {
		t.Errorf("JSONNumber.IsNumber() = %v, want %v", num.IsNumber(), true)
	}
	if num.IsString() {
		t.Errorf("JSONNumber.IsString() = %v, want %v", num.IsString(), false)
	}
	if num.IsArray() {
		t.Errorf("JSONNumber.IsArray() = %v, want %v", num.IsArray(), false)
	}
	if num.IsObject() {
		t.Errorf("JSONNumber.IsObject() = %v, want %v", num.IsObject(), false)
	}

	// 测试转换方法
	if val, err := num.AsBoolean(); err != nil || val != true {
		t.Errorf("JSONNumber(123.45).AsBoolean() = %v, %v, want %v, nil", val, err, true)
	}

	if val, err := num.AsNumber(); err != nil || val != 123.45 {
		t.Errorf("JSONNumber(123.45).AsNumber() = %v, %v, want %v, nil", val, err, 123.45)
	}

	if val, err := num.AsString(); err != nil || val != "123.45" {
		t.Errorf("JSONNumber(123.45).AsString() = %v, %v, want %v, nil", val, err, "123.45")
	}

	if _, err := num.AsArray(); err == nil {
		t.Errorf("JSONNumber.AsArray() should return error")
	}
	if _, err := num.AsObject(); err == nil {
		t.Errorf("JSONNumber.AsObject() should return error")
	}

	// 测试GetValue方法
	if num.GetValue() != 123.45 {
		t.Errorf("JSONNumber(123.45).GetValue() = %v, want %v", num.GetValue(), 123.45)
	}
}

// 测试JSONString
func TestJSONString(t *testing.T) {
	str := NewJSONString("hello")

	// 测试类型
	if str.Type() != "string" {
		t.Errorf("JSONString.Type() = %v, want %v", str.Type(), "string")
	}

	// 测试字符串表示
	if str.String() != `"hello"` {
		t.Errorf(`JSONString("hello").String() = %v, want %v`, str.String(), `"hello"`)
	}

	// 测试类型检查
	if str.IsNull() {
		t.Errorf("JSONString.IsNull() = %v, want %v", str.IsNull(), false)
	}
	if str.IsBoolean() {
		t.Errorf("JSONString.IsBoolean() = %v, want %v", str.IsBoolean(), false)
	}
	if str.IsNumber() {
		t.Errorf("JSONString.IsNumber() = %v, want %v", str.IsNumber(), false)
	}
	if !str.IsString() {
		t.Errorf("JSONString.IsString() = %v, want %v", str.IsString(), true)
	}
	if str.IsArray() {
		t.Errorf("JSONString.IsArray() = %v, want %v", str.IsArray(), false)
	}
	if str.IsObject() {
		t.Errorf("JSONString.IsObject() = %v, want %v", str.IsObject(), false)
	}

	// 测试转换方法
	if _, err := str.AsBoolean(); err == nil {
		t.Errorf("JSONString(\"hello\").AsBoolean() should return error")
	}

	// 测试特殊字符串转布尔值
	trueStr := NewJSONString("true")
	falseStr := NewJSONString("false")
	if val, err := trueStr.AsBoolean(); err != nil || val != true {
		t.Errorf("JSONString(\"true\").AsBoolean() = %v, %v, want %v, nil", val, err, true)
	}
	if val, err := falseStr.AsBoolean(); err != nil || val != false {
		t.Errorf("JSONString(\"false\").AsBoolean() = %v, %v, want %v, nil", val, err, false)
	}

	// 测试数字字符串转数字
	numStr := NewJSONString("123.45")
	if val, err := numStr.AsNumber(); err != nil || val != 123.45 {
		t.Errorf("JSONString(\"123.45\").AsNumber() = %v, %v, want %v, nil", val, err, 123.45)
	}

	if val, err := str.AsString(); err != nil || val != "hello" {
		t.Errorf("JSONString(\"hello\").AsString() = %v, %v, want %v, nil", val, err, "hello")
	}

	if _, err := str.AsArray(); err == nil {
		t.Errorf("JSONString.AsArray() should return error")
	}
	if _, err := str.AsObject(); err == nil {
		t.Errorf("JSONString.AsObject() should return error")
	}

	// 测试GetValue方法
	if str.GetValue() != "hello" {
		t.Errorf("JSONString(\"hello\").GetValue() = %v, want %v", str.GetValue(), "hello")
	}
}
