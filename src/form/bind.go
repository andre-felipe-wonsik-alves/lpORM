package form

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/andre-felipe-wonsik-alves/lpORM/src/schema"
)

// reads form fields and converts them to the approriate type in Go
func Bind[T any](r *http.Request, m *schema.Model) (*T, map[string][]string) {
	errs := map[string][]string{}
	if err := r.ParseForm(); err != nil {
		errs["_form"] = []string{err.Error()}
		return nil, errs
	}

	out := new(T)
	v := reflect.ValueOf(out).Elem()

	for _, f := range m.Fields {
		if f.Auto {
			continue
		}

		raw := r.PostForm.Get(f.Name)
		fv := v.FieldByName(f.Name)

		if err := setValue(fv, raw, r.PostForm.Has(f.Name)); err != nil {
			errs[f.Name] = append(errs[f.Name], err.Error())
			continue
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return out, nil
}

// casts the value to the appropriate Go type
func setValue(fv reflect.Value, raw string, present bool) error {
	s := strings.TrimSpace(raw)

	switch fv.Kind() {
	case reflect.String:
		fv.SetString(raw)

	case reflect.Bool:
		// checkbox not sent if unset
		fv.SetBool(present && raw != "false" && raw != "0")

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if s == "" {
			return nil
		}

		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil || fv.OverflowInt(n) {
			return fmt.Errorf("número inteiro inválido")
		}
		fv.SetInt(n)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if s == "" {
			return nil
		}

		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil || fv.OverflowUint(n) {
			return fmt.Errorf("número inteiro inválido")
		}

		fv.SetUint(n)

	case reflect.Float32, reflect.Float64:
		if s == "" {
			return nil
		}

		n, err := strconv.ParseFloat(s, 64)
		if err != nil || fv.OverflowFloat(n) {
			return fmt.Errorf("número inválido")
		}

		fv.SetFloat(n)

	default:
		return fmt.Errorf("tipo %s não suportado", fv.Type())
	}

	return nil
}
