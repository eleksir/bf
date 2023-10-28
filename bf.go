package main

import (
	"bufio"
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

// Prints similarity score for given phrase.
func checkPhrase(s string) error {
	classifier, err := bayesian.NewClassifierFromFile(dataFile)

	if err != nil {
		return fmt.Errorf("unable to open %s: %w", dataFile, err)
	}

	s = nString(s)

	scores, _, _ := classifier.LogScores(strings.Split(s, " "))

	_, err = fmt.Printf("Score %v\n", scores[0])

	return err
}

// Learns phrases from predefined dictionary files.
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

// Feeds strings data from given file to given class.
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

		line = nString(line)

		if err != nil {
			// Handle EOF and trailing data
			if err == io.EOF {
				if line != "" {
					stuff := strings.Split(line, " ")
					c.Learn(stuff, bc)
				}

				break
			}

			return fmt.Errorf("unable to read %s: %w", filename, err)
		}

		if line != "" {
			stuff := strings.Split(line, " ")
			c.Learn(stuff, bc)
		}
	}

	return nil
}
