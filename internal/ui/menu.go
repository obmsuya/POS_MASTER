package ui

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"github.com/obmsuya/POS_MASTER/internal/api"
	"github.com/obmsuya/POS_MASTER/internal/i18n"
)

type MenuChoice string

const (
	MenuActivate MenuChoice = "activate"
	MenuLookup   MenuChoice = "lookup"
	MenuHistory  MenuChoice = "history"
	MenuPackages MenuChoice = "packages"
	MenuExit     MenuChoice = "exit"
)

func MainMenu(user *api.User) (MenuChoice, error) {
	fmt.Println(SuccessStyle.Render(i18n.T("welcome", user.FullName)))

	options := []huh.Option[MenuChoice]{
		huh.NewOption(i18n.T("menu_activate"), MenuActivate),
		huh.NewOption(i18n.T("menu_lookup"), MenuLookup),
		huh.NewOption(i18n.T("menu_history"), MenuHistory),
	}
	if user.IsSuperuser {
		options = append(options, huh.NewOption(i18n.T("menu_packages"), MenuPackages))
	}
	options = append(options, huh.NewOption(i18n.T("menu_exit"), MenuExit))

	var choice MenuChoice
	err := huh.NewSelect[MenuChoice]().
		Title(i18n.T("main_menu_title")).
		Options(options...).
		Value(&choice).
		Run()
	return choice, err
}
