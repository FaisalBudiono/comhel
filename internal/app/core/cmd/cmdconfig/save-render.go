package cmdconfig

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"

	"github.com/FaisalBudiono/comhel/internal/app/core/util/log"
	"github.com/FaisalBudiono/comhel/internal/app/core/util/log/logattr"
	"github.com/FaisalBudiono/comhel/internal/app/core/util/styleutil"
	"github.com/FaisalBudiono/comhel/internal/app/domain"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func (m modelSaver) Init() tea.Cmd {
	return fetchConfigs(
		m.ctx,
		m.configFetcherBroadcast,
		m.configFetcherErrorBroadcast,
	)
}

func (m modelSaver) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	l := log.Logger().With(logattr.Caller("cmdsaver: update"))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		l.Debug("saver update: keypress", slog.String("key", msg.String()))

		switch msg.String() {
		case key1, key2, key3, key4, key5, key6, key7, key8:
			l.Debug("key press to save", slog.String("key", msg.String()))

			return m, saveConfig(m.ctx, msg.String(), m.markedServices)
		case "esc":
			l.Debug("saver update: escaped")
			return m, quit(m.quitBroadcast)
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case configSaved:
		l.Debug("config saved received")

		return m, fetchConfigs(
			m.ctx,
			m.configFetcherBroadcast,
			m.configFetcherErrorBroadcast,
		)
	case fetchConfigSent:
		l.Debug("fetch config sent")

		return m, tea.Batch(
			listenConfigsReceiver(m.configFetcherBroadcast),
			listenError(m.configFetcherErrorBroadcast),
		)
	case errorReceived:
		l.Debug("error received")
		m.err = msg

		return m, nil
	case configsReceived:
		l.Debug("config received")
		m.configs = m.cleanPresets(msg)

		return m, nil
	}

	return m, nil
}

func (m modelSaver) View() string {
	var s string

	s += styleutil.Title().Render("Preset Saver")
	s += "\n\n"

	s += m.renderError()

	s += m.renderTable()
	s += "\n\n"
	s += m.helperText()

	return s
}

func (m modelSaver) renderError() string {
	if m.err == nil {
		return ""
	}

	s := fmt.Sprintf("Failed to parse .comhelconfig.json:\nReason: %s",
		m.err.Error(),
	)

	return styleutil.Error().Render(s) + "\n\n"
}

func (m modelSaver) renderTable() string {
	l := log.Logger().With(logattr.Caller("cmdsaver: renderTable"))

	l.Debug("render: table", logattr.Any("configs", m.configs))

	if m.configs == nil {
		return "Loading"
	}

	t := table.New().
		StyleFunc(m.tableStyling).
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(styleutil.ColorDarkPurple)).
		Headers("KEY", "NO", "SERVICES").
		Rows(m.renderTableRows()...)

	return t.Render()
}

func (m modelSaver) renderTableRows() [][]string {
	mapConfigs := make(map[string]configPreset, len(m.configs))
	for _, c := range m.configs {
		mapConfigs[c.key] = c
	}

	rows := make([][]string, len(validKeys))
	for i, key := range validKeys {
		no := strconv.FormatInt(int64(i+1), 10)
		services := mapConfigs[key].services

		rows[i] = []string{key, no, m.formatServices(services)}
	}

	return rows
}

func (m modelSaver) tableStyling(row, col int) lipgloss.Style {
	switch row {
	case table.HeaderRow:
		return styleutil.Header().Align(lipgloss.Center)
	default:
		return styleutil.Cell().Align(lipgloss.Center)
	}
}

func (m modelSaver) formatServices(services []string) string {
	if len(services) == 0 {
		return styleutil.Disable().Render("<none>")
	}

	stl := lipgloss.NewStyle().Bold(true).Render

	formattedServices := make([]string, len(services))
	for i, s := range services {
		formattedServices[i] = styleutil.Active().Render(
			fmt.Sprintf("[%d]%s", i+1, s),
		)
	}

	return fmt.Sprintf(
		"%s%s%s",
		stl("[ "),
		strings.Join(formattedServices, ", "),
		stl(" ]"),
	)
}

func (m modelSaver) helperText() string {
	mapKey := func(keys []string) []domain.Keymap {
		if len(keys) == 0 {
			return nil
		}

		res := make([]domain.Keymap, len(keys))
		for i, key := range keys {
			res[i] = domain.NewKeymap([]string{key}, fmt.Sprintf("Save to %s", key))
		}

		return res
	}

	helpGroups := [][]domain.Keymap{
		{
			{Keys: []string{"q", "ctrl+c"}, Description: "quit"},
			{Keys: []string{"<esc>"}, Description: "go back"},
		},
		mapKey([]string{key1, key2, key3, key4}),
		mapKey([]string{key5, key6, key7, key8}),
	}

	return styleutil.RenderHelper(helpGroups)
}

func (m modelSaver) cleanPresets(doms []domain.ConfigPreset) []configPreset {
	l := log.Logger().With(logattr.Caller("cmdsaver: cleanPresets"))

	if doms == nil {
		return nil
	}

	if len(doms) == 0 {
		return []configPreset{}
	}

	res := make([]configPreset, 0)

	for _, d := range doms {
		if !slices.Contains(validKeys, d.Key) {
			l.Debug("skipping invalid key", slog.String("key", d.Key))

			continue
		}

		validServices := make([]string, 0)
		for _, s := range d.Services {
			l.Debug("checking service validity")

			if slices.Contains(m.validServices, s) {
				l.Debug("add service", slog.String("service", s))
				validServices = append(validServices, s)
			}
		}

		res = append(res, configPreset{
			key:      d.Key,
			services: validServices,
		})
	}

	return res
}
