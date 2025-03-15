package cmdmain

import (
	"context"

	"github.com/charmbracelet/bubbles/spinner"
)

type model struct {
	ctx context.Context

	spinner spinner.Model

	clNo      int
	clService int
	clStatus  int

	services []string
	states   map[string]renderableService

	reloadBroadcast chan struct{}
}

func New() model {
	spn := spinner.New()
	spn.Spinner = spinner.Points

	return model{
		ctx:     context.Background(),
		spinner: spn,

		clNo:      5,
		clService: 8,
		clStatus:  14,

		states: make(map[string]renderableService),

		reloadBroadcast: make(chan struct{}),
	}
}
