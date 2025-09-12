package store

import (
	"context"

	"github.com/omatheuscaetano/planus-api/internal/auth/model"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

type AuthStore interface {
    FindUserByEmail(c context.Context , email string) (*model.User, *errs.Error)
	CreateUser(c context.Context, person *model.User) (*model.User, *errs.Error)
}
