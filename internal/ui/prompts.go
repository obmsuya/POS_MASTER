package ui

import (
	"errors"

	"github.com/charmbracelet/huh"

	"github.com/chrisostomemataba/faltasi-cli/internal/i18n"
	"github.com/chrisostomemataba/faltasi-cli/internal/validate"
)

// Shared validators — used both by the Ask* one-shot prompts below and by
// screens that pre-fill a huh.Input with an existing value (e.g. editing
// a package) and need the exact same rule applied.

func validateNonEmpty(s string) error {
	if err := validate.NonEmpty(s); err != nil {
		return errors.New(i18n.T("required_field"))
	}
	return nil
}

func validatePhoneText(s string) error {
	if _, err := validate.NormalizeTZPhone(s); err != nil {
		return errors.New(i18n.T("invalid_phone"))
	}
	return nil
}

func validateHardwareIDText(s string) error {
	if _, err := validate.HardwareID(s); err != nil {
		return errors.New(i18n.T("invalid_hwid"))
	}
	return nil
}

func validateAmountText(s string) error {
	if _, err := validate.Amount(s); err != nil {
		return errors.New(i18n.T("invalid_amount"))
	}
	return nil
}

func validatePositiveIntText(s string) error {
	if _, err := validate.PositiveInt(s); err != nil {
		return errors.New(i18n.T("invalid_input_retry"))
	}
	return nil
}

func parseAmount(s string) (float64, error)    { return validate.Amount(s) }
func parsePositiveInt(s string) (int, error)   { return validate.PositiveInt(s) }
func parsePhone(s string) (string, error)      { return validate.NormalizeTZPhone(s) }
func parseHardwareID(s string) (string, error) { return validate.HardwareID(s) }

// AskText prompts for a single free-text value with an optional required-field check.
func AskText(title string, required bool) (string, error) {
	var value string
	input := huh.NewInput().Title(title).Value(&value)
	if required {
		input = input.Validate(validateNonEmpty)
	}
	if err := input.Run(); err != nil {
		return "", err
	}
	return value, nil
}

// AskPassword prompts for a password, masking input.
func AskPassword(title string) (string, error) {
	var value string
	err := huh.NewInput().
		Title(title).
		Value(&value).
		Password(true).
		Validate(validateNonEmpty).
		Run()
	if err != nil {
		return "", err
	}
	return value, nil
}

// AskPhone prompts for a phone number, validating and normalizing it to
// E.164 before returning.
func AskPhone(title string) (string, error) {
	var raw string
	err := huh.NewInput().Title(title).Value(&raw).Validate(validatePhoneText).Run()
	if err != nil {
		return "", err
	}
	normalized, _ := parsePhone(raw)
	return normalized, nil
}

// AskHardwareID prompts for a hardware ID with a basic sanity check.
func AskHardwareID(title string) (string, error) {
	var raw string
	err := huh.NewInput().Title(title).Value(&raw).Validate(validateHardwareIDText).Run()
	if err != nil {
		return "", err
	}
	clean, _ := parseHardwareID(raw)
	return clean, nil
}

// AskAmount prompts for a positive payment amount.
func AskAmount(title string) (float64, error) {
	var raw string
	err := huh.NewInput().Title(title).Value(&raw).Validate(validateAmountText).Run()
	if err != nil {
		return 0, err
	}
	amount, _ := parseAmount(raw)
	return amount, nil
}

// AskPositiveInt prompts for a whole number greater than zero.
func AskPositiveInt(title string) (int, error) {
	var raw string
	err := huh.NewInput().Title(title).Value(&raw).Validate(validatePositiveIntText).Run()
	if err != nil {
		return 0, err
	}
	value, _ := parsePositiveInt(raw)
	return value, nil
}

// Confirm asks a yes/no question.
func Confirm(title string) (bool, error) {
	var value bool
	err := huh.NewConfirm().Title(title).Value(&value).Run()
	return value, err
}
