package widgets

import (
	"fmt"
	"html/template"
	"strings"
)

type TextInput struct {
	InputType string
}

func (w TextInput) Render(name, label string, value any, errs []string) template.HTML {
	var errHTML strings.Builder
	for _, e := range errs {
		fmt.Fprintf(&errHTML, `<span class="error">%s</span>`, e)
	}

	html := fmt.Sprintf(
		`<div><label for="%s">%s</label><input type="%s" id="%s" name="%s" value="%v">%s</div>`,
		name, label, w.InputType, name, name, w.ParseValue(value), errHTML.String(),
	)
	return template.HTML(html)
}

func (w TextInput) ParseValue(value any) any {
	if value == nil {
		return ""
	}
	return value
}
