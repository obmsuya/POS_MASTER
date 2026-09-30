package validate

import (
	"strings"
	"testing"
)

func TestHardwareID(t *testing.T) {
	desktopID := strings.Repeat("ab12", 16)
	accepted := map[string]string{
		desktopID:                                      desktopID,
		"  " + strings.ToUpper(desktopID) + "\n":       desktopID,
		"cloud-0190f7a2-0000-7000-8000-000000000000":   "cloud-0190f7a2-0000-7000-8000-000000000000",
		" CLOUD-0190F7A2-0000-7000-8000-000000000000 ": "cloud-0190f7a2-0000-7000-8000-000000000000",
	}
	for input, want := range accepted {
		got, err := HardwareID(input)
		if err != nil || got != want {
			t.Errorf("HardwareID(%q) = %q, %v; want %q", input, got, err, want)
		}
	}

	rejected := []string{
		desktopID[:16],
		"abcd1234",
		"not a hardware id at all",
		strings.Repeat("g", 64),
		desktopID + "0",
		"cloud-0190f7a2-0000-7000-8000",
		"cloud-" + desktopID,
		"",
	}
	for _, input := range rejected {
		if _, err := HardwareID(input); err == nil {
			t.Errorf("HardwareID(%q) accepted, want rejected", input)
		}
	}
}
