package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Haugen/bcu/renderer"
)

var version = "dev"

func main() {
	// Handle --version flag
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("bcu version %s\n", version)
		return
	}

	cmdResult, err := exec.Command("git", "branch").CombinedOutput()

	if err != nil {
		fmt.Printf("Error running git command: %s\n", err)
		os.Exit(1)
	}

	output := string(cmdResult)
	if strings.HasPrefix(output, "fatal") {
		fmt.Println(output)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanWords)
	var branches []string
	for scanner.Scan() {
		text := scanner.Text()
		if text != "*" && text != "main" && text != "master" {
			branches = append(branches, text)
		}
	}

	if len(branches) == 0 {
		fmt.Println("No branches to clean up!")
		return
	}

	renderer := renderer.NewRenderer(branches)
	selectedBranches := renderer.Run()

	if len(selectedBranches) == 0 {
		fmt.Println("No branches selected. Exiting.")
		return
	}

	args := append([]string{"branch", "-D"}, selectedBranches...)
	deleteCmd := exec.Command("git", args...)
	deleteCmd.Stdout = os.Stdout
	deleteCmd.Stderr = os.Stderr
	deleteCmd.Run()
}
