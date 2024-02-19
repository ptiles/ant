package main

import "github.com/ptiles/ant/pgrid"

func initialState(field *pgrid.Field) (pgrid.GridPoint, pgrid.GridPoint, pgrid.GridLine, pgrid.GridLine) {
	initialLine := pgrid.GridLine{Axis: pgrid.E, Offset: 0}
	prevLine := pgrid.GridLine{Axis: pgrid.A, Offset: 0}
	currLine := pgrid.GridLine{Axis: pgrid.B, Offset: 0}

	prevPoint := field.MakeGridPoint(initialLine, prevLine)
	currPoint := field.MakeGridPoint(prevLine, currLine)

	return prevPoint, currPoint, prevLine, currLine
}
