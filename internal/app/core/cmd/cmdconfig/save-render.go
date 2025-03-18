package cmdconfig

import (
	"fmt"
	"log/slog"

	"github.com/FaisalBudiono/comhel/internal/app/core/util/log"
	"github.com/FaisalBudiono/comhel/internal/app/core/util/log/logattr"
	"github.com/FaisalBudiono/comhel/internal/app/core/util/styleutil"
	"github.com/FaisalBudiono/comhel/internal/app/domain"
	tea "github.com/charmbracelet/bubbletea"
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
		m.configs = cleanPresets(msg, m.validServices)

		return m, nil
	}

	return m, nil
}

func (m modelSaver) View() string {
	var s string

	s += styleutil.Title().Render("Preset Saver")
	s += "\n\n"

	s += renderError(m.err)

	s += renderTable(m.configs)
	s += "\n\n"
	s += m.helperText()

	return s
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
