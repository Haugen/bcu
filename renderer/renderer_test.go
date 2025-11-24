package renderer

import (
	"reflect"
	"testing"
)

// Helper function to create test branches
func testBranches(names ...string) []Branch {
	branches := make([]Branch, len(names))
	for i, name := range names {
		branches[i] = Branch{Name: name, IsActive: false}
	}
	return branches
}

func TestNewState(t *testing.T) {
	branches := []Branch{
		{Name: "feature-1", IsActive: false},
		{Name: "feature-2", IsActive: false},
		{Name: "bugfix-1", IsActive: false},
	}
	state := NewState(branches)

	if state.cursorPos != 0 {
		t.Errorf("Expected initial cursor position 0, got %d", state.cursorPos)
	}

	if len(state.list) != 3 {
		t.Errorf("Expected 3 branches, got %d", len(state.list))
	}

	if len(state.selected) != 0 {
		t.Errorf("Expected no selections initially, got %d", len(state.selected))
	}
}

func TestMoveCursorDown(t *testing.T) {
	state := NewState(testBranches("branch1", "branch2", "branch3"))

	// Move down from 0 to 1
	moved := state.MoveCursorDown()
	if !moved {
		t.Error("Expected cursor to move down")
	}
	if state.cursorPos != 1 {
		t.Errorf("Expected cursor at 1, got %d", state.cursorPos)
	}

	// Move down from 1 to 2
	moved = state.MoveCursorDown()
	if !moved {
		t.Error("Expected cursor to move down")
	}
	if state.cursorPos != 2 {
		t.Errorf("Expected cursor at 2, got %d", state.cursorPos)
	}

	// Try to move down from last position (should not move)
	moved = state.MoveCursorDown()
	if moved {
		t.Error("Expected cursor NOT to move down from last position")
	}
	if state.cursorPos != 2 {
		t.Errorf("Expected cursor to stay at 2, got %d", state.cursorPos)
	}
}

func TestMoveCursorUp(t *testing.T) {
	state := NewState(testBranches("branch1", "branch2", "branch3"))
	state.cursorPos = 2 // Start at last position

	// Move up from 2 to 1
	moved := state.MoveCursorUp()
	if !moved {
		t.Error("Expected cursor to move up")
	}
	if state.cursorPos != 1 {
		t.Errorf("Expected cursor at 1, got %d", state.cursorPos)
	}

	// Move up from 1 to 0
	moved = state.MoveCursorUp()
	if !moved {
		t.Error("Expected cursor to move up")
	}
	if state.cursorPos != 0 {
		t.Errorf("Expected cursor at 0, got %d", state.cursorPos)
	}

	// Try to move up from first position (should not move)
	moved = state.MoveCursorUp()
	if moved {
		t.Error("Expected cursor NOT to move up from first position")
	}
	if state.cursorPos != 0 {
		t.Errorf("Expected cursor to stay at 0, got %d", state.cursorPos)
	}
}

func TestMoveCursorUpDownSequence(t *testing.T) {
	// This test specifically addresses the duplication bug scenario
	state := NewState(testBranches("branch1", "branch2", "branch3"))

	// Move down, then back up to position 0
	state.MoveCursorDown() // 0 -> 1
	state.MoveCursorDown() // 1 -> 2
	state.MoveCursorUp()   // 2 -> 1
	state.MoveCursorUp()   // 1 -> 0

	if state.cursorPos != 0 {
		t.Errorf("Expected cursor at 0, got %d", state.cursorPos)
	}

	// Verify we can move down again from 0
	moved := state.MoveCursorDown()
	if !moved {
		t.Error("Expected to be able to move down from position 0")
	}
	if state.cursorPos != 1 {
		t.Errorf("Expected cursor at 1, got %d", state.cursorPos)
	}
}

func TestToggleSelection(t *testing.T) {
	state := NewState(testBranches("branch1", "branch2", "branch3"))

	// Initially nothing selected
	if state.selected[0] {
		t.Error("Expected branch 0 to not be selected initially")
	}

	// Select branch 0
	state.ToggleSelection()
	if !state.selected[0] {
		t.Error("Expected branch 0 to be selected")
	}

	// Deselect branch 0
	state.ToggleSelection()
	if state.selected[0] {
		t.Error("Expected branch 0 to be deselected")
	}
}

func TestToggleMultipleSelections(t *testing.T) {
	state := NewState(testBranches("branch1", "branch2", "branch3"))

	// Select branch 0
	state.ToggleSelection()

	// Move to branch 1 and select
	state.MoveCursorDown()
	state.ToggleSelection()

	// Move to branch 2 and select
	state.MoveCursorDown()
	state.ToggleSelection()

	// All three should be selected
	if !state.selected[0] || !state.selected[1] || !state.selected[2] {
		t.Error("Expected all branches to be selected")
	}

	// Deselect branch 1
	state.MoveCursorUp() // Back to position 1
	state.ToggleSelection()

	if state.selected[1] {
		t.Error("Expected branch 1 to be deselected")
	}
	if !state.selected[0] || !state.selected[2] {
		t.Error("Expected branches 0 and 2 to still be selected")
	}
}

func TestGetSelectedBranches(t *testing.T) {
	state := NewState(testBranches("feature-1", "feature-2", "bugfix-1", "hotfix-1"))

	// Select branches at indices 0, 2
	state.selected[0] = true
	state.selected[2] = true

	result := state.GetSelectedBranches()
	expected := []string{"feature-1", "bugfix-1"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetSelectedBranchesEmpty(t *testing.T) {
	state := NewState(testBranches("feature-1", "feature-2"))

	result := state.GetSelectedBranches()

	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %v", result)
	}
}

func TestGetOutputLines(t *testing.T) {
	state := NewState(testBranches("branch1", "branch2"))

	lines := state.GetOutputLines()

	// Should have header + empty line + 2 branches = 4 lines
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(lines))
	}

	// First line should be the header
	if lines[0] == "" {
		t.Error("Expected non-empty header line")
	}

	// Second line should be empty
	if lines[1] != "" {
		t.Error("Expected empty second line")
	}

	// Third line should show cursor on first branch
	if lines[2] != "> [ ] branch1" {
		t.Errorf("Expected '> [ ] branch1', got '%s'", lines[2])
	}

	// Fourth line should not show cursor
	if lines[3] != "  [ ] branch2" {
		t.Errorf("Expected '  [ ] branch2', got '%s'", lines[3])
	}
}

func TestGetOutputLinesWithSelection(t *testing.T) {
	state := NewState(testBranches("branch1", "branch2"))
	state.cursorPos = 1
	state.selected[1] = true

	lines := state.GetOutputLines()

	// Cursor on branch2, which is selected
	if lines[3] != "> [x] branch2" {
		t.Errorf("Expected '> [x] branch2', got '%s'", lines[3])
	}

	// branch1 not selected, no cursor
	if lines[2] != "  [ ] branch1" {
		t.Errorf("Expected '  [ ] branch1', got '%s'", lines[2])
	}
}

func TestGetOutputLinesCount(t *testing.T) {
	// This test verifies that the line count is always consistent
	// which would catch the duplication bug
	state := NewState(testBranches("b1", "b2", "b3"))

	// Get initial line count
	expectedCount := 3 + 2 // branches + header + empty line

	if len(state.GetOutputLines()) != expectedCount {
		t.Errorf("Expected %d lines initially, got %d", expectedCount, len(state.GetOutputLines()))
	}

	// Move cursor around and verify line count stays the same
	state.MoveCursorDown()
	if len(state.GetOutputLines()) != expectedCount {
		t.Error("Line count changed after moving cursor down")
	}

	state.MoveCursorUp()
	if len(state.GetOutputLines()) != expectedCount {
		t.Error("Line count changed after moving cursor up")
	}

	state.ToggleSelection()
	if len(state.GetOutputLines()) != expectedCount {
		t.Error("Line count changed after toggling selection")
	}
}

func TestToggleSelectionOnActiveBranch(t *testing.T) {
	// Test that active branches cannot be selected
	branches := []Branch{
		{Name: "branch1", IsActive: false},
		{Name: "branch2", IsActive: true}, // Active branch
		{Name: "branch3", IsActive: false},
	}
	state := NewState(branches)

	// Try to select an active branch (position 1)
	state.cursorPos = 1
	state.ToggleSelection()

	// Should not be selected
	if state.selected[1] {
		t.Error("Expected active branch to not be selectable")
	}

	// Regular branch should be selectable
	state.cursorPos = 0
	state.ToggleSelection()
	if !state.selected[0] {
		t.Error("Expected non-active branch to be selectable")
	}
}

func TestGetOutputLinesWithActiveBranch(t *testing.T) {
	branches := []Branch{
		{Name: "branch1", IsActive: false},
		{Name: "branch2", IsActive: true},
	}
	state := NewState(branches)

	lines := state.GetOutputLines()

	// Active branch should have (checked out) indicator and no checkbox
	if lines[3] != "      branch2 (checked out)" {
		t.Errorf("Expected '      branch2 (checked out)', got '%s'", lines[3])
	}

	// Regular branch should not have indicator
	if lines[2] != "> [ ] branch1" {
		t.Errorf("Expected '> [ ] branch1', got '%s'", lines[2])
	}
}
