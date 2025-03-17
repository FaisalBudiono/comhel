package cmdmain

import "github.com/charmbracelet/lipgloss"

func cellStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(0, 1)
}

func activeStyle() lipgloss.Style {
	return cellStyle().
		Foreground(lipgloss.Color("#ca07ce"))
}

func numberCellStyle() lipgloss.Style {
	return cellStyle().
		Align(lipgloss.Center)
}

func numberActiveStyle() lipgloss.Style {
	return numberCellStyle().
		Foreground(lipgloss.Color("#ca07ce"))
}

func headerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(0, 2).
		Bold(true).
		Align(lipgloss.Center)
}

func helperStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))
}
