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

type composeAllSent bool

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

		return composeAllSent(false)
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

		return composeAllSent(false)
	}
}

type composeMarkedSent []string

func composeUpMarked(
	ctx context.Context, services []string, b chan<- []string,
) tea.Cmd {
	return func() tea.Msg {
		go func() {
			err := compose.UpByService(ctx, services...)
			if err != nil {
				if !errors.Is(err, compose.ErrNoService) {
					log.Logger().Error("failed compose up manually", slog.String("err", fmt.Sprintf("%#v", err)))
					panic(err)
				}
			}

			b <- services
		}()

		return composeMarkedSent(services)
	}
}

func composeDownMarked(
	ctx context.Context, services []string, b chan<- []string,
) tea.Cmd {
	return func() tea.Msg {
		go func() {
			err := compose.DownByService(ctx, services...)
			if err != nil {
				if !errors.Is(err, compose.ErrNoService) {
					log.Logger().Error("failed compose down manually", slog.String("err", fmt.Sprintf("%#v", err)))
					panic(err)
				}
			}

			b <- services
		}()

		return composeMarkedSent(services)
	}
}

type refetchedAllCalled bool

func refetchAll(b <-chan struct{}) tea.Cmd {
	return func() tea.Msg {
		log.Logger().Debug("cmd: START refetch all")
		<-b

		log.Logger().Debug("cmd: END refetch all")
		return refetchedAllCalled(false)
	}
}

type refetchedMarkedCalled []string

func refetchMarked(b <-chan []string) tea.Cmd {
	return func() tea.Msg {
		log.Logger().Debug("cmd: START refetch marked")
		res := <-b

		log.Logger().Debug("cmd: END refetch marked")

		return refetchedMarkedCalled(res)
	}
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

type queueReset struct{}

func resetKeyQueues(ctx context.Context) tea.Cmd {
	log.Logger().Debug("resetKeyQueue: called")

	return func() tea.Msg {
		select {
		case <-time.After(time.Second):
			log.Logger().Debug("resetKeyQueue: finish")
			return queueReset(struct{}{})
		case <-ctx.Done():
			log.Logger().Debug("resetKeyQueue: canceled")
			return nil
		}
	}
}

type subModelQuitConfirmed struct{}

func waitSubModelQuit(b <-chan struct{}) tea.Cmd {
	return func() tea.Msg {
		log.Logger().Debug("waitSubModelQuit: START")
		<-b
		log.Logger().Debug("waitSubModelQuit: END")

		return subModelQuitConfirmed(struct{}{})
	}
}
