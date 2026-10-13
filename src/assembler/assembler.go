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
	for _, f := range m.Fields {
		col := fmt.Sprintf("%s %s", a.d.Quote(f.Column), a.d.SQLType(f))

		if f.PK {
			col += " PRIMARY KEY"
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
		"CREATE TABLE IF NOT EXISTS %s (\n%s\n);",
		a.d.Quote(m.Table),
		strings.Join(cols, ",\n"),
	)

	fmt.Printf("[*] Executando DDL para %q:\n%s\n\n", m.Table, ddl)

	_, err := a.db.Exec(ddl)
	return err
}
