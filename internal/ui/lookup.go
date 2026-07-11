package ui

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"github.com/chrisostomemataba/faltasi-cli/internal/api"
	"github.com/chrisostomemataba/faltasi-cli/internal/i18n"
)

func LookupFlow(client *api.Client) error {
	var method string
	err := huh.NewSelect[string]().
		Title(i18n.T("lookup_prompt")).
		Options(
			huh.NewOption(i18n.T("lookup_prompt"), "hwid"),
			huh.NewOption(i18n.T("lookup_or_phone"), "phone"),
		).
		Value(&method).
		Run()
	if err != nil {
		return err
	}

	var hardwareID, phone string
	if method == "hwid" {
		hardwareID, err = AskHardwareID(i18n.T("lookup_prompt"))
	} else {
		phone, err = AskPhone(i18n.T("lookup_or_phone"))
	}
	if err != nil {
		return err
	}

	var result *api.LicenseInfo
	err = WithSpinner(i18n.T("loading"), func() error {
		var lookupErr error
		result, lookupErr = client.Lookup(hardwareID, phone)
		return lookupErr
	})
	if err != nil {
		return err
	}

	if result == nil {
		fmt.Println(SubtleStyle.Render(i18n.T("customer_not_found")))
	} else {
		fmt.Println(SuccessStyle.Render(i18n.T("customer_found")))
		fmt.Println(RenderTable(
			[]string{i18n.T("phone_label"), i18n.T("pkg_name_label"), i18n.T("expires_label")},
			[][]string{{result.Phone, result.Package, FormatDate(result.ExpiresAt)}},
		))
	}

	_, _ = AskText(i18n.T("press_enter"), false)
	return nil
}
