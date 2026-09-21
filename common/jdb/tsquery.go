package jdb

import (
	"regexp"
	"strings"
)

var reToTSQuerySpecialChars = regexp.MustCompile(`[&|!(){}\[\]^"'~*?:\\<>]`)

// ToTSQueryInput tokenizes and sanitizes a raw search query to be used as input for to_tsquery
func ToTSQueryInput(rawSearchQuery string) string {
	// Remove special characters used by to_tsquery (see: https://www.postgresql.org/docs/12/textsearch-controls.html#TEXTSEARCH-PARSING-QUERIES)
	// Note that there's no risk of injection, the goal here is to avoid returning user facing errors
	rawSearchQuery = reToTSQuerySpecialChars.ReplaceAllString(rawSearchQuery, " ")

	rawSearchTerms := strings.Fields(rawSearchQuery)

	searchTerms := make([]string, 0, len(rawSearchTerms))
	for _, term := range rawSearchTerms {
		// Append :* to each term to match partial words
		searchTerms = append(searchTerms, term+":*")
	}

	searchQuery := strings.Join(searchTerms, " & ")

	return searchQuery
}
