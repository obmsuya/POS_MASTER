package validate

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var digitsOnly = regexp.MustCompile(`[^0-9]`)

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

func NonEmpty(raw string) error {
	if strings.TrimSpace(raw) == "" {
		return errors.New("required")
	}
	return nil
}
