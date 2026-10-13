package form

import (
	"fmt"
	"html/template"
	"reflect"
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

func New(model *schema.Model, opts ...Option) (*Form, error) {
	f := &Form{
		Model:  model,
		Values: map[string]any{},
		Errors: map[string][]string{},
		method: "POST",
	}

	for _, opt := range opts {
		opt(f)
	}

	return f, nil
}

func (f *Form) Render() template.HTML {
	var buf strings.Builder
	fmt.Fprintf(&buf, `<h3>%s</h3><form action="%s" method="%s">`, f.Model.Name, f.action, f.method)
	for _, field := range f.Model.Fields {

		w := resolveWidget(field)

		if w == nil {
			continue
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
	// Se for chave primária com auto incremento, renderiza como campo oculto
	if f.PK && f.Auto {
		return widgets.HiddenInput{}
	}

	formType := ""
	if f.Form != nil {
		formType = f.Form["type"]
	}

	switch formType {
	case "checkbox":
		return widgets.CheckboxInput{InputType: "checkbox"}
	case "select":
		return widgets.SelectInput{Options: widgets.ParseOptions(f.Form["options"])}
	case "radio":
		return widgets.RadioInput{Options: widgets.ParseOptions(f.Form["options"])}
	case "textarea":
		return widgets.TextareaInput{}
	case "hidden":
		return widgets.HiddenInput{}
	case "password":
		return widgets.TextInput{InputType: "password"}
	case "email":
		return widgets.TextInput{InputType: "email"}
	case "date":
		return widgets.TextInput{InputType: "date"}
	case "number":
		return widgets.TextInput{InputType: "number"}
	case "text":
		return widgets.TextInput{InputType: "text"}
	case "":
		// Se não foi especificado tipo de formulário, infere a partir do tipo Go
		if f.GoType != nil {
			switch f.GoType.Kind() {
			case reflect.Bool:
				return widgets.CheckboxInput{InputType: "checkbox"}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Float32, reflect.Float64:
				return widgets.TextInput{InputType: "number"}
			default:
				return widgets.TextInput{InputType: "text"}
			}
		}
		return widgets.TextInput{InputType: "text"}
	default:
		return widgets.TextInput{InputType: formType}
	}
}
