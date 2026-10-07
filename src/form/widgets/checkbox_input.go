package widgets

import (
	"fmt"
	"html/template"
	"strings"
)

type CheckboxInput struct {
	InputType string
}

func (w CheckboxInput) Render(name, label string, value any, errs []string) template.HTML {
	var errHTML strings.Builder
	for _, e := range errs {
		fmt.Fprintf(&errHTML, `<span class="error">%s</span>`, e)
	}

	html := fmt.Sprintf(
		`<div><label for="%s">%s</label><input type="%s" id="%s" name="%s" %s>%s</div>`,
		name, label, w.InputType, name, name, w.ParseValue(value), errHTML.String(),
	)
	return template.HTML(html)
}

func (w CheckboxInput) ParseValue(value any) any {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case bool:
		if v {
			return "checked"
		}
	case string:
		if v == "true" || v == "on" || v == "1" {
			return "checked"
		}
	case int, int64:
		if v == 1 {
			return "checked"
		}
	}

	return ""
}
