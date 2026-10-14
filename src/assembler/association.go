package assembler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/andre-felipe-wonsik-alves/lpORM/src/schema"
)

// builds and registers an association model like a normal model so it get's a page and handler
func (a *Assembler) registerAssociation(m *schema.Model, assoc schema.Field) error {
	ass, err := a.buildAssociationModel(m, assoc)
	if err != nil {
		return err
	}

	return a.registerModel(ass, a.associationHandler(ass))
}

func (a *Assembler) buildAssociationModel(m *schema.Model, assoc schema.Field) (*schema.Model, error) {
	pk, err := getPK(m)
	if err != nil {
		return nil, err
	}

	otherTable, otherColumn, _ := strings.Cut(assoc.Association, ".")
	otherModel, other := a.findField(otherTable, otherColumn)

	if other == nil {
		return nil, fmt.Errorf("assembler: %s: n-n=%s: %s.%s não encontrado (reigstre esse model antes)",
			m.Name, assoc.Association, otherTable, otherColumn)
	}

	thisCol := m.Table + "_" + pk.Column
	otherCol := otherTable + "_" + otherColumn

	assModel := &schema.Model{
		Name:  m.Name + otherModel.Name,
		Table: m.Table + "_" + otherTable,
		Fields: []schema.Field{
			{
				Name:      thisCol,
				Column:    thisCol,
				GoType:    pk.GoType,
				PK:        true,
				Reference: m.Table + "." + pk.Column,
				Form:      map[string]string{"type": "number", "label": m.Name + " (id)"},
			},
			{
				Name:      otherCol,
				Column:    otherCol,
				GoType:    other.GoType,
				PK:        true,
				Reference: otherTable + "." + otherColumn,
				Form:      map[string]string{"type": "number", "label": otherModel.Name + " (id)"},
			},
		},
	}

	return assModel, nil
}

func (a *Assembler) associationHandler(m *schema.Model) http.HandlerFunc {
	cols := make([]string, len(m.Fields))
	values := make([]string, len(m.Fields))

	for i, f := range m.Fields {
		cols[i] = a.d.Quote(f.Column)
		values[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		a.d.Quote(m.Table), strings.Join(cols, ", "), strings.Join(values, ", "))

	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "%s\n", err)
			return
		}

		args := make([]any, len(m.Fields))
		for i, f := range m.Fields {
			n, err := strconv.ParseInt(strings.TrimSpace(r.PostForm.Get(f.Name)), 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				fmt.Fprintf(w, "%s: número inteiro inválido\n", f.Name)
				return
			}

			args[i] = n
		}

		if _, err := a.db.ExecContext(r.Context(), query, args...); err != nil {
			log.Printf("[!] insert em %s: %v", m.Table, err)
			// debug
			http.Error(w, "erro ao salvar: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, "salvo:")
		for i, f := range m.Fields {
			fmt.Fprintf(w, " %s=%d", f.Column, args[i])
		}

		fmt.Fprintln(w)
	}
}

// returns a model's pk
func getPK(m *schema.Model) (*schema.Field, error) {
	for _, f := range m.Fields {
		if f.PK {
			return &f, nil
		}
	}

	return nil, fmt.Errorf("assembler: %s: n-n exige uma chave primária", m.Name)
}

func (a *Assembler) findField(table, column string) (*schema.Model, *schema.Field) {
	for _, m := range a.models {
		if m.Table != table {
			continue
		}

		for i := range m.Fields {
			if m.Fields[i].Column == column {
				return m, &m.Fields[i]
			}
		}
	}

	return nil, nil
}
