package cmdconfig

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"

	"github.com/FaisalBudiono/comhel/internal/app/core/util/log"
	"github.com/FaisalBudiono/comhel/internal/app/core/util/log/logattr"
	"github.com/FaisalBudiono/comhel/internal/app/core/util/styleutil"
	"github.com/FaisalBudiono/comhel/internal/app/domain"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func renderTable(presets []configPreset) string {
	if presets == nil {
		return "Loading"
	}

	t := table.New().
		StyleFunc(func(row, col int) lipgloss.Style {
			switch row {
			case table.HeaderRow:
				return styleutil.Header().Align(lipgloss.Center)
			default:
				return styleutil.Cell().Align(lipgloss.Center)
			}
		}).
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(styleutil.ColorDarkPurple)).
		Headers("NO", "KEY", "SERVICES").
		Rows(renderTableRows(presets)...)

	return t.Render()
}

func renderTableRows(presets []configPreset) [][]string {
	mapConfigs := make(map[string]configPreset, len(presets))
	for _, c := range presets {
		mapConfigs[c.key] = c
	}

	rows := make([][]string, len(validKeys))
	for i, key := range validKeys {
		no := strconv.FormatInt(int64(i+1), 10)
		services := mapConfigs[key].services

		rows[i] = []string{no, key, formatServices(services)}
	}

	return rows
}

func formatServices(services []string) string {
	if len(services) == 0 {
		return styleutil.Disable().Render("<none>")
	}

	stl := lipgloss.NewStyle().Bold(true).Render

	formattedServices := make([]string, len(services))
	for i, s := range services {
		formattedServices[i] = styleutil.Active().Render(
			fmt.Sprintf("[%d]%s", i+1, s),
		)
	}

	return fmt.Sprintf(
		"%s%s%s",
		stl("[ "),
		strings.Join(formattedServices, ", "),
		stl(" ]"),
	)
}

func renderError(err error) string {
	if err == nil {
		return ""
	}

	s := fmt.Sprintf("Failed to parse .comhelconfig.json:\nReason: %s",
		err.Error(),
	)

	return styleutil.Error().Render(s) + "\n\n"
}

func cleanPresets(doms []domain.ConfigPreset, validServices []string) []configPreset {
	l := log.Logger().With(logattr.Caller("cmdconfig: cleanPresets"))

	if doms == nil {
		return nil
	}

	if len(doms) == 0 {
		return []configPreset{}
	}

	res := make([]configPreset, 0)
	for _, d := range doms {
		if !slices.Contains(validKeys, d.Key) {
			l.Debug("skipping invalid key", slog.String("key", d.Key))

			continue
		}

		ackServices := make([]string, 0)
		for _, s := range d.Services {
			l.Debug("checking service validity")

			if slices.Contains(validServices, s) {
				l.Debug("add service", slog.String("service", s))
				ackServices = append(ackServices, s)
			}
		}

		res = append(res, configPreset{
			key:      d.Key,
			services: ackServices,
		})
	}

	return res
}
