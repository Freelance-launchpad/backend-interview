package jutils

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type currencyValue interface {
	~int | ~int16 | ~int32 | ~int64 | ~float32 | ~float64 | ~string
}

var (
	// ErrInvalidCurrencyValue is returned when an invalid currency value is provided.
	ErrInvalidCurrencyValue = errors.New("INVALID_CURRENCY_VALUE")
)

// ToCents returns int64 cents based on a given string or number currency value.
func ToCents[T currencyValue](value T) (int64, error) {
	return currencyStringToCents(fmt.Sprintf("%v", value))
}

func currencyStringToCents(value string) (int64, error) {
	value = strings.ReplaceAll(value, ",", ".")

	f, ok := big.NewFloat(0).SetPrec(128).SetString(value)
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrInvalidCurrencyValue, value)
	}

	f.Mul(f, big.NewFloat(100))
	i, _ := f.Int(big.NewInt(0))

	return i.Int64(), nil
}

// ToEuros returns a human-readable euro string based on a given cents value.
func ToEuros(amount int64) string {
	return CentsToCurrency(amount, "€")
}

func CentsToCurrency(amount int64, currency string) string {
	value := fmt.Sprintf("%d", amount/100)

	cents := amount % 100
	value += fmt.Sprintf(",%02d", cents)

	return value + " " + currency
}

// TODO: better localization support
func CentsToCurrencyEN(amount int64, currency string) string {
	value := fmt.Sprintf("%d", amount/100)

	cents := amount % 100
	value += fmt.Sprintf(".%02d", cents)

	return currency + value
}

// EuroStringToCents returns a int64 in cents (2 decimals), will add decimal if necessary.
func EuroStringToCents(value string) (int64, error) {
	errMsg := "EuroStringToCents has failed with value: " + value

	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, ",", "")

	parts := strings.Split(value, ".")
	if len(parts) > 2 {
		return 0, fmt.Errorf("%s: INVALID_FORMAT", errMsg)
	}

	intPart := parts[0]
	decimalPart := "00"
	if len(parts) == 2 {
		decimalPart = parts[1]
	}

	// set the right number of decimals
	switch len(decimalPart) {
	case 0:
		decimalPart = "00"
	case 1:
		decimalPart = decimalPart + "0"
	default:
		decimalPart = decimalPart[:2]
	}

	// parse the final result
	amount, err := strconv.ParseInt(intPart+decimalPart, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s: ParseInt err: %w", errMsg, err)
	}

	return amount, nil
}
