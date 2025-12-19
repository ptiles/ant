package pgrid

var newOffsetDeltasTestCases = []struct {
	ax0    uint8
	ax1    uint8
	axT    uint8
	expect offsetDeltas
}{{
	ax0: 0, ax1: 1, axT: 2,
	expect: offsetDeltas{
		zeroZero: -1.381966011250105,
		ax0Delta: -1,
		ax1Delta: 0.6180339887498949,
	},
}, {
	ax0: 0, ax1: 1, axT: 3,
	expect: offsetDeltas{
		zeroZero: -2.23606797749979,
		ax0Delta: -0.6180339887498949,
		ax1Delta: -0.6180339887498949,
	},
}, {
	ax0: 0, ax1: 1, axT: 4,
	expect: offsetDeltas{
		zeroZero: -1.381966011250105,
		ax0Delta: 0.6180339887498949,
		ax1Delta: -1,
	},
}, {
	ax0: 0, ax1: 2, axT: 1,
	expect: offsetDeltas{
		zeroZero: 2.23606797749979,
		ax0Delta: 1.618033988749895,
		ax1Delta: 1.618033988749895,
	},
}, {
	ax0: 0, ax1: 2, axT: 3,
	expect: offsetDeltas{
		zeroZero: -3.618033988749895,
		ax0Delta: -1.618033988749895,
		ax1Delta: -1,
	},
}, {
	ax0: 0, ax1: 2, axT: 4,
	expect: offsetDeltas{
		zeroZero: -3.618033988749895,
		ax0Delta: -1,
		ax1Delta: -1.618033988749895,
	},
}, {
	ax0: 0, ax1: 3, axT: 1,
	expect: offsetDeltas{
		zeroZero: -3.618033988749895,
		ax0Delta: -1,
		ax1Delta: -1.618033988749895,
	},
}, {
	ax0: 0, ax1: 3, axT: 2,
	expect: offsetDeltas{
		zeroZero: -3.618033988749895,
		ax0Delta: -1.618033988749895,
		ax1Delta: -1,
	},
}, {
	ax0: 0, ax1: 3, axT: 4,
	expect: offsetDeltas{
		zeroZero: 2.23606797749979,
		ax0Delta: 1.618033988749895,
		ax1Delta: 1.618033988749895,
	},
}, {
	ax0: 0, ax1: 4, axT: 1,
	expect: offsetDeltas{
		zeroZero: -1.381966011250105,
		ax0Delta: 0.6180339887498949,
		ax1Delta: -1,
	},
}, {
	ax0: 0, ax1: 4, axT: 2,
	expect: offsetDeltas{
		zeroZero: -2.23606797749979,
		ax0Delta: -0.6180339887498949,
		ax1Delta: -0.6180339887498949,
	},
}, {
	ax0: 0, ax1: 4, axT: 3,
	expect: offsetDeltas{
		zeroZero: -1.381966011250105,
		ax0Delta: -1,
		ax1Delta: 0.6180339887498949,
	},
}, {
	ax0: 1, ax1: 0, axT: 2,
	expect: offsetDeltas{
		zeroZero: -1.381966011250105,
		ax0Delta: 0.6180339887498949,
		ax1Delta: -1,
	},
}, {
	ax0: 1, ax1: 0, axT: 3,
	expect: offsetDeltas{
		zeroZero: -2.23606797749979,
		ax0Delta: -0.6180339887498949,
		ax1Delta: -0.6180339887498949,
	},
}, {
	ax0: 1, ax1: 0, axT: 4,
	expect: offsetDeltas{
		zeroZero: -1.381966011250105,
		ax0Delta: -1,
		ax1Delta: 0.6180339887498949,
	},
}, {
	ax0: 1, ax1: 2, axT: 0,
	expect: offsetDeltas{
		zeroZero: -1.381966011250105,
		ax0Delta: 0.6180339887498949,
		ax1Delta: -1,
	},
}, {
	ax0: 1, ax1: 2, axT: 3,
	expect: offsetDeltas{
		zeroZero: -1.381966011250105,
		ax0Delta: -1,
		ax1Delta: 0.6180339887498949,
	},
}, {
	ax0: 1, ax1: 2, axT: 4,
	expect: offsetDeltas{
		zeroZero: -2.23606797749979,
		ax0Delta: -0.6180339887498949,
		ax1Delta: -0.6180339887498949,
	},
}, {
	ax0: 1, ax1: 3, axT: 0,
	expect: offsetDeltas{
		zeroZero: -3.618033988749895,
		ax0Delta: -1,
		ax1Delta: -1.618033988749895,
	},
}, {
	ax0: 1, ax1: 3, axT: 2,
	expect: offsetDeltas{
		zeroZero: 2.23606797749979,
		ax0Delta: 1.618033988749895,
		ax1Delta: 1.618033988749895,
	},
}, {
	ax0: 1, ax1: 3, axT: 4,
	expect: offsetDeltas{
		zeroZero: -3.618033988749895,
		ax0Delta: -1.618033988749895,
		ax1Delta: -1,
	},
}, {
	ax0: 1, ax1: 4, axT: 0,
	expect: offsetDeltas{
		zeroZero: 2.23606797749979,
		ax0Delta: 1.618033988749895,
		ax1Delta: 1.618033988749895,
	},
}, {
	ax0: 1, ax1: 4, axT: 2,
	expect: offsetDeltas{
		zeroZero: -3.618033988749895,
		ax0Delta: -1,
		ax1Delta: -1.618033988749895,
	},
}, {
	ax0: 1, ax1: 4, axT: 3,
	expect: offsetDeltas{
		zeroZero: -3.618033988749895,
		ax0Delta: -1.618033988749895,
		ax1Delta: -1,
	},
}, {
	ax0: 2, ax1: 0, axT: 1,
	expect: offsetDeltas{
		zeroZero: 2.23606797749979,
		ax0Delta: 1.618033988749895,
		ax1Delta: 1.618033988749895,
	},
}, {
	ax0: 2, ax1: 0, axT: 3,
	expect: offsetDeltas{
		zeroZero: -3.618033988749895,
		ax0Delta: -1,
		ax1Delta: -1.618033988749895,
	},
}, {
	ax0: 2, ax1: 0, axT: 4,
	expect: offsetDeltas{
		zeroZero: -3.618033988749895,
		ax0Delta: -1.618033988749895,
		ax1Delta: -1,
	},
}, {
	ax0: 2, ax1: 1, axT: 0,
	expect: offsetDeltas{
		zeroZero: -1.381966011250105,
		ax0Delta: -1,
		ax1Delta: 0.6180339887498949,
	},
}, {
	ax0: 2, ax1: 1, axT: 3,
	expect: offsetDeltas{
		zeroZero: -1.381966011250105,
		ax0Delta: 0.6180339887498949,
		ax1Delta: -1,
	},
}, {
	ax0: 2, ax1: 1, axT: 4,
	expect: offsetDeltas{
		zeroZero: -2.23606797749979,
		ax0Delta: -0.6180339887498949,
		ax1Delta: -0.6180339887498949,
	},
}, {
	ax0: 2, ax1: 3, axT: 0,
	expect: offsetDeltas{
		zeroZero: -2.23606797749979,
		ax0Delta: -0.6180339887498949,
		ax1Delta: -0.6180339887498949,
	},
}, {
	ax0: 2, ax1: 3, axT: 1,
	expect: offsetDeltas{
		zeroZero: -1.381966011250105,
		ax0Delta: 0.6180339887498949,
		ax1Delta: -1,
	},
}, {
	ax0: 2, ax1: 3, axT: 4,
	expect: offsetDeltas{
		zeroZero: -1.381966011250105,
		ax0Delta: -1,
		ax1Delta: 0.6180339887498949,
	},
}, {
	ax0: 2, ax1: 4, axT: 0,
	expect: offsetDeltas{
		zeroZero: -3.618033988749895,
		ax0Delta: -1.618033988749895,
		ax1Delta: -1,
	},
}, {
	ax0: 2, ax1: 4, axT: 1,
	expect: offsetDeltas{
		zeroZero: -3.618033988749895,
		ax0Delta: -1,
		ax1Delta: -1.618033988749895,
	},
}, {
	ax0: 2, ax1: 4, axT: 3,
	expect: offsetDeltas{
		zeroZero: 2.23606797749979,
		ax0Delta: 1.618033988749895,
		ax1Delta: 1.618033988749895,
	},
}, {
	ax0: 3, ax1: 0, axT: 1,
	expect: offsetDeltas{
		zeroZero: -3.618033988749895,
		ax0Delta: -1.618033988749895,
		ax1Delta: -1,
	},
}, {
	ax0: 3, ax1: 0, axT: 2,
	expect: offsetDeltas{
		zeroZero: -3.618033988749895,
		ax0Delta: -1,
		ax1Delta: -1.618033988749895,
	},
}, {
	ax0: 3, ax1: 0, axT: 4,
	expect: offsetDeltas{
		zeroZero: 2.23606797749979,
		ax0Delta: 1.618033988749895,
		ax1Delta: 1.618033988749895,
	},
}, {
	ax0: 3, ax1: 1, axT: 0,
	expect: offsetDeltas{
		zeroZero: -3.618033988749895,
		ax0Delta: -1.618033988749895,
		ax1Delta: -1,
	},
}, {
	ax0: 3, ax1: 1, axT: 2,
	expect: offsetDeltas{
		zeroZero: 2.23606797749979,
		ax0Delta: 1.618033988749895,
		ax1Delta: 1.618033988749895,
	},
}, {
	ax0: 3, ax1: 1, axT: 4,
	expect: offsetDeltas{
		zeroZero: -3.618033988749895,
		ax0Delta: -1,
		ax1Delta: -1.618033988749895,
	},
}, {
	ax0: 3, ax1: 2, axT: 0,
	expect: offsetDeltas{
		zeroZero: -2.23606797749979,
		ax0Delta: -0.6180339887498949,
		ax1Delta: -0.6180339887498949,
	},
}, {
	ax0: 3, ax1: 2, axT: 1,
	expect: offsetDeltas{
		zeroZero: -1.381966011250105,
		ax0Delta: -1,
		ax1Delta: 0.6180339887498949,
	},
}, {
	ax0: 3, ax1: 2, axT: 4,
	expect: offsetDeltas{
		zeroZero: -1.381966011250105,
		ax0Delta: 0.6180339887498949,
		ax1Delta: -1,
	},
}, {
	ax0: 3, ax1: 4, axT: 0,
	expect: offsetDeltas{
		zeroZero: -1.381966011250105,
		ax0Delta: -1,
		ax1Delta: 0.6180339887498949,
	},
}, {
	ax0: 3, ax1: 4, axT: 1,
	expect: offsetDeltas{
		zeroZero: -2.23606797749979,
		ax0Delta: -0.6180339887498949,
		ax1Delta: -0.6180339887498949,
	},
}, {
	ax0: 3, ax1: 4, axT: 2,
	expect: offsetDeltas{
		zeroZero: -1.381966011250105,
		ax0Delta: 0.6180339887498949,
		ax1Delta: -1,
	},
}, {
	ax0: 4, ax1: 0, axT: 1,
	expect: offsetDeltas{
		zeroZero: -1.381966011250105,
		ax0Delta: -1,
		ax1Delta: 0.6180339887498949,
	},
}, {
	ax0: 4, ax1: 0, axT: 2,
	expect: offsetDeltas{
		zeroZero: -2.23606797749979,
		ax0Delta: -0.6180339887498949,
		ax1Delta: -0.6180339887498949,
	},
}, {
	ax0: 4, ax1: 0, axT: 3,
	expect: offsetDeltas{
		zeroZero: -1.381966011250105,
		ax0Delta: 0.6180339887498949,
		ax1Delta: -1,
	},
}, {
	ax0: 4, ax1: 1, axT: 0,
	expect: offsetDeltas{
		zeroZero: 2.23606797749979,
		ax0Delta: 1.618033988749895,
		ax1Delta: 1.618033988749895,
	},
}, {
	ax0: 4, ax1: 1, axT: 2,
	expect: offsetDeltas{
		zeroZero: -3.618033988749895,
		ax0Delta: -1.618033988749895,
		ax1Delta: -1,
	},
}, {
	ax0: 4, ax1: 1, axT: 3,
	expect: offsetDeltas{
		zeroZero: -3.618033988749895,
		ax0Delta: -1,
		ax1Delta: -1.618033988749895,
	},
}, {
	ax0: 4, ax1: 2, axT: 0,
	expect: offsetDeltas{
		zeroZero: -3.618033988749895,
		ax0Delta: -1,
		ax1Delta: -1.618033988749895,
	},
}, {
	ax0: 4, ax1: 2, axT: 1,
	expect: offsetDeltas{
		zeroZero: -3.618033988749895,
		ax0Delta: -1.618033988749895,
		ax1Delta: -1,
	},
}, {
	ax0: 4, ax1: 2, axT: 3,
	expect: offsetDeltas{
		zeroZero: 2.23606797749979,
		ax0Delta: 1.618033988749895,
		ax1Delta: 1.618033988749895,
	},
}, {
	ax0: 4, ax1: 3, axT: 0,
	expect: offsetDeltas{
		zeroZero: -1.381966011250105,
		ax0Delta: 0.6180339887498949,
		ax1Delta: -1,
	},
}, {
	ax0: 4, ax1: 3, axT: 1,
	expect: offsetDeltas{
		zeroZero: -2.23606797749979,
		ax0Delta: -0.6180339887498949,
		ax1Delta: -0.6180339887498949,
	},
}, {
	ax0: 4, ax1: 3, axT: 2,
	expect: offsetDeltas{
		zeroZero: -1.381966011250105,
		ax0Delta: -1,
		ax1Delta: 0.6180339887498949,
	},
}}

var threeAxesOffsetTestCases = []struct {
	ax0    uint8
	ax1    uint8
	axT    uint8
	expect float64
}{
	{ax0: 0, ax1: 1, axT: 2, expect: -1},
	{ax0: 0, ax1: 1, axT: 3, expect: -0.6180339887498949},
	{ax0: 0, ax1: 1, axT: 4, expect: 0.6180339887498949},
	{ax0: 0, ax1: 2, axT: 1, expect: 1.618033988749895},
	{ax0: 0, ax1: 2, axT: 3, expect: -1.618033988749895},
	{ax0: 0, ax1: 2, axT: 4, expect: -1},
	{ax0: 0, ax1: 3, axT: 1, expect: -1},
	{ax0: 0, ax1: 3, axT: 2, expect: -1.618033988749895},
	{ax0: 0, ax1: 3, axT: 4, expect: 1.618033988749895},
	{ax0: 0, ax1: 4, axT: 1, expect: 0.6180339887498949},
	{ax0: 0, ax1: 4, axT: 2, expect: -0.6180339887498949},
	{ax0: 0, ax1: 4, axT: 3, expect: -1},
	{ax0: 1, ax1: 0, axT: 2, expect: 0.6180339887498949},
	{ax0: 1, ax1: 0, axT: 3, expect: -0.6180339887498949},
	{ax0: 1, ax1: 0, axT: 4, expect: -1},
	{ax0: 1, ax1: 2, axT: 0, expect: 0.6180339887498949},
	{ax0: 1, ax1: 2, axT: 3, expect: -1},
	{ax0: 1, ax1: 2, axT: 4, expect: -0.6180339887498949},
	{ax0: 1, ax1: 3, axT: 0, expect: -1},
	{ax0: 1, ax1: 3, axT: 2, expect: 1.618033988749895},
	{ax0: 1, ax1: 3, axT: 4, expect: -1.618033988749895},
	{ax0: 1, ax1: 4, axT: 0, expect: 1.618033988749895},
	{ax0: 1, ax1: 4, axT: 2, expect: -1},
	{ax0: 1, ax1: 4, axT: 3, expect: -1.618033988749895},
	{ax0: 2, ax1: 0, axT: 1, expect: 1.618033988749895},
	{ax0: 2, ax1: 0, axT: 3, expect: -1},
	{ax0: 2, ax1: 0, axT: 4, expect: -1.618033988749895},
	{ax0: 2, ax1: 1, axT: 0, expect: -1},
	{ax0: 2, ax1: 1, axT: 3, expect: 0.6180339887498949},
	{ax0: 2, ax1: 1, axT: 4, expect: -0.6180339887498949},
	{ax0: 2, ax1: 3, axT: 0, expect: -0.6180339887498949},
	{ax0: 2, ax1: 3, axT: 1, expect: 0.6180339887498949},
	{ax0: 2, ax1: 3, axT: 4, expect: -1},
	{ax0: 2, ax1: 4, axT: 0, expect: -1.618033988749895},
	{ax0: 2, ax1: 4, axT: 1, expect: -1},
	{ax0: 2, ax1: 4, axT: 3, expect: 1.618033988749895},
	{ax0: 3, ax1: 0, axT: 1, expect: -1.618033988749895},
	{ax0: 3, ax1: 0, axT: 2, expect: -1},
	{ax0: 3, ax1: 0, axT: 4, expect: 1.618033988749895},
	{ax0: 3, ax1: 1, axT: 0, expect: -1.618033988749895},
	{ax0: 3, ax1: 1, axT: 2, expect: 1.618033988749895},
	{ax0: 3, ax1: 1, axT: 4, expect: -1},
	{ax0: 3, ax1: 2, axT: 0, expect: -0.6180339887498949},
	{ax0: 3, ax1: 2, axT: 1, expect: -1},
	{ax0: 3, ax1: 2, axT: 4, expect: 0.6180339887498949},
	{ax0: 3, ax1: 4, axT: 0, expect: -1},
	{ax0: 3, ax1: 4, axT: 1, expect: -0.6180339887498949},
	{ax0: 3, ax1: 4, axT: 2, expect: 0.6180339887498949},
	{ax0: 4, ax1: 0, axT: 1, expect: -1},
	{ax0: 4, ax1: 0, axT: 2, expect: -0.6180339887498949},
	{ax0: 4, ax1: 0, axT: 3, expect: 0.6180339887498949},
	{ax0: 4, ax1: 1, axT: 0, expect: 1.618033988749895},
	{ax0: 4, ax1: 1, axT: 2, expect: -1.618033988749895},
	{ax0: 4, ax1: 1, axT: 3, expect: -1},
	{ax0: 4, ax1: 2, axT: 0, expect: -1},
	{ax0: 4, ax1: 2, axT: 1, expect: -1.618033988749895},
	{ax0: 4, ax1: 2, axT: 3, expect: 1.618033988749895},
	{ax0: 4, ax1: 3, axT: 0, expect: 0.6180339887498949},
	{ax0: 4, ax1: 3, axT: 1, expect: -0.6180339887498949},
	{ax0: 4, ax1: 3, axT: 2, expect: -1},
}
