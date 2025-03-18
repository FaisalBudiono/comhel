package cmdconfig

import "github.com/FaisalBudiono/comhel/internal/app/port/portout"

var configRepo portout.ConfigRepo

func BindDeps(configFetcher portout.ConfigRepo) {
	configRepo = configFetcher
}
