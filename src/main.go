package main

import (
	"fmt"
	"io"
	"os"

	"oczko/src/impl/interpreter"
	"oczko/src/interface/runner"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "usage: oczko [program.ocz]")
		os.Exit(2)
	}
	var source []byte
	var err error
	if len(os.Args) == 2 {
		source, err = os.ReadFile(os.Args[1])
	} else {
		source, err = io.ReadAll(os.Stdin)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "oczko:", err)
		os.Exit(1)
	}
	var app runner.Runner = interpreter.New(os.Stdout)
	if err := app.Run(string(source)); err != nil {
		fmt.Fprintln(os.Stderr, "oczko:", err)
		os.Exit(1)
	}
}
