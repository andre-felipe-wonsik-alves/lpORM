package widgets

import (
	"fmt"
	"html"
	"html/template"
)

type TextInput struct {
	InputType string
}

func (w TextInput) Render(name, label string, value any, errs []string) template.HTML {
	inputType := w.InputType
	if inputType == "" {
		inputType = "text"
	}

	return template.HTML(fmt.Sprintf(
		`<div><label for="%s">%s</label><input type="%s" id="%s" name="%s" value="%s">%s</div>`,
		html.EscapeString(name),
		html.EscapeString(label),
		html.EscapeString(inputType),
		html.EscapeString(name),
		html.EscapeString(name),
		html.EscapeString(stringify(value)),
		errorsHTML(errs),
	))
}

func (w TextInput) ParseValue(value any) any {
	if value == nil {
		return ""
	}
	return value
}
