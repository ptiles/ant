package main

import (
	"flag"
	"fmt"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/step"
	"github.com/ptiles/ant/utils"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
)

var (
	programName    = filepath.Base(os.Args[0])
	usageTextShort = "\nFor usage run: %s -h\n"
)

func main() {
	commonFlags := &utils.CommonFlags{}
	commonFlags.CommonFlagsSetup(pgrid.GridLinesTotal)
	flag.Parse()
	commonFlags.ParseArgs()

	fmt.Print("\n", commonFlags.String())

	utils.StartCPUProfile(commonFlags.Cpuprofile)
	defer utils.StopCPUProfile()

	rules, err := utils.GetRules(commonFlags.AntName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid name.  Should consist of at least two letters R L r l.")
		fmt.Fprintf(os.Stderr, usageTextShort, programName)
		os.Exit(1)
	}

	field := pgrid.New(commonFlags.Radius, rules, commonFlags.InitialPoint)

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
