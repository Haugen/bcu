package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Haugen/bcu/renderer"
)

func main() {
	cmdResult, err := exec.Command("git", "branch").CombinedOutput()

	if string(cmdResult)[0:5] == "fatal" {
		fmt.Println(string(cmdResult))
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Something went wrong with running \"git branch\": %s", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(cmdResult)))
	scanner.Split(bufio.ScanWords)
	var branches []string
	for scanner.Scan() {
		text := scanner.Text()
		if text != "*" && text != "main" && text != "master" {
			branches = append(branches, text)
		}
	}

	if len(branches) > 0 {
		renderer := renderer.NewRenderer(branches)
		selected := renderer.Run()
		fmt.Println(selected)
	}
}
