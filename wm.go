package main

import (
	"fmt"
	"os"
	"os/exec"
)

type WindowState int

const (
	Unknown WindowState = iota
	Stopped
	Unfocused
	Focused
)

type WinErr string

func (e WinErr) Error() string {
	return string(e)
}

func parseWindowState(data string) (WindowState, error) {
	if len(data) == 0 {
		return Stopped, nil
	}

	if data == "false" {
		return Unfocused, nil
	}

	if data == "true" {
		return Focused, nil
	}

	return Unknown, WinErr(fmt.Sprintf("Unable to parse window state from %q", data))
}

func jsonWindowSelector(selector, id string) string {
	query := ""
	if selector == "class" {
		query = ".window_properties?"
	}

	return fmt.Sprintf(".. | select(%s.%s? == \"%s\") | .focused", query, selector, id)
}

func HideWindow(config Config) {
	fmt.Println("Hide", config.id)
	predicate := fmt.Sprintf("[%s=%s]", config.selector, config.id)
	sway2 := exec.Command("swaymsg", predicate, "scratchpad", "show")
	if err := sway2.Run(); err != nil {
		fmt.Printf("Failed to send %s to scratchpad: %s\n", config.id, err)
		os.Exit(1)
	}
}

func FocusWindow(config Config) {

	fmt.Println("Focus", config.selector, config.id)
	predicate := fmt.Sprintf("[%s=%s]", config.selector, config.id)
	sway2 := exec.Command("swaymsg", predicate, "focus")
	if out, err := sway2.Output(); err != nil {
		fmt.Printf("Failed to focus %s:\n\t%s\n\t%s\n", config.id, out, err)
		os.Exit(1)
	}
}

func getWindowState(config Config) (WindowState, error) {

	jsonWindowState := pipeline(
		exec.Command("swaymsg", "-t", "get_tree", "-r"),
		exec.Command("jq", jsonWindowSelector(config.selector, config.id)))

	state, err := parseWindowState(jsonWindowState)
	if err != nil {
		return Unknown, err
	}

	return state, nil
}
