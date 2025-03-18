package cmdconfig

import (
	"context"

	"github.com/FaisalBudiono/comhel/internal/app/domain"
)

type modelLoader struct {
	ctx context.Context

	markedServices []string
	validServices  []string
	configs        []configPreset

	err error

	quitAndLoadBroadcast        chan<- []string
	configFetcherBroadcast      chan []domain.ConfigPreset
	configFetcherErrorBroadcast chan error
}

func NewLoader(
	ctx context.Context,
	quitAndLoadBoradcast chan<- []string,
	markedServices []string,
	validServices []string,
) modelLoader {
	return modelLoader{
		ctx: ctx,

		markedServices: markedServices,
		validServices:  validServices,

		quitAndLoadBroadcast:        quitAndLoadBoradcast,
		configFetcherBroadcast:      make(chan []domain.ConfigPreset),
		configFetcherErrorBroadcast: make(chan error),
	}
}
