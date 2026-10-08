package web

import (
	"fmt"

	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-ddns-updater/app"
	"github.com/mitoteam/mtweb"
)

type PageBuilder struct {
	mtweb.PageBuilderBase

	toolbar *mtweb.BtnPanelElement
}

func NewPageBuilder(ctx *mbr.MbrContext) *PageBuilder {
	p := &PageBuilder{
		PageBuilderBase: mtweb.NewPageBuilderBase(ctx),
	}

	p.BuildHeadTitleF = func() string {
		title := app.App.AppName

		if p.GetTitle() != "" {
			title = p.GetTitle() + " | " + title
		}

		return title
	}

	p.RenderF = p.render

	return p
}

// Creates a mbr.RouterHandleFunc that builds a page using PageBuilder
func CreatePageBuilderRouteHandler(buildPageF func(*PageBuilder) any) mbr.RouterHandleFunc {
	return func(ctx *mbr.MbrContext) any {
		p := NewPageBuilder(ctx)

		out := buildPageF(p)

		if err, ok := out.(error); ok { // error happen, return it as-is
			return err
		} else if p.Ctx.IsRedirect() { //redirect is already set, so we don't need to return content
			return nil
		} else { // render the page as html
			html, err := p.Render()

			if err != nil {
				return err
			}

			return html
		}
	}
}

// returns BtnPanelElement to be used as page toolbar
func (p *PageBuilder) Toolbar() *mtweb.BtnPanelElement {
	if p.toolbar == nil {
		p.toolbar = mtweb.NewBtnPanel()
	}

	return p.toolbar
}

func (p *PageBuilder) render() error {
	document := p.GetDocument().
		Icon("/favicon.ico").
		Stylesheet("/assets/css/style.css")

	container := dhtml.Div().Class("container my-3")

	container.Append(p.renderHeader())

	// H1 page title
	title := p.GetTitle()
	if title != "" {
		container.Append(dhtml.NewTag("h1").Append(title))
	}

	if p.toolbar != nil && !p.toolbar.IsEmpty() {
		p.toolbar.Class("mb-3").Class("page-toolbar")

		container.Append(p.toolbar)
	}

	container.Append(dhtml.Div().Class("region-main").Append(p.GetMain()))

	container.Append(p.renderFooter())

	document.Body().Append(container)

	//script at the bottom
	document.Body().
		Append(dhtml.NewTag("script").Attribute("src", "/assets/script.min.js"))

	return nil
}

func (p *PageBuilder) renderHeader() (out dhtml.HtmlPiece) {
	header := dhtml.Div().Class("region-header border bg-light p-3 mb-3").Attribute("role", "header")

	header_left := dhtml.Div().
		Append(dhtml.Div().Append(dhtml.NewLink(mbr.Url(RootCtl.Home)).Label(app.App.AppName).Class("text-decoration-none")))

	header_right := dhtml.Div().Class("text-end")

	if p.IsAuthenticated() {
		header_right.Append(
			mtweb.NewSmBtn(mbr.Url(RootCtl.Logout), "arrow-right-from-bracket").
				Label(dhtml.Span().Class("ms-1").Append("Logout")),
		)
	}

	header.Append(dhtmlbs.NewJustifiedLR().L(header_left).R(header_right))

	out.Append(header)
	return out
}

func (p *PageBuilder) renderFooter() (out dhtml.HtmlPiece) {
	out.Append(dhtml.Div().Class("region-footer border bg-light p-3 mt-3").Append(
		dhtmlbs.NewJustifiedLR().
			L(
				fmt.Sprintf("This instance: v%s", app.App.Version),
				dhtml.Span().Class("small text-muted ms-2").Append(
					mtweb.Icon(mtweb.IconTimestamp).Label(app.App.BuildTime),
				),
				dhtml.Div().Append(
					dhtml.NewLink("https://github.com/mitoteam/mt-ddns-updater").Label(
						dhtml.UnsafeText("<img alt=\"GitHub Release\" src=\"https://img.shields.io/github/v/release/mitoteam/mt-ddns-updater?style=flat-square&logo=github&label=latest%20version\">"),
					),
				),
				dhtml.Div().Class("mt-1").Append(
					mtweb.NewSmBtn(
						"https://github.com/mitoteam/mt-ddns-updater/issues/new?template=bug_report.md", "bug",
					).Label(dhtml.Span().Class("ms-1").Append("Report a Bug")).Target("blank"),
					mtweb.NewSmBtn(
						"https://github.com/mitoteam/mt-ddns-updater/issues/new?template=feature_request.md", "lightbulb-on",
					).Label(dhtml.Span().Class("ms-1").Append("Suggest a Feature")).Target("blank"),
				),
			).
			R(
				dhtml.Div().Class("small text-end").Append(
					app.App.AppName+" by ",
					dhtml.NewLink("https://www.mito-team.com").Label("MiTo Team").Target("blank"),
				),
				dhtml.Div().Class("small text-muted text-end").Append(goapp.MOTTO),
			),
	))
	return out
}

// Helper function to check if the user is authenticated in the current page context
func (p *PageBuilder) IsAuthenticated() bool {
	return p.Ctx.Get(isAuthenticatedCtxField) == true
}
