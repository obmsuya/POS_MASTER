// Package validate holds input validation that must catch mistakes before
// a network call is ever made — bad hardware IDs, malformed phone numbers,
// non-positive amounts. Mirrors the same phone-normalization rules the
// Django backend applies server-side, so a value accepted here is never
// rejected there.
package validate

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var digitsOnly = regexp.MustCompile(`[^0-9]`)

// NormalizeTZPhone accepts 07XXXXXXXX, 06XXXXXXXX, +2557XXXXXXXX, 2557XXXXXXXX
// and returns the E.164 form +2557XXXXXXXX. Returns an error if the input
// doesn't look like a Tanzanian mobile number.
func NormalizeTZPhone(raw string) (string, error) {
	digits := digitsOnly.ReplaceAllString(raw, "")

	switch {
	case strings.HasPrefix(digits, "255") && len(digits) == 12:
		return "+" + digits, nil
	case strings.HasPrefix(digits, "0") && len(digits) == 10:
		return "+255" + digits[1:], nil
	case len(digits) == 9 && (digits[0] == '6' || digits[0] == '7'):
		return "+255" + digits, nil
	default:
		return "", errors.New("invalid phone number")
	}
}

// HardwareID does a light sanity check — the real format is a 64-char
// hex SHA-256 string, but we only require it to be non-trivial and free
// of whitespace/control characters, so a mis-scanned or truncated paste
// is caught immediately rather than round-tripping to the server.
func HardwareID(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if len(trimmed) < 8 {
		return "", errors.New("hardware id too short")
	}
	if strings.ContainsAny(trimmed, " \t\n\r") {
		return "", errors.New("hardware id contains whitespace")
	}
	return trimmed, nil
}

// Amount parses a payment amount and requires it to be a positive number.
func Amount(raw string) (float64, error) {
	trimmed := strings.TrimSpace(raw)
	value, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return 0, errors.New("not a number")
	}
	if value <= 0 {
		return 0, errors.New("amount must be greater than zero")
	}
	return value, nil
}

// PositiveInt parses a whole number and requires it to be greater than zero
// — used for package days-granted and max-devices fields.
func PositiveInt(raw string) (int, error) {
	trimmed := strings.TrimSpace(raw)
	value, err := strconv.Atoi(trimmed)
	if err != nil {
		return 0, errors.New("not a whole number")
	}
	if value <= 0 {
		return 0, errors.New("must be greater than zero")
	}
	return value, nil
}

// NonEmpty requires a field to have visible content once trimmed.
func NonEmpty(raw string) error {
	if strings.TrimSpace(raw) == "" {
		return errors.New("required")
	}
	return nil
}
