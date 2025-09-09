package repositories

import (
	"database/sql"

	"github.com/omatheuscaetano/planus-api/internal/company/models"
	"github.com/omatheuscaetano/planus-api/internal/shared/dto"
	"github.com/omatheuscaetano/planus-api/internal/shared/errs"
)

type CompanyRepository interface {
    Paginate(offset, limit int) (dto.Paginated[models.Company], *errs.Error)
    All() ([]models.Company, *errs.Error)
	Find(id int) (*models.Company, *errs.Error)
    Create(company *models.Company) *errs.Error
}

type companyRepository struct {
	db *sql.DB
}

func NewCompanyRepository(db *sql.DB) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) Paginate(offset, limit int) (dto.Paginated[models.Company], *errs.Error) {
    var total int
    countQuery := "SELECT COUNT(*) FROM companies"
    err := r.db.QueryRow(countQuery).Scan(&total)
    if err != nil {
        return dto.Paginated[models.Company]{}, errs.From(err)
    }

    query := "SELECT * FROM companies ORDER BY id LIMIT $1 OFFSET $2"
    rows, err := r.db.Query(query, limit, offset)
    if err != nil {
        return dto.Paginated[models.Company]{}, errs.From(err)
    }
    defer rows.Close()

    var companies []models.Company
    for rows.Next() {
        var company models.Company
        err := rows.Scan(&company.ID, &company.Name, &company.CNPJ, &company.CreatedAt, &company.UpdatedAt)
        if err != nil {
            return dto.Paginated[models.Company]{}, errs.From(err)
        }
        companies = append(companies, company)
    }

    if err = rows.Err(); err != nil {
        return dto.Paginated[models.Company]{}, errs.From(err)
    }

    currentPage := (offset / limit) + 1
    lastPage := (total + limit - 1) / limit

    paginated := dto.Paginated[models.Company]{
        Meta: dto.PaginationMeta{
            Total:       total,
            PerPage:     limit,
            CurrentPage: currentPage,
            LastPage:    lastPage,
            FirstPage:   1,
            SortBy:      []dto.PaginationSortBy{{Key: "id", Direction: "asc"}},
            Where:       []dto.PaginationWhere{},
        },
        Data: companies,
    }

    return paginated, nil
}

func (r *companyRepository) All() ([]models.Company, *errs.Error) {
    var companies []models.Company
    query := "SELECT * FROM companies"
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, errs.From(err)
    }
    defer rows.Close()

    for rows.Next() {
        var company models.Company
        err := rows.Scan(&company.ID, &company.Name, &company.CNPJ, &company.CreatedAt, &company.UpdatedAt)
        if err != nil {
            return nil, errs.From(err)
        }
        companies = append(companies, company)
    }

    if err = rows.Err(); err != nil {
        return nil, errs.From(err)
    }

    return companies, nil
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
