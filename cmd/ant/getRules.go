package main

import "errors"

func getRules(antName string) ([]bool, uint8, error) {
	limit := uint8(len(antName))
	var nameInvalid = limit < 2
	rules := make([]bool, limit)
	for i, letter := range antName {
		if letter != 'R' && letter != 'r' && letter != 'L' && letter != 'l' {
			nameInvalid = true
			break
		}
		rules[i] = letter == 'R' || letter == 'r'
	}
	if nameInvalid {
		return rules, limit, errors.New("invalid name")
	}
	return rules, limit, nil
}
