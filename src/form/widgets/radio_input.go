package widgets

import (
	"fmt"
	"html/template"
	"strings"
)

type RadioInput struct {
	Options []Option
}

func (w RadioInput) Render(name, label string, value any, errs []string) template.HTML {
	var errHTML strings.Builder
	for _, e := range errs {
		fmt.Fprintf(&errHTML, `<span class="error">%s</span>`, e)
	}

	currentVal := fmt.Sprintf("%v", w.ParseValue(value))

	var radiosHTML strings.Builder
	for i, opt := range w.Options {
		radioID := fmt.Sprintf("%s_%d", name, i)
		checked := ""
		if currentVal != "" && currentVal == opt.Value {
			checked = " checked"
		}

		fmt.Fprintf(&radiosHTML,
			`<label for="%s"><input type="radio" id="%s" name="%s" value="%s"%s> %s</label> `,
			radioID, radioID, name, opt.Value, checked, opt.Label,
		)
	}

	html := fmt.Sprintf(
		`<div><label>%s</label><div class="radio-group">%s</div>%s</div>`,
		label, radiosHTML.String(), errHTML.String(),
	)
	return template.HTML(html)
}

func (w RadioInput) ParseValue(value any) any {
	if value == nil {
		return ""
	}
	return value
}
