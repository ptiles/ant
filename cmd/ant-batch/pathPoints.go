package main

import (
	"fmt"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/utils"
	"iter"
	"strings"
)

func (fl *Flags) PathPoints(debug *strings.Builder) iter.Seq[string] {
	return func(yield func(string) bool) {
		if fl.initialPointPath == "" {
			return
		}
		if fl.debug {
			debug.WriteString("\nPathPoints: ")
		}

		pathFlags := utils.CommonFlags{}
		pathFlags.ParseShorthand(fl.initialPointPath)
		field := pgrid.New(pathFlags.Pattern, pathFlags.AntRules, pathFlags.InitialPoint)

		steps := pathFlags.Steps
		for i, turn := range field.RunTurns(pathFlags.Steps.Max) {
			if i >= steps.Min && steps.Inc > 0 && i%steps.Inc == 0 {
				turnString := turn.String()
				if !yield(fmt.Sprintf(" -i %s", turnString)) {
					return
				}
				if fl.debug {
					debug.WriteString(" ")
					debug.WriteString(turnString)
				}
			}
		}
	}
}
