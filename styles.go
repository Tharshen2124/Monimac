package main

import "github.com/charmbracelet/lipgloss"

var (
	// titleStyle is the top header bar.
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#5C5FE0")).
			Padding(0, 2).
			Width(0) // width set dynamically

	// sectionTitleStyle labels each section.
	sectionTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#9CB5FE"))

	// selectedRowStyle highlights the currently selected container row.
	selectedRowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#1A1A2E")).
				Background(lipgloss.Color("#9CB5FE")).
				Bold(true)

	// errorStyle renders errors in a warm red.
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F5F")).
			Bold(true)

	// dimStyle renders muted / empty-state text.
	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	// confirmBoxStyle wraps the stop-confirmation dialog.
	confirmBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#9CB5FE")).
			Padding(1, 3).
			Foreground(lipgloss.Color("#FAFAFA"))

	// barFilledStyle colours the filled portion of a progress bar.
	barFilledStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5C5FE0"))

	// barEmptyStyle colours the empty portion of a progress bar.
	barEmptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3A3A4A"))

	// footerStyle renders the key-binding hint bar at the bottom.
	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)

	// spinnerStyle applies to the loading spinner.
	spinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CB5FE"))
)
