package repository

import (
	"context"

	"github.com/omatheuscaetano/planus-api/internal/model"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)


type UserRepository interface {
	Find(ctx context.Context, id model.ID) (*model.User, *errs.Error)
	Delete(ctx context.Context, id model.ID) *errs.Error
}
