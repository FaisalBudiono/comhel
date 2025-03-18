package main

import (
	"github.com/FaisalBudiono/comhel/internal/app/adapter/doccom"
	"github.com/FaisalBudiono/comhel/internal/app/adapter/jsonconfig"
	logadapter "github.com/FaisalBudiono/comhel/internal/app/adapter/log"
	"github.com/FaisalBudiono/comhel/internal/app/core/cmd/cmdmain"
	"github.com/FaisalBudiono/comhel/internal/app/core/cmd/cmdsaver"
	"github.com/FaisalBudiono/comhel/internal/app/core/util/env"
	"github.com/FaisalBudiono/comhel/internal/app/core/util/log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	err := env.Bind()
	if err != nil {
		panic(err)
	}

	l, err := logadapter.New()
	if err != nil {
		panic(err)
	}
	log.SetDefault(l)

	cmdsaver.BindDeps(jsonconfig.New())
	cmdmain.BindDeps(doccom.New())

	p := tea.NewProgram(cmdmain.New())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
