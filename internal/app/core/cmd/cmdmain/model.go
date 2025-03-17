package cmdmain

import (
	"context"

	"github.com/charmbracelet/bubbles/spinner"
)

type model struct {
	ctx context.Context

	spinner spinner.Model

	services     []string
	states       map[string]renderableService
	cursor       int
	activeStates map[int]bool

	reloadBroadcast  chan struct{}
	serviceBroadcast chan []string
}

func New() model {
	spn := spinner.New()
	spn.Spinner = spinner.Points

	return model{
		ctx:     context.Background(),
		spinner: spn,

		states:       make(map[string]renderableService),
		activeStates: make(map[int]bool),

		reloadBroadcast:  make(chan struct{}),
		serviceBroadcast: make(chan []string),
	}
}
