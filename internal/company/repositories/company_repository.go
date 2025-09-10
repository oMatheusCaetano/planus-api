package repositories

import (
	"database/sql"

	"github.com/omatheuscaetano/planus-api/internal/company/models"
	"github.com/omatheuscaetano/planus-api/internal/db"

	"github.com/omatheuscaetano/planus-api/internal/shared/dto"
	"github.com/omatheuscaetano/planus-api/internal/shared/errs"
)

type CompanyRepository interface {
    Paginate(props dto.PaginationProps) (*dto.Paginated[models.Company], *errs.Error)
    All(props dto.ListingProps) (*[]models.Company, *errs.Error)
	Find(id int) (*models.Company, *errs.Error)
    Create(company *models.Company) *errs.Error
}

type companyRepository struct {
}

func NewCompanyRepository() CompanyRepository {
	return &companyRepository{}
}

func (r *companyRepository) Paginate(props dto.PaginationProps) (*dto.Paginated[models.Company], *errs.Error) {
    counterQb := db.From("companies")

    if len(props.Where) > 0 {
        counterQb.WhereFromLogicBlock(props.Where)
    }

    total, err := counterQb.Count()
    if err != nil {
        return nil, err
    }

    if (props.PerPage == 0) {
        props.PerPage = 10
    }

    if (props.Page == 0) {
        props.Page = 1
    }
    
    query := counterQb.
        Duplicate().
        Limit(props.PerPage).
        Offset((props.Page - 1) * props.PerPage)

    for _, sort := range props.SortBy {
        if (sort.Direction != "asc" && sort.Direction != "desc") {
            sort.Direction = "asc"
        }

        query.SortBy(sort.Key, sort.Direction)
    }

    var companies []models.Company
    scanError := query.ScanMany(companies, func (rows *sql.Rows) error {
        var model models.Company
        err := rows.Scan(&model.ID, &model.Name, &model.CNPJ, &model.CreatedAt, &model.UpdatedAt)
        if err != nil {
            return err
        }
        companies = append(companies, model)
        return nil
    })


    if scanError != nil {
        return nil, scanError
    }

    paginated := &dto.Paginated[models.Company]{
        Meta: dto.PaginationMeta{
            Total:       total,
            PerPage:     props.PerPage,
            Page:        props.Page,
            LastPage:    (total + props.PerPage - 1) / props.PerPage,
            FirstPage:   1,
            SortBy:      []db.SortBy{{Key: "id", Direction: "asc"}},
            Where:       []db.Where{},
        },
        Data: companies,
    }

    return paginated, nil
}

func (r *companyRepository) All(props dto.ListingProps) (*[]models.Company, *errs.Error) {
    
    var companies []models.Company
    query := db.From("companies");

    if len(props.Where) > 0 {
        query.WhereFromLogicBlock(props.Where)
    }

    for _, sort := range props.SortBy {
        if (sort.Direction != "asc" && sort.Direction != "desc") {
            sort.Direction = "asc"
        }

        query.SortBy(sort.Key, sort.Direction)
    }

    err := query.ScanMany(companies, func (rows *sql.Rows) error {
        var model models.Company
        err := rows.Scan(&model.ID, &model.Name, &model.CNPJ, &model.CreatedAt, &model.UpdatedAt)
        if err != nil {
            return err
        }
        companies = append(companies, model)
        return nil
    })


    if err != nil {
        return nil, err
    }
    return &companies, nil
}

func (r *companyRepository) Find(id int) (*models.Company, *errs.Error) {
    var company models.Company

    err := db.From("companies").Where("id", "=", id).
        Scan(&company.ID, &company.Name, &company.CNPJ, &company.CreatedAt, &company.UpdatedAt)

    if err != nil {
        return nil, err
    }
    return &company, nil
}

func (r *companyRepository) Create(company *models.Company) *errs.Error {
    err := db.
        Insert("companies").
        Values(map[string]interface{}{
            "name": company.Name,
            "cnpj": company.CNPJ,
            "created_at": company.CreatedAt,
            "updated_at": company.UpdatedAt,
        }).
        Returning("id").
        Scan(&company.ID)

    if err != nil {
        return errs.From(err)
    }
    return nil
}
