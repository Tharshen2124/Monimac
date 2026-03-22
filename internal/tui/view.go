package tui

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View is the root rendering function.
func (m Model) View() string {
	if m.state == stateLoading {
		return fmt.Sprintf("\n  %s Loading monimac…\n", m.spinner.View())
	}

	if m.state == stateConfirm {
		return m.confirmView()
	}

	sections := []string{
		m.titleView(),
		m.systemView(),
		m.dockerView(),
		m.footerView(),
	}
	return strings.Join(sections, "\n")
}

// titleView renders the top header bar.
func (m Model) titleView() string {
	hostname, _ := os.Hostname()
	label := fmt.Sprintf(" monimac  —  %s ", hostname)
	style := titleStyle.Width(m.width)
	return style.Render(label)
}

// systemView renders CPU and memory metrics.
func (m Model) systemView() string {
	var sb strings.Builder
	sb.WriteString(sectionTitleStyle.Render("SYSTEM METRICS"))
	sb.WriteString("\n")

	if m.metricsErr != nil {
		sb.WriteString(errorStyle.Render(fmt.Sprintf("  metrics error: %v", m.metricsErr)))
		sb.WriteString("\n")
		return sb.String()
	}

	barWidth := 20
	cpuBar := renderBar(m.cpu.Percent, barWidth)
	cpuLine := fmt.Sprintf("  CPU  [%s]  %5.1f%%", cpuBar, m.cpu.Percent)
	sb.WriteString(cpuLine)
	sb.WriteString("\n")

	memBar := renderBar(m.mem.Percent, barWidth)
	memLine := fmt.Sprintf("  MEM  [%s]  %5.1f%%  %s / %s",
		memBar, m.mem.Percent,
		bytesToGB(m.mem.Used),
		bytesToGB(m.mem.Total),
	)
	sb.WriteString(memLine)
	sb.WriteString("\n")

	return sb.String()
}

// dockerView renders the container table or an error/empty message.
func (m Model) dockerView() string {
	var sb strings.Builder
	sb.WriteString(sectionTitleStyle.Render("DOCKER CONTAINERS"))
	sb.WriteString("\n")

	if m.dockerErr != nil {
		sb.WriteString(errorStyle.Render(fmt.Sprintf("  Docker unavailable: %v", m.dockerErr)))
		sb.WriteString("\n")
		return sb.String()
	}

	if len(m.containers) == 0 {
		sb.WriteString(dimStyle.Render("  No containers running"))
		sb.WriteString("\n")
		return sb.String()
	}

	// Column widths.
	const (
		colName   = 24
		colImage  = 28
		colStatus = 20
		colCPU    = 10
		colMem    = 22
	)

	header := fmt.Sprintf("  %-*s  %-*s  %-*s  %-*s  %-*s",
		colName, "NAME",
		colImage, "IMAGE",
		colStatus, "STATUS",
		colCPU, "CPU",
		colMem, "MEM",
	)
	sb.WriteString(dimStyle.Render(header))
	sb.WriteString("\n")

	for i, c := range m.containers {
		name := truncate(c.Name, colName)
		image := truncate(c.Image, colImage)
		status := truncate(c.Status, colStatus)
		cpuPerc := truncate(c.CPU, colCPU)
		memUsage := truncate(c.Mem, colMem)

		row := fmt.Sprintf("  %-*s  %-*s  %-*s  %-*s  %-*s",
			colName, name,
			colImage, image,
			colStatus, status,
			colCPU, cpuPerc,
			colMem, memUsage,
		)

		if i == m.selected {
			sb.WriteString(selectedRowStyle.Render(row))
		} else {
			sb.WriteString(row)
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// footerView renders the key-binding hint line.
func (m Model) footerView() string {
	hint := "  q quit   ↑/↓ select   s/enter stop   r refresh"
	return footerStyle.Render(hint)
}

// confirmView renders the stop-confirmation overlay.
func (m Model) confirmView() string {
	content := fmt.Sprintf(
		"Stop container %s?\n\n  [y] yes   [n] no",
		lipgloss.NewStyle().Bold(true).Render(m.stopTarget.Name),
	)
	box := confirmBoxStyle.Render(content)

	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
	}
	return "\n" + box + "\n"
}

// renderBar returns a string of filled/empty block characters representing percent.
func renderBar(percent float64, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := int(math.Round(percent / 100.0 * float64(width)))
	if filled > width {
		filled = width
	}
	filledStr := barFilledStyle.Render(strings.Repeat("█", filled))
	emptyStr := barEmptyStyle.Render(strings.Repeat("░", width-filled))
	return filledStr + emptyStr
}

// bytesToGB formats a byte count as a human-readable GB string.
func bytesToGB(bytes uint64) string {
	gb := float64(bytes) / (1024 * 1024 * 1024)
	return fmt.Sprintf("%.1f GB", gb)
}

// truncate shortens s to at most n runes, appending "…" if truncated.
func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	if n <= 1 {
		return string(runes[:n])
	}
	return string(runes[:n-1]) + "…"
}
