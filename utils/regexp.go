package utils

import "regexp"

func NamedMatches(r, s string) map[string]string {
	expr := regexp.MustCompile(r)
	match := expr.FindStringSubmatch(s)
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
