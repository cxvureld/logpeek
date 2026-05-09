package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cxvureld/logpeek/internal/analyzer"
)

type tickMsg time.Time

type Model struct {
	stats  *analyzer.Stats
	table  table.Model
	ready  bool
	ticker *time.Ticker
}

func NewModel(stats *analyzer.Stats) *Model {
	columns := []table.Column{
		{Title: "Metric", Width: 35},
		{Title: "Value", Width: 35},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(false),
		table.WithHeight(20),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	t.SetStyles(s)

	return &Model{
		stats: stats,
		table: t,
	}
}

func (m *Model) Init() tea.Cmd {
	m.ticker = time.NewTicker(500 * time.Millisecond)
	return tea.Batch(tickCmd(m.ticker))
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		snap := m.stats.Snapshot()
		m.updateTable(snap)
		return m, tickCmd(m.ticker)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.ticker.Stop()
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *Model) View() string {
	if !m.ready {
		return "Initializing...\n"
	}

	snap := m.stats.Snapshot()

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		PaddingLeft(1).
		PaddingRight(1).
		Render("📊 LogPeek — Real-time Log Analyzer")

	statusBar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575")).
		Render(fmt.Sprintf(
			"Requests: %d | Bytes: %d | 5xx Errors: %d",
			snap.TotalRequests,
			snap.BytesSent,
			snap.StatusCodes[500]+snap.StatusCodes[502]+snap.StatusCodes[503],
		))

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render("q • quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		m.table.View(),
		"",
		statusBar,
		help,
	)
}

func (m *Model) updateTable(snap *analyzer.Stats) {
	var rows []table.Row

	// ------ STATUS CODES ------
	rows = append(rows, table.Row{"─── Status Codes ───", ""})
	rows = append(rows, table.Row{"  2xx ✅ Success", fmt.Sprintf("%d", snap.StatusCodes[200]+snap.StatusCodes[201]+snap.StatusCodes[204])})
	rows = append(rows, table.Row{"  3xx 🔄 Redirect", fmt.Sprintf("%d", snap.StatusCodes[301]+snap.StatusCodes[302])})
	rows = append(rows, table.Row{"  4xx ⚠️ Client Error", fmt.Sprintf("%d", snap.StatusCodes[400]+snap.StatusCodes[401]+snap.StatusCodes[403]+snap.StatusCodes[404])})
	rows = append(rows, table.Row{"  5xx 🔥 Server Error", fmt.Sprintf("%d", snap.StatusCodes[500]+snap.StatusCodes[502]+snap.StatusCodes[503])})

	// ------ TOP ENDPOINTS ------
	if len(snap.TopEndpoints) > 0 {
		rows = append(rows, table.Row{"", ""})
		rows = append(rows, table.Row{"─── Top Endpoints ───", ""})

		type epCount struct {
			path  string
			count int64
		}
		var endpoints []epCount
		for path, count := range snap.TopEndpoints {
			endpoints = append(endpoints, epCount{path, count})
		}

		for i := 0; i < len(endpoints); i++ {
			for j := i + 1; j < len(endpoints); j++ {
				if endpoints[j].count > endpoints[i].count {
					endpoints[i], endpoints[j] = endpoints[j], endpoints[i]
				}
			}
		}

		limit := 5
		if len(endpoints) < limit {
			limit = len(endpoints)
		}
		for i := 0; i < limit; i++ {
			rows = append(rows, table.Row{
				fmt.Sprintf("  %s", endpoints[i].path),
				fmt.Sprintf("%d hits", endpoints[i].count),
			})
		}
	}

	// ------ TOP IP ADDRESSES ------
	if len(snap.TopIPs) > 0 {
		rows = append(rows, table.Row{"", ""})
		rows = append(rows, table.Row{"─── Top IP Addresses ───", ""})

		type ipCount struct {
			ip    string
			count int64
		}
		var ips []ipCount
		for ip, count := range snap.TopIPs {
			ips = append(ips, ipCount{ip, count})
		}

		for i := 0; i < len(ips); i++ {
			for j := i + 1; j < len(ips); j++ {
				if ips[j].count > ips[i].count {
					ips[i], ips[j] = ips[j], ips[i]
				}
			}
		}

		limit := 5
		if len(ips) < limit {
			limit = len(ips)
		}
		for i := 0; i < limit; i++ {
			rows = append(rows, table.Row{
				fmt.Sprintf("  %s", ips[i].ip),
				fmt.Sprintf("%d requests", ips[i].count),
			})
		}
	}

	m.table.SetRows(rows)
	m.ready = true
}

func tickCmd(ticker *time.Ticker) tea.Cmd {
	return func() tea.Msg {
		return tickMsg(<-ticker.C)
	}
}
