package widgets

import (
	"html/template"
	"strings"
)

type Widgets interface {
	Render(name, label string, value any, errs []string) template.HTML
	ParseValue(value any) any
}

type Option struct {
	Value string
	Label string
}

// ParseOptions converte uma string de opções como "PR,SP,SC" ou "M:Masculino,F:Feminino"
// em uma lista de structs Option.
func ParseOptions(raw string) []Option {
	if raw == "" {
		return nil
	}

	var options []Option
	for _, item := range strings.Split(raw, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		parts := strings.SplitN(item, ":", 2)
		if len(parts) == 2 {
			options = append(options, Option{
				Value: strings.TrimSpace(parts[0]),
				Label: strings.TrimSpace(parts[1]),
			})
		} else {
			options = append(options, Option{
				Value: item,
				Label: item,
			})
		}
	}

	return options
}

