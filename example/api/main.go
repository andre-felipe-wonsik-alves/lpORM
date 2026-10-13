package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/andre-felipe-wonsik-alves/lpORM/input"
	"github.com/andre-felipe-wonsik-alves/lpORM/src/assembler"
)

const dbCongig = "host=localhost port=5432 user=user password=password dbname=appdb sslmode=disable connect_timeout=5"

func main() {
	db, err := sql.Open("postgres", dbCongig)
	if err != nil {
		log.Fatal("[!] Erro ao abrir conexão: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("[!] Banco inacessível: ", err)
	}

	app := assembler.New(db)

	if err := assembler.Register[input.Cliente](app); err != nil {
		log.Fatal(err)
	}

	if err := assembler.Register[input.Job](app); err != nil {
		log.Fatal(err)
	}

	fmt.Println("[*] Servidor em http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", app.Handler()))
}
