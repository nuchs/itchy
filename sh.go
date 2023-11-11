package main

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

func Pipeline(commands ...*exec.Cmd) (string, error) {
	var tailPipe io.Reader
	var err error
	for _, cmd := range commands {
		if tailPipe != nil {
			cmd.Stdin = tailPipe
		}

		if tailPipe, err = cmd.StdoutPipe(); err != nil {
			return "", fmt.Errorf("Failed to pipe into %s: %s\n", cmd.String(), err)
		}
	}

	for _, cmd := range commands {
		if err = cmd.Start(); err != nil {
			return "", fmt.Errorf("Failed to start %s: %s\n", cmd.String(), err)
		}
	}

	result, err := io.ReadAll(tailPipe)
	if err != nil {
		return "", fmt.Errorf("Failed to read from pipeline: %s", err)
	}

	for _, cmd := range commands {
		if err = cmd.Wait(); err != nil {
			return "", fmt.Errorf("%s failed to complete: %s\n", cmd.String(), err)
		}
	}

	return strings.TrimSpace(string(result)), nil
}

func StartApp(config Config) error {
	program := exec.Command(config.command, config.args...)
	err := program.Start()
	if err != nil {
		return fmt.Errorf("%s failed to start: %s\n", config.command, err)
	}

	err = program.Process.Release()
	if err != nil {
		return fmt.Errorf("%s failed to detach: %s\n", config.command, err)
	}

	time.Sleep(config.startTime)

	return nil
}
