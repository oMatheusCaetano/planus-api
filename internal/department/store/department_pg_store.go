package store

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/omatheuscaetano/planus-api/internal/department/model"
	dbDto "github.com/omatheuscaetano/planus-api/pkg/db/dto"
	"github.com/omatheuscaetano/planus-api/pkg/db/function"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

type DepartmentPgStore struct {
    db        *sql.DB
    tableName string
}

func NewDepartmentPgStore(db *sql.DB) *DepartmentPgStore {
    return &DepartmentPgStore{
        db:        db,
        tableName: "departments",
    }
}

func (r *DepartmentPgStore) Paginate(ctx context.Context, props *dbDto.Paginate) (*dbDto.PaginatedData[model.Department], *errs.Error) {
	return function.Paginate(ctx, &function.PaginateProps[model.Department]{
		DB:        r.db,
		Props:       props,
		TableName: r.tableName,
		Select: func() []string {
			return []string{"id", "name", "created_at", "updated_at"}
		},
		Scan: func(rows *sql.Rows) (*model.Department, *errs.Error) {
			var model model.Department
			if err := rows.Scan(&model.ID, &model.Name, &model.CreatedAt, &model.UpdatedAt); err != nil {
				return nil, errs.From(err)
			}
			return &model, nil
		},
	})
}

func (r *DepartmentPgStore) All(c context.Context, props *dbDto.List) ([]*model.Department, *errs.Error) {
	return function.All(c, &function.AllProps[model.Department]{
		DB:        r.db,
		Props:     props,
		TableName: r.tableName,
		Select: func() []string {
			return []string{"id", "name", "created_at", "updated_at"}
		},
		Scan: func(rows *sql.Rows) (*model.Department, *errs.Error) {
			var model model.Department
			if err := rows.Scan(&model.ID, &model.Name, &model.CreatedAt, &model.UpdatedAt); err != nil {
				return nil, errs.From(err)
			}
			return &model, nil
		},
	})
}

func (r *DepartmentPgStore) Find(c context.Context, id int) (*model.Department, *errs.Error) {
	var model model.Department

	err := function.Find(c, &function.FindProps{
		DB:        r.db,
		TableName: r.tableName,
		Select: func() []string {
			return []string{"id", "name", "created_at", "updated_at"}
		},
		Where: func(b sq.SelectBuilder) sq.SelectBuilder {
			return b.Where(sq.Eq{"id": id})
		},
		Scan: func() []any {
			return []any{&model.ID, &model.Name, &model.CreatedAt, &model.UpdatedAt}
		},
	})

	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (r *DepartmentPgStore) Create(ctx context.Context, model *model.Department) (*model.Department, *errs.Error) {
	err := function.Create(ctx, &function.CreateProps{
		DB: 	   r.db,
		TableName: r.tableName,
		Columns: func() []string {
			return []string{"name", "created_at", "updated_at"}
		},
		Values: func() []any {
			return []any{model.Name, model.CreatedAt, model.UpdatedAt}
		},
		Returning: func() []string {
			return []string{"id", "name", "created_at", "updated_at"}
		},
		Scan: func() []any {
			return []any{&model.ID, &model.Name, &model.CreatedAt, &model.UpdatedAt}
		},
	})

	if err != nil {
		return nil, err
	}

	return model, nil
}

func (r *DepartmentPgStore) Update(c context.Context, id int, model *model.Department) (*model.Department, *errs.Error) {
	err := function.Update(c, &function.UpdateProps{
		DB:        r.db,
		TableName: r.tableName,
		Set: func(b sq.UpdateBuilder) sq.UpdateBuilder {
			return b.Set("name", model.Name).Set("updated_at", model.UpdatedAt)
		},
		Where: func(b sq.UpdateBuilder) sq.UpdateBuilder {
			return b.Where(sq.Eq{"id": id})
		},
		Returning: func() []string {
			return []string{"id", "name", "created_at", "updated_at"}
		},
		Scan: func() []any {
			return []any{&model.ID, &model.Name, &model.CreatedAt, &model.UpdatedAt}
		},
	})

	if (err != nil) {
		return nil, err
	}

	return model, nil
}

func (r *DepartmentPgStore) Delete(c context.Context, id int) *errs.Error {
	return function.Delete(c, &function.DeleteProps{
        DB:        r.db,
        TableName: r.tableName,
        Where: func (qb sq.DeleteBuilder) sq.DeleteBuilder {
            return qb.Where(sq.Eq{"id": id})
        },
    })
}

