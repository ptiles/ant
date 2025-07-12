package utils

import (
	"strconv"
)

func ParseRangeStr(rangeStr string) (minimum, maximum, delta int, err error) {
	expr := `((?P<min>\d+)-)?(?P<max>\d+)((?P<mode>[%W])((?P<row>\d+)-)?(?P<delta>\d+))?`
	result := NamedStringMatches(expr, rangeStr)
	//fmt.Println(matches)

	//re := regexp.MustCompile(`((\d+)-)?(\d+)(%(\d+))?`)
	//result := re.FindStringSubmatch(rangeStr)

	maximum, err = strconv.Atoi(result["max"])
	if err != nil {
		return
	}

	delta, err = strconv.Atoi(result["delta"])
	if err != nil {
		return
	}

	if result["min"] == "" {
		return 0, maximum, delta, nil
	}

	minimum, err = strconv.Atoi(result["min"])
	if err != nil {
		return
	}

	return minimum, maximum, delta, nil
}
