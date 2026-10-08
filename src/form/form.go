package form

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/andre-felipe-wonsik-alves/lpORM/src/form/widgets"
	"github.com/andre-felipe-wonsik-alves/lpORM/src/schema"
)

type Form struct {
	Model  *schema.Model
	Values map[string]any
	Errors map[string][]string

	action string
	method string
}

func New(model *schema.Model, opts ...any) (*Form, error) {
	f := &Form{Model: model, Values: map[string]any{}, Errors: map[string][]string{}}

	return f, nil
}

func (f *Form) Render() template.HTML {
	var buf strings.Builder
	fmt.Fprintf(&buf, `<h3>%s</h3><form action="%s" method="%s">`, f.Model.Name, f.action, f.method)
	for _, field := range f.Model.Fields {

		w := resolveWidget(field)

		if w == nil {
			continue // widget ainda não implementado
		}

		label := field.Form["label"]
		if label == "" {
			label = field.Name
		}
		buf.WriteString(string(w.Render(field.Name, label, f.Values[field.Name], f.Errors[field.Name])))
	}

	submit := "Enviar"

	fmt.Fprintf(&buf, `<button type="submit">%s</button></form><hr>`, submit)
	return template.HTML(buf.String())
}

func resolveWidget(f schema.Field) widgets.Widgets {
	switch f.Form["type"] {
	case "checkbox":
		println("[*] Processando input do tipo checkbox")
		return widgets.CheckboxInput{InputType: f.Form["type"]}
	case "select":
		// TODO return widgets.SelectInput{Options: opts}
		println("[!] Select não implementado!")
		return nil
	default:
		println("[*] Processando input do tipo texto")
		return widgets.TextInput{InputType: f.Form["type"]}
	}
}
