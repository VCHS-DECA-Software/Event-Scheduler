package utils

import "strings"

func SplitAndTrim(str string, sep string) []string {
	split := strings.Split(str, sep)
	for i := 0; i < len(split); i++ {
		split[i] = strings.Trim(split[i], " ")
	}
	return split
}
