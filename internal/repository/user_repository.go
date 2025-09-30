package repository

import (
	"context"

	"github.com/omatheuscaetano/planus-api/internal/model"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

type UserRepository interface {
	Repository[model.User]
	FindByUsername(ctx context.Context, username string) (*model.User, *errs.Error)
}
