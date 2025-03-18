package cmdsaver

import (
	"context"

	"github.com/FaisalBudiono/comhel/internal/app/domain"
)

type configPreset struct {
	key      string
	services []string
}

const (
	key1 string = "1"
	key2 string = "2"
	key3 string = "3"
	key4 string = "4"
	key5 string = "z"
	key6 string = "x"
	key7 string = "c"
	key8 string = "v"
)

var validKeys = []string{key1, key2, key3, key4, key5, key6, key7, key8}

type model struct {
	ctx context.Context

	markedServices []string
	validServices  []string
	configs        []configPreset

	quitBroadcast          chan<- struct{}
	configFetcherBroadcast chan []domain.ConfigPreset
}

func New(
	ctx context.Context,
	quitBroadcast chan<- struct{},
	markedServices []string,
	validServices []string,
) model {
	return model{
		ctx: ctx,

		markedServices: markedServices,
		validServices:  validServices,

		quitBroadcast:          quitBroadcast,
		configFetcherBroadcast: make(chan []domain.ConfigPreset),
	}
}
