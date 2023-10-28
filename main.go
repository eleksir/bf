package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var (
		helpFlag  = flag.Bool("help", false, "displays help message")
		learnFlag = flag.Bool("learn", false, "learns data from sample good and bad dictionaries")
		phrase    = flag.String("phrase", "", "check phrase for similarity criterion")
	)

	flag.Parse()

	switch {
	case *helpFlag:
		printHelp()
		os.Exit(0)

	case *learnFlag:
		if err := learn(); err != nil {
			log.Fatalf("Error: %s", err)
		}

	case *phrase != "":
		if err := checkPhrase(*phrase); err != nil {
			log.Fatalf("Error: %s", err)
		}

	default:
		printHelp()
		os.Exit(1)
	}
}

// printHelp Prints help message.
func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("\tbf <option> [arg]")
	fmt.Println("")
	fmt.Println("where option can be:")
	fmt.Println("\t--help - show help")
	fmt.Println("\t--learn - learn data from dictionary data/data.txt")
	fmt.Println("\t--phrase 'text' - check test for similarity to dictionary data")
}
