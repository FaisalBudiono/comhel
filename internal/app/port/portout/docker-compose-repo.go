package portout

import (
	"context"
	"errors"

	"github.com/FaisalBudiono/comhel/internal/app/domain"
)

var (
	ErrNoService = errors.New("no service provided")
	ErrNotFound  = errors.New("service not found")
)

type DockerComposePort interface {
	List(ctx context.Context) ([]string, error)
	Service(ctx context.Context, serviceName string) (domain.Service, error)
	Down(ctx context.Context) error
	DownByService(ctx context.Context, services ...string) error
	Up(ctx context.Context) error
	UpByService(ctx context.Context, services ...string) error
}
