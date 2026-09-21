package jutils

import (
	"regexp"
	"strings"
	"unicode"
)

// PointerEqual returns true if pointers are deeply equal.
func PointerEqual[T comparable](first *T, second *T) bool {
	if (first == nil && second != nil) ||
		(second == nil && first != nil) ||
		(first != nil && second != nil &&
			*first != *second) {
		return false
	}

	return true
}

var (
	reSpaces = regexp.MustCompile(`\s+`)
	rePoint  = regexp.MustCompile(`\.\s*`)
	reComma  = regexp.MustCompile(`,\s*`)
)

// CleanupText cleanup an input string:
// - Remove leading and trailing spaces
// - Replace multiple spaces by a single space
// - Add spaces after commas and dots
// - Capitalize the first letter of the string
// - Add a dot at the end if not already there
func CleanupText(input string) string {
	// Remove leading and trailing spaces
	formattedText := strings.TrimSpace(input)

	// Remove multiple spaces
	formattedText = reSpaces.ReplaceAllString(formattedText, " ")

	// Fix spaces after dots
	formattedText = rePoint.ReplaceAllString(formattedText, ". ")

	// Fix spaces after commas
	formattedText = reComma.ReplaceAllString(formattedText, ", ")

	// Split into sentences while keeping the dots
	sentences := strings.SplitAfter(formattedText, ".")
	for i, sentence := range sentences {
		trimmedSentence := strings.TrimSpace(sentence)
		if len(trimmedSentence) > 0 && unicode.IsLetter(rune(trimmedSentence[0])) {
			sentences[i] = capitalizeFirstLetter(trimmedSentence)
		} else {
			sentences[i] = trimmedSentence
		}
	}

	// Recompose the text
	formattedText = strings.Join(sentences, " ")

	// Remove any extra space before the final dot
	formattedText = strings.TrimSpace(formattedText)

	// Add ending dot if missing
	if !strings.HasSuffix(formattedText, ".") {
		formattedText += "."
	}

	return formattedText
}

func capitalizeFirstLetter(s string) string {
	return strings.ToUpper(string(s[0])) + s[1:]
}
