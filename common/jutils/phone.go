package jutils

import (
	"fmt"
	"regexp"

	"github.com/nyaruka/phonenumbers/v2"
)

var (
	// reSequentialDigits verifies if there is a sequence of 5 or more consecutive increasing or decreasing digits
	reSequentialDigits = regexp.MustCompile(`(?:012345|123456|234567|345678|456789|567890|098765|987654|876543|765432|654321|543210)`)
)

type Phone struct {
	Prefix string `json:"prefix"`
	Number string `json:"number"`
}

func (p *Phone) IsValid(advancedValidation bool) bool {
	if p == nil {
		return false
	}

	if p.Prefix == "" || p.Number == "" {
		return false
	}

	number, err := phonenumbers.Parse(p.String(), "")
	if err != nil {
		return false
	}

	if !phonenumbers.IsValidNumber(number) {
		return false
	}

	if advancedValidation {
		if reSequentialDigits.MatchString(fmt.Sprintf("%d", number.GetNationalNumber())) {
			return false
		}

		if containsRepeatedSequence(fmt.Sprintf("%d", number.GetNationalNumber())) {
			return false
		}
	}

	return true
}

func (p Phone) String() string {
	return fmt.Sprintf("%s %s", p.Prefix, p.Number)
}

// containsRepeatedSequence checks if the input string contains sequences of repeated patterns:
// - A single digit repeated 7 or more times ("1111111").
// - A 2-digit pattern repeated 4 or more times ("12121212").
// - A 3-digit pattern repeated 3 or more times ("123123123").
// - A 4-digit pattern repeated 2 or more times ("12341234").
func containsRepeatedSequence(s string) bool {
	rules := map[int]int{
		// sequence size -> max nb repetitions
		1: 7,
		2: 4,
		3: 3,
		4: 2,
	}

	for size, maxReps := range rules {
		for start := range size {
			count := 1
			for i := start; i <= len(s)-size*2; i += size {
				if s[i:i+size] == s[i+size:i+size*2] {
					count++
					if count >= maxReps {
						return true
					}
				} else {
					count = 1
				}
			}
		}
	}

	return false
}
