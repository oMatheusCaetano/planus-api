package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"github.com/omatheuscaetano/planus-api/pkg/app"
	"github.com/omatheuscaetano/planus-api/pkg/db/dto"
)

var (
	con  *sql.DB
	once sync.Once
)

func GetDB() *sql.DB {
	return con
}

func Init() {
	once.Do(func() {
		var err error
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
            app.DBHost(), app.DBUser(), app.DBPassword(), app.DBName(), app.DBPort(),
        )
		con, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Fatal(err)
		}
		if err := con.Ping(); err != nil {
			log.Fatal(err)
		}
	})
}

func WhereLogic(query sq.SelectBuilder, block []*dto.WhereLogicBlock) sq.SelectBuilder {
	if len(block) == 0 {
		return query
	}

	sqlizer := buildWhereSqlizer(block)
	if sqlizer != nil {
		query = query.Where(sqlizer)
	}

	return query
}

func buildWhereSqlizer(block []*dto.WhereLogicBlock) sq.Sqlizer {
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

