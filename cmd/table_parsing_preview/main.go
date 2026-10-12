package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/andre-felipe-wonsik-alves/lpORM/input"
	"github.com/andre-felipe-wonsik-alves/lpORM/src/dialect"
	"github.com/andre-felipe-wonsik-alves/lpORM/src/schema"
	_ "github.com/lib/pq"
)

const configInfo = "host=localhost port=5432 user=user password=password dbname=appdb sslmode=disable connect_timeout=5"

func main() {
	db, err := sql.Open("postgres", configInfo)
	if err != nil {
		log.Fatal("[!] Erro ao abrir conexão: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("[!] Banco inacessível: ", err)
	}
	fmt.Println("[*] Conectado ao PostgreSQL com sucesso!")

	d := dialect.Postgres{}
	for _, entry := range []any{input.Cliente{}, input.Job{}} {
		model, _ := schema.Parse(entry)

		var cols []string
		for _, field := range model.Fields {
			col := fmt.Sprintf("%s %s", d.Quote(field.Column), d.SQLType(field))
			if field.PK {
				col += " PRIMARY KEY"
			}
			if field.Unique && !field.PK {
				col += " UNIQUE"
			}
			if !field.PK {
				col += " NOT NULL"
			}
			cols = append(cols, "\t"+col)
		}

		ddl := fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS %s (\n%s\n);",
			d.Quote(model.Table),
			strings.Join(cols, ",\n"),
		)

		fmt.Printf("[*] DDL para %q:\n%s\n\n", model.Table, ddl)

		if _, err := db.Exec(ddl); err != nil {
			log.Fatalf("[!] Erro ao criar tabela %q: %v", model.Table, err)
		}
		fmt.Printf("[*] Tabela %q criada!\n\n", model.Table)
	}
}
