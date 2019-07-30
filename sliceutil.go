package main

import "strings"

// StringContains AAA
func StringContains(slice *[]string, item string) bool {
	for _, it := range *slice {
		if strings.Compare(it, item) == 0 {
			return true
		}

	}
	return false
}
