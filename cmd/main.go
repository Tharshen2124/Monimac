package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"tharshen.xyz/monimac/internal/tui"
)

func main() {
	p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
