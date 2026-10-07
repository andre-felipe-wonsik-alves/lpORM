package widgets

import "html/template"

type Widgets interface {
	Render(name, label string, value any, errs []string) template.HTML
	ParseValue(value any) any
}
