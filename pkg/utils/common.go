package utils

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func ListContains(subList, list []string) bool {
	if len(subList) > len(list) {
		return false
	}
	for _, element := range subList {
		found := false
		for _, b := range list {
			if element == b {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func IsElementExists(list []string, element string) bool {
	for _, item := range list {
		if item == element {
			return true
		}
	}
	return false
}

func IsElementMatch(element string, elements []string) bool {
	if IsElementExists(elements, "*") {
		return true
	}

	return IsElementExists(elements, element)
}

func GetClientConfig(kubeConfig string) (*rest.Config, error) {
	if kubeConfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeConfig)
	}
	return rest.InClusterConfig()
}

func AddPrefixToList(list []string, prefix string) []string {
	result := make([]string, len(list))

	for i, item := range list {
		result[i] = prefix + item
	}

	return result
}

func StringToMap(str string) map[string]string {
	pairs := strings.Split(str, ",")
	m := make(map[string]string)

	for _, pair := range pairs {
		keyValue := strings.Split(pair, ":")
		if len(keyValue) == 2 {
			key := strings.TrimSpace(keyValue[0])
			value := strings.TrimSpace(keyValue[1])
			m[key] = value
		}
	}

	return m
}

func ConvertYAMLListToJSONList(yamlString string) ([]byte, error) {
	// Unmarshal YAML into a map[string]interface{}
	yamlData := make([]map[string]interface{}, 0)
	err := yaml.Unmarshal([]byte(yamlString), &yamlData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %v", err)
	}

	// Marshal the YAML data as JSON
	jsonBytes, err := json.Marshal(&yamlData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	return jsonBytes, nil
}

func ConvertYAMLToJSON(yamlString []byte) ([]byte, error) {
	// Unmarshal YAML into a map[string]interface{}
	yamlData := make(map[string]interface{})
	err := yaml.Unmarshal(yamlString, &yamlData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %v", err)
	}

	// Marshal the YAML data as JSON
	jsonBytes, err := json.Marshal(&yamlData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	return jsonBytes, nil
}

func SPtr(str string) *string {
	return &str
}

func BPtr(b bool) *bool {
	return &b
}

func IPtr(i int64) *int64 {
	return &i
}

func ValidateHTTPFormat(input string) bool {
	regex := `^(https?://)([\w-]+(\.[\w-]+)*)(:\d+)?(/[\w-./?%&=]*)?$`
	match, _ := regexp.MatchString(regex, input)
	return match
}

func TrimString(s string, maxLength int) string {
	if maxLength >= len(s) {
		return s
	}
	return s[:maxLength]
}

func StringToInt64(input string) int64 {
	h := fnv.New64a()
	h.Write([]byte(input))
	hashValue := h.Sum64()

	// Convert the hash value to int64
	int64Value := int64(hashValue)
	// Make sure the value is positive (int64 can represent only non-negative values)
	if int64Value < 0 {
		int64Value = int64Value * -1
	}

	return int64Value
}

func RemoveBraces(input string) string {
	output := strings.ReplaceAll(input, "{", "")
	output = strings.ReplaceAll(output, "}", "")
	return output
}

func ExtractStringsBetweenTags(input string) []string {
	re := regexp.MustCompile(`<([^>]+)>`)
	matches := re.FindAllStringSubmatch(input, -1)

	var result []string
	for _, match := range matches {
		result = append(result, match[1])
	}

	return result
}
