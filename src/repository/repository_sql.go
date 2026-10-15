package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/andre-felipe-wonsik-alves/lpORM/src/dialect"
	"github.com/andre-felipe-wonsik-alves/lpORM/src/schema"
)

type repository[T any] struct {
	db     *sql.DB
	model  *schema.Model
	query  string
	fields []string
}

func NewRepository[T any](db *sql.DB, d dialect.Dialect) (Repository[T], error) {
	var zero T
	model, err := schema.Parse(zero)
	if err != nil {
		return nil, err
	}

	r := &repository[T]{db: db, model: model}

	var cols, values []string
	for _, f := range model.Fields {
		if f.Auto {
			continue
		}

		r.fields = append(r.fields, f.Name)
		cols = append(cols, d.Quote(f.Column))
		values = append(values, fmt.Sprintf("$%d", len(r.fields))) // psotgres uses "$" to represent a value in VALUES(...)
	}

	r.query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		d.Quote(model.Table), strings.Join(cols, ", "), strings.Join(values, ", "))

	return r, nil
}

func (r *repository[T]) Create(ctx context.Context, e *T) error {
	if e == nil {
		return errors.New("lporm: Create recebeu ponteiro nil")
	}

	v := reflect.ValueOf(e).Elem()
	args := make([]any, len(r.fields))
	for i, name := range r.fields {
		args[i] = v.FieldByName(name).Interface()
	}

	// maybe change to return the inserted tuple's id?
	_, err := r.db.ExecContext(ctx, r.query, args...)
	return err
}
