package jutils

import (
	"regexp"
	"strings"
)

var reNIN = regexp.MustCompile(`^[A-Z]{2}[0-9]{6}[A-D]$`)

// ValidateNIN checks if the provided string is a valid UK National Insurance Number (NIN).
func ValidateNIN(nin string) bool {
	return reNIN.MatchString(strings.ToUpper(nin))
}
