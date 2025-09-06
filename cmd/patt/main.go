package main

import (
	"context"
	"os"
	"os/signal"
	"patt"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	exitIfErr(patt.RunCLI(ctx, os.Args, os.Stdin, os.Stdout))
}

func exitIfErr(err error) {
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
