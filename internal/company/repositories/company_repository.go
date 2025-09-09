package repositories

import (
	"database/sql"

	"github.com/omatheuscaetano/planus-api/internal/company/models"
	"github.com/omatheuscaetano/planus-api/internal/shared/errs"
)

type CompanyRepository interface {
	Find(id int) (*models.Company, *errs.Error)
}

type companyRepository struct {
	db *sql.DB
}

func NewCompanyRepository(db *sql.DB) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) Find(id int) (*models.Company, *errs.Error) {
    var company models.Company
    query := "SELECT * FROM companies WHERE id = $1"
    err := r.db.QueryRow(query, id).Scan(&company.ID, &company.Name, &company.CNPJ, &company.CreatedAt, &company.UpdatedAt)
    if err != nil {
        return nil, errs.From(err)
    }
    return &company, nil
}
