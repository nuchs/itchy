package main

import (
	"flag"
	"fmt"
	"time"
)

const (
	ErrSelector  = ConfigErr("you must specify exactly one of -a or -c")
	ErrStartTime = ConfigErr("start time must be positive")
	ErrCommand   = ConfigErr("you must provide a command to run")
)

var noConfig = Config{}

type ConfigErr string

func (c ConfigErr) Error() string {
	return string(c)
}

type Config struct {
	id        string
	selector  string
	startTime time.Duration
	command   string
	args      []string
}

func LoadConfig() (Config, error) {
	app_id := flag.String("a", "", "'app_id' of window to manage")
	class := flag.String("c", "", "'class' of window to manage")
	startTime := flag.Int("t", 20, "How many ms to wait after starting app before trying to find its window")
	flag.Usage = Usage

	flag.Parse()
	if err := validateCommandLine(*app_id, *class, *startTime); err != nil {
		return noConfig, err
	}

	return newConfig(*app_id, *class, *startTime), nil
}

func Usage() {
	fmt.Fprintf(flag.CommandLine.Output(), `
itch - Utility for managing multiple sway scratchpads

Searches for the window identified by the provided app_id (-a) or class (-c)
and either sends it to the scratchpad or displays it in the active
workspace. If no such window can be found the provided command is run (the
assumption being it will create a window with the specified id)

Usage:
  itch -a|-c <id> <command> [command arguments]...
    
`)
	flag.PrintDefaults()
	fmt.Fprintln(flag.CommandLine.Output())
}

func validateCommandLine(app_id, class string, startTime int) error {
	if (app_id != "" && class != "") || (app_id == "" && class == "") {
		return ErrSelector
	}

	if startTime < 0 {
		return ErrStartTime
	}

	if len(flag.Args()) < 1 {
		return ErrCommand
	}

	return nil
}

func newConfig(app_id, class string, startTime int) Config {
	var id string
	var selector string

	if app_id != "" {
		selector = "app_id"
		id = app_id
	} else {
		selector = "class"
		id = class
	}

	return Config{
		id,
		selector,
		time.Duration(time.Duration(startTime) * time.Millisecond),
		flag.Arg(0),
		flag.Args()[1:],
	}
}
