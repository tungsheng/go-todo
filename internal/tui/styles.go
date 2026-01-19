package tui

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/tungsheng/go-todo/internal/model"
)

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#424242"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	statusColors = map[model.Status]lipgloss.AdaptiveColor{
		model.StatusPending:    {Light: "#888888", Dark: "#626262"},
		model.StatusInProgress: {Light: "#F1C40F", Dark: "#F1C40F"},
		model.StatusDone:       {Light: "#2ECC71", Dark: "#2ECC71"},
		model.StatusClosed:     {Light: "#E74C3C", Dark: "#E74C3C"},
	}

	headerStyle = lipgloss.NewStyle().
			Height(3)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Padding(0, 1).
			MarginLeft(1)

	timeTagStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#3498DB")).
			Padding(0, 1).
			MarginLeft(3)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(highlight).
				Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(subtle).
			MarginTop(1).
			MarginLeft(1)

	inputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Padding(0, 1).
			MarginLeft(1)

	confirmStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E74C3C")).
			Bold(true)

	separatorStyle = lipgloss.NewStyle().
			Foreground(subtle).
			MarginTop(1).
			MarginBottom(0).
			PaddingLeft(2)
)

func StatusStyle(status model.Status) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(statusColors[status])
}
