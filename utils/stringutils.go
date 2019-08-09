package utils

import "regexp"

func IsBlankString(str string) bool {
	match, _ := regexp.MatchString("^(\\s*)$", str)
	return match
}

func IsNotBlankString(str string) bool {
	return !IsBlankString(str)
}
