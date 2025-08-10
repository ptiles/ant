package main

import (
	"fmt"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/utils"
	"iter"
	"strings"
)

type path struct {
	flags   utils.CommonFlags
	present bool
}

func (p *path) parser() flagParser {
	return func(path string) error {
		if path == "" {
			return nil
		}

		p.flags.ParseShorthand(path)
		p.present = true

		return nil
	}
}

func (p *path) skip() bool {
	return !p.present
}

func (p *path) seq(debug *strings.Builder) iter.Seq[string] {
	return func(yield func(string) bool) {
		if p.skip() {
			return
		}

		if debug != nil {
			debug.WriteString("\nPathPoints: ")
		}

		pathFlags := p.flags
		field := pgrid.New(pathFlags.Pattern, pathFlags.AntRules, pathFlags.InitialPoint)

		steps := pathFlags.Steps
		for i, turn := range field.RunTurns(pathFlags.Steps.Max) {
			if i >= steps.Min && steps.Inc > 0 && i%steps.Inc == 0 {
				turnString := turn.String()
				if !yield(fmt.Sprintf(" -i %s", turnString)) {
					return
				}
				if debug != nil {
					debug.WriteString(" ")
					debug.WriteString(turnString)
				}
			}
		}
	}
}
