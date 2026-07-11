package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func FormatDate(iso string) string {
	if iso == "" {
		return "-"
	}
	parsed, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return iso
	}
	return parsed.Format("2 Jan 2006")
}

func FormatMoney(raw string) string {
	whole := strings.SplitN(raw, ".", 2)[0]
	negative := strings.HasPrefix(whole, "-")
	if negative {
		whole = whole[1:]
	}
	if _, err := strconv.Atoi(whole); err != nil {
		return raw
	}

	var grouped []string
	for len(whole) > 3 {
		grouped = append([]string{whole[len(whole)-3:]}, grouped...)
		whole = whole[:len(whole)-3]
	}
	grouped = append([]string{whole}, grouped...)

	result := strings.Join(grouped, ",")
	if negative {
		result = "-" + result
	}
	return result
}

func FormatMoneyFloat(amount float64) string {
	return FormatMoney(fmt.Sprintf("%.0f", amount))
}
