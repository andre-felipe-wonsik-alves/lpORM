package main

import (
	"fmt"
	"log"
	"os"

	"github.com/andre-felipe-wonsik-alves/lpORM/input"
	"github.com/andre-felipe-wonsik-alves/lpORM/src/form"
	"github.com/andre-felipe-wonsik-alves/lpORM/src/schema"
)

func main() {
	m, err := schema.Parse(input.Cliente{})
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}
	fmt.Printf("Model: %s  (tabela: %s)\n\n", m.Name, m.Table)
	for _, f := range m.Fields {
		fmt.Printf("  Campo:  %s\n", f.Name)
		fmt.Printf("  Coluna: %s\n", f.Column)
		fmt.Printf("  Tipo:   %s\n", f.GoType)
		fmt.Printf("  PK=%v  Auto=%v  Unique=%v  Size=%d\n", f.PK, f.Auto, f.Unique, f.Size)
		fmt.Printf("  Form:   %v\n", f.Form)
		fmt.Printf("  Rules:  %v\n\n", f.Rules)
	}

	form_gerado, _ := form.New(input.Cliente{})

	html := form_gerado.Render()
	htmlCompleto := `
	<!DOCTYPE html>
		<html lang="pt-BR">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Minha página</title>
		</head>
		<body>
		` + html + `
		</body>
	</html>`

	err = os.WriteFile("./out/FORM.html", []byte(htmlCompleto), 0o644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Arquivo gerado com sucesso!\n")
}
