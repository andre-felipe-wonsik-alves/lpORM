package assembler

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andre-felipe-wonsik-alves/lpORM/src/dialect"
	"github.com/andre-felipe-wonsik-alves/lpORM/src/form"
	"github.com/andre-felipe-wonsik-alves/lpORM/src/schema"
	_ "github.com/lib/pq"
)

type Assembler struct {
	declarations []any
	db           *sql.DB
}

func (a *Assembler) New(declarations []any, host string, dbname string) {
	connection := fmt.Sprintf("host=%s dbname=%s connect_timeout=5", host, dbname)
	fmt.Println("[*] Conectando com: ", connection)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal("[!] Erro Estabelecendo conexão com o banco: ", err)
	}

	a.db = db
	a.declarations = declarations
}

func (a *Assembler) Assemble() {
	var htmlCompleto strings.Builder
	fmt.Fprintf(&htmlCompleto, `
	<!DOCTYPE html>
		<html lang="pt-BR">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Minha página</title>
			<style>%s</style>
		</head>
		<body>`,
		form.CSS(),
	)

	for _, element := range a.declarations {
		m, err := schema.Parse(element)
		if err != nil {
			fmt.Println("[!] Erro:", err)
			return
		}

		if err := a.migrate(m); err != nil {
			log.Fatalf("[!] Erro ao migrar tabela %q: %v", m.Table, err)
		}

		form_gerado, _ := form.New(m)

		fmt.Fprint(&htmlCompleto, form_gerado.Render())
	}

	fmt.Fprint(&htmlCompleto, `</body></html>`)

	outputDir := "./out"

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("Erro ao criar o diretório: %v", err)
	}

	err := os.WriteFile("./out/FORM.html", []byte(htmlCompleto.String()), 0o644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[*] Arquivo gerado com sucesso!\n")
	err = a.db.Close()
	if err != nil {
		log.Fatal("[!] Erro fechando a conexão com o banco: ", err)
	}
}

func (a *Assembler) migrate(m *schema.Model) error {
	d := dialect.Postgres{}

	var cols []string
	for _, f := range m.Fields {
		col := fmt.Sprintf("%s %s", d.Quote(f.Column), d.SQLType(f))

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
		d.Quote(m.Table),
		strings.Join(cols, ",\n"),
	)

	fmt.Printf("[*] Executando DDL para %q:\n%s\n\n", m.Table, ddl)

	_, err := a.db.Exec(ddl)
	return err
}
