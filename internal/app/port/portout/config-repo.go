package portout

import (
	"context"
	"fmt"

	"github.com/FaisalBudiono/comhel/internal/app/domain"
)

type ConfigErr struct {
	CmdName string
	Msg     string
}

func (c *ConfigErr) Error() string {
	return fmt.Sprintf("%s: %s", c.CmdName, c.Msg)
}

type ConfigRepo interface {
	Save(ctx context.Context, p domain.ConfigPreset) (domain.ConfigPreset, error)
	Fetch(ctx context.Context) ([]domain.ConfigPreset, error)
}
