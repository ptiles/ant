package utils

import (
	"regexp"
	"strconv"
)

func NamedStringMatches(rx, str string) map[string]string {
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

func NamedIntMatches(expr *regexp.Regexp, str string) map[string]int {
	if !expr.MatchString(str) {
		return nil
	}

	match := expr.FindStringSubmatch(str)
	result := make(map[string]int)
	matchLen := len(match)

	var err error
	for i, name := range expr.SubexpNames() {
		if i > matchLen {
			break
		}
		if i != 0 && name != "" {
			result[name], err = strconv.Atoi(match[i])
			if err != nil {
				return nil
			}
		}
	}

	return result
}
