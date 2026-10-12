package widgets

import (
	"fmt"
	"html"
	"html/template"
	"strings"
)

type SelectInput struct {
	Options     []Option
	Placeholder string
}

func (w SelectInput) Render(name, label string, value any, errs []string) template.HTML {
	currentVal := stringify(value)

	placeholder := w.Placeholder
	if placeholder == "" {
		placeholder = "Selecione..."
	}

	var optionsHTML strings.Builder
	fmt.Fprintf(&optionsHTML, `<option value="">%s</option>`, html.EscapeString(placeholder))

	for _, opt := range w.Options {
		selected := ""
		if currentVal != "" && currentVal == opt.Value {
			selected = " selected"
		}
		fmt.Fprintf(&optionsHTML, `<option value="%s"%s>%s</option>`,
			html.EscapeString(opt.Value), selected, html.EscapeString(opt.Label))
	}

	return template.HTML(fmt.Sprintf(
		`<div><label for="%s">%s</label><select id="%s" name="%s">%s</select>%s</div>`,
		html.EscapeString(name),
		html.EscapeString(label),
		html.EscapeString(name),
		html.EscapeString(name),
		optionsHTML.String(),
		errorsHTML(errs),
	))
}

func (w SelectInput) ParseValue(value any) any {
	if value == nil {
		return ""
	}
	return value
}
