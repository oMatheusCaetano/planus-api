package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/omatheuscaetano/planus-api/internal/shared/errs"
)

//! READ

type SortBy struct {
    Key       string `json:"key"`
    Direction string `json:"direction"`
}

type Where struct {
	Key      string      `json:"key"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type WhereLogicBlock struct {
	Operator  string      `json:"operator"` // "and" or "or"
	Condition interface{} `json:"condition"` // Where or []LogicBlock
}

type SelectQb struct {
    tableName    string
    selectFields []string
    limit        int
    offset       int
    where        []WhereLogicBlock
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

func (s *SelectQb) WhereFromLogicBlock(blocks []WhereLogicBlock) *SelectQb {
    s.where = append(s.where, blocks...)
    return s
}

func (s *SelectQb) Where(column string, operator string, value any) *SelectQb {
    s.where = append(s.where, WhereLogicBlock{
        Operator:  "and",
        Condition: Where{Key: column, Operator: operator, Value: value},
    })
    return s
}

func (s *SelectQb) Or(column string, operator string, value any) *SelectQb {
    s.where = append(s.where, WhereLogicBlock{
        Operator:  "or",
        Condition: Where{Key: column, Operator: operator, Value: value},
    })
    return s
}


func (s *SelectQb) WhereSub(callback func (subQueryBuilder *SelectQb) *SelectQb) *SelectQb {
    subQb := callback(&SelectQb{})
    s.where = append(s.where, WhereLogicBlock{
        Operator:  "and",
        Condition: subQb.where,
    })
    return s
}

func (s *SelectQb) OrSub(callback func (subQueryBuilder *SelectQb) *SelectQb) *SelectQb {
    subQb := callback(&SelectQb{})
    s.where = append(s.where, WhereLogicBlock{
        Operator:  "or",
        Condition: subQb.where,
    })
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
        where:        append([]WhereLogicBlock{}, s.where...),
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
        var processCondition func(cond interface{}) (string, error)
        processCondition = func(cond interface{}) (string, error) {
            switch c := cond.(type) {
            case Where:
                c.Operator, c.Value = mapOperatorsToSql(c.Operator, c.Value)
                args = append(args, c.Value)
                return fmt.Sprintf("%s %s $%d", c.Key, c.Operator, len(args)), nil
            case WhereLogicBlock:
                if logicCond, ok := c.Condition.([]WhereLogicBlock); ok {
                    var subConditions []string
                    for _, subCond := range logicCond {
                        subCondStr, err := processCondition(subCond)
                        if err != nil {
                            return "", err
                        }
                        op := "AND"
                        if subCond.Operator == "or" {
                            op = "OR"
                        }
                        subConditions = append(subConditions, op + " " + subCondStr)
                    }
                    return "(" + strings.Join(subConditions, " ") + ")", nil
                } else {
                    return processCondition(c.Condition)
                }
            default:
                return "", fmt.Errorf("invalid condition type")
            }
        }

        var conditions []string

        for _, block := range s.where {
            condStr, err := processCondition(block)
            if err != nil {
                // Handle error appropriately, here we just return an empty query and args
                return "", nil
            }
            op := "AND"
            if block.Operator == "or" {
                op = "OR"
            }
            conditions = append(conditions, op + " " + condStr)
        }

        q := (" WHERE " + strings.Join(conditions, " "))
        q = strings.ReplaceAll(q, "WHERE AND", "WHERE")
        q = strings.ReplaceAll(q, "AND (AND", "AND (")
        q = strings.ReplaceAll(q, "OR (AND", "OR (")
        query += q
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


    log.Printf("\n\n\n\n\n\n")
    log.Printf(query)
    log.Println(args)
    log.Printf("\n\n\n\n\n\n")
    return query, args
}

func mapOperatorsToSql(operator string, value any) (string, any) {
    switch operator {
        case "contain":
        case "contains":
            return "ILIKE", fmt.Sprintf("%%%v%%", value)
        case "startWith":
        case "startsWith":
            return "ILIKE", fmt.Sprintf("%v%%", value)
        case "endWith":
        case "endsWith":
            return "ILIKE", fmt.Sprintf("%%%v", value)
        default:
            return operator, value
    }

    return operator, value
}




//! INSERT
type InsertQb struct {
    tableName    string
    columns     []string
    values      []any
    returning   []string
}

func Insert(tableName string) *InsertQb {
    return &InsertQb{tableName: tableName}
}

func (q *InsertQb) Values(values map[string]any) *InsertQb {
    for col, val := range values {
        q.columns = append(q.columns, col)
        q.values = append(q.values, val)
    }
    return q
}

func (q *InsertQb) Returning(columns ...string) *InsertQb {
    q.returning = columns
    return q
}

func (q *InsertQb) ToSql() (string, []any) {
    var args []any
    var placeholders []string

    for i, val := range q.values {
        args = append(args, val)
        placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
    }

    query := fmt.Sprintf(
        "INSERT INTO %s (%s) VALUES (%s)",
        q.tableName,
        strings.Join(q.columns, ", "),
        strings.Join(placeholders, ", "),
    )

    if len(q.returning) > 0 {
        query += fmt.Sprintf(" RETURNING %s", strings.Join(q.returning, ", "))
    }

    return query, args
}

func (q *InsertQb) Scan(dest ...any) *errs.Error {
    query, args := q.ToSql()
    err := con.QueryRow(query, args...).Scan(dest...)
    if err != nil {
        return errs.From(err)
    }
    return nil
}

func (q *InsertQb) Run() (*sql.Rows, *errs.Error) {
    query, args := q.ToSql()
    rows, err := con.Query(query, args...)
    if err != nil {
        return nil, errs.From(err)
    }
    return rows, nil
}


//! UPDATE
type UpdateQb struct {
    tableName string
    sets      map[string]any
    where     []WhereLogicBlock
    returning []string
}

func Update(tableName string) *UpdateQb {
    return &UpdateQb{
        tableName: tableName,
        sets:      make(map[string]any),
    }
}

func (q *UpdateQb) Set(column string, value any) *UpdateQb {
    q.sets[column] = value
    return q
}

func (q *UpdateQb) WhereFromLogicBlock(blocks []WhereLogicBlock) *UpdateQb {
    q.where = append(q.where, blocks...)
    return q
}

func (q *UpdateQb) Where(column string, operator string, value any) *UpdateQb {
    q.where = append(q.where, WhereLogicBlock{
        Operator:  "and",
        Condition: Where{Key: column, Operator: operator, Value: value},
    })
    return q
}

func (q *UpdateQb) Or(column string, operator string, value any) *UpdateQb {
    q.where = append(q.where, WhereLogicBlock{
        Operator:  "or",
        Condition: Where{Key: column, Operator: operator, Value: value},
    })
    return q
}

func (q *UpdateQb) WhereSub(callback func (subQueryBuilder *SelectQb) *SelectQb) *UpdateQb {
    subQb := callback(&SelectQb{})
    q.where = append(q.where, WhereLogicBlock{
        Operator:  "and",
        Condition: subQb.where,
    })
    return q
}

func (q *UpdateQb) OrSub(callback func (subQueryBuilder *SelectQb) *SelectQb) *UpdateQb {
    subQb := callback(&SelectQb{})
    q.where = append(q.where, WhereLogicBlock{
        Operator:  "or",
        Condition: subQb.where,
    })
    return q
}

func (q *UpdateQb) Returning(columns ...string) *UpdateQb {
    q.returning = columns
    return q
}

func (q *UpdateQb) ToSql() (string, []any) {
    var args []any
    var setClauses []string

    i := 1
    for col, val := range q.sets {
        args = append(args, val)
        setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, i))
        i++
    }

    query := fmt.Sprintf("UPDATE %s SET %s", q.tableName, strings.Join(setClauses, ", "))

    if len(q.where) > 0 {
        var processCondition func(cond interface{}) (string, error)
        processCondition = func(cond interface{}) (string, error) {
            switch c := cond.(type) {
            case Where:
                c.Operator, c.Value = mapOperatorsToSql(c.Operator, c.Value)
                args = append(args, c.Value)
                return fmt.Sprintf("%s %s $%d", c.Key, c.Operator, len(args)), nil
            case WhereLogicBlock:
                if logicCond, ok := c.Condition.([]WhereLogicBlock); ok {
                    var subConditions []string
                    for _, subCond := range logicCond {
                        subCondStr, err := processCondition(subCond)
                        if err != nil {
                            return "", err
                        }
                        op := "AND"
                        if subCond.Operator == "or" {
                            op = "OR"
                        }
                        subConditions = append(subConditions, op + " " + subCondStr)
                    }
                    return "(" + strings.Join(subConditions, " ") + ")", nil
                } else {
                    return processCondition(c.Condition)
                }
            default:
                return "", fmt.Errorf("invalid condition type")
            }
        }

        var conditions []string

        for _, block := range q.where {
            condStr, err := processCondition(block)
            if err != nil {
                // Handle error appropriately, here we just return an empty query and args
                return "", nil
            }
            op := "AND"
            if block.Operator == "or" {
                op = "OR"
            }
            conditions = append(conditions, op + " " + condStr)
        }

        condQuery := (" WHERE " + strings.Join(conditions, " "))
        condQuery = strings.ReplaceAll(condQuery, "WHERE AND", "WHERE")
        condQuery = strings.ReplaceAll(condQuery, "AND (AND", "AND (")
        condQuery = strings.ReplaceAll(condQuery, "OR (AND", "OR (")
        query += condQuery
    }

    if len(q.returning) > 0 {
        query += fmt.Sprintf(" RETURNING %s", strings.Join(q.returning, ", "))
    }

    return query, args
}

func (q *UpdateQb) Scan(dest ...any) *errs.Error {
    query, args := q.ToSql()
    err := con.QueryRow(query, args...).Scan(dest...)
    if err != nil {
        return errs.From(err)
    }
    return nil
}

func (q *UpdateQb) Run() (*sql.Rows, *errs.Error) {
    query, args := q.ToSql()
    rows, err := con.Query(query, args...)
    if err != nil {
        return nil, errs.From(err)
    }
    return rows, nil
}
