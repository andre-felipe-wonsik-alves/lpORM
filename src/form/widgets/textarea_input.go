package widgets

import (
	"fmt"
	"html"
	"html/template"
)

type TextareaInput struct {
	Rows int
	Cols int
}

func (w TextareaInput) Render(name, label string, value any, errs []string) template.HTML {
	rows := w.Rows
	if rows <= 0 {
		rows = 4
	}

	colsAttr := ""
	if w.Cols > 0 {
		colsAttr = fmt.Sprintf(` cols="%d"`, w.Cols)
	}

	return template.HTML(fmt.Sprintf(
		`<div><label for="%s">%s</label><textarea id="%s" name="%s" rows="%d"%s>%s</textarea>%s</div>`,
		html.EscapeString(name),
		html.EscapeString(label),
		html.EscapeString(name),
		html.EscapeString(name),
		rows,
		colsAttr,
		html.EscapeString(stringify(value)),
		errorsHTML(errs),
	))
}

func (w TextareaInput) ParseValue(value any) any {
	if value == nil {
		return ""
	}
	return value
}
