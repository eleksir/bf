package main

import (
	"regexp"
	"strings"
)

var (
	pMarks   = []string{".", ",", "!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "{", "}", "<", ">", "[", "]", "\\"}
	pMarks2  = []string{"-", "_", "+", "=", ":", ";", "'", "`", "~", "\""}
	newLines = []string{"\n", "\r", "\n\r", "\r\n"}
)

// Normalizes text buffer. Remove newlines, leading and trailing spaces, punctuation marks and replace repeated spaces
// with single one.
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
