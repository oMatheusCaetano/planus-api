package store

import (
	"context"
	"database/sql"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/omatheuscaetano/planus-api/internal/person/dto"
	"github.com/omatheuscaetano/planus-api/internal/person/model"
	dbDto "github.com/omatheuscaetano/planus-api/pkg/db/dto"
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

func whereLogic(query sq.SelectBuilder, block []*dbDto.WhereLogicBlock) sq.SelectBuilder {
	if len(block) == 0 {
		return query
	}

	sqlizer := buildWhereSqlizer(block)
	if sqlizer != nil {
		query = query.Where(sqlizer)
	}

	return query
}

func buildWhereSqlizer(block []*dbDto.WhereLogicBlock) sq.Sqlizer {
	if len(block) == 0 {
		return nil
	}

	var parts []sq.Sqlizer

	for i, b := range block {
		var expr sq.Sqlizer
		if b.Condition != nil {
			switch strings.ToUpper(b.Condition.Operator) {
			case "=":
				expr = sq.Eq{b.Condition.Key: b.Condition.Value}
			case "!=":
				expr = sq.NotEq{b.Condition.Key: b.Condition.Value}
			case "<":
				expr = sq.Lt{b.Condition.Key: b.Condition.Value}
			case "<=":
				expr = sq.LtOrEq{b.Condition.Key: b.Condition.Value}
			case ">":
				expr = sq.Gt{b.Condition.Key: b.Condition.Value}
			case ">=":
				expr = sq.GtOrEq{b.Condition.Key: b.Condition.Value}
			case "LIKE":
				expr = sq.Like{b.Condition.Key: b.Condition.Value}
			case "ILIKE", "CONTAIN", "CONTAINS", "CONTAINING":
				expr = sq.ILike{b.Condition.Key: b.Condition.Value}
			case "STARTWITH", "STARTSWITH":
				expr = sq.ILike{b.Condition.Key: b.Condition.Value.(string) + "%"}
			case "ENDWITH", "ENDSWITH":
				expr = sq.ILike{b.Condition.Key: "%" + b.Condition.Value.(string)}
			case "IN":
				expr = sq.Eq{b.Condition.Key: b.Condition.Value}
			default:
				continue
			}
		} else if b.Sub != nil {
			expr = buildWhereSqlizer(b.Sub)
			if expr == nil {
				continue
			}
		} else {
			continue
		}

		if i == 0 {
			parts = append(parts, expr)
			continue
		}

		op := strings.ToLower(b.Operator)
		if op == "or" {
			prev := parts[len(parts)-1]
			if prevOr, ok := prev.(sq.Or); ok {
				parts[len(parts)-1] = sq.Or(append([]sq.Sqlizer(prevOr), expr))
			} else {
				parts[len(parts)-1] = sq.Or{prev, expr}
			}
		} else {
			parts = append(parts, expr)
		}
	}

	if len(parts) == 0 {
		return nil
	}
	if len(parts) == 1 {
		return parts[0]
	}
	return sq.And(parts)
}

func (r *PersonPgStore) Paginate(ctx context.Context, props *dto.PaginatePerson) (*dto.PaginatedPerson, *errs.Error) {
	if props.Page == 0 { props.Page = 1 }
	if props.PerPage == 0 { props.PerPage = 10 }

	var total int
	countQuery := r.psql.Select("COUNT(1) AS total").From(r.tableName)
	countQuery = whereLogic(countQuery, props.Where)

	countSqlStr, countArgs, err := countQuery.ToSql()
	if err != nil { return nil, errs.From(err) }
	if err := r.db.QueryRowContext(ctx, countSqlStr, countArgs...).Scan(&total); err != nil {
		return nil, errs.From(err)
	}

	query := r.psql.
		Select("id", "name", "created_at", "updated_at").
		From(r.tableName).
		Limit(uint64(props.PerPage)).
		Offset(uint64((props.Page - 1) * props.PerPage))
	query = whereLogic(query, props.Where)

	if len(props.SortBy) > 0 {
		for _, sort := range props.SortBy {
			direction := strings.ToUpper(sort.Direction)
			if direction != "ASC" && direction != "DESC" {
				direction = "ASC"
			}
			query = query.OrderBy(sort.Key + " " + direction)
		}
	} else {
		query = query.OrderBy("id DESC")
	}

	sqlStr, args, err := query.ToSql()
	if err != nil { return nil, errs.From(err) }

	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil { return nil, errs.From(err) }
	defer rows.Close()

	var people []*model.Person
	for rows.Next() {
		var person model.Person
		if err := rows.Scan(&person.ID, &person.Name, &person.CreatedAt, &person.UpdatedAt); err != nil {
			return nil, errs.From(err)
		}
		people = append(people, &person)
	}

	if err := rows.Err(); err != nil { return nil, errs.From(err) }

	return &dto.PaginatedPerson{
		Data: people,
		Meta: &dto.PaginatedPersonMeta{
			Page:     props.Page,
			PerPage:  props.PerPage,
			LastPage: (total + props.PerPage - 1) / props.PerPage,
			Total:    total,
			SortBy:   props.SortBy,
			Where:    props.Where,
		},
	}, nil
}

func (r *PersonPgStore) All(ctx context.Context, dto *dto.ListPerson) ([]*model.Person, *errs.Error) {
	query := r.psql.
		Select("id", "name", "created_at", "updated_at").
		From(r.tableName)
		
	query = whereLogic(query, dto.Where)

	if len(dto.SortBy) > 0 {
		for _, sort := range dto.SortBy {
			direction := strings.ToUpper(sort.Direction)
			if direction != "ASC" && direction != "DESC" {
				direction = "ASC"
			}
			query = query.OrderBy(sort.Key + " " + direction)
		}
	}  else {
		query = query.OrderBy("id DESC")
	}

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
