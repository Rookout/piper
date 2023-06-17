package utils

import (
	"reflect"
	"testing"
)

func TestListContains(t *testing.T) {
	// Test Case 1: subList is empty
	list1 := []string{"apple", "banana", "orange"}
	subList1 := []string{}
	expectedResult1 := true
	if result := ListContains(subList1, list1); result != expectedResult1 {
		t.Errorf("Expected %v, but got %v", expectedResult1, result)
	}

	// Test Case 2: subList is a subset of list
	list2 := []string{"apple", "banana", "orange"}
	subList2 := []string{"banana"}
	expectedResult2 := true
	if result := ListContains(subList2, list2); result != expectedResult2 {
		t.Errorf("Expected %v, but got %v", expectedResult2, result)
	}

	// Test Case 3: subList is not a subset of list
	list3 := []string{"apple", "banana", "orange"}
	subList3 := []string{"banana", "mango"}
	expectedResult3 := false
	if result := ListContains(subList3, list3); result != expectedResult3 {
		t.Errorf("Expected %v, but got %v", expectedResult3, result)
	}

	// Test Case 4: subList is longer than list
	list4 := []string{"apple", "banana", "orange"}
	subList4 := []string{"apple", "banana", "orange", "mango"}
	expectedResult4 := false
	if result := ListContains(subList4, list4); result != expectedResult4 {
		t.Errorf("Expected %v, but got %v", expectedResult4, result)
	}
}

func TestIsElementExists(t *testing.T) {
	// Test Case 1: Element exists in the list
	list1 := []string{"apple", "banana", "orange"}
	element1 := "banana"
	expectedResult1 := true
	if result := IsElementExists(list1, element1); result != expectedResult1 {
		t.Errorf("Expected %v, but got %v", expectedResult1, result)
	}

	// Test Case 2: Element does not exist in the list
	list2 := []string{"apple", "banana", "orange"}
	element2 := "mango"
	expectedResult2 := false
	if result := IsElementExists(list2, element2); result != expectedResult2 {
		t.Errorf("Expected %v, but got %v", expectedResult2, result)
	}

	// Test Case 3: Empty list
	list3 := []string{}
	element3 := "apple"
	expectedResult3 := false
	if result := IsElementExists(list3, element3); result != expectedResult3 {
		t.Errorf("Expected %v, but got %v", expectedResult3, result)
	}
}

func TestIsElementMatch(t *testing.T) {
	// Test Case 1: Element matches "*" wildcard
	elements1 := []string{"apple", "banana", "orange", "*"}
	element1 := "mango"
	expectedResult1 := true
	if result := IsElementMatch(element1, elements1); result != expectedResult1 {
		t.Errorf("Expected %v, but got %v", expectedResult1, result)
	}

	// Test Case 2: Element matches a specific element in the list
	elements2 := []string{"apple", "banana", "orange"}
	element2 := "banana"
	expectedResult2 := true
	if result := IsElementMatch(element2, elements2); result != expectedResult2 {
		t.Errorf("Expected %v, but got %v", expectedResult2, result)
	}

	// Test Case 3: Element does not match any element in the list
	elements3 := []string{"apple", "banana", "orange"}
	element3 := "mango"
	expectedResult3 := false
	if result := IsElementMatch(element3, elements3); result != expectedResult3 {
		t.Errorf("Expected %v, but got %v", expectedResult3, result)
	}

	// Test Case 4: Element matches "*" wildcard but is not present in the list
	elements4 := []string{"apple", "banana", "orange", "*"}
	element4 := "grape"
	expectedResult4 := true
	if result := IsElementMatch(element4, elements4); result != expectedResult4 {
		t.Errorf("Expected %v, but got %v", expectedResult4, result)
	}

	// Test Case 5: Empty list
	elements5 := []string{}
	element5 := "apple"
	expectedResult5 := false
	if result := IsElementMatch(element5, elements5); result != expectedResult5 {
		t.Errorf("Expected %v, but got %v", expectedResult5, result)
	}
}

func TestAddPrefixToList(t *testing.T) {
	// Test Case 1: Add prefix to each item in the list
	list1 := []string{"apple", "banana", "orange"}
	prefix1 := "fruit_"
	expectedResult1 := []string{"fruit_apple", "fruit_banana", "fruit_orange"}
	result1 := AddPrefixToList(list1, prefix1)
	if !reflect.DeepEqual(result1, expectedResult1) {
		t.Errorf("Expected %v, but got %v", expectedResult1, result1)
	}

	// Test Case 2: Add empty prefix to each item in the list
	list2 := []string{"apple", "banana", "orange"}
	prefix2 := ""
	expectedResult2 := []string{"apple", "banana", "orange"}
	result2 := AddPrefixToList(list2, prefix2)
	if !reflect.DeepEqual(result2, expectedResult2) {
		t.Errorf("Expected %v, but got %v", expectedResult2, result2)
	}

	// Test Case 3: Empty list
	list3 := []string{}
	prefix3 := "prefix_"
	expectedResult3 := []string{}
	result3 := AddPrefixToList(list3, prefix3)
	if !reflect.DeepEqual(result3, expectedResult3) {
		t.Errorf("Expected %v, but got %v", expectedResult3, result3)
	}
}

func TestStringToMap(t *testing.T) {
	// Test Case 1: Valid string with key-value pairs
	str1 := "key1:value1, key2:value2, key3:value3"
	expectedResult1 := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	result1 := StringToMap(str1)
	for key, value := range expectedResult1 {
		if result1[key] != value {
			t.Errorf("Expected value '%v' for key '%v', but got '%v'", value, key, result1[key])
		}
	}

	// Test Case 2: Valid string with empty key-value pairs
	str2 := "key1:value1,, key2:value2, :value3"
	expectedResult2 := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"":     "value3",
	}
	result2 := StringToMap(str2)
	for key, value := range expectedResult2 {
		if result2[key] != value {
			t.Errorf("Expected value '%v' for key '%v', but got '%v'", value, key, result2[key])
		}
	}

	// Test Case 3: Empty string
	str3 := ""
	result3 := StringToMap(str3)
	if len(result3) != 0 {
		t.Errorf("Expected an empty map, but got '%v'", result3)
	}

}

func TestConvertYAMLListToJSONList(t *testing.T) {
	// Test Case 1: Valid YAML list
	yamlString1 := `
- name: John
  age: 30
- name: Jane
  age: 25
`
	expectedJSON1 := `[{"age":30,"name":"John"},{"age":25,"name":"Jane"}]`
	resultJSON1, err := ConvertYAMLListToJSONList(yamlString1)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
	if string(resultJSON1) != expectedJSON1 {
		t.Errorf("Expected JSON: %s\nGot JSON: %s", expectedJSON1, string(resultJSON1))
	}

	// Test Case 2: Empty YAML list
	yamlString2 := `[]`
	expectedJSON2 := `[]`
	resultJSON2, err := ConvertYAMLListToJSONList(yamlString2)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
	if string(resultJSON2) != expectedJSON2 {
		t.Errorf("Expected JSON: %s\nGot JSON: %s", expectedJSON2, string(resultJSON2))
	}

	// Test Case 3: Invalid YAML format
	yamlString3 := `
name: John
age: 30
`
	resultJSON3, err := ConvertYAMLListToJSONList(yamlString3)
	if err == nil {
		t.Errorf("Expected  to get error, but got %v", err)
	}
	if resultJSON3 != nil {
		t.Errorf("Expected JSON to be nil, but got: %s", string(resultJSON3))
	}
}

func TestConvertYAMToJSON(t *testing.T) {
	// Test Case 1: Valid YAML
	yamlString1 := []byte(`
name: John
age: 30
`)
	expectedJSON1 := `{"age":30,"name":"John"}`
	resultJSON1, err := ConvertYAMToJSON(yamlString1)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
	if string(resultJSON1) != expectedJSON1 {
		t.Errorf("Expected JSON: %s\nGot JSON: %s", expectedJSON1, string(resultJSON1))
	}

	// Test Case 2: Empty YAML
	yamlString2 := []byte("")
	expectedJSON2 := `{}`
	resultJSON2, err := ConvertYAMToJSON(yamlString2)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
	if string(resultJSON2) != expectedJSON2 {
		t.Errorf("Expected JSON: %s\nGot JSON: %s", expectedJSON2, string(resultJSON2))
	}

	// Test Case 3: Invalid YAML format
	yamlString3 := []byte(`
- name: John
  age: 30
`)
	resultJSON3, err := ConvertYAMToJSON(yamlString3)
	if err == nil {
		t.Errorf("Expected  to get error, but got %v", err)
	}
	if resultJSON3 != nil {
		t.Errorf("Expected JSON to be nil, but got: %s", string(resultJSON3))
	}
}
