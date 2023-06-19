package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListContains(t *testing.T) {
	// Test Case 1: subList is empty
	list1 := []string{"apple", "banana", "orange"}
	subList1 := []string{}
	expectedResult1 := true
	result := ListContains(subList1, list1)
	assert.Equal(t, expectedResult1, result)

	// Test Case 2: subList is a subset of list
	list2 := []string{"apple", "banana", "orange"}
	subList2 := []string{"banana"}
	expectedResult2 := true
	result = ListContains(subList2, list2)
	assert.Equal(t, expectedResult2, result)

	// Test Case 3: subList is not a subset of list
	list3 := []string{"apple", "banana", "orange"}
	subList3 := []string{"banana", "mango"}
	expectedResult3 := false
	result = ListContains(subList3, list3)
	assert.Equal(t, expectedResult3, result)

	// Test Case 4: subList is longer than list
	list4 := []string{"apple", "banana", "orange"}
	subList4 := []string{"apple", "banana", "orange", "mango"}
	expectedResult4 := false
	result = ListContains(subList4, list4)
	assert.Equal(t, expectedResult4, result)
}

func TestIsElementExists(t *testing.T) {
	// Test Case 1: Element exists in the list
	list1 := []string{"apple", "banana", "orange"}
	element1 := "banana"
	expectedResult1 := true
	result := IsElementExists(list1, element1)
	assert.Equal(t, expectedResult1, result)

	// Test Case 2: Element does not exist in the list
	list2 := []string{"apple", "banana", "orange"}
	element2 := "mango"
	expectedResult2 := false
	result = IsElementExists(list2, element2)
	assert.Equal(t, expectedResult2, result)

	// Test Case 3: Empty list
	list3 := []string{}
	element3 := "apple"
	expectedResult3 := false
	result = IsElementExists(list3, element3)
	assert.Equal(t, expectedResult3, result)
}

func TestIsElementMatch(t *testing.T) {
	// Test Case 1: Element matches "*" wildcard
	elements1 := []string{"apple", "banana", "orange", "*"}
	element1 := "mango"
	expectedResult1 := true
	result := IsElementMatch(element1, elements1)
	assert.Equal(t, expectedResult1, result)

	// Test Case 2: Element matches a specific element in the list
	elements2 := []string{"apple", "banana", "orange"}
	element2 := "banana"
	expectedResult2 := true
	result = IsElementMatch(element2, elements2)
	assert.Equal(t, expectedResult2, result)

	// Test Case 3: Element does not match any element in the list
	elements3 := []string{"apple", "banana", "orange"}
	element3 := "mango"
	expectedResult3 := false
	result = IsElementMatch(element3, elements3)
	assert.Equal(t, expectedResult3, result)

	// Test Case 4: Element matches "*" wildcard but is not present in the list
	elements4 := []string{"apple", "banana", "orange", "*"}
	element4 := "grape"
	expectedResult4 := true
	result = IsElementMatch(element4, elements4)
	assert.Equal(t, expectedResult4, result)

	// Test Case 5: Empty list
	elements5 := []string{}
	element5 := "apple"
	expectedResult5 := false
	result = IsElementMatch(element5, elements5)
	assert.Equal(t, expectedResult5, result)
}

func TestAddPrefixToList(t *testing.T) {
	// Test Case 1: Add prefix to each item in the list
	list1 := []string{"apple", "banana", "orange"}
	prefix1 := "fruit_"
	expectedResult1 := []string{"fruit_apple", "fruit_banana", "fruit_orange"}
	result1 := AddPrefixToList(list1, prefix1)
	assert.Equal(t, expectedResult1, result1)

	// Test Case 2: Add empty prefix to each item in the list
	list2 := []string{"apple", "banana", "orange"}
	prefix2 := ""
	expectedResult2 := []string{"apple", "banana", "orange"}
	result2 := AddPrefixToList(list2, prefix2)
	assert.Equal(t, expectedResult2, result2)

	// Test Case 3: Empty list
	list3 := []string{}
	prefix3 := "prefix_"
	expectedResult3 := []string{}
	result3 := AddPrefixToList(list3, prefix3)
	assert.Equal(t, expectedResult3, result3)
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
		assert.Equal(t, value, result1[key])
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
		assert.Equal(t, value, result2[key])
	}

	// Test Case 3: Empty string
	str3 := ""
	result3 := StringToMap(str3)
	assert.Empty(t, result3)
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
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON1, string(resultJSON1))

	// Test Case 2: Empty YAML list
	yamlString2 := `[]`
	expectedJSON2 := `[]`
	resultJSON2, err := ConvertYAMLListToJSONList(yamlString2)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON2, string(resultJSON2))

	// Test Case 3: Invalid YAML format
	yamlString3 := `
name: John
age: 30
`
	resultJSON3, err := ConvertYAMLListToJSONList(yamlString3)
	assert.Error(t, err)
	assert.Nil(t, resultJSON3)
}

func TestConvertYAMToJSON(t *testing.T) {
	// Test Case 1: Valid YAML
	yamlString1 := []byte(`
name: John
age: 30
`)
	expectedJSON1 := `{"age":30,"name":"John"}`
	resultJSON1, err := ConvertYAMToJSON(yamlString1)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON1, string(resultJSON1))

	// Test Case 2: Empty YAML
	yamlString2 := []byte("")
	expectedJSON2 := `{}`
	resultJSON2, err := ConvertYAMToJSON(yamlString2)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON2, string(resultJSON2))

	// Test Case 3: Invalid YAML format
	yamlString3 := []byte(`
- name: John
  age: 30
`)
	resultJSON3, err := ConvertYAMToJSON(yamlString3)
	assert.Error(t, err)
	assert.Nil(t, resultJSON3)
}
