package portout

import (
	"context"

	"github.com/FaisalBudiono/comhel/internal/app/domain"
)

type ConfigRepo interface {
	Save(ctx context.Context, p domain.ConfigPreset) (domain.ConfigPreset, error)
	Fetch(ctx context.Context) ([]domain.ConfigPreset, error)
}
