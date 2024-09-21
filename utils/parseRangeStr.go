package utils

import (
	"regexp"
	"strconv"
)

func ParseRangeStr(rangeStr string) (minimum, maximum int, err error) {
	re := regexp.MustCompile(`((\d+)-)?(\d+)`)
	result := re.FindStringSubmatch(rangeStr)

	maximum, err = strconv.Atoi(result[3])
	if err != nil {
		return
	}

	if result[2] == "" {
		return 0, maximum, nil
	}

	minimum, err = strconv.Atoi(result[2])
	if err != nil {
		return
	}

	return minimum, maximum, nil
}
