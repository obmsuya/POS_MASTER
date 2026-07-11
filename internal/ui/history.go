package ui

import (
	"fmt"

	"github.com/chrisostomemataba/faltasi-cli/internal/api"
	"github.com/chrisostomemataba/faltasi-cli/internal/i18n"
)

func HistoryFlow(client *api.Client) error {
	var payments []api.PaymentRecord
	err := WithSpinner(i18n.T("loading"), func() error {
		var listErr error
		payments, listErr = client.PaymentHistory("")
		return listErr
	})
	if err != nil {
		return err
	}

	fmt.Println(TitleStyle.Render(i18n.T("history_title")))

	if len(payments) == 0 {
		fmt.Println(SubtleStyle.Render(i18n.T("history_empty")))
	} else {
		rows := make([][]string, len(payments))
		for i, payment := range payments {
			rows[i] = []string{
				payment.Phone,
				payment.Package,
				"TZS " + FormatMoney(payment.Amount),
				methodLabel(payment.Method),
				payment.ProcessedBy,
				FormatDate(payment.CreatedAt),
			}
		}
		fmt.Println(RenderTable(
			[]string{i18n.T("phone_label"), i18n.T("pkg_name_label"), i18n.T("amount_label"), i18n.T("method_label"), i18n.T("recorded_by"), i18n.T("expires_label")},
			rows,
		))
	}

	_, _ = AskText(i18n.T("press_enter"), false)
	return nil
}
