package ui

import (
	"errors"

	"github.com/charmbracelet/huh"

	"github.com/chrisostomemataba/faltasi-cli/internal/i18n"
	"github.com/chrisostomemataba/faltasi-cli/internal/validate"
)

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

func AskPhone(title string) (string, error) {
	var raw string
	err := huh.NewInput().Title(title).Value(&raw).Validate(validatePhoneText).Run()
	if err != nil {
		return "", err
	}
	normalized, _ := parsePhone(raw)
	return normalized, nil
}

func AskHardwareID(title string) (string, error) {
	var raw string
	err := huh.NewInput().Title(title).Value(&raw).Validate(validateHardwareIDText).Run()
	if err != nil {
		return "", err
	}
	clean, _ := parseHardwareID(raw)
	return clean, nil
}

func AskAmount(title string) (float64, error) {
	var raw string
	err := huh.NewInput().Title(title).Value(&raw).Validate(validateAmountText).Run()
	if err != nil {
		return 0, err
	}
	amount, _ := parseAmount(raw)
	return amount, nil
}

func AskPositiveInt(title string) (int, error) {
	var raw string
	err := huh.NewInput().Title(title).Value(&raw).Validate(validatePositiveIntText).Run()
	if err != nil {
		return 0, err
	}
	value, _ := parsePositiveInt(raw)
	return value, nil
}

func Confirm(title string) (bool, error) {
	var value bool
	err := huh.NewConfirm().Title(title).Value(&value).Run()
	return value, err
}
