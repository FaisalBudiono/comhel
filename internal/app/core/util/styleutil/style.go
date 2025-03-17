package styleutil

import "github.com/charmbracelet/lipgloss"

func Cell() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(0, 1)
}

func Active() lipgloss.Style {
	return Cell().
		Foreground(ColorActive)
}

func NumberCell() lipgloss.Style {
	return Cell().
		Align(lipgloss.Center)
}

func NumberActive() lipgloss.Style {
	return NumberCell().
		Foreground(ColorActive)
}

func Header() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(0, 2).
		Bold(true).
		Align(lipgloss.Center)
}

func ActiveHeader() lipgloss.Style {
	return Header().
		Background(ColorLightGray)
}

func ActiveHeaderMarker() lipgloss.Style {
	return ActiveHeader().
		Padding(0, 1)
}

func Helper() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))
}
