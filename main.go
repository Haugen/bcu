package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
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
	fmt.Println(branches)
}
