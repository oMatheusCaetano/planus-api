package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/omatheuscaetano/planus-api/internal/shared/errs"
)

type SelectQb struct {
    tableName    string
    selectFields []string
    limit        int
    offset       int
    where        []Where
    sortBy       []SortBy
}

func From(tableName string) *SelectQb {
    qb := &SelectQb{}
    qb.SetTableName(tableName)
    return qb
}

func (s *SelectQb) SetTableName(tableName string) *SelectQb {
    s.tableName = tableName
    return s
}

func (s *SelectQb) Select(dest ...string) *SelectQb {
    s.selectFields = dest
    return s
}

func (s *SelectQb) Limit(limit int) *SelectQb {
    s.limit = limit
    return s
}

func (s *SelectQb) Offset(offset int) *SelectQb {
    s.offset = offset
    return s
}

func (s *SelectQb) Where(column string, operator string, value any) *SelectQb {
    s.where = append(s.where, Where{Key: column, Operator: operator, Value: value})
    return s
}

func (s *SelectQb) SortBy(column string, direction string) *SelectQb {
    s.sortBy = append(s.sortBy, SortBy{Key: column, Direction: direction})
    return s
}

func (s *SelectQb) Count() (int, *errs.Error) {
    var count int
    query, args := s.ToSql()
    countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count_table", query)
    err := con.QueryRow(countQuery, args...).Scan(&count)
    if err != nil {
        return 0, errs.From(err)
    }
    return count, nil
}

func (s *SelectQb) Exists() (bool, *errs.Error) {
    query, args := s.ToSql()
    var exists bool
    err := con.QueryRow(query, args...).Scan(&exists)
    if err != nil {
        return false, errs.From(err)
    }
    return exists, nil
}

func (s *SelectQb) Duplicate() *SelectQb {
    return &SelectQb{
        tableName:    s.tableName,
        selectFields: append([]string{}, s.selectFields...),
        limit:        s.limit,
        offset:       s.offset,
        where:        append([]Where{}, s.where...),
        sortBy:       append([]SortBy{}, s.sortBy...),
    }
}

func (s *SelectQb) Scan(dest ...any) *errs.Error {
	query, args := s.ToSql()
	err := con.QueryRow(query, args...).Scan(dest...)
    if err != nil {
        return errs.From(err)
    }
    return nil
}

func (s *SelectQb) ScanMany(list any, eachScan func(rows *sql.Rows) error) *errs.Error {
    query, args := s.ToSql()
	rows, err := con.Query(query, args...)
	if err != nil {
		return errs.From(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := eachScan(rows)
		if err != nil {
			return errs.From(err)
		}
	}

	if err := rows.Err(); err != nil {
		return errs.From(err)
	}

	return nil
}

func (s *SelectQb) ToSql() (string, []any) {
    var args []any

    //! Select From
	if len(s.selectFields) == 0 {
		s.selectFields = []string{"*"}
	}
    query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(s.selectFields, ", "), s.tableName)

    //! Where
    if len(s.where) > 0 {
        var conditions []string
        for i, where := range s.where {
            conditions = append(conditions, fmt.Sprintf("%s %s $%d", where.Key, where.Operator, i+1))
            args = append(args, where.Value)
        }
        query += " WHERE " + strings.Join(conditions, " AND ")
    }

    //! Sort By
    if len(s.sortBy) > 0 {
        var sortClauses []string
        for _, sort := range s.sortBy {
            sortClauses = append(sortClauses, fmt.Sprintf("%s %s", sort.Key, sort.Direction))
        }
        query += " ORDER BY " + strings.Join(sortClauses, ", ")
    }

    //! Limit
    if s.limit > 0 {
        query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
        args = append(args, s.limit)
    }

    //! Offset
    if s.offset > 0 {
        query += fmt.Sprintf(" OFFSET $%d", len(args)+1)
        args = append(args, s.offset)
    }

    //! Returning
    // if len(qb.Returning) > 0 {
    //     query += " RETURNING " + strings.Join(qb.Returning, ", ")
    // }

    return query, args
}
