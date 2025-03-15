package cmdmain

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/FaisalBudiono/comhel/internal/app/adapter/log"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
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

	s += m.actionCommandText()

	return s
}

func (m model) actionCommandText() string {
	return "q: quit | R: refresh"
}

func (m model) renderTable() string {
	var s string
	if m.services == nil {
		return fmt.Sprintf("%s Loading", m.spinner.View())
	}

	serviceFmt := fmt.Sprintf("| %%%ds | %%%ds | %%%ds |\n",
		m.clNo,
		m.clService,
		m.clStatus,
	)

	s += m.dash()
	s += fmt.Sprintf(serviceFmt, "No", "Service", "Status")
	s += m.dash()

	for i, name := range m.services {
		no := strconv.FormatInt(int64(i+1), 10)
		s += fmt.Sprintf(serviceFmt, no, name, m.renderStatus(name))
	}

	s += m.dash()

	return s
}

func (m model) renderStatus(serviceName string) string {
	s, found := m.states[serviceName]
	if !found {
		return m.spinner.View()
	}

	return s.status
}

func (m model) dash() string {
	total := m.clNo + m.clService + m.clStatus + 10

	return strings.Repeat("-", total) + "\n"
}
