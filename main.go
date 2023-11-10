package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

func main() {
	config := loadConfig()

	sway := exec.Command("swaymsg", "-t", "get_tree", "-r")

	jsonQuery := fmt.Sprintf(".. | select(.%s? == \"%s\") | .focused", config.selector, config.id)
	jq := exec.Command("jq", jsonQuery)

	var err error
	if jq.Stdin, err = sway.StdoutPipe(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to pipe between swaymsg and jq:", err)
		os.Exit(1)
	}

	var pipe io.Reader
	if pipe, err = jq.StdoutPipe(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect pipe to jq:", err)
		os.Exit(1)
	}

	if err = sway.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to start swaymsg:", err)
		os.Exit(1)
	}

	if err = jq.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to start jq:", err)
		os.Exit(1)
	}

	var result []byte
	if result, err = io.ReadAll(pipe); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read from jq pipe:", err)
		os.Exit(1)
	}

	if err = jq.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "Jq failed to complete:", err)
		os.Exit(1)
	}

	if err = sway.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "Swaymsg failed to complete:", err)
		os.Exit(1)
	}

	if len(result) == 0 {
		fmt.Println("Start", config.id)
		program := exec.Command(config.command, config.args...)
		err = program.Start()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s failed to start: %s\n", config.command, err)
			os.Exit(1)
		}
		err = program.Process.Release()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s failed to detach: %s\n", config.command, err)
			os.Exit(1)
		}

		time.Sleep(config.startTime)
	}

	if string(result) == "true" {
		fmt.Println("Hide", config.id)
		predicate := fmt.Sprintf("[%s=%s]", config.selector, config.id)
		sway2 := exec.Command("swaymsg", predicate, "scratchpad", "show")
		if err = sway2.Run(); err != nil {
			fmt.Printf("Failed to send %s to scratchpad: %s\n", config.id, err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Focus", config.selector, config.id)
		predicate := fmt.Sprintf("[%s=%s]", config.selector, config.id)
		sway2 := exec.Command("swaymsg", predicate, "focus")
		if out, err := sway2.Output(); err != nil {
			fmt.Printf("Failed to focus %s:\n\t%s\n\t%s\n", config.id, out, err)
			os.Exit(1)
		}
	}
}
