package doccom

import (
	"context"
	"encoding/json"
	"os/exec"
	"slices"
	"strings"

	"github.com/FaisalBudiono/comhel/internal/app/domain"
	"github.com/FaisalBudiono/comhel/internal/app/port/portout"
)

type service struct {
	Name  string `json:"Service"`
	State string `json:"State"`
}

type dockerCompose struct{}

func (repo *dockerCompose) List(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(
		ctx,
		"docker", "compose", "config", "--services",
	)

	o, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	services := strings.Split(string(o), "\n")

	var uniqueServices []string
	for _, s := range services {
		if s != "" {
			uniqueServices = append(uniqueServices, s)
		}
	}
	slices.Sort(uniqueServices)

	return uniqueServices, nil
}

func (repo *dockerCompose) Service(ctx context.Context, serviceName string) (domain.Service, error) {
	cmd := exec.CommandContext(
		ctx,
		"docker", "compose", "ps", "-a", "--format", "json", serviceName,
	)

	o, err := cmd.Output()
	if err != nil {
		return domain.Service{}, err
	}

	stringOutput := strings.Trim(string(o), " ")
	if stringOutput == "" {
		return domain.Service{}, portout.ErrNotFound
	}

	splitOutputs := strings.Split(string(o), "\n")
	if len(splitOutputs) == 0 {
		return domain.Service{}, portout.ErrNotFound
	}

	output := splitOutputs[0]

	var s service
	err = json.Unmarshal([]byte(output), &s)
	if err != nil {
		return domain.Service{}, err
	}

	return domain.NewService(s.Name, domain.StatusFrom(s.State)), nil
}

func (repo *dockerCompose) Up(ctx context.Context) error {
	cmd := exec.CommandContext(
		ctx,
		"docker", "compose", "up", "-d",
	)

	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func (repo *dockerCompose) UpByService(ctx context.Context, services ...string) error {
	if len(services) == 0 {
		return portout.ErrNoService
	}

	args := []string{"compose", "up", "-d"}
	for _, s := range services {
		args = append(args, s)
	}

	cmd := exec.CommandContext(ctx, "docker", args...)

	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func (repo *dockerCompose) Down(ctx context.Context) error {
	cmd := exec.CommandContext(
		ctx,
		"docker", "compose", "down", "--remove-orphans",
	)

	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func (repo *dockerCompose) DownByService(ctx context.Context, services ...string) error {
	if len(services) == 0 {
		return portout.ErrNoService
	}

	args := []string{"compose", "down", "--remove-orphans"}
	for _, s := range services {
		args = append(args, s)
	}

	cmd := exec.CommandContext(ctx, "docker", args...)

	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func New() *dockerCompose {
	return &dockerCompose{}
}
