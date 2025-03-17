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

func noCellStyle() lipgloss.Style {
	return cellStyle().
		Align(lipgloss.Center)
}

func noActiveStyle() lipgloss.Style {
	return noCellStyle().
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
