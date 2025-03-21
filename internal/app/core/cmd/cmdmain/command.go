package cmdmain

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/FaisalBudiono/comhel/internal/app/core/util/log"
	"github.com/FaisalBudiono/comhel/internal/app/core/util/log/logattr"
	"github.com/FaisalBudiono/comhel/internal/app/port/portout"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	fetchedListNames []string
	errorReceived    error
)

func fetchList(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		names, err := composeRepo.List(ctx)
		if err != nil {
			switch cusErr := err.(type) {
			case *portout.ConfigErr:
				return errorReceived(cusErr)
			}

			panic(err)
		}

		return fetchedListNames(names)
	}
}

type composeAllSent bool

func composeDown(
	ctx context.Context, b chan<- struct{}, bErr chan<- error,
) tea.Cmd {
	return func() tea.Msg {
		go func() {
			err := composeRepo.Down(ctx)
			b <- struct{}{}
			bErr <- err
		}()

		return composeAllSent(false)
	}
}

func composeUp(
	ctx context.Context, b chan<- struct{}, bErr chan<- error,
) tea.Cmd {
	return func() tea.Msg {
		go func() {
			err := composeRepo.Up(ctx)
			b <- struct{}{}
			bErr <- err
		}()

		return composeAllSent(false)
	}
}

type composeMarkedSent []string

func composeUpMarked(
	ctx context.Context, services []string, b chan<- []string, bErr chan<- error,
) tea.Cmd {
	return func() tea.Msg {
		go func() {
			err := composeRepo.UpByService(ctx, services...)

			b <- services
			bErr <- err
		}()

		return composeMarkedSent(services)
	}
}

func composeDownMarked(
	ctx context.Context, services []string, b chan<- []string, bErr chan<- error,
) tea.Cmd {
	return func() tea.Msg {
		go func() {
			err := composeRepo.DownByService(ctx, services...)

			b <- services
			bErr <- err
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

func listenError(b <-chan error) tea.Cmd {
	return func() tea.Msg {
		err := <-b

		var errMsg string
		if err != nil {
			errMsg = err.Error()
		}
		log.Logger().Debug("error received",
			logattr.Caller("cmdmain: command: listen Error"),
			slog.String("err", errMsg),
			logattr.Any("errComplete", err),
		)

		if errors.Is(err, portout.ErrNoService) {
			return nil
		}

		return errorReceived(err)
	}
}

type (
	fetchedService         renderableService
	fetchedServiceNotFound string
)

func fetchService(ctx context.Context, serviceName string) tea.Cmd {
	return func() tea.Msg {
		s, err := composeRepo.Service(ctx, serviceName)
		if err != nil {
			if errors.Is(err, portout.ErrNotFound) {
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

type loadedServiceQuitConfirmed []string

func waitModelLoaderQuit(b <-chan []string) tea.Cmd {
	return func() tea.Msg {
		res := <-b

		return loadedServiceQuitConfirmed(res)
	}
}
