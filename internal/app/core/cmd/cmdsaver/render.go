package cmdsaver

import (
	"log/slog"

	"github.com/FaisalBudiono/comhel/internal/app/adapter/log"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		log.Logger().Debug("saver update: keypress", slog.String("key", msg.String()))

		switch msg.String() {
		case "esc":
			log.Logger().Debug("saver update: escaped")
			return m, quit(m.quitBroadcast)
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	return "hello"
}
