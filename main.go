package main

import (
	"github.com/FaisalBudiono/comhel/internal/app/adapter/doccom"
	"github.com/FaisalBudiono/comhel/internal/app/adapter/env"
	"github.com/FaisalBudiono/comhel/internal/app/adapter/log"
	"github.com/FaisalBudiono/comhel/internal/app/core/cmd/cmdmain"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	err := env.Bind()
	if err != nil {
		panic(err)
	}

	l, err := log.New()
	if err != nil {
		panic(err)
	}
	log.SetDefault(l)

	cmdmain.BindDeps(doccom.New())

	p := tea.NewProgram(cmdmain.New())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
