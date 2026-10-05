package repository

import "context"

type Repository[T any] interface {
	Create(ctx context.Context, e *T) error
}
