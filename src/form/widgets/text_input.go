package widgets

import (
	"fmt"
	"html/template"
)

type TextInput struct {
	InputType string
}

func (w TextInput) Render(name, label string, value any, errs []string) template.HTML {
	errHTML := ""
	for _, e := range errs {
		errHTML += fmt.Sprintf(`<span class="error">%s</span>`, e)
	}
	html := fmt.Sprintf(
		`<div><label for="%s">%s</label><input type="%s" id="%s" name="%s" value="%v">%s</div>`,
		name, label, w.InputType, name, name, value, errHTML,
	)
	return template.HTML(html)
}
