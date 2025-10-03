package utils

import (
	"regexp"
	"strconv"
)

func NamedStringMatches(expr *regexp.Regexp, str string) map[string]string {
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

func NamedIntMatches(expr *regexp.Regexp, str string) map[string]int {
	if !expr.MatchString(str) {
		return nil
	}

	match := expr.FindStringSubmatch(str)
	result := make(map[string]int)
	matchLen := len(match)

	for i, name := range expr.SubexpNames() {
		if i > matchLen {
			break
		}
		if i != 0 && name != "" {
			result[name], _ = strconv.Atoi(match[i])
		}
	}

	return result
}
