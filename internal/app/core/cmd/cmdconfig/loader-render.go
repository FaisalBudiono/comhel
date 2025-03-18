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

func (m modelLoader) Init() tea.Cmd {
	return fetchConfigs(
		m.ctx,
		m.configFetcherBroadcast,
		m.configFetcherErrorBroadcast,
	)
}

func (m modelLoader) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	l := log.Logger().With(logattr.Caller("cmdloader: update"))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		l.Debug("loader update: keypress", slog.String("key", msg.String()))

		switch msg.String() {
		case key1, key2, key3, key4, key5, key6, key7, key8:
			key := msg.String()

			l.Debug("key press to load", slog.String("key", key))

			return m, quitAndLoad(m.quitAndLoadBroadcast, m.serviceByKey(key))
		case "esc":
			l.Debug("loader update: escaped")
			return m, quitAndLoad(m.quitAndLoadBroadcast, nil)
		case "q", "ctrl+c":
			return m, tea.Quit
		}
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

func (m modelLoader) View() string {
	var s string

	s += styleutil.Title().Render("Preset Loader")
	s += "\n\n"

	s += renderError(m.err)

	s += renderTable(m.configs)
	s += "\n\n"
	s += m.helperText()

	return s
}

func (m modelLoader) serviceByKey(key string) []string {
	for _, s := range m.configs {
		if s.key == key {
			return s.services
		}
	}

	return nil
}

func (m modelLoader) helperText() string {
	mapKey := func(keys []string) []domain.Keymap {
		if len(keys) == 0 {
			return nil
		}

		res := make([]domain.Keymap, len(keys))
		for i, key := range keys {
			res[i] = domain.NewKeymap([]string{key}, fmt.Sprintf("Load from %s", key))
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
