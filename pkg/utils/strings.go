package utils

import "strings"

func StringToList(input string) []string {
	return strings.Split(strings.ReplaceAll(input, " ", ""), ",")
}
