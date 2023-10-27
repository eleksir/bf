/*
It maybe useful to apply some stemmers from https://github.com/blevesearch/snowballstem
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
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

var (
	pMarks   = []string{".", ",", "!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "{", "}", "<", ">", "[", "]", "\\"}
	pMarks2  = []string{"-", "_", "+", "=", ":", ";", "'", "`", "~", "\""}
	newLines = []string{"\n", "\r", "\n\r", "\r\n"}
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

// Learns phrases form given dictionary file.
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
		return fmt.Errorf("unable to open file %s: %w", filename, err)
	}

	defer func(fh *os.File) {
		if err := fh.Close(); err != nil {
			log.Printf("Unable to close %s cleanly: %s", filename, err)
		}
	}(fh)

	reader := bufio.NewReader(fh)

	for {
		line, err := reader.ReadString('\n')

		fmt.Printf("Line before norm: %s", line)
		line = nString(line)
		fmt.Printf("Line after norm: %s\n", line)

		if err != nil {
			if err == io.EOF {
				if line != "" {
					stuff := []string{line}
					c.Learn(stuff, bc)
				}

				break
			}

			return fmt.Errorf("unable to read %s: %w", filename, err)
		}

		if line != "" {
			stuff := []string{line}
			c.Learn(stuff, bc)
		}
	}

	return nil
}

// Prints similarity score for given phrase.
func checkPhrase(s string) error {
	classifier, err := bayesian.NewClassifierFromFile(dataFile)

	if err != nil {
		return fmt.Errorf("unable to open %s: %w", dataFile, err)
	}

	s = nString(s)

	scores, _, _ := classifier.LogScores([]string{s})

	_, err = fmt.Printf("Score %v\n", scores[0])

	return err
}

// Normalizes text buffer.
func nString(buf string) string {
	buf = strings.Trim(buf, "\n\r\t ")

	// Remove punctuation marks.
	for _, pMark := range pMarks {
		buf = strings.ReplaceAll(buf, pMark, "")
	}

	for _, pMark := range pMarks2 {
		buf = strings.ReplaceAll(buf, pMark, "")
	}

	// Replace newline sequences with space, even erroneous ones.
	for _, newline := range newLines {
		buf = strings.ReplaceAll(buf, newline, " ")
	}

	// Also replace any kind of space sequences with single space.
	buf = regexp.MustCompile(`\s+`).ReplaceAllString(buf, " ")

	// Transform all letters to lowercase
	buf = strings.ToLower(buf)

	return buf
}
