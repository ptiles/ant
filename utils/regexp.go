package utils

import (
	"regexp"
)

func NamedMatches(rx, str string) map[string]string {
	expr := regexp.MustCompile(rx)
	match := expr.FindStringSubmatch(str)
	result := make(map[string]string)
	matchLen := len(match)

	for i, name := range expr.SubexpNames() {
		if i > matchLen {
			break
		}
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	return result
}
