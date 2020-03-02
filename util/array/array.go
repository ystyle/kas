package array

import "strings"

func IncludesString(strs []string, item string) bool {
	for _, str := range strs {
		if str == item {
			return true
		}
	}
	return false
}

func IncludesFromString(strs []string, item string) bool {
	for _, str := range strs {
		if strings.Contains(item, str) {
			return true
		}
	}
	return false
}
