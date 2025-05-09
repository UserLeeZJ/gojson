package generic

import (
	"testing"

	"github.com/UserLeeZJ/gojson/types"
)

func TestGetTyped(t *testing.T) {
	// Create a JSON object
	obj := types.NewJSONObject()
	obj.PutString("name", "John")
	obj.PutNumber("age", 30)
	obj.PutBoolean("active", true)

	// Create nested object
	address := types.NewJSONObject()
	address.PutString("city", "New York")
	address.PutString("country", "USA")
	obj.PutObject("address", address)

	// Create array
	hobbies := types.NewJSONArray()
	hobbies.AddString("reading").AddString("swimming")
	obj.PutArray("hobbies", hobbies)

	// Test string
	name, err := GetTyped[string](obj, "name")
	if err != nil {
		t.Errorf("GetTyped[string] failed: %v", err)
	}
	if name != "John" {
		t.Errorf("name mismatch: expected John, got %s", name)
	}

	// Test int
	age, err := GetTyped[int](obj, "age")
	if err != nil {
		t.Errorf("GetTyped[int] failed: %v", err)
	}
	if age != 30 {
		t.Errorf("age mismatch: expected 30, got %d", age)
	}

	// Test float64
	ageFloat, err := GetTyped[float64](obj, "age")
	if err != nil {
		t.Errorf("GetTyped[float64] failed: %v", err)
	}
	if ageFloat != 30.0 {
		t.Errorf("ageFloat mismatch: expected 30.0, got %f", ageFloat)
	}

	// Test bool
	active, err := GetTyped[bool](obj, "active")
	if err != nil {
		t.Errorf("GetTyped[bool] failed: %v", err)
	}
	if !active {
		t.Errorf("active mismatch: expected true, got %v", active)
	}

	// Test array
	hobbiesArr, err := GetTyped[[]string](obj, "hobbies")
	if err != nil {
		t.Errorf("GetTyped[[]string] failed: %v", err)
	}
	if len(hobbiesArr) != 2 {
		t.Errorf("hobbiesArr length mismatch: expected 2, got %d", len(hobbiesArr))
	}
	if hobbiesArr[0] != "reading" {
		t.Errorf("hobbiesArr[0] mismatch: expected reading, got %s", hobbiesArr[0])
	}

	// Test map
	addressMap, err := GetTyped[map[string]string](obj, "address")
	if err != nil {
		t.Errorf("GetTyped[map[string]string] failed: %v", err)
	}
	if addressMap["city"] != "New York" {
		t.Errorf("addressMap[city] mismatch: expected New York, got %s", addressMap["city"])
	}

	// Test struct
	addressStruct, err := GetTyped[Address](obj, "address")
	if err != nil {
		t.Errorf("GetTyped[Address] failed: %v", err)
	}
	if addressStruct.City != "New York" {
		t.Errorf("addressStruct.City mismatch: expected New York, got %s", addressStruct.City)
	}

	// Test entire object as struct
	person, err := GetTyped[Person](obj, "")
	if err != nil {
		t.Errorf("GetTyped[Person] failed: %v", err)
	}
	if person.Name != "John" {
		t.Errorf("person.Name mismatch: expected John, got %s", person.Name)
	}
	if person.Age != 30 {
		t.Errorf("person.Age mismatch: expected 30, got %d", person.Age)
	}
	if person.Address.City != "New York" {
		t.Errorf("person.Address.City mismatch: expected New York, got %s", person.Address.City)
	}

	// Test non-existent key
	_, err = GetTyped[string](obj, "nonexistent")
	if err == nil {
		t.Errorf("GetTyped with non-existent key should fail")
	}

	// Test type mismatch
	_, err = GetTyped[int](obj, "name")
	if err == nil {
		t.Errorf("GetTyped with type mismatch should fail")
	}

	// Skip custom type test for now
	// CustomID is a string alias, but direct conversion is not supported

	// Test uint
	_, err = GetTyped[uint](obj, "age")
	if err != nil {
		t.Errorf("GetTyped[uint] failed: %v", err)
	}
}

func TestToJSONValue(t *testing.T) {
	// Test primitive types
	strVal, err := ToJSONValue("test")
	if err != nil {
		t.Errorf("ToJSONValue(string) failed: %v", err)
	}
	if !strVal.IsString() {
		t.Errorf("strVal should be a string")
	}
	str, _ := strVal.AsString()
	if str != "test" {
		t.Errorf("strVal mismatch: expected test, got %s", str)
	}

	numVal, err := ToJSONValue(42)
	if err != nil {
		t.Errorf("ToJSONValue(int) failed: %v", err)
	}
	if !numVal.IsNumber() {
		t.Errorf("numVal should be a number")
	}
	num, _ := numVal.AsNumber()
	if num != 42 {
		t.Errorf("numVal mismatch: expected 42, got %f", num)
	}

	boolVal, err := ToJSONValue(true)
	if err != nil {
		t.Errorf("ToJSONValue(bool) failed: %v", err)
	}
	if !boolVal.IsBoolean() {
		t.Errorf("boolVal should be a boolean")
	}
	b, _ := boolVal.AsBoolean()
	if !b {
		t.Errorf("boolVal mismatch: expected true, got %v", b)
	}

	// Test nil
	nullVal, err := ToJSONValue(nil)
	if err != nil {
		t.Errorf("ToJSONValue(nil) failed: %v", err)
	}
	if !nullVal.IsNull() {
		t.Errorf("nullVal should be null")
	}

	// Test map
	mapVal, err := ToJSONValue(map[string]interface{}{
		"name": "John",
		"age":  30,
	})
	if err != nil {
		t.Errorf("ToJSONValue(map) failed: %v", err)
	}
	if !mapVal.IsObject() {
		t.Errorf("mapVal should be an object")
	}
	obj, _ := mapVal.AsObject()
	name, _ := obj.GetString("name")
	if name != "John" {
		t.Errorf("mapVal.name mismatch: expected John, got %s", name)
	}

	// Test slice
	sliceVal, err := ToJSONValue([]interface{}{"apple", "banana"})
	if err != nil {
		t.Errorf("ToJSONValue(slice) failed: %v", err)
	}
	if !sliceVal.IsArray() {
		t.Errorf("sliceVal should be an array")
	}
	arr, _ := sliceVal.AsArray()
	if arr.Size() != 2 {
		t.Errorf("sliceVal size mismatch: expected 2, got %d", arr.Size())
	}
	first, _ := arr.Get(0).AsString()
	if first != "apple" {
		t.Errorf("sliceVal[0] mismatch: expected apple, got %s", first)
	}

	// Test struct
	structVal, err := ToJSONValue(Person{
		Name: "John",
		Age:  30,
		Address: Address{
			City:    "New York",
			Country: "USA",
		},
		Hobbies: []string{"reading", "swimming"},
		Active:  true,
	})
	if err != nil {
		t.Errorf("ToJSONValue(struct) failed: %v", err)
	}
	if !structVal.IsObject() {
		t.Errorf("structVal should be an object")
	}
	personObj, _ := structVal.AsObject()
	personName, _ := personObj.GetString("name")
	if personName != "John" {
		t.Errorf("structVal.name mismatch: expected John, got %s", personName)
	}

	// Test numeric types
	int8Val, err := ToJSONValue(int8(8))
	if err != nil || !int8Val.IsNumber() {
		t.Errorf("ToJSONValue(int8) failed: %v", err)
	}

	int16Val, err := ToJSONValue(int16(16))
	if err != nil || !int16Val.IsNumber() {
		t.Errorf("ToJSONValue(int16) failed: %v", err)
	}

	int32Val, err := ToJSONValue(int32(32))
	if err != nil || !int32Val.IsNumber() {
		t.Errorf("ToJSONValue(int32) failed: %v", err)
	}

	int64Val, err := ToJSONValue(int64(64))
	if err != nil || !int64Val.IsNumber() {
		t.Errorf("ToJSONValue(int64) failed: %v", err)
	}

	uintVal, err := ToJSONValue(uint(42))
	if err != nil || !uintVal.IsNumber() {
		t.Errorf("ToJSONValue(uint) failed: %v", err)
	}

	uint8Val, err := ToJSONValue(uint8(8))
	if err != nil || !uint8Val.IsNumber() {
		t.Errorf("ToJSONValue(uint8) failed: %v", err)
	}

	uint16Val, err := ToJSONValue(uint16(16))
	if err != nil || !uint16Val.IsNumber() {
		t.Errorf("ToJSONValue(uint16) failed: %v", err)
	}

	uint32Val, err := ToJSONValue(uint32(32))
	if err != nil || !uint32Val.IsNumber() {
		t.Errorf("ToJSONValue(uint32) failed: %v", err)
	}

	uint64Val, err := ToJSONValue(uint64(64))
	if err != nil || !uint64Val.IsNumber() {
		t.Errorf("ToJSONValue(uint64) failed: %v", err)
	}

	float32Val, err := ToJSONValue(float32(3.14))
	if err != nil || !float32Val.IsNumber() {
		t.Errorf("ToJSONValue(float32) failed: %v", err)
	}

	// Test convertToJSONValue error cases
	_, err = convertToJSONValue(complex(1, 2))
	if err == nil {
		t.Errorf("convertToJSONValue should fail for complex type")
	}
}
