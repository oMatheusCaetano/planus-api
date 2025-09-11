package store

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/omatheuscaetano/planus-api/internal/company/model"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

type CompanyRepository struct {
	db *sql.DB
	psql sq.StatementBuilderType
}

func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *CompanyRepository) FindById(ctx context.Context, id int) (*model.Company, *errs.Error) {
	query := r.psql.
		Select("id", "name", "cnpj").
		From("companies").
		Where(sq.Eq{"id": id})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, errs.From(err)
	}

	row := r.db.QueryRowContext(ctx, sqlStr, args...)
	var company model.Company
	if err := row.Scan(&company.ID, &company.Name, &company.CNPJ); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errs.From(err)
	}

	return &company, nil
}

func (r *CompanyRepository) Create(ctx context.Context, company *model.Company) *errs.Error {
	query := r.psql.
		Insert("companies").
		Columns("name", "cnpj").
		Values(company.Name, company.CNPJ).
		Suffix("RETURNING id")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return errs.From(err)
	}

	return errs.From(r.db.QueryRowContext(ctx, sqlStr, args...).Scan(&company.ID))
}

func (r *CompanyRepository) Update(ctx context.Context, company *model.Company) *errs.Error {
	query := r.psql.
		Update("companies").
		Set("name", company.Name).
		Set("cnpj", company.CNPJ).
		Where(sq.Eq{"id": company.ID})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return errs.From(err)
	}

	res, err := r.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return errs.From(err)
	}

	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return errs.From(fmt.Errorf("company with id %d not found", company.ID))
	}

	return nil
}

func (r *CompanyRepository) Delete(ctx context.Context, id int) *errs.Error {
	sqlStr, args, err := r.psql.Delete("companies").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return errs.From(err)
	}

	res, err := r.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return errs.From(err)
	}

	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return errs.From(fmt.Errorf("company with id %d not found", id))
	}

	return nil
}