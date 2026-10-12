package widgets

import (
	"fmt"
	"html/template"
	"strings"
)

type SelectInput struct {
	Options     []Option
	Placeholder string
}

func (w SelectInput) Render(name, label string, value any, errs []string) template.HTML {
	var errHTML strings.Builder
	for _, e := range errs {
		fmt.Fprintf(&errHTML, `<span class="error">%s</span>`, e)
	}

	currentVal := fmt.Sprintf("%v", w.ParseValue(value))

	var optionsHTML strings.Builder
	placeholder := w.Placeholder
	if placeholder == "" {
		placeholder = "Selecione..."
	}
	fmt.Fprintf(&optionsHTML, `<option value="">%s</option>`, placeholder)

	for _, opt := range w.Options {
		selected := ""
		if currentVal != "" && currentVal == opt.Value {
			selected = " selected"
		}
		fmt.Fprintf(&optionsHTML, `<option value="%s"%s>%s</option>`, opt.Value, selected, opt.Label)
	}

	html := fmt.Sprintf(
		`<div><label for="%s">%s</label><select id="%s" name="%s">%s</select>%s</div>`,
		name, label, name, name, optionsHTML.String(), errHTML.String(),
	)
	return template.HTML(html)
}

func (w SelectInput) ParseValue(value any) any {
	if value == nil {
		return ""
	}
	return value
}
