package generic

import (
	"testing"

	"github.com/UserLeeZJ/gojson/types"
)

// Test structs
type Person struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Address Address  `json:"address"`
	Hobbies []string `json:"hobbies"`
	Active  bool     `json:"active"`
}

type Address struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Zip     string `json:"zip,omitempty"`
}

// Custom type for testing
type CustomID string

func TestJSONObject(t *testing.T) {
	// Create a generic JSON object
	obj := NewJSONObject[map[string]interface{}]()
	obj.PutString("name", "John")
	obj.PutNumber("age", 30)

	// Create nested object
	address := types.NewJSONObject()
	address.PutString("city", "New York")
	address.PutString("country", "USA")
	obj.PutObject("address", address)

	// Create array
	hobbies := types.NewJSONArray()
	hobbies.AddString("reading").AddString("swimming")
	obj.PutArray("hobbies", hobbies)

	// Test boolean
	obj.PutBoolean("active", true)

	// Test null
	obj.Put("data", types.NewJSONNull())

	// Test Get methods
	name, err := obj.GetString("name")
	if err != nil {
		t.Errorf("Failed to get name: %v", err)
	}
	if name != "John" {
		t.Errorf("name mismatch: expected John, got %s", name)
	}

	age, err := obj.GetNumber("age")
	if err != nil {
		t.Errorf("Failed to get age: %v", err)
	}
	if age != 30 {
		t.Errorf("age mismatch: expected 30, got %f", age)
	}

	active, err := obj.GetBoolean("active")
	if err != nil {
		t.Errorf("Failed to get active: %v", err)
	}
	if !active {
		t.Errorf("active mismatch: expected true, got %v", active)
	}

	// Test GetTyped
	ageAny, err := obj.GetTyped("age")
	if err != nil {
		t.Errorf("Failed to get age: %v", err)
	}
	ageFloat, ok := ageAny.(float64)
	if !ok {
		t.Errorf("age is not a float64")
	}
	if ageFloat != 30 {
		t.Errorf("ageFloat mismatch: expected 30, got %f", ageFloat)
	}

	// Test Value
	value := obj.Value()

	// Check values
	if value["name"] != "John" {
		t.Errorf("value[name] mismatch: expected John, got %v", value["name"])
	}

	if value["age"].(float64) != 30 {
		t.Errorf("value[age] mismatch: expected 30, got %v", value["age"])
	}

	// Check nested object
	addressMap := value["address"].(map[string]interface{})
	if addressMap["city"] != "New York" {
		t.Errorf("addressMap[city] mismatch: expected New York, got %v", addressMap["city"])
	}

	// Check array
	hobbiesArr := value["hobbies"].([]interface{})
	if len(hobbiesArr) != 2 {
		t.Errorf("hobbies length mismatch: expected 2, got %d", len(hobbiesArr))
	}
	if hobbiesArr[0] != "reading" {
		t.Errorf("hobbies[0] mismatch: expected reading, got %v", hobbiesArr[0])
	}

	// Test Keys and Has
	keys := obj.Keys()
	if len(keys) != 6 {
		t.Errorf("keys length mismatch: expected 6, got %d", len(keys))
	}

	if !obj.Has("name") {
		t.Errorf("Has(name) should return true")
	}

	if obj.Has("nonexistent") {
		t.Errorf("Has(nonexistent) should return false")
	}

	// Test Size
	if obj.Size() != 6 {
		t.Errorf("Size mismatch: expected 6, got %d", obj.Size())
	}

	// Test Remove
	obj.Remove("data")
	if obj.Has("data") {
		t.Errorf("Remove(data) failed")
	}

	// Test ToMap
	m := obj.ToMap()
	if m["name"] != "John" {
		t.Errorf("ToMap()[name] mismatch: expected John, got %v", m["name"])
	}

	// Test ForEach
	count := 0
	obj.ForEach(func(key string, value types.JSONValue) {
		count++
	})
	if count != 5 {
		t.Errorf("ForEach count mismatch: expected 5, got %d", count)
	}

	// Test Clone
	clone := obj.Clone()
	if clone.Size() != 5 {
		t.Errorf("Clone size mismatch: expected 5, got %d", clone.Size())
	}

	// Test PutTyped
	obj2 := NewJSONObject[map[string]interface{}]()
	_, err = obj2.PutTyped("person", Person{
		Name: "Alice",
		Age:  25,
		Address: Address{
			City:    "Boston",
			Country: "USA",
		},
		Hobbies: []string{"coding", "gaming"},
		Active:  true,
	})
	if err != nil {
		t.Errorf("PutTyped failed: %v", err)
	}

	personObj, err := obj2.GetObject("person")
	if err != nil {
		t.Errorf("GetObject failed: %v", err)
	}

	personName, err := personObj.GetString("name")
	if err != nil {
		t.Errorf("GetString failed: %v", err)
	}
	if personName != "Alice" {
		t.Errorf("personName mismatch: expected Alice, got %s", personName)
	}

	// Test Merge
	obj3 := NewJSONObject[map[string]interface{}]()
	obj3.PutString("city", "San Francisco")
	obj3.PutNumber("population", 884363)
	
	obj4 := NewJSONObject[map[string]interface{}]()
	obj4.PutString("city", "Los Angeles")
	obj4.PutNumber("founded", 1781)
	
	obj3.Merge(obj4)
	
	city, err := obj3.GetString("city")
	if err != nil || city != "Los Angeles" {
		t.Errorf("Merge failed: city should be Los Angeles, got %s", city)
	}
	
	pop, err := obj3.GetNumber("population")
	if err != nil || pop != 884363 {
		t.Errorf("Merge failed: population should be 884363, got %f", pop)
	}
	
	founded, err := obj3.GetNumber("founded")
	if err != nil || founded != 1781 {
		t.Errorf("Merge failed: founded should be 1781, got %f", founded)
	}
}
