package store

import (
	"context"

	"github.com/omatheuscaetano/planus-api/internal/company/model"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

type CompanyStore interface {
    FindById(c context.Context, id int) (*model.Company, *errs.Error)
    Create(c context.Context, company *model.Company) *errs.Error
    Update(c context.Context, company *model.Company) *errs.Error
    Delete(c context.Context, id int) *errs.Error
}
