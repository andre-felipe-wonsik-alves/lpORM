package repository

import (
	"context"
	"database/sql"

	"github.com/andre-felipe-wonsik-alves/lpORM/src/schema"
)

type repository[T any] struct {
	db    *sql.DB
	model *schema.Model
}

func NewRepository[T any](db *sql.DB) (Repository[T], error) {
	var zero T
	model, err := schema.Parse(zero) // however your schema package builds a Model
	if err != nil {
		return nil, err
	}

	return &repository[T]{db: db, model: model}, nil
}

func (r *repository[T]) Create(ctx context.Context, e *T) error {
	print("Método não implementado!")

	return nil
}
