package widgets

import (
	"fmt"
	"html/template"
	"strings"
)

type TextareaInput struct {
	Rows int
	Cols int
}

func (w TextareaInput) Render(name, label string, value any, errs []string) template.HTML {
	var errHTML strings.Builder
	for _, e := range errs {
		fmt.Fprintf(&errHTML, `<span class="error">%s</span>`, e)
	}

	//
	rows := w.Rows
	if rows <= 0 {
		rows = 4
	}

	colsAttr := ""
	if w.Cols > 0 {
		colsAttr = fmt.Sprintf(` cols="%d"`, w.Cols)
	}

	valStr := fmt.Sprintf("%v", w.ParseValue(value))

	html := fmt.Sprintf(
		`<div><label for="%s">%s</label><textarea id="%s" name="%s" rows="%d"%s>%s</textarea>%s</div>`,
		name, label, name, name, rows, colsAttr, valStr, errHTML.String(),
	)
	return template.HTML(html)
}

func (w TextareaInput) ParseValue(value any) any {
	if value == nil {
		return ""
	}
	return value
}
