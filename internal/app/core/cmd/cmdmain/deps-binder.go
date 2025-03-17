package cmdmain

import "github.com/FaisalBudiono/comhel/internal/app/port/portout"

var composeRepo portout.DockerComposePort

func BindDeps(dockerComposePort portout.DockerComposePort) {
	composeRepo = dockerComposePort
}
