package pgrid

import "fmt"

const GridLinesTotal = uint8(5)

func init() {
	if GridLinesTotal%2 == 0 || GridLinesTotal < 5 || GridLinesTotal > 25 {
		fmt.Println("GridLinesTotal should be odd number between 5 and 25")
	}
}
