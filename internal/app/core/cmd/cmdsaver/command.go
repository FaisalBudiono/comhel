package cmdsaver

import (
	"context"

	"github.com/FaisalBudiono/comhel/internal/app/domain"
	tea "github.com/charmbracelet/bubbletea"
)

func quit(b chan<- struct{}) tea.Cmd {
	return func() tea.Msg {
		b <- struct{}{}
		return nil
	}
}

type fetchConfigSent struct{}

func fetchConfigs(ctx context.Context, b chan<- []domain.ConfigPreset) tea.Cmd {
	return func() tea.Msg {
		go func() {
			res, err := configRepo.Fetch(ctx)
			if err != nil {
				panic(err)
			}

			b <- res
		}()

		return fetchConfigSent(struct{}{})
	}
}

type configsReceived []domain.ConfigPreset

func receiveConfigs(b <-chan []domain.ConfigPreset) tea.Cmd {
	return func() tea.Msg {
		res := <-b

		return configsReceived(res)
	}
}
