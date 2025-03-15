package compose

import (
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"slices"
	"strings"

	"github.com/FaisalBudiono/comhel/internal/app/domain"
)

func List(ctx context.Context) ([]string, error) {
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

var ErrNotFound = errors.New("service not found")

type service struct {
	Name  string `json:"Service"`
	State string `json:"State"`
}

func Service(ctx context.Context, serviceName string) (domain.Service, error) {
	cmd := exec.CommandContext(
		ctx,
		"docker", "compose", "ps", "--format", "json", serviceName,
	)

	o, err := cmd.Output()
	if err != nil {
		return domain.Service{}, err
	}

	stringOutput := strings.Trim(string(o), " ")
	if stringOutput == "" {
		return domain.Service{}, ErrNotFound
	}

	splitOutputs := strings.Split(string(o), "\n")
	if len(splitOutputs) == 0 {
		return domain.Service{}, ErrNotFound
	}

	output := splitOutputs[0]

	var s service
	err = json.Unmarshal([]byte(output), &s)
	if err != nil {
		return domain.Service{}, err
	}

	return domain.NewService(s.Name, domain.StatusFrom(s.State)), nil
}

func Up(ctx context.Context) error {
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

func Down(ctx context.Context) error {
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
