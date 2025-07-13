package pgrid

import (
	"fmt"
	"math"
)

func printOffsets(name string, offsets *allOffsetDeltas) {
	fmt.Printf("\n\n\n// ---------- %s ----------\n\n", name)
	for i0, o0 := range offsets {
		for i1, o1 := range o0 {
			if o1[0].targetAx == 0 && o1[1].targetAx == 0 && o1[2].targetAx == 0 {
				for i2 := range o1 {
					fmt.Printf(
						"// %s[%d][%d][%d] = offsetDeltas{}\n",
						name, i0, i1, i2,
					)
				}
			} else {
				for i2, o2 := range o1 {
					fmt.Printf(
						"%s[%d][%d][%d].ax0Delta = %6.3f // %20.17f\n",
						name, i0, i1, i2, o2.ax0Delta, o2.ax0Delta,
					)
					fmt.Printf(
						"%s[%d][%d][%d].ax1Delta = %6.3f // %20.17f\n",
						name, i0, i1, i2, o2.ax1Delta, o2.ax1Delta,
					)
				}
			}
		}
	}
}

// TODO: calculate this properly from angles instead of heuristic
// delta == -1.618 =>  -math.Phi
// delta == -0.618 => 1-math.Phi
// delta ==  0.618 =>   math.Phi-1
// delta ==  1.618 =>   math.Phi
func updateOffsetsToFirst(offsetsToFirst *allOffsetDeltas) {
	// offsetsToFirst[0][0][0] = offsetDeltas{}
	// offsetsToFirst[0][0][1] = offsetDeltas{}
	// offsetsToFirst[0][0][2] = offsetDeltas{}
	offsetsToFirst[0][1][0].ax0Delta = math.Phi - 1 //  0.61803398874989501
	offsetsToFirst[0][1][0].ax1Delta = -1           // -1.00000000000000022
	offsetsToFirst[0][1][1].ax0Delta = -1           // -0.99999999999999922
	offsetsToFirst[0][1][1].ax1Delta = -math.Phi    // -1.61803398874989424
	offsetsToFirst[0][1][2].ax0Delta = math.Phi     //  1.61803398874989868
	offsetsToFirst[0][1][2].ax1Delta = math.Phi     //  1.61803398874989823
	offsetsToFirst[0][2][0].ax0Delta = -1           // -1.00000000000000022
	offsetsToFirst[0][2][0].ax1Delta = math.Phi - 1 //  0.61803398874989501
	offsetsToFirst[0][2][1].ax0Delta = 1 - math.Phi // -0.61803398874989468
	offsetsToFirst[0][2][1].ax1Delta = 1 - math.Phi // -0.61803398874989501
	offsetsToFirst[0][2][2].ax0Delta = -math.Phi    // -1.61803398874989357
	offsetsToFirst[0][2][2].ax1Delta = -1           // -0.99999999999999856
	offsetsToFirst[0][3][0].ax0Delta = -math.Phi    // -1.61803398874989401
	offsetsToFirst[0][3][0].ax1Delta = -1           // -0.99999999999999922
	offsetsToFirst[0][3][1].ax0Delta = 1 - math.Phi // -0.61803398874989501
	offsetsToFirst[0][3][1].ax1Delta = 1 - math.Phi // -0.61803398874989468
	offsetsToFirst[0][3][2].ax0Delta = -1           // -1.00000000000000067
	offsetsToFirst[0][3][2].ax1Delta = math.Phi - 1 //  0.61803398874989501
	offsetsToFirst[0][4][0].ax0Delta = math.Phi     //  1.61803398874989823
	offsetsToFirst[0][4][0].ax1Delta = math.Phi     //  1.61803398874989868
	offsetsToFirst[0][4][1].ax0Delta = -1           // -0.99999999999999856
	offsetsToFirst[0][4][1].ax1Delta = -math.Phi    // -1.61803398874989357
	offsetsToFirst[0][4][2].ax0Delta = math.Phi - 1 //  0.61803398874989501
	offsetsToFirst[0][4][2].ax1Delta = -1           // -1.00000000000000067
	offsetsToFirst[1][0][0].ax0Delta = math.Phi     //  1.61803398874989424
	offsetsToFirst[1][0][0].ax1Delta = math.Phi     //  1.61803398874989490
	offsetsToFirst[1][0][1].ax0Delta = -1           // -1.00000000000000089
	offsetsToFirst[1][0][1].ax1Delta = -math.Phi    // -1.61803398874989579
	offsetsToFirst[1][0][2].ax0Delta = math.Phi - 1 //  0.61803398874989324
	offsetsToFirst[1][0][2].ax1Delta = -1           // -0.99999999999999978
	// offsetsToFirst[1][1][0] = offsetDeltas{}
	// offsetsToFirst[1][1][1] = offsetDeltas{}
	// offsetsToFirst[1][1][2] = offsetDeltas{}
	offsetsToFirst[1][2][0].ax0Delta = math.Phi     //  1.61803398874989490
	offsetsToFirst[1][2][0].ax1Delta = math.Phi     //  1.61803398874989424
	offsetsToFirst[1][2][1].ax0Delta = math.Phi - 1 //  0.61803398874989535
	offsetsToFirst[1][2][1].ax1Delta = -1           // -1.00000000000000000
	offsetsToFirst[1][2][2].ax0Delta = -1           // -0.99999999999999678
	offsetsToFirst[1][2][2].ax1Delta = -math.Phi    // -1.61803398874989224
	offsetsToFirst[1][3][0].ax0Delta = -math.Phi    // -1.61803398874989579
	offsetsToFirst[1][3][0].ax1Delta = -1           // -1.00000000000000067
	offsetsToFirst[1][3][1].ax0Delta = -1           // -1.00000000000000000
	offsetsToFirst[1][3][1].ax1Delta = math.Phi - 1 //  0.61803398874989535
	offsetsToFirst[1][3][2].ax0Delta = 1 - math.Phi // -0.61803398874989401
	offsetsToFirst[1][3][2].ax1Delta = 1 - math.Phi // -0.61803398874989546
	offsetsToFirst[1][4][0].ax0Delta = -1           // -0.99999999999999978
	offsetsToFirst[1][4][0].ax1Delta = math.Phi - 1 //  0.61803398874989335
	offsetsToFirst[1][4][1].ax0Delta = -math.Phi    // -1.61803398874989224
	offsetsToFirst[1][4][1].ax1Delta = -1           // -0.99999999999999678
	offsetsToFirst[1][4][2].ax0Delta = 1 - math.Phi // -0.61803398874989546
	offsetsToFirst[1][4][2].ax1Delta = 1 - math.Phi // -0.61803398874989401
	offsetsToFirst[2][0][0].ax0Delta = -1           // -0.99999999999999956
	offsetsToFirst[2][0][0].ax1Delta = math.Phi - 1 //  0.61803398874989479
	offsetsToFirst[2][0][1].ax0Delta = -math.Phi    // -1.61803398874989490
	offsetsToFirst[2][0][1].ax1Delta = -1           // -1.00000000000000044
	offsetsToFirst[2][0][2].ax0Delta = 1 - math.Phi // -0.61803398874989546
	offsetsToFirst[2][0][2].ax1Delta = 1 - math.Phi // -0.61803398874989468
	offsetsToFirst[2][1][0].ax0Delta = math.Phi - 1 //  0.61803398874989479
	offsetsToFirst[2][1][0].ax1Delta = -1           // -0.99999999999999956
	offsetsToFirst[2][1][1].ax0Delta = math.Phi     //  1.61803398874989379
	offsetsToFirst[2][1][1].ax1Delta = math.Phi     //  1.61803398874989379
	offsetsToFirst[2][1][2].ax0Delta = -1           // -1.00000000000000355
	offsetsToFirst[2][1][2].ax1Delta = -math.Phi    // -1.61803398874989801
	// offsetsToFirst[2][2][0] = offsetDeltas{}
	// offsetsToFirst[2][2][1] = offsetDeltas{}
	// offsetsToFirst[2][2][2] = offsetDeltas{}
	offsetsToFirst[2][3][0].ax0Delta = -1           // -1.00000000000000044
	offsetsToFirst[2][3][0].ax1Delta = -math.Phi    // -1.61803398874989446
	offsetsToFirst[2][3][1].ax0Delta = math.Phi     //  1.61803398874989357
	offsetsToFirst[2][3][1].ax1Delta = math.Phi     //  1.61803398874989379
	offsetsToFirst[2][3][2].ax0Delta = math.Phi - 1 //  0.61803398874989601
	offsetsToFirst[2][3][2].ax1Delta = -1           // -1.00000000000000022
	offsetsToFirst[2][4][0].ax0Delta = 1 - math.Phi // -0.61803398874989468
	offsetsToFirst[2][4][0].ax1Delta = 1 - math.Phi // -0.61803398874989546
	offsetsToFirst[2][4][1].ax0Delta = -math.Phi    // -1.61803398874989801
	offsetsToFirst[2][4][1].ax1Delta = -1           // -1.00000000000000355
	offsetsToFirst[2][4][2].ax0Delta = -1           // -1.00000000000000022
	offsetsToFirst[2][4][2].ax1Delta = math.Phi - 1 //  0.61803398874989601
	offsetsToFirst[3][0][0].ax0Delta = 1 - math.Phi // -0.61803398874989512
	offsetsToFirst[3][0][0].ax1Delta = 1 - math.Phi // -0.61803398874989479
	offsetsToFirst[3][0][1].ax0Delta = -math.Phi    // -1.61803398874989424
	offsetsToFirst[3][0][1].ax1Delta = -1           // -0.99999999999999967
	offsetsToFirst[3][0][2].ax0Delta = -1           // -0.99999999999999911
	offsetsToFirst[3][0][2].ax1Delta = math.Phi - 1 //  0.61803398874989446
	offsetsToFirst[3][1][0].ax0Delta = 1 - math.Phi // -0.61803398874989468
	offsetsToFirst[3][1][0].ax1Delta = 1 - math.Phi // -0.61803398874989512
	offsetsToFirst[3][1][1].ax0Delta = -1           // -1.00000000000000000
	offsetsToFirst[3][1][1].ax1Delta = math.Phi - 1 //  0.61803398874989546
	offsetsToFirst[3][1][2].ax0Delta = -math.Phi    // -1.61803398874989757
	offsetsToFirst[3][1][2].ax1Delta = -1           // -1.00000000000000266
	offsetsToFirst[3][2][0].ax0Delta = -1           // -0.99999999999999967
	offsetsToFirst[3][2][0].ax1Delta = -math.Phi    // -1.61803398874989424
	offsetsToFirst[3][2][1].ax0Delta = math.Phi - 1 //  0.61803398874989546
	offsetsToFirst[3][2][1].ax1Delta = -1           // -1.00000000000000000
	offsetsToFirst[3][2][2].ax0Delta = math.Phi     //  1.61803398874989246
	offsetsToFirst[3][2][2].ax1Delta = math.Phi     //  1.61803398874989246
	// offsetsToFirst[3][3][0] = offsetDeltas{}
	// offsetsToFirst[3][3][1] = offsetDeltas{}
	// offsetsToFirst[3][3][2] = offsetDeltas{}
	offsetsToFirst[3][4][0].ax0Delta = math.Phi - 1 //  0.61803398874989446
	offsetsToFirst[3][4][0].ax1Delta = -1           // -0.99999999999999933
	offsetsToFirst[3][4][1].ax0Delta = -1           // -1.00000000000000266
	offsetsToFirst[3][4][1].ax1Delta = -math.Phi    // -1.61803398874989757
	offsetsToFirst[3][4][2].ax0Delta = math.Phi     //  1.61803398874989246
	offsetsToFirst[3][4][2].ax1Delta = math.Phi     //  1.61803398874989246
	offsetsToFirst[4][0][0].ax0Delta = math.Phi - 1 //  0.61803398874989346
	offsetsToFirst[4][0][0].ax1Delta = -1           // -1.00000000000000044
	offsetsToFirst[4][0][1].ax0Delta = -1           // -1.00000000000000111
	offsetsToFirst[4][0][1].ax1Delta = -math.Phi    // -1.61803398874989535
	offsetsToFirst[4][0][2].ax0Delta = math.Phi     //  1.61803398874989468
	offsetsToFirst[4][0][2].ax1Delta = math.Phi     //  1.61803398874989623
	offsetsToFirst[4][1][0].ax0Delta = -1           // -1.00000000000000022
	offsetsToFirst[4][1][0].ax1Delta = math.Phi - 1 //  0.61803398874989346
	offsetsToFirst[4][1][1].ax0Delta = 1 - math.Phi // -0.61803398874989579
	offsetsToFirst[4][1][1].ax1Delta = 1 - math.Phi // -0.61803398874989368
	offsetsToFirst[4][1][2].ax0Delta = -math.Phi    // -1.61803398874989357
	offsetsToFirst[4][1][2].ax1Delta = -1           // -0.99999999999999745
	offsetsToFirst[4][2][0].ax0Delta = -math.Phi    // -1.61803398874989535
	offsetsToFirst[4][2][0].ax1Delta = -1           // -1.00000000000000111
	offsetsToFirst[4][2][1].ax0Delta = 1 - math.Phi // -0.61803398874989368
	offsetsToFirst[4][2][1].ax1Delta = 1 - math.Phi // -0.61803398874989579
	offsetsToFirst[4][2][2].ax0Delta = -1           // -0.99999999999999956
	offsetsToFirst[4][2][2].ax1Delta = math.Phi - 1 //  0.61803398874989579
	offsetsToFirst[4][3][0].ax0Delta = math.Phi     //  1.61803398874989623
	offsetsToFirst[4][3][0].ax1Delta = math.Phi     //  1.61803398874989424
	offsetsToFirst[4][3][1].ax0Delta = -1           // -0.99999999999999722
	offsetsToFirst[4][3][1].ax1Delta = -math.Phi    // -1.61803398874989357
	offsetsToFirst[4][3][2].ax0Delta = math.Phi - 1 //  0.61803398874989579
	offsetsToFirst[4][3][2].ax1Delta = -1           // -0.99999999999999956
	// offsetsToFirst[4][4][0] = offsetDeltas{}
	// offsetsToFirst[4][4][1] = offsetDeltas{}
	// offsetsToFirst[4][4][2] = offsetDeltas{}
}

// TODO: calculate this properly from angles instead of heuristic
// delta == -1.618 =>  -math.Phi
// delta == -0.618 => 1-math.Phi
// delta ==  0.618 =>   math.Phi-1
// delta ==  1.618 =>   math.Phi
func updateOffsetsToLast(offsetsToLast *allOffsetDeltas) {
	// offsetsToLast[0][0][0] = offsetDeltas{}
	// offsetsToLast[0][0][1] = offsetDeltas{}
	// offsetsToLast[0][0][2] = offsetDeltas{}
	offsetsToLast[0][1][0].ax0Delta = -1           // -0.99999999999999956
	offsetsToLast[0][1][0].ax1Delta = math.Phi - 1 //  0.61803398874989479
	offsetsToLast[0][1][1].ax0Delta = 1 - math.Phi // -0.61803398874989512
	offsetsToLast[0][1][1].ax1Delta = 1 - math.Phi // -0.61803398874989479
	offsetsToLast[0][1][2].ax0Delta = math.Phi - 1 //  0.61803398874989346
	offsetsToLast[0][1][2].ax1Delta = -1           // -1.00000000000000044
	offsetsToLast[0][2][0].ax0Delta = math.Phi     //  1.61803398874989424
	offsetsToLast[0][2][0].ax1Delta = math.Phi     //  1.61803398874989490
	offsetsToLast[0][2][1].ax0Delta = -math.Phi    // -1.61803398874989424
	offsetsToLast[0][2][1].ax1Delta = -1           // -0.99999999999999967
	offsetsToLast[0][2][2].ax0Delta = -1           // -1.00000000000000111
	offsetsToLast[0][2][2].ax1Delta = -math.Phi    // -1.61803398874989535
	offsetsToLast[0][3][0].ax0Delta = -1           // -1.00000000000000089
	offsetsToLast[0][3][0].ax1Delta = -math.Phi    // -1.61803398874989579
	offsetsToLast[0][3][1].ax0Delta = -math.Phi    // -1.61803398874989490
	offsetsToLast[0][3][1].ax1Delta = -1           // -1.00000000000000044
	offsetsToLast[0][3][2].ax0Delta = math.Phi     //  1.61803398874989468
	offsetsToLast[0][3][2].ax1Delta = math.Phi     //  1.61803398874989623
	offsetsToLast[0][4][0].ax0Delta = math.Phi - 1 //  0.61803398874989324
	offsetsToLast[0][4][0].ax1Delta = -1           // -0.99999999999999978
	offsetsToLast[0][4][1].ax0Delta = 1 - math.Phi // -0.61803398874989546
	offsetsToLast[0][4][1].ax1Delta = 1 - math.Phi // -0.61803398874989468
	offsetsToLast[0][4][2].ax0Delta = -1           // -0.99999999999999911
	offsetsToLast[0][4][2].ax1Delta = math.Phi - 1 //  0.61803398874989446
	offsetsToLast[1][0][0].ax0Delta = math.Phi - 1 //  0.61803398874989479
	offsetsToLast[1][0][0].ax1Delta = -1           // -0.99999999999999956
	offsetsToLast[1][0][1].ax0Delta = 1 - math.Phi // -0.61803398874989468
	offsetsToLast[1][0][1].ax1Delta = 1 - math.Phi // -0.61803398874989512
	offsetsToLast[1][0][2].ax0Delta = -1           // -1.00000000000000022
	offsetsToLast[1][0][2].ax1Delta = math.Phi - 1 //  0.61803398874989346
	// offsetsToLast[1][1][0] = offsetDeltas{}
	// offsetsToLast[1][1][1] = offsetDeltas{}
	// offsetsToLast[1][1][2] = offsetDeltas{}
	offsetsToLast[1][2][0].ax0Delta = math.Phi - 1 //  0.61803398874989501
	offsetsToLast[1][2][0].ax1Delta = -1           // -1.00000000000000022
	offsetsToLast[1][2][1].ax0Delta = -1           // -1.00000000000000000
	offsetsToLast[1][2][1].ax1Delta = math.Phi - 1 //  0.61803398874989546
	offsetsToLast[1][2][2].ax0Delta = 1 - math.Phi // -0.61803398874989579
	offsetsToLast[1][2][2].ax1Delta = 1 - math.Phi // -0.61803398874989368
	offsetsToLast[1][3][0].ax0Delta = -1           // -0.99999999999999922
	offsetsToLast[1][3][0].ax1Delta = -math.Phi    // -1.61803398874989424
	offsetsToLast[1][3][1].ax0Delta = math.Phi     //  1.61803398874989379
	offsetsToLast[1][3][1].ax1Delta = math.Phi     //  1.61803398874989379
	offsetsToLast[1][3][2].ax0Delta = -math.Phi    // -1.61803398874989357
	offsetsToLast[1][3][2].ax1Delta = -1           // -0.99999999999999745
	offsetsToLast[1][4][0].ax0Delta = math.Phi     //  1.61803398874989868
	offsetsToLast[1][4][0].ax1Delta = math.Phi     //  1.61803398874989823
	offsetsToLast[1][4][1].ax0Delta = -1           // -1.00000000000000355
	offsetsToLast[1][4][1].ax1Delta = -math.Phi    // -1.61803398874989801
	offsetsToLast[1][4][2].ax0Delta = -math.Phi    // -1.61803398874989757
	offsetsToLast[1][4][2].ax1Delta = -1           // -1.00000000000000266
	offsetsToLast[2][0][0].ax0Delta = math.Phi     //  1.61803398874989490
	offsetsToLast[2][0][0].ax1Delta = math.Phi     //  1.61803398874989424
	offsetsToLast[2][0][1].ax0Delta = -1           // -0.99999999999999967
	offsetsToLast[2][0][1].ax1Delta = -math.Phi    // -1.61803398874989424
	offsetsToLast[2][0][2].ax0Delta = -math.Phi    // -1.61803398874989535
	offsetsToLast[2][0][2].ax1Delta = -1           // -1.00000000000000111
	offsetsToLast[2][1][0].ax0Delta = -1           // -1.00000000000000022
	offsetsToLast[2][1][0].ax1Delta = math.Phi - 1 //  0.61803398874989501
	offsetsToLast[2][1][1].ax0Delta = math.Phi - 1 //  0.61803398874989546
	offsetsToLast[2][1][1].ax1Delta = -1           // -1.00000000000000000
	offsetsToLast[2][1][2].ax0Delta = 1 - math.Phi // -0.61803398874989368
	offsetsToLast[2][1][2].ax1Delta = 1 - math.Phi // -0.61803398874989579
	// offsetsToLast[2][2][0] = offsetDeltas{}
	// offsetsToLast[2][2][1] = offsetDeltas{}
	// offsetsToLast[2][2][2] = offsetDeltas{}
	offsetsToLast[2][3][0].ax0Delta = 1 - math.Phi // -0.61803398874989468
	offsetsToLast[2][3][0].ax1Delta = 1 - math.Phi // -0.61803398874989501
	offsetsToLast[2][3][1].ax0Delta = math.Phi - 1 //  0.61803398874989535
	offsetsToLast[2][3][1].ax1Delta = -1           // -1.00000000000000000
	offsetsToLast[2][3][2].ax0Delta = -1           // -0.99999999999999956
	offsetsToLast[2][3][2].ax1Delta = math.Phi - 1 //  0.61803398874989579
	offsetsToLast[2][4][0].ax0Delta = -math.Phi    // -1.61803398874989357
	offsetsToLast[2][4][0].ax1Delta = -1           // -0.99999999999999856
	offsetsToLast[2][4][1].ax0Delta = -1           // -0.99999999999999678
	offsetsToLast[2][4][1].ax1Delta = -math.Phi    // -1.61803398874989224
	offsetsToLast[2][4][2].ax0Delta = math.Phi     //  1.61803398874989246
	offsetsToLast[2][4][2].ax1Delta = math.Phi     //  1.61803398874989246
	offsetsToLast[3][0][0].ax0Delta = -math.Phi    // -1.61803398874989579
	offsetsToLast[3][0][0].ax1Delta = -1           // -1.00000000000000067
	offsetsToLast[3][0][1].ax0Delta = -1           // -1.00000000000000044
	offsetsToLast[3][0][1].ax1Delta = -math.Phi    // -1.61803398874989446
	offsetsToLast[3][0][2].ax0Delta = math.Phi     //  1.61803398874989623
	offsetsToLast[3][0][2].ax1Delta = math.Phi     //  1.61803398874989424
	offsetsToLast[3][1][0].ax0Delta = -math.Phi    // -1.61803398874989401
	offsetsToLast[3][1][0].ax1Delta = -1           // -0.99999999999999922
	offsetsToLast[3][1][1].ax0Delta = math.Phi     //  1.61803398874989357
	offsetsToLast[3][1][1].ax1Delta = math.Phi     //  1.61803398874989379
	offsetsToLast[3][1][2].ax0Delta = -1           // -0.99999999999999722
	offsetsToLast[3][1][2].ax1Delta = -math.Phi    // -1.61803398874989357
	offsetsToLast[3][2][0].ax0Delta = 1 - math.Phi // -0.61803398874989501
	offsetsToLast[3][2][0].ax1Delta = 1 - math.Phi // -0.61803398874989468
	offsetsToLast[3][2][1].ax0Delta = -1           // -1.00000000000000000
	offsetsToLast[3][2][1].ax1Delta = math.Phi - 1 //  0.61803398874989535
	offsetsToLast[3][2][2].ax0Delta = math.Phi - 1 //  0.61803398874989579
	offsetsToLast[3][2][2].ax1Delta = -1           // -0.99999999999999956
	// offsetsToLast[3][3][0] = offsetDeltas{}
	// offsetsToLast[3][3][1] = offsetDeltas{}
	// offsetsToLast[3][3][2] = offsetDeltas{}
	offsetsToLast[3][4][0].ax0Delta = -1           // -1.00000000000000067
	offsetsToLast[3][4][0].ax1Delta = math.Phi - 1 //  0.61803398874989501
	offsetsToLast[3][4][1].ax0Delta = 1 - math.Phi // -0.61803398874989401
	offsetsToLast[3][4][1].ax1Delta = 1 - math.Phi // -0.61803398874989546
	offsetsToLast[3][4][2].ax0Delta = math.Phi - 1 //  0.61803398874989601
	offsetsToLast[3][4][2].ax1Delta = -1           // -1.00000000000000022
	offsetsToLast[4][0][0].ax0Delta = -1           // -0.99999999999999978
	offsetsToLast[4][0][0].ax1Delta = math.Phi - 1 //  0.61803398874989335
	offsetsToLast[4][0][1].ax0Delta = 1 - math.Phi // -0.61803398874989468
	offsetsToLast[4][0][1].ax1Delta = 1 - math.Phi // -0.61803398874989546
	offsetsToLast[4][0][2].ax0Delta = math.Phi - 1 //  0.61803398874989446
	offsetsToLast[4][0][2].ax1Delta = -1           // -0.99999999999999933
	offsetsToLast[4][1][0].ax0Delta = math.Phi     //  1.61803398874989823
	offsetsToLast[4][1][0].ax1Delta = math.Phi     //  1.61803398874989868
	offsetsToLast[4][1][1].ax0Delta = -math.Phi    // -1.61803398874989801
	offsetsToLast[4][1][1].ax1Delta = -1           // -1.00000000000000355
	offsetsToLast[4][1][2].ax0Delta = -1           // -1.00000000000000266
	offsetsToLast[4][1][2].ax1Delta = -math.Phi    // -1.61803398874989757
	offsetsToLast[4][2][0].ax0Delta = -1           // -0.99999999999999856
	offsetsToLast[4][2][0].ax1Delta = -math.Phi    // -1.61803398874989357
	offsetsToLast[4][2][1].ax0Delta = -math.Phi    // -1.61803398874989224
	offsetsToLast[4][2][1].ax1Delta = -1           // -0.99999999999999678
	offsetsToLast[4][2][2].ax0Delta = math.Phi     //  1.61803398874989246
	offsetsToLast[4][2][2].ax1Delta = math.Phi     //  1.61803398874989246
	offsetsToLast[4][3][0].ax0Delta = math.Phi - 1 //  0.61803398874989501
	offsetsToLast[4][3][0].ax1Delta = -1           // -1.00000000000000067
	offsetsToLast[4][3][1].ax0Delta = 1 - math.Phi // -0.61803398874989546
	offsetsToLast[4][3][1].ax1Delta = 1 - math.Phi // -0.61803398874989401
	offsetsToLast[4][3][2].ax0Delta = -1           // -1.00000000000000022
	offsetsToLast[4][3][2].ax1Delta = math.Phi - 1 //  0.61803398874989601
	// offsetsToLast[4][4][0] = offsetDeltas{}
	// offsetsToLast[4][4][1] = offsetDeltas{}
	// offsetsToLast[4][4][2] = offsetDeltas{}
}
