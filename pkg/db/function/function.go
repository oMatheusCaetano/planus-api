package function

import (
	"context"
	"database/sql"
	"strings"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"github.com/omatheuscaetano/planus-api/pkg/db"
	"github.com/omatheuscaetano/planus-api/pkg/db/dto"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

type PaginateProps[T any] struct {
	DB        *sql.DB
	Props     *dto.Paginate
	TableName string
	Select    func() []string
	Scan      func(*sql.Rows) (*T, *errs.Error)
}

func Paginate[T any](c context.Context, p *PaginateProps[T]) (*dto.PaginatedData[T], *errs.Error) {
	type paginateResult struct {
		total  int
		list   []*T
		err    *errs.Error
	}

	if p.Props.Page == 0 {
		p.Props.Page = 1
	}

	if p.Props.PerPage == 0 {
		p.Props.PerPage = 10
	}

	resultCh := make(chan paginateResult, 2)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		countQuery := QueryBuilder().Select("COUNT(1) AS total").From(p.TableName)
		countQuery = db.WhereLogic(countQuery, p.Props.Where)

		countSqlStr, countArgs, err := countQuery.ToSql()

		if err != nil {
			resultCh <- paginateResult{err: errs.From(err)}
			return
		}

		var total int
		if err := p.DB.QueryRowContext(c, countSqlStr, countArgs...).Scan(&total); err != nil {
			resultCh <- paginateResult{err: errs.From(err)}
			return
		}

		resultCh <- paginateResult{total: total}
	}()

	go func() {
		defer wg.Done()

		query := QueryBuilder().
			Select(p.Select()...).
			From(p.TableName).
			Limit(uint64(p.Props.PerPage)).
			Offset(uint64((p.Props.Page - 1) * p.Props.PerPage))

		query = db.WhereLogic(query, p.Props.Where)

		if len(p.Props.SortBy) > 0 {
			for _, sort := range p.Props.SortBy {
				direction := strings.ToUpper(sort.Direction)
				if direction != "ASC" && direction != "DESC" {
					direction = "ASC"
				}
				query = query.OrderBy(sort.Key + " " + direction)
			}
		}

		sqlStr, args, err := query.ToSql()

		if err != nil {
			resultCh <- paginateResult{err: errs.From(err)}
			return
		}

		rows, err := p.DB.QueryContext(c, sqlStr, args...)

		if err != nil {
			resultCh <- paginateResult{err: errs.From(err)}
			return
		}

		defer rows.Close()

		var list []*T
		for rows.Next() {
			model, err := p.Scan(rows)

			if err != nil {
				resultCh <- paginateResult{err: errs.From(err)}
				return
			}

			list = append(list, model)
		}

		if err := rows.Err(); err != nil {
			resultCh <- paginateResult{err: errs.From(err)}
			return
		}

		resultCh <- paginateResult{list: list}
	}()

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var (
		total  int
		list   []*T
	)

	for res := range resultCh {
		if res.err != nil {
			return nil, res.err
		}
		if res.total != 0 {
			total = res.total
		}
		if res.list != nil {
			list = res.list
		}
	}

	return &dto.PaginatedData[T]{
		Data: list,
		Meta: &dto.PaginationMeta{
			Page:     p.Props.Page,
			PerPage:  p.Props.PerPage,
			LastPage: (total + p.Props.PerPage - 1) / p.Props.PerPage,
			Total:    total,
			SortBy:   p.Props.SortBy,
			Where:    p.Props.Where,
		},
	}, nil
}

type AllProps[T any] struct {
	DB        *sql.DB
	Props     *dto.List
	TableName string
	Select    func() []string
	Scan      func(*sql.Rows) (*T, *errs.Error)
}

func All[T any](c context.Context, p *AllProps[T]) ([]*T, *errs.Error) {
	qb := QueryBuilder()

	query := qb.Select(p.Select()...).From(p.TableName)
	query = db.WhereLogic(query, p.Props.Where)

	if len(p.Props.SortBy) > 0 {
		for _, sort := range p.Props.SortBy {

			direction := strings.ToUpper(sort.Direction)

			if direction != "ASC" && direction != "DESC" {
				direction = "ASC"
			}

			query = query.OrderBy(sort.Key + " " + direction)

		}
	}

	sqlStr, args, err := query.ToSql()

	if err != nil {
		return nil, errs.From(err)
	}

	rows, err := p.DB.QueryContext(c, sqlStr, args...)

	if err != nil {
		return nil, errs.From(err)
	}

	defer rows.Close()

	var list []*T
	for rows.Next() {
		model, err := p.Scan(rows)

		if err != nil {
			return nil, err
		}

		list = append(list, model)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.From(err)
	}

	return list, nil
}

type FindProps struct {
	DB        *sql.DB
	TableName string
	Select    func() []string
	Where     func (sq.SelectBuilder) sq.SelectBuilder
	Scan 	  func() []any
}

func Find(c context.Context, p *FindProps) *errs.Error {
	qb := QueryBuilder()

	query := qb.Select(p.Select()...).From(p.TableName)
	query = p.Where(query)

	sqlStr, args, err := query.ToSql()

	if err != nil {
		return errs.From(err)
	}

	if err := p.DB.QueryRowContext(c, sqlStr, args...).Scan(p.Scan()...); err != nil {
		return errs.From(err)
	}

	return nil
}

type CreateProps struct {
	DB        *sql.DB
	TableName string
	Columns   func() []string
	Values    func() []any
	Returning func() []string
	Scan 	  func() []any
}

func Create(c context.Context, p *CreateProps) *errs.Error {
	qb := QueryBuilder()

	query := qb.Insert(p.TableName).Columns(p.Columns()...).Values(p.Values()...)
	query = query.Suffix("RETURNING " + strings.Join(p.Returning(), ", "))

	sqlStr, args, err := query.ToSql()

	if err != nil {
		return errs.From(err)
	}

	if err := p.DB.QueryRowContext(c, sqlStr, args...).Scan(p.Scan()...); err != nil {
		return errs.From(err)
	}

	return nil
}

type UpdateProps struct {
	DB 	      *sql.DB
	TableName string
	Set	      func(sq.UpdateBuilder) sq.UpdateBuilder
	Where	  func(sq.UpdateBuilder) sq.UpdateBuilder
	Returning func() []string
	Scan 	  func() []any
}

func  Update(ctx context.Context, p *UpdateProps) *errs.Error {
	qb := QueryBuilder()

	query := qb.Update(p.TableName)
	query = p.Set(query)
	query = p.Where(query)
	query = query.Suffix("RETURNING " + strings.Join(p.Returning(), ", "))

	sqlStr, args, err := query.ToSql()

	if err != nil {
		return errs.From(err)
	}

	if err := p.DB.QueryRowContext(ctx, sqlStr, args...).Scan(p.Scan()...); err != nil {
		return errs.From(err)
	}

	return nil
}

type DeleteProps struct {
	DB        *sql.DB
	TableName string
	Where     func (sq.DeleteBuilder) sq.DeleteBuilder
}

func Delete(ctx context.Context, p *DeleteProps) *errs.Error {
	qb := QueryBuilder()

	sqlStr, args, err := p.Where(qb.Delete(p.TableName)).ToSql()

	if err != nil {
		return errs.From(err)
	}

	_, err = p.DB.ExecContext(ctx, sqlStr, args...)

	if err != nil {
		return errs.From(err)
	}

	return nil
}

func QueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
