package tea

import (
	"github.com/charmbracelet/lipgloss"
)

type eventStyle struct {
	Screen lipgloss.Style

	Title       lipgloss.Style
	List        lipgloss.Style
	Description lipgloss.Style
	Pagination  lipgloss.Style
	Help        lipgloss.Style

	NormalItem   lipgloss.Style
	CurrentItem  lipgloss.Style
	SelectedItem lipgloss.Style

	ActivePaginationDot   lipgloss.Style
	InactivePaginationDot lipgloss.Style
}

func newEventStyle() (s eventStyle) {
	s.Screen = lipgloss.NewStyle().
		Padding(1, 2)

	s.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("#9254de")).
		Foreground(lipgloss.Color("#fafafa")).
		Padding(0, 1)

	s.List = lipgloss.NewStyle().
		MarginTop(1).
		Border(lipgloss.RoundedBorder(), true, false)

	s.Description = lipgloss.NewStyle()

	s.Pagination = lipgloss.NewStyle().
		PaddingTop(1)

	s.Help = lipgloss.NewStyle().
		PaddingTop(1)

	s.NormalItem = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#d9d9d9"))

	s.CurrentItem = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#69c0ff"))

	s.SelectedItem = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#95de64"))

	s.ActivePaginationDot = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f5f5f5")).
		SetString("•")

	s.InactivePaginationDot = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#595959")).
		SetString("•")

	return
}
