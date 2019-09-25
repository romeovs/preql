package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/romeovs/preql"
)

var (
	help = flag.Bool("help", false, "Show help")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("preql: ")

	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if len(args) != 1 {
		usage()
		os.Exit(1)
	}

	if *help {
		usage()
		os.Exit(0)
	}

	pkg, err := preql.Load(args[0])
	if err != nil {
		log.Fatal(err)
	}

	err = pkg.Generate()
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <package name>\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Builds scanners and queries for all relevant types and methods in the specified package")
	flag.PrintDefaults()
}
