package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"tharshen.xyz/monimac/internal/docker"
	"tharshen.xyz/monimac/internal/metrics"
)

type state int

const (
	stateLoading state = iota
	stateNormal
	stateConfirm
)

// systemMetricsMsg carries the result of a metrics fetch.
type systemMetricsMsg struct {
	cpu metrics.CPUStats
	mem metrics.MemStats
	err error
}

// dockerMsg carries the result of a docker container listing.
type dockerMsg struct {
	containers []docker.Container
	err        error
}

// containerStoppedMsg carries the result of a stop command.
type containerStoppedMsg struct {
	id  string
	err error
}

// tickMsg triggers a periodic refresh.
type tickMsg time.Time

// Model is the root BubbleTea model.
type Model struct {
	state      state
	cpu        metrics.CPUStats
	mem        metrics.MemStats
	metricsErr error
	containers []docker.Container
	dockerErr  error
	selected   int
	stopTarget docker.Container
	width      int
	height     int
	spinner    spinner.Model
}

// NewModel constructs the initial model.
func NewModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle
	return Model{
		state:   stateLoading,
		spinner: s,
	}
}

func tick() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func fetchSystemMetrics() tea.Cmd {
	return func() tea.Msg {
		cpuStats, err := metrics.GetCPUStats()
		if err != nil {
			return systemMetricsMsg{err: err}
		}
		memStats, err := metrics.GetMemStats()
		if err != nil {
			return systemMetricsMsg{cpu: cpuStats, err: err}
		}
		return systemMetricsMsg{cpu: cpuStats, mem: memStats}
	}
}

func fetchDockerContainers() tea.Cmd {
	return func() tea.Msg {
		containers, err := docker.ListContainers()
		return dockerMsg{containers: containers, err: err}
	}
}

func stopContainer(id string) tea.Cmd {
	return func() tea.Msg {
		err := docker.StopContainer(id)
		return containerStoppedMsg{id: id, err: err}
	}
}

// Init starts the spinner, fetches metrics, and schedules the first tick.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		fetchSystemMetrics(),
		fetchDockerContainers(),
		tick(),
	)
}

// Update handles all incoming messages and key events.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case spinner.TickMsg:
		if m.state == stateLoading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil

	case systemMetricsMsg:
		m.metricsErr = msg.err
		if msg.err == nil {
			m.cpu = msg.cpu
			m.mem = msg.mem
		}
		if m.state == stateLoading {
			// Stay in loading until we also have a docker response.
			// We'll transition when we get dockerMsg.
		}
		return m, nil

	case dockerMsg:
		m.dockerErr = msg.err
		if msg.err == nil {
			m.containers = msg.containers
		}
		if m.state == stateLoading {
			m.state = stateNormal
		}
		// Clamp selection.
		if m.selected >= len(m.containers) {
			m.selected = max(0, len(m.containers)-1)
		}
		return m, nil

	case containerStoppedMsg:
		// Trigger a refresh after stopping.
		return m, tea.Batch(fetchDockerContainers(), fetchSystemMetrics())

	case tickMsg:
		return m, tea.Batch(
			fetchSystemMetrics(),
			fetchDockerContainers(),
			tick(),
		)

	case tea.KeyMsg:
		switch m.state {
		case stateNormal:
			return m.handleNormalKey(msg)
		case stateConfirm:
			return m.handleConfirmKey(msg)
		}
	}

	return m, nil
}

func (m Model) handleNormalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "up", "k":
		if m.selected > 0 {
			m.selected--
		}

	case "down", "j":
		if m.selected < len(m.containers)-1 {
			m.selected++
		}

	case "enter", "s":
		if len(m.containers) > 0 {
			m.stopTarget = m.containers[m.selected]
			m.state = stateConfirm
		}

	case "r":
		return m, tea.Batch(fetchSystemMetrics(), fetchDockerContainers())
	}

	return m, nil
}

func (m Model) handleConfirmKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		id := m.stopTarget.ID
		m.state = stateNormal
		return m, stopContainer(id)

	case "n", "esc":
		m.state = stateNormal

	case "q", "ctrl+c":
		return m, tea.Quit
	}

	return m, nil
}
