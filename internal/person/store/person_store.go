package store

import (
	"context"

	"github.com/omatheuscaetano/planus-api/internal/person/model"
	dbDto "github.com/omatheuscaetano/planus-api/pkg/db/dto"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

type PersonStore interface {
	Paginate(c context.Context, dto *dbDto.Paginate) (*dbDto.PaginatedData[model.Person], *errs.Error)
	All(c context.Context, dto *dbDto.List) ([]*model.Person, *errs.Error)
	Find(c context.Context, id int) (*model.Person, *errs.Error)
    Create(c context.Context, model *model.Person) (*model.Person, *errs.Error)
    Update(c context.Context, id int, model *model.Person) (*model.Person, *errs.Error)
	Delete(c context.Context, id int) *errs.Error
}
