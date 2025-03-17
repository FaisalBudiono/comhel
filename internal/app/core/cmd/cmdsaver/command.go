package cmdsaver

import tea "github.com/charmbracelet/bubbletea"

func quit(b chan<- struct{}) tea.Cmd {
	return func() tea.Msg {
		b <- struct{}{}
		return nil
	}
}
