package widgets

import (
	"fmt"
	"html/template"
)

type HiddenInput struct{}

func (w HiddenInput) Render(name, label string, value any, errs []string) template.HTML {
	valStr := fmt.Sprintf("%v", w.ParseValue(value))
	if valStr == "<nil>" {
		valStr = ""
	}

	html := fmt.Sprintf(`<input type="hidden" id="%s" name="%s" value="%s">`, name, name, valStr)
	return template.HTML(html)
}

func (w HiddenInput) ParseValue(value any) any {
	if value == nil {
		return ""
	}
	return value
}
