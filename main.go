package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

func main() {
	fmt.Println("Hello")

	var selector string
	if os.Args[1] == "-a" {
		selector = "app_id"
	} else {
		selector = "class"
	}
	id := os.Args[2]
	jsonQuery := fmt.Sprintf(".. | select(.%s? == \"%s\") | .focused", selector, id)
	fmt.Println("query = ", jsonQuery)

	sway := exec.Command("swaymsg", "-t", "get_tree", "-r")
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
		fmt.Println("Start", id)
		var program *exec.Cmd
		var startupTime time.Duration
		if id == "notes" {
			program = exec.Command("foot", "-a", id, "nvim", "-S", "~/.local/share/nvim/session/notes")
			startupTime = time.Millisecond * 20
		} else {
			program = exec.Command("spotify")
			startupTime = time.Millisecond * 375
		}
		err = program.Start()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Nvim failed to start:", err)
			os.Exit(1)
		}
		err = program.Process.Release()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Nvim failed to detach:", err)
			os.Exit(1)
		}
		time.Sleep(startupTime)
	}
	if string(result) == "true" {
		fmt.Println("Hide", id)
		predicate := fmt.Sprintf("[%s=%s]", selector, id)
		sway2 := exec.Command("swaymsg", predicate, "scratchpad", "show")
		if err = sway2.Run(); err != nil {
			fmt.Printf("Failed to send %s to scratchpad: %s\n", id, err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Focus", selector, id)
		predicate := fmt.Sprintf("[%s=%s]", selector, id)
		sway2 := exec.Command("swaymsg", predicate, "focus")
		if out, err := sway2.Output(); err != nil {
			fmt.Printf("Failed to focus %s:\n\t%s\n\t%s\n", id, out, err)
			os.Exit(1)
		}
	}

	fmt.Println("Bye bye")
}
