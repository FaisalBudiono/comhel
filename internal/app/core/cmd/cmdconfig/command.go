package cmdconfig

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

func quitAndLoad(b chan<- []string, loadServices []string) tea.Cmd {
	return func() tea.Msg {
		b <- loadServices

		return nil
	}
}

type fetchConfigSent struct{}

func fetchConfigs(
	ctx context.Context,
	cRes chan<- []domain.ConfigPreset,
	cErr chan<- error,
) tea.Cmd {
	return func() tea.Msg {
		go func() {
			res, err := configRepo.Fetch(ctx)
			cRes <- res
			cErr <- err
		}()

		return fetchConfigSent(struct{}{})
	}
}

type configsReceived []domain.ConfigPreset

func listenConfigsReceiver(b <-chan []domain.ConfigPreset) tea.Cmd {
	return func() tea.Msg {
		res := <-b

		log.Logger().Debug("configs received",
			logattr.Caller("cmdconfig: command: listen configs"),
			logattr.Any("res", res),
		)

		return configsReceived(res)
	}
}

type errorReceived error

func listenError(b <-chan error) tea.Cmd {
	return func() tea.Msg {
		err := <-b

		var errMsg string
		if err != nil {
			errMsg = err.Error()
		}
		log.Logger().Debug("error received",
			logattr.Caller("cmdconfig: command: listen Error"),
			slog.String("err", errMsg),
			logattr.Any("errComplete", err),
		)

		return errorReceived(err)
	}
}

type configSaved struct{}

func saveConfig(ctx context.Context, key string, services []string) tea.Cmd {
	return func() tea.Msg {
		l := log.Logger().With(logattr.Caller("cmdconfig: command: saveConfig"))

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
			return errorReceived(err)
		}

		return configSaved(struct{}{})
	}
}
