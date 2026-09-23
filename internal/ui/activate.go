package ui

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"

	"github.com/obmsuya/POS_MASTER/internal/api"
	"github.com/obmsuya/POS_MASTER/internal/i18n"
)

func ActivateFlow(client *api.Client) error {
	hardwareID, err := AskHardwareID(i18n.T("lookup_prompt"))
	if err != nil {
		return err
	}

	var existing *api.LicenseInfo
	err = WithSpinner(i18n.T("loading"), func() error {
		var lookupErr error
		existing, lookupErr = client.Lookup(hardwareID, "")
		return lookupErr
	})
	if err != nil {
		return err
	}

	var phone string
	if existing != nil {
		fmt.Println(SuccessStyle.Render(i18n.T("customer_found")))
		fmt.Println(RenderTable(
			[]string{i18n.T("phone_label"), i18n.T("pkg_name_label"), i18n.T("expires_label")},
			[][]string{{existing.Phone, existing.Package, FormatDate(existing.ExpiresAt)}},
		))
	} else {
		fmt.Println(SubtleStyle.Render(i18n.T("customer_not_found")))
		phone, err = AskPhone(i18n.T("new_customer_phone"))
		if err != nil {
			return err
		}
	}

	var packages []api.Package
	err = WithSpinner(i18n.T("loading"), func() error {
		var listErr error
		packages, listErr = client.ListPackages()
		return listErr
	})
	if err != nil {
		return err
	}
	if len(packages) == 0 {
		return errors.New(i18n.T("no_packages"))
	}

	options := make([]huh.Option[int], len(packages))
	for i, pkg := range packages {
		label := fmt.Sprintf("%s — TZS %s (%d %s)", pkg.Name, FormatMoney(pkg.Price), pkg.DaysGranted, i18n.T("days_suffix"))
		options[i] = huh.NewOption(label, pkg.ID)
	}

	var packageID int
	err = huh.NewSelect[int]().
		Title(i18n.T("pick_package")).
		Options(options...).
		Value(&packageID).
		Run()
	if err != nil {
		return err
	}

	amount, err := AskAmount(i18n.T("amount_label"))
	if err != nil {
		return err
	}

	var method string
	err = huh.NewSelect[string]().
		Title(i18n.T("method_label")).
		Options(
			huh.NewOption(i18n.T("method_cash"), "cash"),
			huh.NewOption(i18n.T("method_mobile_money"), "mobile_money"),
			huh.NewOption(i18n.T("method_bank"), "bank_transfer"),
		).
		Value(&method).
		Run()
	if err != nil {
		return err
	}

	reference, err := AskText(i18n.T("reference_label"), false)
	if err != nil {
		return err
	}

	fmt.Println(RenderTable(
		[]string{i18n.T("amount_label"), i18n.T("method_label")},
		[][]string{{"TZS " + FormatMoneyFloat(amount), methodLabel(method)}},
	))
	confirmed, err := Confirm(i18n.T("confirm_prompt"))
	if err != nil {
		return err
	}
	if !confirmed {
		return nil
	}

	var result *api.LicenseInfo
	err = WithSpinner(i18n.T("activating"), func() error {
		var recordErr error
		result, recordErr = client.RecordPayment(api.RecordPaymentInput{
			HardwareID: hardwareID,
			Phone:      phone,
			PackageID:  packageID,
			Amount:     amount,
			Method:     method,
			Reference:  reference,
		})
		return recordErr
	})
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			return errors.New(i18n.T("error_generic", apiErr.Message))
		}
		return errors.New(i18n.T("error_network"))
	}

	fmt.Println(SuccessStyle.Render("✅ " + i18n.T("activated_title")))
	fmt.Println(i18n.T("activated_body", result.Package, result.DaysGranted))
	fmt.Println(RenderTable(
		[]string{i18n.T("phone_label"), i18n.T("expires_label")},
		[][]string{{result.Phone, FormatDate(result.ExpiresAt)}},
	))
	_, _ = AskText(i18n.T("press_enter"), false)
	return nil
}

func methodLabel(method string) string {
	switch method {
	case "cash":
		return i18n.T("method_cash")
	case "mobile_money":
		return i18n.T("method_mobile_money")
	case "bank_transfer":
		return i18n.T("method_bank")
	default:
		return method
	}
}
