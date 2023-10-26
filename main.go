package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jbrukh/bayesian"
)

const (
	Bad                bayesian.Class = "Bad"
	Good               bayesian.Class = "Good"
	goodDictionaryFile string         = "./data/good_dictionary.txt"
	badDictionaryFile  string         = "./data/bad_dictionary.txt"
	dataFile           string         = "./data/data.bin"
)

func main() {
	var (
		helpFlag  = flag.Bool("help", false, "displays help message")
		learnFlag = flag.Bool("learn", false, "learns data from sample good and bad dictionaries")
		phrase    = flag.String("phrase", "", "check phrase for similarity criterion")
	)

	flag.Parse()

	switch {
	case *helpFlag == true:
		printHelp()
		os.Exit(0)

	case *learnFlag == true:
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

// printHelp Prints help message
func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("\tbf <option> [arg]")
	fmt.Println("")
	fmt.Println("where option can be:")
	fmt.Println("\t--help - show help")
	fmt.Println("\t--learn - learn data from dictionary data/data.txt")
	fmt.Println("\t--phrase 'text' - check test for similarity to dictionary data")
}

// Learns phrases form ./data/dictionary.txt
func learn() error {
	classifier := bayesian.NewClassifier(Bad, Good)

	if err := feed(classifier, goodDictionaryFile, Good); err != nil {
		return err
	}

	if err := feed(classifier, badDictionaryFile, Bad); err != nil {
		return err
	}

	if err := classifier.WriteToFile(dataFile); err != nil {
		return fmt.Errorf("unable to save data to file %s: %w", dataFile, err)
	}

	return nil
}

func feed(c *bayesian.Classifier, filename string, bc bayesian.Class) error {
	fh, err := os.Open(filename)

	if err != nil {
		log.Fatalf("Unable to open file %s: %s\n", filename, err)
	}

	defer func(fh *os.File) {
		err := fh.Close()

		if err != nil {
			log.Printf("Unable to close %s cleanly: %s", filename, err)
		}
	}(fh)

	reader := bufio.NewReader(fh)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				value := strings.Trim(line, "\n\r\t ")

				if value != "" {
					// stuff := regexp.MustCompile(`\s+`).Split(value, -1)
					stuff := []string{value}
					c.Learn(stuff, bc)
				}

				break
			}

			return fmt.Errorf("unable to read %s: %w", filename, err)
		}

		value := strings.Trim(line, "\n\r\t ")

		if value != "" {
			// stuff := regexp.MustCompile(`\s+`).Split(value, -1)
			stuff := []string{value}
			c.Learn(stuff, bc)
		}
	}

	return nil
}

// Prints similarity score for given phrase
func checkPhrase(s string) error {
	classifier, err := bayesian.NewClassifierFromFile(dataFile)

	if err != nil {
		return fmt.Errorf("unable to open %s: %w", dataFile, err)
	}

	s = strings.Trim(s, "\n\r\t ")

	// classifier.ConvertTermsFreqToTfIdf()

	// Split string into words and feed it to classifier
	// scores, _, _ := classifier.LogScores(regexp.MustCompile(`\s+`).Split(s, -1))
	scores, _, _ := classifier.LogScores([]string{s})

	_, err = fmt.Printf("Score %v\n", scores[0])

	return err
}
