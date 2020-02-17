package array

import "strings"

func IncludesString(strs []string, item string) bool {
	for _, str := range strs {
		if strings.Contains(str, item) {
			return true
		}
	}
	return false
}
