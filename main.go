package main

import (
	"fmt"
	"os"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		handleError(err)
	}

	state, err := GetWindowState(config)
	if err != nil {
		handleError(err)
	}

	switch state {
	case Stopped:
		if err = StartApp(config); err != nil {
			handleError(err)
		}
		if err = FocusWindow(config); err != nil {
			handleError(err)
		}

	case Unfocused:
		if err = FocusWindow(config); err != nil {
			handleError(err)
		}

	case Focused:
		if err = HideWindow(config); err != nil {
			handleError(err)
		}

	default:
		handleError(fmt.Errorf("ERROR: Unable to determine state for %s\n", config.id))
	}
}

func handleError(err error) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	os.Exit(1)
}
