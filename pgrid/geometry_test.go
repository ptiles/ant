package pgrid

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

func approxEq(a, b float64) bool {
	//return a == b
	return math.Abs(a-b) < 8.882e-16
}

func TestPrepareRotation(t *testing.T) {
	if GridLinesTotal != 5 {
		t.Skip("This test is only valid when GridLinesTotal == 5")
	}

	testCase := [5][5]bool{
		{true, true, true, false, false},
		{false, true, true, true, false},
		{false, false, true, true, true},
		{true, false, false, true, true},
		{true, true, false, false, true},
	}

	g := Geometry{}

	g.prepareRotation()

	for ax0 := range GridLinesTotal {
		for ax1 := range GridLinesTotal {
			if g[ax0][ax1].rotation != testCase[ax0][ax1] {
				t.Errorf("g[%d][%d].rotation expected %5v got %5v",
					ax0, ax1, testCase[ax0][ax1], g[ax0][ax1].rotation)
			}
		}
	}
}

func TestNewOffsetDeltas(t *testing.T) {
	if GridLinesTotal != 5 {
		t.Skip("This test is only valid when GridLinesTotal == 5")
	}

	radius := 1.0

	for _, tc := range newOffsetDeltasTestCases {
		result := newOffsetDeltas(radius, tc.ax0, tc.ax1, tc.axT)

		if !approxEq(result.zeroZero, tc.expect.zeroZero) {
			t.Errorf("newOffsetDeltas(%d, %d, %d) failed", tc.ax0, tc.ax1, tc.axT)
			t.Errorf("zeroZero expected: %6.3f; got: %6.3f; diff: %v",
				tc.expect.zeroZero, result.zeroZero, tc.expect.zeroZero-result.zeroZero,
			)
		}

		if !approxEq(result.ax0Delta, tc.expect.ax0Delta) {
			t.Errorf("newOffsetDeltas(%d, %d, %d) failed", tc.ax0, tc.ax1, tc.axT)
			t.Errorf("ax0Delta expected: %6.3f; got: %6.3f; diff: %v",
				tc.expect.ax0Delta, result.ax0Delta, tc.expect.ax0Delta-result.ax0Delta,
			)
		}

		if !approxEq(result.ax1Delta, tc.expect.ax1Delta) {
			t.Errorf("newOffsetDeltas(%d, %d, %d) failed", tc.ax0, tc.ax1, tc.axT)
			t.Errorf("ax1Delta expected: %6.3f; got: %6.3f; diff: %v",
				tc.expect.ax1Delta, result.ax1Delta, tc.expect.ax1Delta-result.ax1Delta,
			)
		}
	}
}

func TestGenNewOffsetDeltas(t *testing.T) {
	t.SkipNow()

	var sb strings.Builder

	radius := 1.0

	sb.WriteString("\nvar newOffsetDeltasTestCases = []struct {\n\tax0    uint8\n\tax1    uint8\n\taxT    uint8\n\texpect offsetDeltas\n}{{\n")

	for ax0, ax1 := range AxesAll() {
		for _, axT := range otherAxes(ax0, ax1) {
			deltas := newOffsetDeltas(radius, ax0, ax1, axT)

			sb.WriteString(fmt.Sprintf("\tax0: %d, ax1: %d, axT: %d,\n", ax0, ax1, axT))
			sb.WriteString("\texpect: offsetDeltas{\n")
			sb.WriteString(fmt.Sprintf("\t\tzeroZero: %v,\n", deltas.zeroZero))
			sb.WriteString(fmt.Sprintf("\t\tax0Delta: %v,\n", deltas.ax0Delta))
			sb.WriteString(fmt.Sprintf("\t\tax1Delta: %v,\n", deltas.ax1Delta))
			sb.WriteString("\t},\n")
			sb.WriteString("}, {\n")
		}
	}

	sb.WriteString("}}\n")

	t.Log(sb.String())
}

func TestThreeAxesOffset(t *testing.T) {
	if GridLinesTotal != 5 {
		t.Skip("This test is only valid when GridLinesTotal == 5")
	}

	for _, tc := range threeAxesOffsetTestCases {
		result := threeAxesOffset(tc.ax0, tc.ax1, tc.axT)

		if !approxEq(result, tc.expect) {
			t.Errorf("threeAxesOffset(%d, %d, %d) expected: %6.3f; got: %6.3f; diff: %v",
				tc.ax0, tc.ax1, tc.axT, tc.expect, result, tc.expect-result,
			)
		}
	}

}

func TestGenThreeAxesOffset(t *testing.T) {
	t.SkipNow()

	var sb strings.Builder

	sb.WriteString("\nvar threeAxesOffsetTestCases = []struct {\n\tax0    uint8\n\tax1    uint8\n\taxT    uint8\n\texpect float64\n}{\n")

	for ax0, ax1 := range AxesAll() {
		for _, axT := range otherAxes(ax0, ax1) {
			sb.WriteString(fmt.Sprintf("\t{ax0: %d, ax1: %d, axT: %d, expect: %v},\n",
				ax0, ax1, axT, threeAxesOffset(ax0, ax1, axT)))
		}
	}

	sb.WriteString("}\n")

	t.Log(sb.String())
}
