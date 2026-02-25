// Package main provides a command-line interface for the resolution theorem prover.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/thxrsxm/res/internal/clause"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: res [options] <clause1> <clause2> ...\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "A resolution theorem prover for propositional logic.\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  <clause>    A clause in the format: A,B,-C (comma-separated literals)\n")
		fmt.Fprintf(os.Stderr, "              Each literal is a single letter (A-Z) optionally prefixed with '-'\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Output:\n")
		fmt.Fprintf(os.Stderr, "  [ ]         The clause set is unsatisfiable (contradiction found)\n")
		fmt.Fprintf(os.Stderr, "  [x]         The clause set is satisfiable (no contradiction found)\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  res a,-a\n")
		fmt.Fprintf(os.Stderr, "  res \"a,b\" \"-a,c\" \"-b,c\" \"-c\"\n")
		fmt.Fprintf(os.Stderr, "  res a,b,-c -a,b,c -b,c -c\n")
		fmt.Fprintf(os.Stderr, "  res -- -a,b,-c -a,b,c -b,c -c\n")
	}
	// Stop flag parsing after first non-flag argument
	flag.CommandLine.SetOutput(os.Stderr)
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(0)
	}
	set := []clause.Clause{}
	for i := range flag.Args() {
		c, err := clause.Parse(flag.Args()[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing clause %q: %v\n", flag.Args()[i], err)
			os.Exit(1)
		}
		set = append(set, *c)
	}
	result := clause.Res(set, 0)
	if result {
		fmt.Println("[ ]")
	} else {
		fmt.Println("[x]")
	}
}
