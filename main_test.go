package main

import (
	"reflect"
	"testing"
)

func TestParseBranches(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name: "typical git branch output with current branch",
			input: `* main
  feature-1
  feature-2
  bugfix-123`,
			expected: []string{"feature-1", "feature-2", "bugfix-123"},
		},
		{
			name: "output with master branch",
			input: `  develop
* master
  hotfix-456`,
			expected: []string{"develop", "hotfix-456"},
		},
		{
			name: "output with both main and master",
			input: `* main
  master
  feature-a
  feature-b`,
			expected: []string{"feature-a", "feature-b"},
		},
		{
			name:     "only main branch exists",
			input:    "* main",
			expected: []string{},
		},
		{
			name:     "empty output",
			input:    "",
			expected: []string{},
		},
		{
			name: "branches with special characters",
			input: `* main
  feature/new-ui
  bugfix/JIRA-123
  release-1.0.0`,
			expected: []string{"feature/new-ui", "bugfix/JIRA-123", "release-1.0.0"},
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
