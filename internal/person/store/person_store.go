package store

import (
	"context"

	"github.com/omatheuscaetano/planus-api/internal/person/model"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

type PersonStore interface {
	All(c context.Context) ([]*model.Person, *errs.Error)
	Find(c context.Context, id int) (*model.Person, *errs.Error)
    Create(c context.Context, person *model.Person) (*model.Person, *errs.Error)
    Update(c context.Context, id int, person *model.Person) (*model.Person, *errs.Error)
	Delete(c context.Context, id int) *errs.Error
}
