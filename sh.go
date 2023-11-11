package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func pipeline(commands ...*exec.Cmd) string {
	var tailPipe io.Reader
	var err error
	for _, cmd := range commands {
		if tailPipe != nil {
			cmd.Stdin = tailPipe
		}

		if tailPipe, err = cmd.StdoutPipe(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to pipe into %s: %s\n", cmd.String(), err)
			os.Exit(1)
		}
	}

	for _, cmd := range commands {
		if err = cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to start %s: %s\n", cmd.String(), err)
			os.Exit(1)
		}
	}

	result, err := io.ReadAll(tailPipe)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read from pipeline:", err)
		os.Exit(1)
	}

	for _, cmd := range commands {
		if err = cmd.Wait(); err != nil {
			fmt.Fprintf(os.Stderr, "%s failed to complete: %s\n", cmd.String(), err)
			os.Exit(1)
		}
	}

	return strings.TrimSpace(string(result))
}

func StartApp(config Config) {

	fmt.Println("Start", config.id)
	program := exec.Command(config.command, config.args...)
	err := program.Start()
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
