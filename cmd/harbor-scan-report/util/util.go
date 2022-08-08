package util

import "strings"

func IsStringEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func IsStringPresent(s string) bool {
	return !IsStringEmpty(s)
}
