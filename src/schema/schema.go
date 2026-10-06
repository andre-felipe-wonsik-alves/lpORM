package schema

import (
	"reflect"
)

type Field struct {
	Name             string
	Column           string
	GoType           reflect.Type
	PK, Auto, Unique bool
	Size             int
	Form             map[string]string // opções de formulário (label, type, options)
	Rules            []string
}

type Model struct {
	Name   string
	Table  string
	Fields []Field
}
