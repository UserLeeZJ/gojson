package types

import (
	"testing"
)

func TestJSONObject(t *testing.T) {
	// 创建空对象
	obj := NewJSONObject()

	// 测试类型
	if obj.Type() != "object" {
		t.Errorf("JSONObject.Type() = %v, want %v", obj.Type(), "object")
	}

	// 测试空对象的字符串表示
	if obj.String() != "{}" {
		t.Errorf("Empty JSONObject.String() = %v, want %v", obj.String(), "{}")
	}

	// 测试类型检查
	if obj.IsNull() {
		t.Errorf("JSONObject.IsNull() = %v, want %v", obj.IsNull(), false)
	}
	if obj.IsBoolean() {
		t.Errorf("JSONObject.IsBoolean() = %v, want %v", obj.IsBoolean(), false)
	}
	if obj.IsNumber() {
		t.Errorf("JSONObject.IsNumber() = %v, want %v", obj.IsNumber(), false)
	}
	if obj.IsString() {
		t.Errorf("JSONObject.IsString() = %v, want %v", obj.IsString(), false)
	}
	if obj.IsArray() {
		t.Errorf("JSONObject.IsArray() = %v, want %v", obj.IsArray(), false)
	}
	if !obj.IsObject() {
		t.Errorf("JSONObject.IsObject() = %v, want %v", obj.IsObject(), true)
	}

	// 测试转换方法
	if _, err := obj.AsBoolean(); err == nil {
		t.Errorf("JSONObject.AsBoolean() should return error")
	}
	if _, err := obj.AsNumber(); err == nil {
		t.Errorf("JSONObject.AsNumber() should return error")
	}
	if val, err := obj.AsString(); err != nil || val != "{}" {
		t.Errorf("JSONObject.AsString() = %v, %v, want %v, nil", val, err, "{}")
	}
	if _, err := obj.AsArray(); err == nil {
		t.Errorf("JSONObject.AsArray() should return error")
	}
	if val, err := obj.AsObject(); err != nil || val != obj {
		t.Errorf("JSONObject.AsObject() = %v, %v, want %v, nil", val, err, obj)
	}

	// 测试添加属性
	obj.PutBoolean("bool", true).PutNumber("num", 123).PutString("str", "hello").PutNull("null")

	// 测试Size方法
	if obj.Size() != 4 {
		t.Errorf("JSONObject.Size() = %v, want %v", obj.Size(), 4)
	}

	// 测试Has方法
	if !obj.Has("bool") {
		t.Errorf("obj.Has(\"bool\") = %v, want %v", obj.Has("bool"), true)
	}
	if obj.Has("nonexistent") {
		t.Errorf("obj.Has(\"nonexistent\") = %v, want %v", obj.Has("nonexistent"), false)
	}

	// 测试Get方法
	if !obj.Get("bool").IsBoolean() {
		t.Errorf("obj.Get(\"bool\").IsBoolean() = %v, want %v", obj.Get("bool").IsBoolean(), true)
	}
	if !obj.Get("num").IsNumber() {
		t.Errorf("obj.Get(\"num\").IsNumber() = %v, want %v", obj.Get("num").IsNumber(), true)
	}
	if !obj.Get("str").IsString() {
		t.Errorf("obj.Get(\"str\").IsString() = %v, want %v", obj.Get("str").IsString(), true)
	}
	if !obj.Get("null").IsNull() {
		t.Errorf("obj.Get(\"null\").IsNull() = %v, want %v", obj.Get("null").IsNull(), true)
	}
	if !obj.Get("nonexistent").IsNull() { // 不存在的键应返回null
		t.Errorf("obj.Get(\"nonexistent\").IsNull() = %v, want %v", obj.Get("nonexistent").IsNull(), true)
	}

	// 测试GetBoolean方法
	if val, err := obj.GetBoolean("bool"); err != nil || val != true {
		t.Errorf("obj.GetBoolean(\"bool\") = %v, %v, want %v, nil", val, err, true)
	}
	if _, err := obj.GetBoolean("nonexistent"); err == nil {
		t.Errorf("obj.GetBoolean(\"nonexistent\") should return error")
	}

	// 测试GetNumber方法
	if val, err := obj.GetNumber("num"); err != nil || val != 123 {
		t.Errorf("obj.GetNumber(\"num\") = %v, %v, want %v, nil", val, err, 123)
	}
	if _, err := obj.GetNumber("nonexistent"); err == nil {
		t.Errorf("obj.GetNumber(\"nonexistent\") should return error")
	}

	// 测试GetString方法
	if val, err := obj.GetString("str"); err != nil || val != "hello" {
		t.Errorf("obj.GetString(\"str\") = %v, %v, want %v, nil", val, err, "hello")
	}
	if _, err := obj.GetString("nonexistent"); err == nil {
		t.Errorf("obj.GetString(\"nonexistent\") should return error")
	}

	// 测试嵌套对象和数组
	nestedObj := NewJSONObject()
	nestedObj.PutString("nested", "value")
	obj.PutObject("obj", nestedObj)

	nestedArr := NewJSONArray()
	nestedArr.AddNumber(1).AddNumber(2).AddNumber(3)
	obj.PutArray("arr", nestedArr)

	// 测试GetObject方法
	if val, err := obj.GetObject("obj"); err != nil {
		t.Errorf("obj.GetObject(\"obj\") error = %v", err)
	} else {
		if nestedVal, err := val.GetString("nested"); err != nil || nestedVal != "value" {
			t.Errorf("obj.GetObject(\"obj\").GetString(\"nested\") = %v, %v, want %v, nil", nestedVal, err, "value")
		}
	}
	if _, err := obj.GetObject("nonexistent"); err == nil {
		t.Errorf("obj.GetObject(\"nonexistent\") should return error")
	}

	// 测试GetArray方法
	if val, err := obj.GetArray("arr"); err != nil {
		t.Errorf("obj.GetArray(\"arr\") error = %v", err)
	} else {
		if val.Size() != 3 {
			t.Errorf("obj.GetArray(\"arr\").Size() = %v, want %v", val.Size(), 3)
		}
	}
	if _, err := obj.GetArray("nonexistent"); err == nil {
		t.Errorf("obj.GetArray(\"nonexistent\") should return error")
	}

	// 测试Keys方法
	keys := obj.Keys()
	if len(keys) != 6 {
		t.Errorf("obj.Keys() length = %v, want %v", len(keys), 6)
	}

	// 测试SortedKeys方法
	sortedKeys := obj.SortedKeys()
	if len(sortedKeys) != 6 {
		t.Errorf("obj.SortedKeys() length = %v, want %v", len(sortedKeys), 6)
	}
	// 检查是否已排序
	for i := 0; i < len(sortedKeys)-1; i++ {
		if sortedKeys[i] > sortedKeys[i+1] {
			t.Errorf("obj.SortedKeys() not sorted: %v > %v", sortedKeys[i], sortedKeys[i+1])
		}
	}

	// 测试Remove方法
	obj.Remove("num")
	if obj.Has("num") {
		t.Errorf("After Remove, obj.Has(\"num\") = %v, want %v", obj.Has("num"), false)
	}
	if obj.Size() != 5 {
		t.Errorf("After Remove, obj.Size() = %v, want %v", obj.Size(), 5)
	}

	// 测试ToMap方法
	m := obj.ToMap()
	if len(m) != 5 {
		t.Errorf("obj.ToMap() length = %v, want %v", len(m), 5)
	}
	if m["bool"] != true {
		t.Errorf("obj.ToMap()[\"bool\"] = %v, want %v", m["bool"], true)
	}
	if m["str"] != "hello" {
		t.Errorf("obj.ToMap()[\"str\"] = %v, want %v", m["str"], "hello")
	}

	// 测试ForEach方法
	count := 0
	obj.ForEach(func(key string, value JSONValue) {
		count++
	})
	if count != 5 {
		t.Errorf("ForEach count = %v, want %v", count, 5)
	}

	// 测试Merge方法
	other := NewJSONObject()
	other.PutString("newKey", "newValue")
	other.PutString("str", "overridden") // 覆盖现有键
	obj.Merge(other)
	if !obj.Has("newKey") {
		t.Errorf("After Merge, obj.Has(\"newKey\") = %v, want %v", obj.Has("newKey"), true)
	}
	if val, _ := obj.GetString("str"); val != "overridden" {
		t.Errorf("After Merge, obj.GetString(\"str\") = %v, want %v", val, "overridden")
	}

	// 测试Clone方法
	clone := obj.Clone()
	if clone.Size() != obj.Size() {
		t.Errorf("clone.Size() = %v, want %v", clone.Size(), obj.Size())
	}
	// 修改克隆不应影响原对象
	clone.PutString("cloneOnly", "value")
	if obj.Has("cloneOnly") {
		t.Errorf("After modifying clone, obj.Has(\"cloneOnly\") = %v, want %v", obj.Has("cloneOnly"), false)
	}
}
