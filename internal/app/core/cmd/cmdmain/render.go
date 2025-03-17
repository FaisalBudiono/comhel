package cmdmain

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/charmbracelet/lipgloss/table"

	"github.com/FaisalBudiono/comhel/internal/app/adapter/log"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		fetchList(m.ctx),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		log.Logger().Debug("update: keypress", slog.String("key", msg.String()))

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}

			return m, nil
		case "j", "down":
			maxServiceCursor := len(m.services) - 1
			if m.cursor < maxServiceCursor {
				m.cursor++
			}

			return m, nil
		case "U":
			return m, tea.Batch(
				composeUp(m.ctx, m.reloadBroadcast),
				refetchListener(m.reloadBroadcast),
			)
		case "D":
			return m, tea.Batch(
				composeDown(m.ctx, m.reloadBroadcast),
				refetchListener(m.reloadBroadcast),
			)
		case "R":
			return m, fetchList(m.ctx)
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)

		return m, cmd
	case composeFinished:
		log.Logger().Debug("update: compose finish")

		m.states = make(map[string]renderableService)

		return m, nil
	case refetchedCalled:
		log.Logger().Debug("update: refetched called")

		return m, fetchList(m.ctx)
	case fetchedListNames:
		log.Logger().Debug("update: fetchlist name", slog.String("list", fmt.Sprintf("%#v", msg)))

		m.states = make(map[string]renderableService)
		m.services = msg

		cmds := make([]tea.Cmd, len(msg))
		for i, sn := range msg {
			cmds[i] = fetchService(m.ctx, sn)
		}

		return m, tea.Batch(cmds...)
	case fetchedService:
		log.Logger().Debug("update: fetchedService", slog.String("service", fmt.Sprintf("%#v", msg)))

		m.states[msg.name] = renderableService(msg)

		return m, nil
	case fetchedServiceNotFound:
		log.Logger().Debug("update: service not found", slog.String("service", string(msg)))

		serviceName := string(msg)
		m.states[serviceName] = offService(serviceName)

		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	var s string

	s += m.renderTable()

	s += "\n\n"

	s += m.helperText()

	return s
}

func (m model) helperText() string {
	var s string
	render := helperStyle().Render

	s += render("k/↑: cursor up | j/↓: cursor down")
	s += "\n"
	s += render("q: quit | R: refresh | U: Up ALL | D: Down ALL")

	return s
}

func (m model) renderTable() string {
	var s string
	if m.services == nil {
		return fmt.Sprintf("%s Loading", m.spinner.View())
	}

	log.Logger().Debug("render: table", slog.Int("cursor", m.cursor))

	rows := make([][]string, len(m.services))
	for i, name := range m.services {
		no := strconv.FormatInt(int64(i+1), 10)

		rows[i] = []string{no, name, m.renderStatus(name)}
	}

	t := table.New().
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle()
			}

			isNumberCol := col == 0
			isActiveRow := row == m.cursor
			if isActiveRow {
				if isNumberCol {
					return noActiveStyle()
				}

				return activeStyle()
			}

			if isNumberCol {
				return noCellStyle()
			}

			return cellStyle()
		}).
		Border(lipgloss.NormalBorder()).
		Headers("No", "Service", "Status").
		Rows(rows...)

	s += t.Render()

	return s
}

func (m model) renderStatus(serviceName string) string {
	s, found := m.states[serviceName]
	if !found {
		return m.spinner.View()
	}

	return s.status
}
