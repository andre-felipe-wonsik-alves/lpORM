package assembler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/andre-felipe-wonsik-alves/lpORM/src/form"
	"github.com/andre-felipe-wonsik-alves/lpORM/src/repository"
	"github.com/andre-felipe-wonsik-alves/lpORM/src/schema"
)

const outputDir = "./out"

// register table T on Assembler instance so the corresponding HTML form is created and linked to index.html
func Register[T any](a *Assembler) error {
	var zero T
	model, err := schema.Parse(zero)
	if err != nil {
		return err
	}

	for _, m := range a.models {
		if m.Table == model.Table {
			return fmt.Errorf("assembler: %s já foi adicionado", model.Name)
		}
	}

	// migrate table
	if err := a.migrate(model); err != nil {
		return fmt.Errorf("assembler: migrar %q: %w", model.Table, err)
	}

	repo, err := repository.NewRepository[T](a.db, a.d)
	if err != nil {
		return err
	}

	// create route with the same name as T
	route := "/" + model.Table

	f, err := form.New(model, form.WithAction(route))
	if err != nil {
		return err
	}

	// build page for T
	html := page(model.Name, string(f.Render()))

	if err := writePage(model.Table, html); err != nil {
		return err
	}

	// handelrs for page T
	handler := insertHandler(model, repo)
	a.mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, html)

		case http.MethodPost:
			handler(w, r)

		default:
			http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		}
	})

	a.models = append(a.models, model)
	return nil
}

// returns a handler for POST /T
func insertHandler[T any](model *schema.Model, repo repository.Repository[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := form.Bind[T](r, model)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "%s\n", err)

			return
		}

		if err := repo.Create(r.Context(), s); err != nil {
			log.Printf("[!] insert em %s: %v", model.Table, err)
			// debgu
			http.Error(w, "erro ao salvar: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "salvo: %+v\n", *s)
	}
}

// persists a html page
func writePage(table, html string) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(outputDir, table+".html"), []byte(html), 0o644)
}
