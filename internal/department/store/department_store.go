package store

import (
	"context"

	"github.com/omatheuscaetano/planus-api/internal/department/model"
	dbDto "github.com/omatheuscaetano/planus-api/pkg/db/dto"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

type DepartmentStore interface {
    Paginate(c context.Context, dto *dbDto.List) (*dbDto.PaginatedData[model.Department], *errs.Error)
	All(c context.Context, dto *dbDto.List) ([]*model.Department, *errs.Error)
	Find(c context.Context, id int) (*model.Department, *errs.Error)
    Create(c context.Context, person *model.Department) (*model.Department, *errs.Error)
    Update(c context.Context, id int, person *model.Department) (*model.Department, *errs.Error)
	Delete(c context.Context, id int) *errs.Error
}
