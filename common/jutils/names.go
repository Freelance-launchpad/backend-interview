package jutils

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// NormalizeName returns the name in uppercase, without accents and whitespace.
func NormalizeName(name string) (string, error) {
	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.NotIn(unicode.Letter)),
		norm.NFC,
	)
	result, _, err := transform.String(t, name)
	if err != nil {
		return "", fmt.Errorf("normalizing name %s: %w", name, err)
	}

	return strings.ToUpper(result), nil
}
