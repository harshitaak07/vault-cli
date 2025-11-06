package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	borderBlue = lipgloss.Color("#3b82f6")
	bgDark     = lipgloss.Color("#0d1117")
	textBright = lipgloss.Color("#e2e8f0")

	leftBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderBlue).
		Padding(1, 2).
		Width(45).
		Align(lipgloss.Left).
		Background(bgDark).
		Foreground(textBright)

	rightBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderBlue).
			Padding(1, 2).
			Width(55).
			Align(lipgloss.Left).
			Background(bgDark).
			Foreground(textBright)
)

type model struct {
	progress int
	message  string
	done     bool
}

func (m model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*80, func(t time.Time) tea.Msg {
		return "tick"
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg {
	case "tick":
		if m.progress < 100 {
			m.progress += 2
			m.message = fmt.Sprintf("Uploading... %d%%", m.progress)
			return m, tick()
		}
		m.done = true
		m.message = "Upload complete! File successfully encrypted & stored."
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	left := leftBox.Render(fmt.Sprintf("Progress\n\n%s", m.message))

	var right string
	if m.done {
		right = rightBox.Render(fmt.Sprintf("Results\n\nFile: sample.txt\nHash: SHA256:xxxxxx\nStorage: S3 Bucket\nTime: %s",
			time.Now().Format("15:04:05")))
	} else {
		right = rightBox.Render("Results\n\nWaiting for upload to complete...")
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func RunTUI() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running TUI:", err)
	}
}

func ShowVaultBanner() {
	fmt.Println(lipgloss.NewStyle().
		Foreground(borderBlue).
		Bold(true).
		Render(`
──────────────────────────────────────────────
  Vault CLI — Secure Encrypted File Manager
──────────────────────────────────────────────`))
	time.Sleep(800 * time.Millisecond)
}
