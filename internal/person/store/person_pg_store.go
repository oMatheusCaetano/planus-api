package store

import (
	"context"
	"database/sql"
	"sync"

	sq "github.com/Masterminds/squirrel"
	userModel "github.com/omatheuscaetano/planus-api/internal/auth/model"
	"github.com/omatheuscaetano/planus-api/internal/person/model"
	dbDto "github.com/omatheuscaetano/planus-api/pkg/db/dto"
	"github.com/omatheuscaetano/planus-api/pkg/db/function"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

type PersonPgStore struct {
	db        *sql.DB
	tableName string
}

func NewPersonPgStore(db *sql.DB) *PersonPgStore {
	return &PersonPgStore{
		db:        db,
		tableName: "people",
	}
}

func (r *PersonPgStore) Paginate(ctx context.Context, props *dbDto.Paginate) (*dbDto.PaginatedData[model.Person], *errs.Error) {
	return function.Paginate(ctx, &function.PaginateProps[model.Person]{
		DB:        r.db,
		Props:     props,
		TableName: r.tableName,
		Select: func() []string {
			return []string{"id", "name", "created_at", "updated_at"}
		},
		Scan: func(rows *sql.Rows) (*model.Person, *errs.Error) {
			var model model.Person
			if err := rows.Scan(&model.ID, &model.Name, &model.CreatedAt, &model.UpdatedAt); err != nil {
				return nil, errs.From(err)
			}
			return &model, nil
		},
	})
}

func (r *PersonPgStore) All(c context.Context, props *dbDto.List) ([]*model.Person, *errs.Error) {
	return function.All(c, &function.AllProps[model.Person]{
		DB:        r.db,
		Props:     props,
		TableName: r.tableName,
		Select: func() []string {
			return []string{"id", "name", "created_at", "updated_at"}
		},
		Scan: func(rows *sql.Rows) (*model.Person, *errs.Error) {
			var model model.Person
			if err := rows.Scan(&model.ID, &model.Name, &model.CreatedAt, &model.UpdatedAt); err != nil {
				return nil, errs.From(err)
			}
			return &model, nil
		},
	})
}

func (r *PersonPgStore) Find(c context.Context, id int) (*model.Person, *errs.Error) {
	type findResult struct {
		person *model.Person
		user   *userModel.User
		err    *errs.Error
	}

    resultCh := make(chan findResult, 2)
    wg := sync.WaitGroup{}
    wg.Add(2)

    //! Fetch person
    go func() {
        defer wg.Done()

        query := function.QueryBuilder().
            Select("id", "name", "created_at", "updated_at").
            From(r.tableName).
            Where(sq.Eq{"id": id})

        sqlStr, args, err := query.ToSql()
        if err != nil {
            resultCh <- findResult{err: errs.From(err)}
            return
        }

        person := &model.Person{}
        if err := r.db.QueryRowContext(c, sqlStr, args...).Scan(
            &person.ID, &person.Name, &person.CreatedAt, &person.UpdatedAt,
        ); err != nil {
            resultCh <- findResult{err: errs.From(err)}
            return
        }

        resultCh <- findResult{person: person}
    }()

    //! Fetch associated user
    go func() {
        defer wg.Done()

        query := function.QueryBuilder().
            Select("email", "created_at", "updated_at").
            From("users").
            Where(sq.Eq{"id": id})

        sqlStr, args, err := query.ToSql()
        if err != nil {
            resultCh <- findResult{err: errs.From(err)}
            return
        }

        user := &userModel.User{}
        if err := r.db.QueryRowContext(c, sqlStr, args...).Scan(&user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
            resultCh <- findResult{err: errs.From(err)}
            return
        }

        resultCh <- findResult{user: user}
    }()

    go func() {
        wg.Wait()
        close(resultCh)
    }()

    var (
        person *model.Person
        user   *userModel.User
    )

    for res := range resultCh {
        if res.err != nil {
            return nil, res.err
        }
        if res.person != nil {
            person = res.person
        }
        if res.user != nil {
            user = res.user
        }
    }

    if person != nil {
        person.User = user
    }

    return person, nil
}

func (r *PersonPgStore) Create(c context.Context, model *model.Person) (*model.Person, *errs.Error) {
	err := function.Create(c, &function.CreateProps{
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

func (r *PersonPgStore) Update(c context.Context, id int, model *model.Person) (*model.Person, *errs.Error) {
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

func (r *PersonPgStore) Delete(c context.Context, id int) *errs.Error {
	return function.Delete(c, &function.DeleteProps{
        DB:        r.db,
        TableName: r.tableName,
        Where: func (qb sq.DeleteBuilder) sq.DeleteBuilder {
            return qb.Where(sq.Eq{"id": id})
        },
    })
}
