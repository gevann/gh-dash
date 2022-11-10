package runssection

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gevann/gh-dash/config"
	"github.com/gevann/gh-dash/data"
	"github.com/gevann/gh-dash/ui/components/section"
	"github.com/gevann/gh-dash/ui/components/table"
	"github.com/gevann/gh-dash/ui/constants"
	"github.com/gevann/gh-dash/ui/context"
	"github.com/gevann/gh-dash/utils"
)

const SectionType = "runs"

type Model struct {
	Runs    []data.RunData
	section section.Model
	error   error
}

func (m *Model) getDimensions() constants.Dimensions {
	return constants.Dimensions{
		Width:  m.section.Ctx.MainContentWidth - containerStyle.GetHorizontalPadding(),
		Height: m.section.Ctx.MainContentHeight - 2,
	}
}

func (m *Model) GetSectionColumns() []table.Column {
	return []table.Column{
		{
			Title: "Title",
			Width: &titleCellWidth,
		},
		{
			Title: "Status",
			Width: &statusCellWidth,
		},
		{
			Title: "Workflow",
			Width: &workflowCellWidth,
		},
		{
			Title: "Branch",
			Width: &branchCellWidth,
		},
	}
}

func NewModel(id int, ctx *context.ProgramContext, config config.SectionConfig) Model {
	m := Model{
		Runs: []data.RunData{},
		section: section.Model{
			Id:        id,
			Config:    config,
			Ctx:       ctx,
			Spinner:   spinner.Model{Spinner: spinner.Dot},
			IsLoading: true,
			Type:      SectionType,
		},
		error: nil,
	}

	m.section.Table = table.NewModel(
		m.getDimensions(),
		m.GetSectionColumns(),
		m.BuildRows(),
		"Run",
		utils.StringPtr(emptyStateStyle.Render(fmt.Sprintf(
			"No runs were found that match the given filters: %s",
			lipgloss.NewStyle().Italic(true).Render(m.section.Config.Filters),
		))),
	)

	return m
}

type SectionRunsFetchedMsg struct {
	SectionId int
	Runs      []data.RunData
	Err       error
}

type Run struct {
	Data  data.RunData
	Width int
}

func (run *Run) renderTitle() string {
	return lipgloss.NewStyle().Render(run.Data.Title)
}

func (run *Run) renderStatus() string {
	return lipgloss.NewStyle().Render(run.Data.Status)
}

func (run *Run) renderWorkflow() string {
	return lipgloss.NewStyle().Render(run.Data.Workflow)
}

func (run *Run) renderBranch() string {
	return lipgloss.NewStyle().Render(run.Data.Branch)
}

func (run *Run) renderEvent() string {
	return lipgloss.NewStyle().Render(run.Data.Event)
}

func (run *Run) renderId() string {
	return lipgloss.NewStyle().Render(run.Data.Id)
}

func (run *Run) renderElapsed() string {
	return lipgloss.NewStyle().Render(run.Data.Elapsed)
}

func (run *Run) renderAge() string {
	return lipgloss.NewStyle().Render(run.Data.Age)
}

func (run *Run) ToTableRow() table.Row {
	return table.Row{
		run.renderTitle(),
		run.renderStatus(),
		run.renderWorkflow(),
		run.renderBranch(),
		run.renderEvent(),
		run.renderId(),
		run.renderElapsed(),
		run.renderAge(),
	}
}

func (m Model) BuildRows() []table.Row {
	rows := make([]table.Row, 0, len(m.Runs))

	return rows
}

func (m *Model) FetchSectionRows() tea.Cmd {
	m.error = nil
	if m == nil {
		return nil
	}
	m.Runs = nil
	m.section.Table.ResetCurrItem()
	m.section.Table.Rows = nil
	m.section.IsLoading = true
	var cmds []tea.Cmd
	cmds = append(cmds, m.section.CreateNextTickCmd(spinner.Tick))

	cmds = append(cmds, func() tea.Msg {
		limit := m.section.Config.Limit
		if limit == nil {
			limit = &m.section.Ctx.Config.Defaults.RunsLimit
		}
		fetchedRuns, err := data.ListRuns(limit)
		if err != nil {
			return SectionRunsFetchedMsg{
				SectionId: m.section.Id,
				Runs:      []data.RunData{},
				Err:       err,
			}
		}

		return SectionRunsFetchedMsg{
			SectionId: m.section.Id,
			Runs:      fetchedRuns,
		}
	})

	return tea.Batch(cmds...)
}

func (m *Model) GetCurrRow() data.RowData {
	if len(m.Runs) == 0 {
		return nil
	}
	run := m.Runs[m.section.Table.GetCurrItem()]
	return &run
}

func (m *Model) NextRow() int {
	return m.section.Table.NextItem()
}

func (m *Model) PrevRow() int {
	return m.section.Table.PrevItem()
}

func (m *Model) FirstItem() int {
	return m.section.Table.FirstItem()
}

func (m *Model) LastItem() int {
	return m.section.Table.LastItem()
}

func (m *Model) GetIsLoading() bool {
	return m.section.IsLoading
}

func (m *Model) Id() int {
	return m.section.Id
}

func (m *Model) NumRows() int {
	return len(m.Runs)
}

func (m *Model) UpdateProgramContext(ctx *context.ProgramContext) {
	oldDimensions := m.getDimensions()
	m.section.Ctx = ctx
	newDimensions := m.getDimensions()
	m.section.Table.SetDimensions(newDimensions)

	if oldDimensions.Height != newDimensions.Height || oldDimensions.Width != newDimensions.Width {
		m.section.Table.SyncViewPortContent()
	}
}

func (m *Model) View() string {
	var spinnerText *string
	if m.section.IsLoading {
		spinnerText = utils.StringPtr(lipgloss.JoinHorizontal(lipgloss.Top,
			spinnerStyle.Copy().Render(m.section.Spinner.View()),
			"Fetching Workflow Runs...",
		))
	}

	if m.error != nil {
		spinnerText = utils.StringPtr(fmt.Sprintf("Error while fetching workflow runs: %v", m.error))
	}

	return containerStyle.Copy().Render(
		m.section.Table.View(spinnerText),
	)
}

func (m Model) Update(msg tea.Msg) (section.Section, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case SectionRunsFetchedMsg:
		m.Runs = msg.Runs
		m.section.IsLoading = false
		m.section.Table.Rows = m.BuildRows()
		m.section.Table.EmptyState = utils.StringPtr(emptyStateStyle.Render(fmt.Sprintf(
			"No runs were found that match the given filters: %s",
			lipgloss.NewStyle().Italic(true).Render(m.section.Config.Filters),
		)))
	case section.SectionTickMsg:
		if !m.section.IsLoading {
			return &m, nil
		}

		var internalTickCmd tea.Cmd
		m.section.Spinner, internalTickCmd = m.section.Spinner.Update(msg.InternalTickMsg)
		cmd = m.section.CreateNextTickCmd(internalTickCmd)
	}

	return &m, cmd
}
