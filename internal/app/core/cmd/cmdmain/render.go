package cmdmain

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss/table"

	"github.com/FaisalBudiono/comhel/internal/app/core/cmd/cmdsaver"
	"github.com/FaisalBudiono/comhel/internal/app/core/util/log"
	"github.com/FaisalBudiono/comhel/internal/app/core/util/styleutil"
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
	switch msg.(type) {
	case subModelQuitConfirmed:
		log.Logger().Debug("update: sub model quit confirm")
		m.subModel = nil

		return m, nil
	}

	if m.subModel != nil {
		log.Logger().Debug("update: enter sub model update")

		var cmd tea.Cmd
		m.subModel, cmd = m.subModel.Update(msg)

		var spinnerCMD tea.Cmd
		m.spinner, spinnerCMD = m.spinner.Update(msg)

		return m, tea.Batch(cmd, spinnerCMD)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		log.Logger().Debug("update: keypress", slog.String("key", msg.String()))

		switch msg.String() {
		case "@":
			m.subModel = cmdsaver.New(m.subModelQuitBroadcast)

			return m, waitSubModelQuit(m.subModelQuitBroadcast)
		case "q", "ctrl+c":
			return m, tea.Quit
		case "g":
			m.keyQueues = append(m.keyQueues, "g")
			keys := strings.Join(m.keyQueues, "")

			log.Logger().Debug("update: g", slog.String("key", msg.String()))

			if keys == "gg" {
				m.cursor = indexStateAll
				m.keyQueues = []string{}

				return m, nil
			}

			m.cancelQueueReset()

			ctx, cancel := context.WithCancel(m.ctx)
			m.cancelQueueReset = cancel

			return m, resetKeyQueues(ctx)
		case "home":
			m.cursor = indexStateAll

			return m, nil
		case "end", "G":
			m.cursor = len(m.services) - 1

			return m, nil
		case "k", "up":
			if m.cursor > indexStateAll {
				m.cursor--
			}

			return m, nil
		case "j", "down":
			maxServiceCursor := len(m.services) - 1
			if m.cursor < maxServiceCursor {
				m.cursor++
			}

			return m, nil
		case " ":
			m := m.toggleMark()

			log.Logger().Debug(
				"update: change active state",
				slog.String("active-states", fmt.Sprintf("%#v", m.activeStates)),
			)

			return m, nil
		case "u":
			markedServices := m.markedServices()

			log.Logger().Debug(
				"update: up marked services",
				slog.String("services", fmt.Sprintf("%#v", markedServices)),
			)

			return m, tea.Batch(
				composeUpMarked(m.ctx, markedServices, m.serviceBroadcast),
				refetchMarked(m.serviceBroadcast),
			)
		case "d":
			markedServices := m.markedServices()

			log.Logger().Debug(
				"update: up marked services",
				slog.String("services", fmt.Sprintf("%#v", markedServices)),
			)

			return m, tea.Batch(
				composeDownMarked(m.ctx, markedServices, m.serviceBroadcast),
				refetchMarked(m.serviceBroadcast),
			)
		case "U":
			return m, tea.Batch(
				composeUp(m.ctx, m.reloadBroadcast),
				refetchAll(m.reloadBroadcast),
			)
		case "D":
			return m, tea.Batch(
				composeDown(m.ctx, m.reloadBroadcast),
				refetchAll(m.reloadBroadcast),
			)
		case "R":
			return m, fetchList(m.ctx)
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)

		return m, cmd
	case composeAllSent:
		log.Logger().Debug("update: compose finish")

		m.states = make(map[string]renderableService)

		return m, nil

	case composeMarkedSent:
		for _, name := range msg {
			delete(m.states, name)
		}

		return m, nil
	case refetchedAllCalled:
		log.Logger().Debug("update: refetched called")

		return m, fetchList(m.ctx)
	case refetchedMarkedCalled:
		log.Logger().Debug("update: refetched marked only")

		cmds := make([]tea.Cmd, len(msg))
		for i, sn := range msg {
			cmds[i] = fetchService(m.ctx, sn)
		}

		return m, tea.Batch(cmds...)
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
	if m.subModel != nil {
		log.Logger().Debug("view: enter sub model view")

		return m.subModel.View()
	}

	var s string

	s += m.renderTable()
	s += "\n\n"
	s += m.helperText()

	return s
}

func (m model) helperText() string {
	var s string
	render := styleutil.Helper().Render

	s += render("k/↑: cursor up | j/↓: cursor down | home/gg: Go top | end/G: Go bottom")
	s += "\n"
	s += render("u: Up marked | d: Down marked | space: mark")
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

		rows[i] = []string{m.renderCursor(i), no, name, m.renderStatus(name)}
	}

	t := table.New().
		StyleFunc(m.tableStyling).
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(styleutil.ColorDarkPurple)).
		Headers(m.headerMark(), "NO", "SERVICE", "STATUS").
		Rows(rows...)

	s += t.Render()

	return s
}

func (m model) headerMark() string {
	if m.allMark() {
		return "[x]"
	}
	return "[ ]"
}

func (m model) tableStyling(row, col int) lipgloss.Style {
	isNumberCol := col == 1

	switch row {
	case table.HeaderRow:
		if m.cursor == table.HeaderRow {
			if col == 0 {
				return styleutil.ActiveHeaderMarker()
			}

			return styleutil.ActiveHeader()
		}

		if col == 0 {
			return styleutil.Cell()
		}

		return styleutil.Header()
	case m.cursor:
		if isNumberCol {
			return styleutil.NumberActive()
		}
		return styleutil.Active()
	default:
		if isNumberCol {
			return styleutil.NumberCell()
		}
		return styleutil.Cell()
	}
}

func (m model) renderStatus(serviceName string) string {
	s, found := m.states[serviceName]
	if !found {
		return m.spinner.View()
	}

	return s.status
}

func (m model) renderCursor(i int) string {
	isActive := m.activeStates[i]
	if isActive {
		return "[x]"
	}

	return "[ ]"
}

func (m model) toggleMark() model {
	isTriggeredAll := m.cursor == indexStateAll
	if !isTriggeredAll {
		activeState := m.activeStates[m.cursor]
		m.activeStates[m.cursor] = !activeState

		log.Logger().Debug(
			"toggleMark: triggered one by one",
			slog.String("activeStates", fmt.Sprintf("%#v", m.activeStates)),
		)

		return m
	}

	firstMarkState := m.activeStates[0]

	for i := range m.services {
		m.activeStates[i] = !firstMarkState
	}

	log.Logger().Debug(
		"toggleMark: trigger all",
		slog.Bool("firstMark", firstMarkState),
		slog.String("services", fmt.Sprintf("%#v", m.services)),
		slog.String("activeStates", fmt.Sprintf("%#v", m.activeStates)),
	)

	return m
}

func (m model) allMark() bool {
	for i := range m.services {
		mark := m.activeStates[i]
		if !mark {
			return false
		}
	}

	return true
}

func (m model) markedServices() []string {
	maxIndex := len(m.services) - 1
	markedServices := make([]string, 0)

	for i, isMarked := range m.activeStates {
		if isMarked && i <= maxIndex {
			markedServices = append(markedServices, m.services[i])
		}
	}

	log.Logger().Debug("fetching markedServices",
		slog.Int("maxIndex", maxIndex),
		slog.String("activeStates", fmt.Sprintf("%#v", m.activeStates)),
		slog.String("markedServices", fmt.Sprintf("%#v", markedServices)),
	)

	return markedServices
}
