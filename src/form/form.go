package form

import (
	"html/template"
	"net/http"

	"github.com/andre-felipe-wonsik-alves/lpORM/src/schema"
)

type Form struct {
	Model  *schema.Model
	Values map[string]any
	Errors map[string][]string

	action      string
	method      string
	submitLabel string
}

func New(model any, opts ...any) (*Form, error)
func (f *Form) Bind(r *http.Request) error
func (f *Form) Validate() bool
func (f *Form) Render() template.HTML
