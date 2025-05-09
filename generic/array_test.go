package generic

import (
	"testing"

	"github.com/UserLeeZJ/gojson/types"
)

func TestJSONArray(t *testing.T) {
	// Create a generic JSON array
	arr := NewJSONArray[string]()
	arr.Add(types.NewJSONString("apple"))
	arr.Add(types.NewJSONString("banana"))
	arr.Add(types.NewJSONString("orange"))

	// Test AddTyped
	_, err := arr.AddTyped("grape")
	if err != nil {
		t.Errorf("AddTyped failed: %v", err)
	}

	// Test Size
	if arr.Size() != 4 {
		t.Errorf("Size mismatch: expected 4, got %d", arr.Size())
	}

	// Test Get
	apple, err := arr.Get(0).AsString()
	if err != nil {
		t.Errorf("Get(0).AsString() failed: %v", err)
	}
	if apple != "apple" {
		t.Errorf("apple mismatch: expected apple, got %s", apple)
	}

	// Test GetTyped
	grape, err := arr.GetTyped(3)
	if err != nil {
		t.Errorf("GetTyped(3) failed: %v", err)
	}
	if grape != "grape" {
		t.Errorf("grape mismatch: expected grape, got %s", grape)
	}

	// Test Value
	value := arr.Value()

	// Check values
	if len(value) != 4 {
		t.Errorf("value length mismatch: expected 4, got %d", len(value))
	}
	if value[0] != "apple" {
		t.Errorf("value[0] mismatch: expected apple, got %s", value[0])
	}
	if value[3] != "grape" {
		t.Errorf("value[3] mismatch: expected grape, got %s", value[3])
	}

	// Test AddString, AddNumber, AddBoolean, AddNull
	arr2 := NewJSONArray[interface{}]()
	arr2.AddString("string")
	arr2.AddNumber(123)
	arr2.AddBoolean(true)
	arr2.AddNull()

	if arr2.Size() != 4 {
		t.Errorf("arr2 size mismatch: expected 4, got %d", arr2.Size())
	}

	// Test Set
	arr.Set(1, types.NewJSONString("pear"))
	pear, err := arr.Get(1).AsString()
	if err != nil {
		t.Errorf("Get(1).AsString() failed: %v", err)
	}
	if pear != "pear" {
		t.Errorf("pear mismatch: expected pear, got %s", pear)
	}

	// Test SetTyped
	_, err = arr.SetTyped(2, "mango")
	if err != nil {
		t.Errorf("SetTyped failed: %v", err)
	}
	mango, err := arr.GetTyped(2)
	if err != nil {
		t.Errorf("GetTyped(2) failed: %v", err)
	}
	if mango != "mango" {
		t.Errorf("mango mismatch: expected mango, got %s", mango)
	}

	// Test Remove
	arr.Remove(0)
	if arr.Size() != 3 {
		t.Errorf("Size after Remove mismatch: expected 3, got %d", arr.Size())
	}
	firstItem, err := arr.GetTyped(0)
	if err != nil {
		t.Errorf("GetTyped(0) failed: %v", err)
	}
	if firstItem != "pear" {
		t.Errorf("firstItem mismatch: expected pear, got %s", firstItem)
	}

	// Test ToArray
	array := arr.ToArray()
	if len(array) != 3 {
		t.Errorf("ToArray length mismatch: expected 3, got %d", len(array))
	}

	// Test ForEach
	count := 0
	arr.ForEach(func(value types.JSONValue, index int) {
		count++
	})
	if count != 3 {
		t.Errorf("ForEach count mismatch: expected 3, got %d", count)
	}

	// Test Map
	mappedArr := arr.Map(func(value types.JSONValue, index int) types.JSONValue {
		str, _ := value.AsString()
		return types.NewJSONString(str + "_mapped")
	})
	mappedFirst, err := mappedArr.GetTyped(0)
	if err != nil {
		t.Errorf("mappedArr.GetTyped(0) failed: %v", err)
	}
	if mappedFirst != "pear_mapped" {
		t.Errorf("mappedFirst mismatch: expected pear_mapped, got %s", mappedFirst)
	}

	// Test Filter
	filteredArr := arr.Filter(func(value types.JSONValue, index int) bool {
		str, _ := value.AsString()
		return str == "mango"
	})
	if filteredArr.Size() != 1 {
		t.Errorf("filteredArr size mismatch: expected 1, got %d", filteredArr.Size())
	}
	filteredItem, err := filteredArr.GetTyped(0)
	if err != nil {
		t.Errorf("filteredArr.GetTyped(0) failed: %v", err)
	}
	if filteredItem != "mango" {
		t.Errorf("filteredItem mismatch: expected mango, got %s", filteredItem)
	}

	// Test Slice
	slicedArr := arr.Slice(1, 3)
	if slicedArr.Size() != 2 {
		t.Errorf("slicedArr size mismatch: expected 2, got %d", slicedArr.Size())
	}
	slicedFirst, err := slicedArr.GetTyped(0)
	if err != nil {
		t.Errorf("slicedArr.GetTyped(0) failed: %v", err)
	}
	if slicedFirst != "mango" {
		t.Errorf("slicedFirst mismatch: expected mango, got %s", slicedFirst)
	}
	
	// Test JSONValue
	jsonValue := arr.JSONValue()
	if !jsonValue.IsArray() {
		t.Errorf("JSONValue should be an array")
	}
	jsonArr, err := jsonValue.AsArray()
	if err != nil {
		t.Errorf("JSONValue.AsArray() failed: %v", err)
	}
	if jsonArr.Size() != 3 {
		t.Errorf("JSONValue size mismatch: expected 3, got %d", jsonArr.Size())
	}
}
