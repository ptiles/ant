package utils

import (
	"strconv"
)

func ParseRangeStr(rangeStr string) (int, int, error) {
	expr := `((?P<min>\d+)-)?(?P<max>\d+)`
	result := NamedStringMatches(expr, rangeStr)

	maximum, err := strconv.Atoi(result["max"])
	if err != nil {
		return 0, 0, err
	}

	if result["min"] == "" {
		return 0, maximum, nil
	}

	minimum, err := strconv.Atoi(result["min"])
	if err != nil {
		return 0, 0, err
	}

	return minimum, maximum, nil
}

func ParseRangeDeltaStr(rangeStr string) (int, int, int, error) {
	expr := `((?P<min>\d+)-)?(?P<max>\d+)((?P<mode>[%W])((?P<row>\d+)-)?(?P<delta>\d+))?`
	result := NamedStringMatches(expr, rangeStr)

	maximum, err := strconv.Atoi(result["max"])
	if err != nil {
		return 0, 0, 0, err
	}

	delta, err := strconv.Atoi(result["delta"])
	if err != nil {
		return 0, 0, 0, err
	}

	if result["min"] == "" {
		return 0, maximum, delta, nil
	}

	minimum, err := strconv.Atoi(result["min"])
	if err != nil {
		return 0, 0, 0, err
	}

	return minimum, maximum, delta, nil
}
