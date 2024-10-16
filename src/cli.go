package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type SelectionMode int

const (
	SingleSelection SelectionMode = iota
	MultiSelection
)

type SelectorModel struct {
	cursor             int
	startIndex         int
	selected           map[int]bool
	options            []string
	selectionMode      SelectionMode
	selectionStatement string
	viewHeight         int
}

// Create a new selector.
func NewSelectorModel(options []string, mode SelectionMode,
	selectionStatement string, viewHeight int) SelectorModel {
	return SelectorModel{
		options:            options,
		selected:           make(map[int]bool),
		selectionMode:      mode,
		selectionStatement: selectionStatement,
		viewHeight:         viewHeight,
	}
}

func (m SelectorModel) Init() tea.Cmd {
	return nil
}

// Update the view.
func (m SelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//nolint: gocritic // the next directive must be in a switch.
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter": // Quit the program
			return m, tea.Quit
		case "up", "k": // Move cursor up
			moveCursorUp(&m)
		case "down", "j": // Move cursor down
			moveCursorDown(&m)
		case " ": // Toggle selection
			toggleSelection(&m)
		case "ctrl+c", "ctrl+q":
			log.Panic("Cancel tool execute")
		}
	}
	return m, nil
}

// Move cursor up, and update the view.
func moveCursorUp(m *SelectorModel) {
	if m.cursor > 0 {
		m.cursor--
	}
	if m.cursor < m.startIndex {
		m.startIndex--
	}
}

// Move cursor down, and update the view.
func moveCursorDown(m *SelectorModel) {
	if m.cursor < len(m.options)-1 {
		m.cursor++
	}
	if m.cursor >= m.startIndex+m.viewHeight {
		m.startIndex++
	}
}

// Toggle selection.
func toggleSelection(m *SelectorModel) {
	if m.selectionMode == SingleSelection {
		// Clear previous selection in single selection mode
		for k := range m.selected {
			delete(m.selected, k)
		}
		m.selected[m.cursor] = true
	} else if m.selectionMode == MultiSelection {
		// Toggle in multi-selection mode
		m.selected[m.cursor] = !m.selected[m.cursor]
	}
}

// View renders the UI for the selector.
func (m SelectorModel) View() string {
	s := "Select an option (use ↑/↓ and press SPACE to select):\n" +
		m.selectionStatement + "\n"

	if m.startIndex > 0 {
		s += "↑ More elements ↑\n"
	}
	// Render the list of options
	for i := m.startIndex; i < len(m.options) &&
		i < m.startIndex+m.viewHeight; i++ {
		cursor := " " // No cursor
		if m.cursor == i {
			cursor = ">" // Add cursor for current selection
		}

		checked := " " // No checkmark
		if m.selected[i] {
			checked = "x" // Add checkmark if selected
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, m.options[i])
	}

	if m.startIndex+m.viewHeight < len(m.options) {
		s += "↓ More elements ↓\n"
	}

	s += "\nPress ENTER to continue\n\n"
	return s
}

// FinalView outputs the selected options when quitting.
func (m SelectorModel) FinalView() (string, []string) {
	var selected []string
	for i, chosen := range m.selected {
		if chosen {
			selected = append(selected, m.options[i])
		}
	}
	if m.selectionMode == SingleSelection {
		if len(selected) == 0 {
			return "", nil
		}
		return selected[0], nil
	}
	return "", selected
}

// RunSelector starts the selector model with Bubble Tea, return two parameters
// depend of the type of selector, if is single selector return value , nil
// if is a multi selector return "", values.
func RunSelector(options []string, mode SelectionMode,
	selectionStatement string, viewHeigh int) (string, []string) {
	p := tea.NewProgram(NewSelectorModel(options, mode, selectionStatement,
		viewHeigh))

	// Start the program
	m, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Output the final view
	model, _ := m.(SelectorModel)
	return model.FinalView()
}
