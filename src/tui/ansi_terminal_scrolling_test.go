package tui

import (
	"testing"
)

func TestScrollHistory(t *testing.T) {
	term := NewAnsiTerminal(10, 5)
	term.historySize = 10

	// Fill screen
	for range 5 {
		term.Write([]byte("Line\n"))
	}

	// Write more lines to force scroll
	term.Write([]byte("Line 6\n"))
	term.Write([]byte("Line 7\n"))

	// History should have 3 lines because filling the screen with 5 lines + newlines triggers one scroll, plus 2 more lines.
	// Actually:
	// 1..4: Y=1..4.
	// 5: Y=5 -> Scroll. Hist=1.
	// 6: Y=5 -> Scroll. Hist=2.
	// 7: Y=5 -> Scroll. Hist=3.
	if len(term.history) != 3 {
		t.Errorf("Expected history length 3, got %d", len(term.history))
	}

	// Check content of history (first line pushed out)
	// It was "Line" + spaces
	cell := term.history[0][0]
	if cell.Char != 'L' {
		t.Errorf("Expected 'L' in history, got %c", cell.Char)
	}
}

func TestViewportAccess(t *testing.T) {
	term := NewAnsiTerminal(10, 3)
	// Screen height 3.
	// Write 1, 2, 3. Screen full.
	// Write 4. 1 goes to history. Screen: 2, 3, 4.

	term.Write([]byte("1\n"))
	term.Write([]byte("2\n"))
	term.Write([]byte("3\n"))
	term.Write([]byte("4\n"))
	// History: ["1", "2"]
	// Screen: ["3", "4", ""]

	term.ScrollViewport(1) // ViewOffset = 1.
	// Row 0: Logical -1 -> Hist[last] -> "2"
	// Row 1: Logical 0 -> Cells[0] -> "3"

	c0 := term.GetCell(0, 0)
	if c0.Char != '2' {
		t.Errorf("Row 0: Expected '2', got %c", c0.Char)
	}

	c1 := term.GetCell(0, 1)
	if c1.Char != '3' {
		t.Errorf("Row 1: Expected '3', got %c", c1.Char)
	}
}

func TestScrollLimits(t *testing.T) {
	term := NewAnsiTerminal(10, 5)
	term.Write([]byte("1\n2\n3\n4\n5\n6\n"))
	// History should have 2 lines ("1", "2")

	term.ScrollViewport(100)
	if term.viewOffset != 2 {
		t.Errorf("Expected viewOffset capped at 2, got %d", term.viewOffset)
	}

	term.ResetViewport()
	if term.viewOffset != 0 {
		t.Errorf("Expected viewOffset 0, got %d", term.viewOffset)
	}
}
