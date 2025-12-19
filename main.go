package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/iosifvieru/goparser/parser"
)

func main() {
	filePtr := flag.String("file", "", "input file path")
	flag.Parse()

	if *filePtr == "" {
		fmt.Println("Error: -file flag is required")
		flag.Usage()
		os.Exit(1)
	}

	err := parser.GenerateReport(*filePtr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate report: %v\n", err)
		os.Exit(1)
	}
}
