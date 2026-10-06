package schema

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

/*
Lê os marcadores e devolve um Model com os metadados.
*/
func Parse(v any) (*Model, error) {
	t := reflect.TypeOf(v)

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("lporm: Parse espera uma struct, recebeu %s", t.Kind())
	}

	model := &Model{
		Name:  t.Name(),
		Table: toSnakeCase(t.Name()) + "s",
	}

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)

		field := Field{
			Name:   sf.Name,
			Column: toSnakeCase(sf.Name),
			GoType: sf.Type,
		}

		if raw := sf.Tag.Get("lporm"); raw != "" {
			for _, part := range strings.Split(raw, ",") {
				part = strings.TrimSpace(part)
				if kv := strings.SplitN(part, "=", 2); len(kv) == 2 {
					switch kv[0] {
					case "size":
						field.Size, _ = strconv.Atoi(kv[1])
					}
				} else {
					switch part {
					case "pk":
						field.PK = true
					case "auto":
						field.Auto = true
					case "unique":
						field.Unique = true
					}
				}
			}
		}

		if raw := sf.Tag.Get("form"); raw != "" {
			field.Form = make(map[string]string)
			for _, pair := range strings.Split(raw, ";") {
				if kv := strings.SplitN(strings.TrimSpace(pair), "=", 2); len(kv) == 2 {
					field.Form[kv[0]] = kv[1]
				}
			}
		}

		if raw := sf.Tag.Get("validate"); raw != "" {
			field.Rules = strings.Split(raw, ",")
		}

		model.Fields = append(model.Fields, field)
	}

	return model, nil
}

func toSnakeCase(s string) string {
	runes := []rune(s)
	var b strings.Builder
	for i, r := range runes {
		upper := unicode.IsUpper(r)
		if upper && i > 0 {
			prevLower := unicode.IsLower(runes[i-1])
			nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if prevLower || nextLower {
				b.WriteByte('_')
			}
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}
