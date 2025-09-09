package repositories

import (
	"database/sql"

	"github.com/omatheuscaetano/planus-api/internal/company/models"
	"github.com/omatheuscaetano/planus-api/internal/shared/errs"
)

type CompanyRepository interface {
	Find(id int) (*models.Company, *errs.Error)
    Create(company *models.Company) *errs.Error
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

func (r *companyRepository) Create(company *models.Company) *errs.Error {
    query := "INSERT INTO companies (name, cnpj, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id"
    err := r.db.QueryRow(query, company.Name, company.CNPJ, company.CreatedAt, company.UpdatedAt).Scan(&company.ID)
    if err != nil {
        return errs.From(err)
    }
    return nil
}
