package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/andre-felipe-wonsik-alves/lpORM/input"
	"github.com/andre-felipe-wonsik-alves/lpORM/src/assembler"
	"github.com/andre-felipe-wonsik-alves/lpORM/src/dialect"
	"github.com/andre-felipe-wonsik-alves/lpORM/src/form"
	"github.com/andre-felipe-wonsik-alves/lpORM/src/repository"
	"github.com/andre-felipe-wonsik-alves/lpORM/src/schema"
)

const dsn = "host=localhost port=5432 user=user password=password dbname=appdb sslmode=disable connect_timeout=5"

func main() {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("[!] Erro ao abrir conexão: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("[!] Banco inacessível: ", err)
	}

	// create table and form for each Model
	a := &assembler.Assembler{}
	a.New([]any{input.Cliente{}}, db)
	a.Assemble()

	// create repository and routes
	mux := http.NewServeMux()
	register[input.Cliente](mux, db, dialect.Postgres{})

	// get every form
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, "out/FORM.html")
	})

	fmt.Println("[*] Servidor em http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func register[T any](mux *http.ServeMux, db *sql.DB, d dialect.Dialect) {
	var zero T
	model, err := schema.Parse(zero)
	if err != nil {
		log.Fatal(err)
	}

	repo, err := repository.NewRepository[T](db, d)
	if err != nil {
		log.Fatal(err)
	}

	route := "/" + model.Table
	mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "use POST", http.StatusMethodNotAllowed)
			return
		}

		e, errs := form.Bind[T](r, model)
		if len(errs) > 0 {
			w.WriteHeader(http.StatusUnprocessableEntity)
			for field, msgs := range errs {
				fmt.Fprintf(w, "%s: %v\n", field, msgs)
			}
			return
		}

		if err := repo.Create(r.Context(), e); err != nil {
			log.Printf("[!] insert em %s: %v", route, err)
			// debug
			http.Error(w, "erro ao salvar: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "salvo: %+v\n", *e)
	})
}
