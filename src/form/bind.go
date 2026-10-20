package form

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/andre-felipe-wonsik-alves/lpORM/src/schema"
)

// reads form fields and builds a struct T instance with the approriate types in Go
func Bind[T any](r *http.Request, m *schema.Model) (*T, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	// create empyt T struct
	s := new(T)

	// get the struct
	v := reflect.ValueOf(s).Elem()

	for _, f := range m.Fields {
		if f.Auto {
			continue
		}

		// set the current struct field to the form value in the request
		fv := r.PostForm.Get(f.Name)
		sf := v.FieldByName(f.Name)

		if err := setValue(sf, fv, r.PostForm.Has(f.Name)); err != nil {
			continue
		}
	}

	return s, nil
}

// parses the value received in the form (fv) to the struct field (sf)
func setValue(sf reflect.Value, fv string, present bool) error {
	s := strings.TrimSpace(fv)

	switch sf.Kind() {
	case reflect.String:
		sf.SetString(fv)

	case reflect.Bool:
		sf.SetBool(present && fv != "false" && fv != "0")

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if s == "" {
			return nil
		}

		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fmt.Errorf("número inteiro inválido")
		}

		sf.SetInt(n)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if s == "" {
			return nil
		}

		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return fmt.Errorf("número inteiro inválido")
		}

		sf.SetUint(n)

	case reflect.Float32, reflect.Float64:
		if s == "" {
			return nil
		}

		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("número inválido")
		}

		sf.SetFloat(n)

	default:
		return fmt.Errorf("tipo %s não suportado", sf.Type())
	}

	return nil
}
