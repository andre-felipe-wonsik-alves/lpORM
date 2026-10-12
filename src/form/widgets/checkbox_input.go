package widgets

import (
	"fmt"
	"html"
	"html/template"
)

type CheckboxInput struct {
	InputType string
}

func (w CheckboxInput) Render(name, label string, value any, errs []string) template.HTML {
	checked := ""
	if isChecked(value) {
		checked = " checked"
	}

	return template.HTML(fmt.Sprintf(
		`<div><label for="%s">%s</label><input type="checkbox" id="%s" name="%s"%s>%s</div>`,
		html.EscapeString(name),
		html.EscapeString(label),
		html.EscapeString(name),
		html.EscapeString(name),
		checked,
		errorsHTML(errs),
	))
}

func (w CheckboxInput) ParseValue(value any) any {
	return isChecked(value)
}

// isChecked decide se o valor representa "marcado".
func isChecked(v any) bool {
	switch x := v.(type) {
	case bool:
		return x
	case string:
		return x == "true" || x == "on" || x == "1"
	case int:
		return x == 1
	case int64:
		return x == 1
	}
	return false
}
