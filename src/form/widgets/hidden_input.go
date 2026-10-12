package widgets

import (
	"fmt"
	"html"
	"html/template"
)

type HiddenInput struct{}

func (w HiddenInput) Render(name, label string, value any, errs []string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<input type="hidden" id="%s" name="%s" value="%s">`,
		html.EscapeString(name),
		html.EscapeString(name),
		html.EscapeString(stringify(value)),
	))
}

func (w HiddenInput) ParseValue(value any) any {
	return stringify(value)
}
