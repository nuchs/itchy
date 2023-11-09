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

	appId := os.Args[1]
	jsonQuery := fmt.Sprintf(".. | select(.app_id? == \"%s\") | .focused", appId)
	fmt.Println("query = ", jsonQuery)

	sway := exec.Command("swaymsg", "-t", "get_tree", "-r")
	jq := exec.Command("jq", jsonQuery)

	var err error
	if jq.Stdin, err = sway.StdoutPipe(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to pipe between swaymsg and jq:", err)
		os.Exit(1)
	}

	var pipe io.ReadCloser
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
		fmt.Println("Start", appId)
		nvim := exec.Command("foot", "-a", appId, "nvim", "-S", "~/.local/share/nvim/session/notes")
		nvim.Stdin = os.Stdin
		nvim.Stdout = os.Stdout
		nvim.Stderr = os.Stderr
		err = nvim.Start()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Nvim failed to start:", err)
			os.Exit(1)
		}
		err = nvim.Process.Release()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Nvim failed to detach:", err)
			os.Exit(1)
		}
		time.Sleep(time.Millisecond * 30)
	}
	if string(result) == "true" {
		fmt.Println("Hide", appId)
		predicate := fmt.Sprintf("[app_id=%s]", appId)
		sway2 := exec.Command("swaymsg", predicate, "scratchpad", "show")
		if err = sway2.Run(); err != nil {
			fmt.Printf("Failed to send %s to scratchpad: %s\n", appId, err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Focus", appId)
		predicate := fmt.Sprintf("[app_id=%s]", appId)
		sway2 := exec.Command("swaymsg", predicate, "focus")
		if out, err := sway2.Output(); err != nil {
			fmt.Printf("Failed to focus %s:\n\t%s\n\t%s\n", appId, out, err)
			os.Exit(1)
		}
	}

	fmt.Println("Bye bye")
}
