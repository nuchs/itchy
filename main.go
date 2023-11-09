package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("Hello")

	appId := os.Args[1]
	jsonQuery := fmt.Sprintf(".. | select(.app_id? == \"%s\")", appId)
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

	fmt.Printf("Focused window app_id = %s\n", result)

	fmt.Println("Bye bye")
}
