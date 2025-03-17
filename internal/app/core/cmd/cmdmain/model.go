package cmdmain

import (
	"context"

	"github.com/charmbracelet/bubbles/spinner"
)

type model struct {
	ctx context.Context

	spinner spinner.Model

	services []string
	states   map[string]renderableService
	cursor   int

	reloadBroadcast chan struct{}
}

func New() model {
	spn := spinner.New()
	spn.Spinner = spinner.Points

	return model{
		ctx:     context.Background(),
		spinner: spn,

		states: make(map[string]renderableService),

		reloadBroadcast: make(chan struct{}),
	}
}
