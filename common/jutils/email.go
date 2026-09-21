package jutils

import (
	"regexp"

	"github.com/bombsimon/tld-validator"
)

var emailRegex = regexp.MustCompile(`^[\w\-.+_]+@[\w\-_.]+\.+\w{2,63}$`)

// IsValidEmail checks if the email is valid.
func IsValidEmail(email string) bool {
	if !emailRegex.MatchString(email) {
		return false
	}

	if !tld.FromDomainName(email).IsValid() {
		return false
	}

	return true
}
