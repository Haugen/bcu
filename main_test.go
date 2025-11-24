package main

import (
	"reflect"
	"testing"

	"github.com/Haugen/bcu/renderer"
)

func TestParseBranches(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []renderer.Branch
	}{
		{
			name: "typical git branch output with current branch",
			input: `* main
  feature-1
  feature-2
  bugfix-123`,
			expected: []renderer.Branch{
				{Name: "feature-1", IsActive: false},
				{Name: "feature-2", IsActive: false},
				{Name: "bugfix-123", IsActive: false},
			},
		},
		{
			name: "output with master branch",
			input: `  develop
* master
  hotfix-456`,
			expected: []renderer.Branch{
				{Name: "develop", IsActive: false},
				{Name: "hotfix-456", IsActive: false},
			},
		},
		{
			name: "output with both main and master",
			input: `* main
  master
  feature-a
  feature-b`,
			expected: []renderer.Branch{
				{Name: "feature-a", IsActive: false},
				{Name: "feature-b", IsActive: false},
			},
		},
		{
			name:     "only main branch exists",
			input:    "* main",
			expected: []renderer.Branch{},
		},
		{
			name:     "empty output",
			input:    "",
			expected: []renderer.Branch{},
		},
		{
			name: "branches with special characters",
			input: `* main
  feature/new-ui
  bugfix/JIRA-123
  release-1.0.0`,
			expected: []renderer.Branch{
				{Name: "feature/new-ui", IsActive: false},
				{Name: "bugfix/JIRA-123", IsActive: false},
				{Name: "release-1.0.0", IsActive: false},
			},
		},
		{
			name: "output with worktree branches",
			input: `* main
  test
  test1
  test2
+ worktree-test1
+ worktree-test2`,
			expected: []renderer.Branch{
				{Name: "test", IsActive: false},
				{Name: "test1", IsActive: false},
				{Name: "test2", IsActive: false},
				{Name: "worktree-test1", IsActive: true},
				{Name: "worktree-test2", IsActive: true},
			},
		},
		{
			name: "mixed worktrees and regular branches",
			input: `  develop
* main
+ worktree-feature
  feature-1
+ worktree-hotfix
  feature-2`,
			expected: []renderer.Branch{
				{Name: "develop", IsActive: false},
				{Name: "worktree-feature", IsActive: true},
				{Name: "feature-1", IsActive: false},
				{Name: "worktree-hotfix", IsActive: true},
				{Name: "feature-2", IsActive: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseBranches(tt.input)

			if len(result) == 0 && len(tt.expected) == 0 {
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseBranches() = %v, want %v", result, tt.expected)
			}
		})
	}
}
