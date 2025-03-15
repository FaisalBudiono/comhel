package cmdmain

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/FaisalBudiono/comhel/internal/app/adapter/log"
	"github.com/FaisalBudiono/comhel/internal/app/core/compose"
	tea "github.com/charmbracelet/bubbletea"
)

type fetchedListNames []string

func fetchList(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		names, err := compose.List(ctx)
		if err != nil {
			panic(err)
		}

		return fetchedListNames(names)
	}
}

type composeFinished bool

func composeDown(ctx context.Context, b chan<- struct{}) tea.Cmd {
	return func() tea.Msg {
		go func() {
			err := compose.Down(ctx)
			if err != nil {
				log.Logger().Error("failed compose down", slog.String("err", fmt.Sprintf("%#v", err)))
				panic(err)
			}

			b <- struct{}{}
		}()

		return composeFinished(false)
	}
}

func composeUp(ctx context.Context, b chan<- struct{}) tea.Cmd {
	return func() tea.Msg {
		go func() {
			err := compose.Up(ctx)
			if err != nil {
				log.Logger().Error("failed compose up", slog.String("err", fmt.Sprintf("%#v", err)))
				panic(err)
			}

			b <- struct{}{}
		}()

		return composeFinished(false)
	}
}

type refetchedCalled bool

func refetchListener(b <-chan struct{}) tea.Cmd {
	return func() tea.Msg {
		log.Logger().Debug("cmd: START refetch listener")
		<-b

		log.Logger().Debug("cmd: END refetch listener")
		return refetchedCalled(false)
	}
}

func waitThen(d time.Duration, cmd tea.Cmd) tea.Cmd {
	t := time.After(d)
	<-t

	return cmd
}

type (
	fetchedService         renderableService
	fetchedServiceNotFound string
)

func fetchService(ctx context.Context, serviceName string) tea.Cmd {
	return func() tea.Msg {
		s, err := compose.Service(ctx, serviceName)
		if err != nil {
			if errors.Is(err, compose.ErrNotFound) {
				return fetchedServiceNotFound(serviceName)
			}

			log.Logger().Error("failed fetching service", slog.String("err", fmt.Sprintf("%#v", err)))
			panic(err)
		}

		return fetchedService(fromDomain(s))
	}
}
