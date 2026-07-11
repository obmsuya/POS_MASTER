package ui

import (
	"github.com/charmbracelet/huh"

	"github.com/chrisostomemataba/faltasi-cli/internal/i18n"
)

// PickLanguage is the one screen allowed to show both languages at once —
// nobody has picked a preference yet. Every screen after this uses only
// the chosen language.
func PickLanguage() (i18n.Lang, error) {
	var choice string
	err := huh.NewSelect[string]().
		Title("Chagua lugha / Choose a language").
		Options(
			huh.NewOption("Kiswahili", "sw"),
			huh.NewOption("English", "en"),
		).
		Value(&choice).
		Run()
	if err != nil {
		return i18n.Swahili, err
	}
	if choice == "en" {
		return i18n.English, nil
	}
	return i18n.Swahili, nil
}
