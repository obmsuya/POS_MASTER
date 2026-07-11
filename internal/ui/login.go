package ui

import (
	"errors"
	"fmt"

	"github.com/chrisostomemataba/faltasi-cli/internal/api"
	"github.com/chrisostomemataba/faltasi-cli/internal/i18n"
)

func Login(client *api.Client) (*api.User, error) {
	fmt.Println(TitleStyle.Render(i18n.T("login_title")))

	phone, err := AskPhone(i18n.T("phone_label"))
	if err != nil {
		return nil, err
	}
	password, err := AskPassword(i18n.T("password_label"))
	if err != nil {
		return nil, err
	}

	var user *api.User
	err = WithSpinner(i18n.T("connecting"), func() error {
		var loginErr error
		user, loginErr = client.Login(phone, password)
		return loginErr
	})
	if err != nil {
		if apiErr, ok := err.(*api.APIError); ok {
			return nil, fmt.Errorf("%s", i18n.T("login_failed", apiErr.Message))
		}
		return nil, errors.New(i18n.T("error_network"))
	}

	if !user.IsSuperuser && user.UserType != "partner" {
		return nil, errors.New(i18n.T("role_not_allowed"))
	}

	return user, nil
}
