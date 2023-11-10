package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type config struct {
	id        string
	selector  string
	startTime time.Duration
	command   string
	args      []string
}

func loadConfig() config {
	app_id := flag.String("a", "", "'app_id' of window to manage")
	class := flag.String("c", "", "'class' of window to manage")
	startTime := flag.Int("t", 20, "How many ms to wait after starting app before trying to find its window")
	flag.Usage = usage
	flag.Parse()

	if (*app_id != "" && *class != "") || (*app_id == "" && *class == "") {
		fmt.Fprint(flag.CommandLine.Output(), "ERROR: you must specify exactly one of -a or -c")
		usage()
		os.Exit(1)
	}

	if *startTime < 0 {
		fmt.Fprint(flag.CommandLine.Output(), "ERROR: start time must be positive")
		usage()
		os.Exit(1)
	}

	if len(flag.Args()) < 1 {
		fmt.Fprint(flag.CommandLine.Output(), "ERROR: you must provide a command to run")
		usage()
		os.Exit(1)
	}

	var id string
	var selector string

	if *app_id != "" {
		selector = "app_id"
		id = *app_id
	} else {
		selector = "class"
		id = *class
	}

	return config{
		id,
		selector,
		time.Duration(time.Duration(*startTime) * time.Millisecond),
		flag.Arg(0),
		flag.Args()[1:],
	}
}

func usage() {
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
