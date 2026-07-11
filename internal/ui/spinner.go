package ui

import (
	"fmt"
	"time"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func WithSpinner(label string, work func() error) error {
	done := make(chan error, 1)
	go func() {
		done <- work()
	}()

	frame := 0
	ticker := time.NewTicker(90 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case err := <-done:
			fmt.Print("\r\033[K")
			return err
		case <-ticker.C:
			fmt.Printf("\r\033[K%s %s", TitleStyle.Render(spinnerFrames[frame%len(spinnerFrames)]), label)
			frame++
		}
	}
}
