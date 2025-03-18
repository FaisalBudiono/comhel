package styleutil

import "github.com/charmbracelet/lipgloss"

func Title() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(ColorDarkPurple).
		Border(lipgloss.DoubleBorder()).
		Width(50).Height(3).
		Align(lipgloss.Center).AlignVertical(lipgloss.Center).
		Bold(true).Italic(true)
}

func Error() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color("#ff0000")).
		Foreground(lipgloss.Color("#ffffff")).
		Border(lipgloss.RoundedBorder()).
		Height(4).
		Padding(0, 2).
		Align(lipgloss.Center).AlignVertical(lipgloss.Center).
		Bold(true).Italic(true)
}

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

func Disable() lipgloss.Style {
	return Helper()
}
