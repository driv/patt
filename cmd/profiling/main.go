package main

import (
	"context"
	"log"
	"os"
	"patt"
	"runtime"
	"runtime/pprof"
)

func main() {

	// CPU profiling
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	args := []string{"patt", "[<_> <_>] [error] <_>", "", "./testdata/Apache_500MB.log"}
	err = patt.RunCLI(context.Background(), args, os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	f, err = os.Create("mem.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal(err)
	}
}
