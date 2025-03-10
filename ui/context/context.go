package context

import "github.com/gevann/gh-dash/config"

type ProgramContext struct {
	ScreenHeight      int
	ScreenWidth       int
	MainContentWidth  int
	MainContentHeight int
	Config            *config.Config
	View              config.ViewType
}

func (ctx *ProgramContext) GetViewSectionsConfig() []config.SectionConfig {
	if ctx.View == config.PRsView {
		return ctx.Config.PRSections
	} else {
		return ctx.Config.IssuesSections
	}
}
