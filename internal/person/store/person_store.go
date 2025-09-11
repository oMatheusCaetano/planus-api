package store

import (
	"context"

	"github.com/omatheuscaetano/planus-api/internal/person/dto"
	"github.com/omatheuscaetano/planus-api/internal/person/model"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

type PersonStore interface {
	Paginate(c context.Context, dto *dto.PaginatePerson) (*dto.PaginatedPerson, *errs.Error)
	All(c context.Context, dto *dto.ListPerson) ([]*model.Person, *errs.Error)
	Find(c context.Context, id int) (*model.Person, *errs.Error)
    Create(c context.Context, person *model.Person) (*model.Person, *errs.Error)
    Update(c context.Context, id int, person *model.Person) (*model.Person, *errs.Error)
	Delete(c context.Context, id int) *errs.Error
}
