package main

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"

	"github.com/obmsuya/POS_MASTER/internal/api"
	"github.com/obmsuya/POS_MASTER/internal/i18n"
	"github.com/obmsuya/POS_MASTER/internal/session"
	"github.com/obmsuya/POS_MASTER/internal/ui"
)

func main() {
	ui.PrintBanner()

	lang, err := ui.PickLanguage()
	if err != nil {
		sayGoodbye()
		return
	}
	i18n.Set(lang)

	client := api.New()
	user, err := authenticate(client)
	if err != nil {
		if !errors.Is(err, huh.ErrUserAborted) {
			fmt.Println(ui.ErrorStyle.Render(err.Error()))
		}
		sayGoodbye()
		return
	}

mainLoop:
	for {
		choice, err := ui.MainMenu(user)
		if err != nil {
			break mainLoop
		}

		var flowErr error
		switch choice {
		case ui.MenuActivate:
			flowErr = ui.ActivateFlow(client)
		case ui.MenuLookup:
			flowErr = ui.LookupFlow(client)
		case ui.MenuHistory:
			flowErr = ui.HistoryFlow(client)
		case ui.MenuPackages:
			flowErr = ui.PackagesFlow(client)
		case ui.MenuExit:
			break mainLoop
		}

		if flowErr != nil && !errors.Is(flowErr, huh.ErrUserAborted) {
			fmt.Println(ui.ErrorStyle.Render(flowErr.Error()))
		}
	}

	sayGoodbye()
}

func sayGoodbye() {
	fmt.Println(ui.SuccessStyle.Render(i18n.T("goodbye")))
}

func authenticate(client *api.Client) (*api.User, error) {
	if cached, err := session.Load(); err == nil && cached != nil {
		client.RestoreSession(cached.AccessToken, cached.RefreshToken)
		if refreshErr := client.RefreshAccessToken(); refreshErr == nil {
			fmt.Println(ui.SuccessStyle.Render(i18n.T("session_restored", cached.FullName)))
			return &api.User{
				FullName:    cached.FullName,
				IsSuperuser: cached.IsSuperuser,
				UserType:    cached.UserType,
			}, nil
		}
		_ = session.Clear()
	}

	for {
		user, err := ui.Login(client)
		if err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return nil, huh.ErrUserAborted
			}
			fmt.Println(ui.ErrorStyle.Render(err.Error()))
			continue
		}

		_ = session.Save(session.Cached{
			AccessToken:  client.AccessToken,
			RefreshToken: client.RefreshToken,
			FullName:     user.FullName,
			IsSuperuser:  user.IsSuperuser,
			UserType:     user.UserType,
		})

		return user, nil
	}
}
