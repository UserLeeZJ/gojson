package types

import (
	"testing"
)

func TestJSONArray(t *testing.T) {
	// 创建空数组
	arr := NewJSONArray()

	// 测试类型
	if arr.Type() != "array" {
		t.Errorf("JSONArray.Type() = %v, want %v", arr.Type(), "array")
	}

	// 测试空数组的字符串表示
	if arr.String() != "[]" {
		t.Errorf("Empty JSONArray.String() = %v, want %v", arr.String(), "[]")
	}

	// 测试类型检查
	if arr.IsNull() {
		t.Errorf("JSONArray.IsNull() = %v, want %v", arr.IsNull(), false)
	}
	if arr.IsBoolean() {
		t.Errorf("JSONArray.IsBoolean() = %v, want %v", arr.IsBoolean(), false)
	}
	if arr.IsNumber() {
		t.Errorf("JSONArray.IsNumber() = %v, want %v", arr.IsNumber(), false)
	}
	if arr.IsString() {
		t.Errorf("JSONArray.IsString() = %v, want %v", arr.IsString(), false)
	}
	if !arr.IsArray() {
		t.Errorf("JSONArray.IsArray() = %v, want %v", arr.IsArray(), true)
	}
	if arr.IsObject() {
		t.Errorf("JSONArray.IsObject() = %v, want %v", arr.IsObject(), false)
	}

	// 测试转换方法
	if _, err := arr.AsBoolean(); err == nil {
		t.Errorf("JSONArray.AsBoolean() should return error")
	}
	if _, err := arr.AsNumber(); err == nil {
		t.Errorf("JSONArray.AsNumber() should return error")
	}
	if val, err := arr.AsString(); err != nil || val != "[]" {
		t.Errorf("JSONArray.AsString() = %v, %v, want %v, nil", val, err, "[]")
	}
	if val, err := arr.AsArray(); err != nil || val != arr {
		t.Errorf("JSONArray.AsArray() = %v, %v, want %v, nil", val, err, arr)
	}
	if _, err := arr.AsObject(); err == nil {
		t.Errorf("JSONArray.AsObject() should return error")
	}

	// 测试添加元素
	arr.AddBoolean(true).AddNumber(123).AddString("hello").AddNull()

	// 测试Size方法
	if arr.Size() != 4 {
		t.Errorf("JSONArray.Size() = %v, want %v", arr.Size(), 4)
	}

	// 测试Get方法
	if !arr.Get(0).IsBoolean() {
		t.Errorf("arr.Get(0).IsBoolean() = %v, want %v", arr.Get(0).IsBoolean(), true)
	}
	if !arr.Get(1).IsNumber() {
		t.Errorf("arr.Get(1).IsNumber() = %v, want %v", arr.Get(1).IsNumber(), true)
	}
	if !arr.Get(2).IsString() {
		t.Errorf("arr.Get(2).IsString() = %v, want %v", arr.Get(2).IsString(), true)
	}
	if !arr.Get(3).IsNull() {
		t.Errorf("arr.Get(3).IsNull() = %v, want %v", arr.Get(3).IsNull(), true)
	}
	if !arr.Get(4).IsNull() { // 超出范围应返回null
		t.Errorf("arr.Get(4).IsNull() = %v, want %v", arr.Get(4).IsNull(), true)
	}

	// 测试GetBoolean方法
	if val, err := arr.GetBoolean(0); err != nil || val != true {
		t.Errorf("arr.GetBoolean(0) = %v, %v, want %v, nil", val, err, true)
	}
	if _, err := arr.GetBoolean(4); err == nil {
		t.Errorf("arr.GetBoolean(4) should return error")
	}

	// 测试GetNumber方法
	if val, err := arr.GetNumber(1); err != nil || val != 123 {
		t.Errorf("arr.GetNumber(1) = %v, %v, want %v, nil", val, err, 123)
	}
	if _, err := arr.GetNumber(4); err == nil {
		t.Errorf("arr.GetNumber(4) should return error")
	}

	// 测试GetString方法
	if val, err := arr.GetString(2); err != nil || val != "hello" {
		t.Errorf("arr.GetString(2) = %v, %v, want %v, nil", val, err, "hello")
	}
	if _, err := arr.GetString(4); err == nil {
		t.Errorf("arr.GetString(4) should return error")
	}

	// 测试Set方法
	arr.Set(1, NewJSONNumber(456))
	if val, err := arr.GetNumber(1); err != nil || val != 456 {
		t.Errorf("After Set, arr.GetNumber(1) = %v, %v, want %v, nil", val, err, 456)
	}

	// 测试Set超出范围自动扩展
	arr.Set(10, NewJSONString("expanded"))
	if arr.Size() != 11 {
		t.Errorf("After Set(10), arr.Size() = %v, want %v", arr.Size(), 11)
	}
	if val, err := arr.GetString(10); err != nil || val != "expanded" {
		t.Errorf("arr.GetString(10) = %v, %v, want %v, nil", val, err, "expanded")
	}

	// 测试Remove方法
	arr.Remove(1)
	if arr.Size() != 10 {
		t.Errorf("After Remove(1), arr.Size() = %v, want %v", arr.Size(), 10)
	}
	if !arr.Get(1).IsString() { // 原来的索引2现在变成了1
		t.Errorf("After Remove(1), arr.Get(1).IsString() = %v, want %v", arr.Get(1).IsString(), true)
	}

	// 测试ToArray方法
	goArr := arr.ToArray()
	if len(goArr) != arr.Size() {
		t.Errorf("arr.ToArray() length = %v, want %v", len(goArr), arr.Size())
	}

	// 测试Join方法
	joined := arr.Join(",")
	if joined == "" {
		t.Errorf("arr.Join(\",\") returned empty string")
	}

	// 测试ForEach方法
	count := 0
	arr.ForEach(func(value JSONValue, index int) {
		count++
	})
	if count != arr.Size() {
		t.Errorf("ForEach count = %v, want %v", count, arr.Size())
	}

	// 测试Map方法
	mapped := arr.Map(func(value JSONValue, index int) JSONValue {
		if value.IsString() {
			str, _ := value.AsString()
			return NewJSONString(str + "_mapped")
		}
		return value
	})
	if mapped.Size() != arr.Size() {
		t.Errorf("mapped.Size() = %v, want %v", mapped.Size(), arr.Size())
	}
	if val, err := mapped.GetString(1); err != nil || val != "hello_mapped" {
		t.Errorf("mapped.GetString(1) = %v, %v, want %v, nil", val, err, "hello_mapped")
	}

	// 测试Filter方法
	filtered := arr.Filter(func(value JSONValue, index int) bool {
		return value.IsString()
	})
	if filtered.Size() == 0 {
		t.Errorf("filtered.Size() = 0, want > 0")
	}
	for i := 0; i < filtered.Size(); i++ {
		if !filtered.Get(i).IsString() {
			t.Errorf("filtered.Get(%d).IsString() = %v, want %v", i, filtered.Get(i).IsString(), true)
		}
	}

	// 测试Slice方法
	sliced := arr.Slice(1, 3)
	if sliced.Size() != 2 {
		t.Errorf("sliced.Size() = %v, want %v", sliced.Size(), 2)
	}
}

func TestNewJSONArrayFromValues(t *testing.T) {
	values := []JSONValue{
		NewJSONBool(true),
		NewJSONNumber(123),
		NewJSONString("hello"),
	}

	arr := NewJSONArrayFromValues(values)

	if arr.Size() != 3 {
		t.Errorf("arr.Size() = %v, want %v", arr.Size(), 3)
	}

	if val, err := arr.GetBoolean(0); err != nil || val != true {
		t.Errorf("arr.GetBoolean(0) = %v, %v, want %v, nil", val, err, true)
	}
	if val, err := arr.GetNumber(1); err != nil || val != 123 {
		t.Errorf("arr.GetNumber(1) = %v, %v, want %v, nil", val, err, 123)
	}
	if val, err := arr.GetString(2); err != nil || val != "hello" {
		t.Errorf("arr.GetString(2) = %v, %v, want %v, nil", val, err, "hello")
	}
}
