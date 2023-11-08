package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("Hello")

	appId := os.Args[1]
	fmt.Println(appId)

	sway, err := exec.Command("swaymsg", "-t", "get_tree", "-r").Output()
	if err != nil {
		fmt.Println("Failed to run swaymsg:", err)
		os.Exit(1)
	}

	var result map[string]string
	json.Unmarshal(sway, &result)

	fmt.Printf("Swaymsg: %s\n", sway)

	fmt.Println("Bye bye")
}
