package main

import (
	"fmt"
	"github.com/ptiles/ant/pgrid"
	"regexp"
	"strconv"
	"strings"
)

func initialState(field *pgrid.Field, initialPoint string) (pgrid.GridPoint, pgrid.GridPoint, pgrid.GridLine, pgrid.GridLine) {
	re := regexp.MustCompile(`([A-E])(-?\d+)([+-]?)([A-E])(-?\d+)`)
	result := re.FindStringSubmatch(initialPoint)

	currAx, currOff, dir, nextAx, nextOff := result[1], result[2], result[3], result[4], result[5]

	currAxis := strings.Index("ABCDE", currAx)
	currOffset, _ := strconv.Atoi(currOff)
	currLine := pgrid.GridLine{Axis: uint8(currAxis), Offset: int16(currOffset)}

	nextAxis := strings.Index("ABCDE", nextAx)
	nextOffset, _ := strconv.Atoi(nextOff)
	nextLine := pgrid.GridLine{Axis: uint8(nextAxis), Offset: int16(nextOffset)}

	currAxIncreasing := dir != "-"
	currPoint := field.MakeGridPoint(currLine, nextLine)

	prevPoint, prevLine := field.NearestNeighbor(currPoint, nextLine, currLine, !currAxIncreasing)

	axisNames := [5]string{"A", "B", "C", "D", "E"}
	fmt.Printf("Initial step: ")
	fmt.Printf("%s%d%s%d=>", axisNames[currLine.Axis], currLine.Offset, axisNames[prevLine.Axis], prevLine.Offset)
	fmt.Printf("%s%d%s%d\n", axisNames[nextLine.Axis], nextLine.Offset, axisNames[currLine.Axis], currLine.Offset)

	return prevPoint, currPoint, prevLine, currLine
}
