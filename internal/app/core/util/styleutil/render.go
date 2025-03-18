package styleutil

import (
	"strings"

	"github.com/FaisalBudiono/comhel/internal/app/domain"
	"github.com/charmbracelet/lipgloss"
)

func RenderHelper(groups [][]domain.Keymap) string {
	separator := "    "

	outs := make([]string, 0)
	for _, g := range groups {
		keys := make([]string, 0)
		descriptions := make([]string, 0)

		for _, h := range g {
			keys = append(keys, strings.Join(h.Keys, "/"))
			descriptions = append(descriptions, h.Description)
		}

		outs = append(outs, lipgloss.JoinHorizontal(lipgloss.Top,
			Helper().Render(strings.Join(keys, "\n")),
			Helper().Render(strings.Repeat(" : \n", len(keys))),
			Helper().Render(strings.Join(descriptions, "\n")),
		), separator)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, outs...)
}
