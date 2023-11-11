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

	state, err := getWindowState(config)
	if err != nil {
		handleError(err)
	}

	switch state {
	case Stopped:
		if err = StartApp(config); err != nil {
			handleError(err)
		}
		FocusWindow(config)
	case Unfocused:
		FocusWindow(config)
	case Focused:
		HideWindow(config)
	default:
		handleError(fmt.Errorf("ERROR: Unable to determine state for %s\n", config.id))
	}
}

func handleError(err error) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	os.Exit(1)
}
