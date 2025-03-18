package cmdsaver

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/FaisalBudiono/comhel/internal/app/core/util/log"
	"github.com/FaisalBudiono/comhel/internal/app/core/util/log/logattr"
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

type configSaved struct{}

func saveConfig(ctx context.Context, key string, services []string) tea.Cmd {
	return func() tea.Msg {
		l := log.Logger().With(logattr.Caller("cmdsaver: command: saveConfig"))

		l.DebugContext(ctx, "key press",
			slog.String("key", key),
			slog.String("services", fmt.Sprintf("%#v", services)),
		)

		cp := domain.NewConfigPreset(
			key,
			services,
		)
		_, err := configRepo.Save(ctx, cp)
		if err != nil {
			panic(err)
		}

		return configSaved(struct{}{})
	}
}
