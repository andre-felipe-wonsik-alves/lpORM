package widgets

import (
	"fmt"
	"html"
	"html/template"
	"strings"
)

type RadioInput struct {
	Options []Option
}

func (w RadioInput) Render(name, label string, value any, errs []string) template.HTML {
	currentVal := stringify(value)

	var radiosHTML strings.Builder
	for i, opt := range w.Options {
		radioID := fmt.Sprintf("%s_%d", name, i)
		checked := ""
		if currentVal != "" && currentVal == opt.Value {
			checked = " checked"
		}

		fmt.Fprintf(&radiosHTML,
			`<label for="%s"><input type="radio" id="%s" name="%s" value="%s"%s> %s</label> `,
			html.EscapeString(radioID),
			html.EscapeString(radioID),
			html.EscapeString(name),
			html.EscapeString(opt.Value),
			checked,
			html.EscapeString(opt.Label),
		)
	}

	return template.HTML(fmt.Sprintf(
		`<div><label>%s</label><div class="radio-group">%s</div>%s</div>`,
		html.EscapeString(label),
		radiosHTML.String(),
		errorsHTML(errs),
	))
}

func (w RadioInput) ParseValue(value any) any {
	if value == nil {
		return ""
	}
	return value
}
