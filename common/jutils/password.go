package jutils

import (
	"crypto/rand"
	"encoding/base64"
	"slices"
	"unicode"
)

// Special chars are shared with the frontend.
var specialChars = []rune{'!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '_', '+', '-', '=', '[', ']', '{', '}', ';', ':', '\'', '"', '\\', '|', ',', '.', '<', '>', '/', '?'}

func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	// Very long passwords can be used for DoS attacks.
	if len(password) > 100 {
		return false
	}

	var containsNumber, containsUpper, containsSpecial bool
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			containsNumber = true
		case unicode.IsUpper(c):
			containsUpper = true
		case slices.Contains(specialChars, c):
			containsSpecial = true
		}

		if containsNumber && containsUpper && containsSpecial {
			return true
		}
	}

	return false
}

func GenerateSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
