package main

import (
	"fmt"
	"os"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		Usage()
		os.Exit(1)
	}

	state, err := getWindowState(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}

	switch state {
	case Stopped:
		StartApp(config)
		FocusWindow(config)
	case Unfocused:
		FocusWindow(config)
	case Focused:
		HideWindow(config)
	default:
		fmt.Printf("ERROR: Unable to determine state for %s\n", config.id)
		os.Exit(1)
	}
}
