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

// parseBranches extracts branch names from git branch output,
// filtering out protected branches (main, master), and marking branches
// that are checked out in a worktree as active
func parseBranches(gitOutput string) []renderer.Branch {
	scanner := bufio.NewScanner(strings.NewReader(gitOutput))
	var branches []renderer.Branch
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		isActive := false

		// Check if branch is checked out in current worktree (*)
		if strings.HasPrefix(line, "* ") {
			isActive = true
			line = strings.TrimPrefix(line, "* ")
			line = strings.TrimSpace(line)
		}

		// Check if branch is checked out in another worktree (+)
		if strings.HasPrefix(line, "+ ") {
			isActive = true
			line = strings.TrimPrefix(line, "+ ")
			line = strings.TrimSpace(line)
		}

		// Skip protected branches
		if line == "main" || line == "master" {
			continue
		}

		branches = append(branches, renderer.Branch{
			Name:     line,
			IsActive: isActive,
		})
	}
	return branches
}

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

	branches := parseBranches(output)

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
	if err := deleteCmd.Run(); err != nil {
		os.Exit(1)
	}
}
