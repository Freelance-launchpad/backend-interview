package jutils

import "regexp"

var auth0IDRegex = regexp.MustCompile(`^[0-9A-Fa-f]{24}$`)

func IsAuth0ID(s string) bool {
	return auth0IDRegex.MatchString(s)
}
