package ui

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"github.com/obmsuya/POS_MASTER/internal/api"
	"github.com/obmsuya/POS_MASTER/internal/i18n"
)

func PackagesFlow(client *api.Client) error {
	var packages []api.Package
	err := WithSpinner(i18n.T("loading"), func() error {
		var listErr error
		packages, listErr = client.ListPackages()
		return listErr
	})
	if err != nil {
		return err
	}

	if len(packages) > 0 {
		rows := make([][]string, len(packages))
		for i, pkg := range packages {
			rows[i] = []string{
				pkg.Name,
				"TZS " + FormatMoney(pkg.Price),
				fmt.Sprintf("%d %s", pkg.DaysGranted, i18n.T("days_suffix")),
				fmt.Sprintf("%d %s", pkg.MaxDevices, i18n.T("devices_suffix")),
			}
		}
		fmt.Println(RenderTable(
			[]string{i18n.T("pkg_name_label"), i18n.T("pkg_price_label"), i18n.T("pkg_days_label"), i18n.T("pkg_devices_label")},
			rows,
		))
	} else {
		fmt.Println(SubtleStyle.Render(i18n.T("no_packages")))
	}

	const (
		actionCreate = "create"
		actionEdit   = "edit"
		actionBack   = "back"
	)
	options := []huh.Option[string]{huh.NewOption(i18n.T("pkg_create"), actionCreate)}
	if len(packages) > 0 {
		options = append(options, huh.NewOption(i18n.T("pkg_edit"), actionEdit))
	}
	options = append(options, huh.NewOption(i18n.T("pkg_back"), actionBack))

	var action string
	if err := huh.NewSelect[string]().Title(i18n.T("packages_title")).Options(options...).Value(&action).Run(); err != nil {
		return err
	}

	switch action {
	case actionCreate:
		return createPackage(client)
	case actionEdit:
		return editPackage(client, packages)
	default:
		return nil
	}
}

func createPackage(client *api.Client) error {
	name, err := AskText(i18n.T("pkg_name_label"), true)
	if err != nil {
		return err
	}
	price, err := AskAmount(i18n.T("pkg_price_label"))
	if err != nil {
		return err
	}
	days, err := AskPositiveInt(i18n.T("pkg_days_label"))
	if err != nil {
		return err
	}
	maxDevices, err := AskPositiveInt(i18n.T("pkg_devices_label"))
	if err != nil {
		return err
	}

	err = WithSpinner(i18n.T("loading"), func() error {
		_, createErr := client.CreatePackage(api.CreatePackageInput{
			Name: name, Price: price, DaysGranted: days, MaxDevices: maxDevices,
		})
		return createErr
	})
	if err != nil {
		return err
	}

	fmt.Println(SuccessStyle.Render("✅ " + i18n.T("pkg_created")))
	_, _ = AskText(i18n.T("press_enter"), false)
	return nil
}

func editPackage(client *api.Client, packages []api.Package) error {
	options := make([]huh.Option[int], len(packages))
	for i, pkg := range packages {
		options[i] = huh.NewOption(pkg.Name, pkg.ID)
	}
	var packageID int
	if err := huh.NewSelect[int]().Title(i18n.T("pkg_edit")).Options(options...).Value(&packageID).Run(); err != nil {
		return err
	}

	var selected api.Package
	for _, pkg := range packages {
		if pkg.ID == packageID {
			selected = pkg
		}
	}

	fmt.Println(SubtleStyle.Render(i18n.T("pkg_edit_notice")))

	name := selected.Name
	priceText := selected.Price
	daysText := fmt.Sprintf("%d", selected.DaysGranted)
	devicesText := fmt.Sprintf("%d", selected.MaxDevices)

	if err := huh.NewInput().Title(i18n.T("pkg_name_label")).Value(&name).
		Validate(validateNonEmpty).Run(); err != nil {
		return err
	}
	if err := huh.NewInput().Title(i18n.T("pkg_price_label")).Value(&priceText).
		Validate(validateAmountText).Run(); err != nil {
		return err
	}
	if err := huh.NewInput().Title(i18n.T("pkg_days_label")).Value(&daysText).
		Validate(validatePositiveIntText).Run(); err != nil {
		return err
	}
	if err := huh.NewInput().Title(i18n.T("pkg_devices_label")).Value(&devicesText).
		Validate(validatePositiveIntText).Run(); err != nil {
		return err
	}

	price, _ := parseAmount(priceText)
	days, _ := parsePositiveInt(daysText)
	maxDevices, _ := parsePositiveInt(devicesText)

	err := WithSpinner(i18n.T("loading"), func() error {
		_, updateErr := client.UpdatePackage(packageID, api.CreatePackageInput{
			Name: name, Price: price, DaysGranted: days, MaxDevices: maxDevices,
		})
		return updateErr
	})
	if err != nil {
		return err
	}

	fmt.Println(SuccessStyle.Render("✅ " + i18n.T("pkg_updated")))
	_, _ = AskText(i18n.T("press_enter"), false)
	return nil
}
