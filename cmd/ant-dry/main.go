package main

import (
	"flag"
	"fmt"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/step"
	"github.com/ptiles/ant/utils"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	commonFlags := &utils.CommonFlags{}
	commonFlags.CommonFlagsSetup(pgrid.GridLinesTotal)
	flag.Parse()
	commonFlags.ParseArgs()

	fmt.Print("\n", commonFlags.String())

	utils.StartCPUProfile(commonFlags.Cpuprofile)
	defer utils.StopCPUProfile()

	field := pgrid.New(commonFlags.Pattern, commonFlags.AntRules, commonFlags.InitialPoint)

	step.DryRunStepper(field, commonFlags.Steps.Max, commonFlags.MaxNoisyDots)
	fmt.Printf(" %s\n", commonFlags.String())

	if commonFlags.Memprofile != "" {
		f, mpErr := os.Create(commonFlags.Memprofile)
		if mpErr != nil {
			log.Fatal("could not create memory profile: ", mpErr)
		}
		defer f.Close()
		runtime.GC()
		if wrErr := pprof.WriteHeapProfile(f); wrErr != nil {
			log.Fatal("could not write memory profile: ", wrErr)
		}
	}
}
