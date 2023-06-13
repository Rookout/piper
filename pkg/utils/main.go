package utils

import (
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
