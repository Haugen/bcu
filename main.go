package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"github.com/Haugen/bcu/tea"
)

func main() {
	cmdResult, err := exec.Command("git", "branch").Output()

	if err != nil {
		fmt.Printf("error %s", err)
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

	selectionResult, err := tea.SelectBranches(branches)

	if err != nil {
		fmt.Printf("error %s", err)
	}

	fmt.Printf("\n %s \n", selectionResult)
}
