package store

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/omatheuscaetano/planus-api/internal/person/model"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

type PersonPgStore struct {
	db        *sql.DB
	psql      sq.StatementBuilderType
	tableName string
}

func NewPersonPgStore(db *sql.DB) *PersonPgStore {
	return &PersonPgStore{
		db:        db,
		tableName: "people",
		psql:      sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *PersonPgStore) All(ctx context.Context) ([]*model.Person, *errs.Error) {
	query := r.psql.
		Select("id", "name", "created_at", "updated_at").
		From(r.tableName)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, errs.From(err)
	}

	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, errs.From(err)
	}
	defer rows.Close()

	var people []*model.Person
	for rows.Next() {
		var person model.Person
		if err := rows.Scan(&person.ID, &person.Name, &person.CreatedAt, &person.UpdatedAt); err != nil {
			return nil, errs.From(err)
		}
		people = append(people, &person)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.From(err)
	}

	return people, nil
}

func (r *PersonPgStore) Find(ctx context.Context, id int) (*model.Person, *errs.Error) {
	query := r.psql.
		Select("id", "name", "created_at", "updated_at").
		From(r.tableName).
		Where(sq.Eq{"id": id})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, errs.From(err)
	}

	var person model.Person
	if err := r.db.QueryRowContext(ctx, sqlStr, args...).Scan(&person.ID, &person.Name, &person.CreatedAt, &person.UpdatedAt); err != nil {
		return nil, errs.From(err)
	}

	return &person, nil
}

func (r *PersonPgStore) Create(ctx context.Context, model *model.Person) (*model.Person, *errs.Error) {
	query := r.psql.
		Insert(r.tableName).
		Columns("name", "created_at", "updated_at").
		Values(model.Name, model.CreatedAt, model.UpdatedAt).
		Suffix("RETURNING id, name, created_at, updated_at")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, errs.From(err)
	}

	if err := r.db.QueryRowContext(ctx, sqlStr, args...).Scan(&model.ID, &model.Name, &model.CreatedAt, &model.UpdatedAt); err != nil {
		return nil, errs.From(err)
	}

	return model, nil
}

func (r *PersonPgStore) Update(ctx context.Context, id int, model *model.Person) (*model.Person, *errs.Error) {
	query := r.psql.
		Update(r.tableName).
		Set("name", model.Name).
		Set("updated_at", model.UpdatedAt).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, name, created_at, updated_at")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, errs.From(err)
	}

	if err := r.db.QueryRowContext(ctx, sqlStr, args...).Scan(&model.ID, &model.Name, &model.CreatedAt, &model.UpdatedAt); err != nil {
		return nil, errs.From(err)
	}

	return model, nil
}

func (r *PersonPgStore) Delete(ctx context.Context, id int) *errs.Error {
	query := r.psql.
		Delete(r.tableName).
		Where(sq.Eq{"id": id})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return errs.From(err)
	}

	_, err = r.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return errs.From(err)
	}

	return nil
}
