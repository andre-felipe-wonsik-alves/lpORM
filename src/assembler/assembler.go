package assembler

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

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
	fmt.Fprint(&htmlCompleto, `
	<!DOCTYPE html>
		<html lang="pt-BR">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Minha página</title>
		</head>
		<body>`,
	)

	for _, element := range a.declarations {
		m, err := schema.Parse(element)
		if err != nil {
			fmt.Println("Erro:", err)
			return
		}

		// TODO: criar o schema no banco

		//fmt.Printf("Model: %s  (tabela: %s)\n\n", m.Name, m.Table)
		//for _, f := range m.Fields {
		//	fmt.Printf("  Campo:  %s\n", f.Name)
		//	fmt.Printf("  Coluna: %s\n", f.Column)
		//	fmt.Printf("  Tipo:   %s\n", f.GoType)
		//	fmt.Printf("  PK=%v  Auto=%v  Unique=%v  Size=%d\n", f.PK, f.Auto, f.Unique, f.Size)
		//	fmt.Printf("  Form:   %v\n", f.Form)
		//	fmt.Printf("  Rules:  %v\n\n", f.Rules)
		//}

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

	fmt.Printf("Arquivo gerado com sucesso!\n")
	err = a.db.Close()
	if err != nil {
		log.Fatal("[!] Erro fechando a conexão com o banco: ", err)
	}
}
