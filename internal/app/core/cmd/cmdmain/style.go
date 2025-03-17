package cmdmain

import "github.com/charmbracelet/lipgloss"

var (
	activeColor = lipgloss.Color("#ca07ce")
	lightGray   = lipgloss.Color("241")
)

func cellStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(0, 1)
}

func activeStyle() lipgloss.Style {
	return cellStyle().
		Foreground(activeColor)
}

func numberCellStyle() lipgloss.Style {
	return cellStyle().
		Align(lipgloss.Center)
}

func numberActiveStyle() lipgloss.Style {
	return numberCellStyle().
		Foreground(activeColor)
}

func headerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(0, 2).
		Bold(true).
		Align(lipgloss.Center)
}

func activeHeaderStyle() lipgloss.Style {
	return headerStyle().
		Background(lightGray)
}

func activeHeaderMarkerStyle() lipgloss.Style {
	return activeHeaderStyle().
		Padding(0, 1)
}

func helperStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))
}
