package repository

import (
	"context"

	"github.com/omatheuscaetano/planus-api/internal/model"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

type Repository[T any] interface {
	Find(ctx context.Context, id model.ID) (*T, *errs.Error)
	Delete(ctx context.Context, id model.ID) *errs.Error
	Create(ctx context.Context, entity *T) *errs.Error
}
