package runssection

import "github.com/charmbracelet/lipgloss"

var (
	titleCellWidth    = lipgloss.Width("vxx.xxx.xxx Release to Production")
	statusCellWidth   = lipgloss.Width("completed success")
	branchCellWidth   = 50
	workflowCellWidth = lipgloss.Width("Test, Build & Deploy to ECS")
	ContainerPadding  = 1

	containerStyle = lipgloss.NewStyle().
			Padding(0, ContainerPadding)

	spinnerStyle = lipgloss.NewStyle().Padding(0, 1)

	emptyStateStyle = lipgloss.NewStyle().
			Faint(true).
			PaddingLeft(1).
			MarginBottom(1)
)
