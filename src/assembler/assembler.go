package assembler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/andre-felipe-wonsik-alves/lpORM/src/dialect"
	"github.com/andre-felipe-wonsik-alves/lpORM/src/form"
	"github.com/andre-felipe-wonsik-alves/lpORM/src/schema"
)

type Assembler struct {
	db     *sql.DB
	d      dialect.Dialect
	mux    *http.ServeMux
	models []*schema.Model
}

func New(db *sql.DB) *Assembler {
	a := &Assembler{db: db, d: dialect.Postgres{}, mux: http.NewServeMux()}
	a.mux.HandleFunc("/", a.index)
	return a
}

// returns the router with every route registered
func (a *Assembler) Handler() http.Handler {
	return a.mux
}

// returns the index page; rebuilt every time
func (a *Assembler) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var links strings.Builder
	for _, m := range a.models {
		fmt.Fprintf(&links, `<li><a href="/%s">%s</a></li>`, m.Table, m.Name)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, page("Formulários", "<h1>Formulários</h1><ul>"+links.String()+"</ul>"))
}

// common page layout
func page(title, body string) string {
	return fmt.Sprintf(
		`<!DOCTYPE html>
<html lang="pt-BR">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>%s</title>
	<style>%s</style>
</head>
<body>%s</body>
</html>`,
		title, form.CSS(), body)
}

func (a *Assembler) migrate(m *schema.Model) error {
	var cols []string
	var foreign_keys []schema.Field
	var associations []schema.Field
	var pk schema.Field

	for i, f := range m.Fields {
		col := fmt.Sprintf("%s %s", a.d.Quote(f.Column), a.d.SQLType(f))

		if f.PK {
			pk = m.Fields[i]
			col += " PRIMARY KEY"
		}
		if f.Reference != "" {
			foreign_keys = append(foreign_keys, f)
		}
		if f.Association != "" {
			associations = append(associations, f)
		}
		if f.Unique && !f.PK {
			col += " UNIQUE"
		}
		if !f.PK {
			col += " NOT NULL"
		}

		cols = append(cols, "\t"+col)
	}

	ddl := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (\n%s",
		a.d.Quote(m.Table),
		strings.Join(cols, ",\n"),
	)

	for _, entry := range foreign_keys {
		references := strings.Split(entry.Reference, ".")
		ddl += fmt.Sprintf(",\n\tFOREIGN KEY (%s)\n\t\tREFERENCES %s (%s)\n", a.d.Quote(entry.Column), references[0], references[1])
	}
	ddl += ");"

	fmt.Printf("[*] Executando DDL para %q:\n%s\n\n", m.Table, ddl)

	if _, err := a.db.Exec(ddl); err != nil {
		return err
	}

	for _, ass := range associations {
		assDdl, err := a.createAssociation(m, pk, ass)
		if err != nil {
			return err
		}

		fmt.Printf("[*] Executando DDL para %q:\n%s\n\n", m.Table, assDdl)

		if _, err := a.db.Exec(ddl); err != nil {
			return err
		}
	}

	return nil
}

func (a *Assembler) createAssociation(m *schema.Model, pk schema.Field, ass schema.Field) (string, error) {
	otherReference := strings.Split(ass.Association, ".")
	otherTable := otherReference[0]
	otherColumn := otherReference[1]
	var other schema.Field

	for _, m := range a.models {
		if m.Table != otherReference[0] {
			continue
		}

		for i := range m.Fields {
			if m.Fields[i].Column == otherReference[1] {
				other = m.Fields[i]
			}
		}
	}

	thisCol := m.Table + "_" + pk.Column
	otherCol := otherTable + "_" + otherColumn

	cols := []string{
		fmt.Sprintf("\t%s %s NOT NULL", a.d.Quote(thisCol), a.d.SQLType(pk)),
		fmt.Sprintf("\t%s %s NOT NULL", a.d.Quote(otherCol), a.d.SQLType(other)),
		fmt.Sprintf("\tPRIMARY KEY (%s, %s)", a.d.Quote(thisCol), a.d.Quote(otherCol)),
		fmt.Sprintf("\tFOREIGN KEY (%s) REFERENCES %s (%s)", a.d.Quote(thisCol), a.d.Quote(m.Table), a.d.Quote(pk.Column)),
		fmt.Sprintf("\tFOREIGN KEY (%s) REFERENCES %s (%s)", a.d.Quote(otherCol), a.d.Quote(otherTable), a.d.Quote(otherColumn)),
	}

	ddl := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (\n%s\n);",
		a.d.Quote(m.Table+"_"+otherReference[0]),
		strings.Join(cols, ",\n"),
	)

	return ddl, nil
}
