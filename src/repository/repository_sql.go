package repository

import (
	"context"
	"database/sql"

	"github.com/andre-felipe-wonsik-alves/lpORM/src/dialect"
	"github.com/andre-felipe-wonsik-alves/lpORM/src/schema"
)

type repository[T any] struct {
	db      *sql.DB
	model   *schema.Model
	dialect dialect.Dialect
}

func NewRepository[T any](
	db *sql.DB,
	d dialect.Dialect,
) (Repository[T], error) {
	return &repository[T]{
		db:      db,
		dialect: d,
	}, nil
}

func (r *repository[T]) Create(ctx context.Context, e *T) error {
	print("Método não implementado!")

	return nil
}
