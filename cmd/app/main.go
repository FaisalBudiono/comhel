package main

import (
	"github.com/FaisalBudiono/comhel/internal/app/adapter/log"
	"github.com/FaisalBudiono/comhel/internal/app/core/cmd/cmdmain"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	l, err := log.NewLogger()
	if err != nil {
		panic(err)
	}
	log.SetDefault(l)

	p := tea.NewProgram(cmdmain.New())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
