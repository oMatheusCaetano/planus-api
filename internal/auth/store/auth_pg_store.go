package store

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/omatheuscaetano/planus-api/internal/auth/model"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

type AuthPgStore struct {
	db         *sql.DB
	psql       sq.StatementBuilderType
	usersTable string
}

func NewAuthPgStore(db *sql.DB) *AuthPgStore {
	return &AuthPgStore{
		db:         db,
		usersTable: "users",
		psql:       sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *AuthPgStore) FindUserByEmail(c context.Context , email string) (*model.User, *errs.Error) {
	query := s.psql.
		Select("id", "email", "password", "created_at", "updated_at").
		From(s.usersTable).
		Where(sq.Eq{"email": email})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, errs.From(err)
	}

	var model model.User
	if err := s.db.QueryRowContext(c, sqlStr, args...).Scan(&model.ID, &model.Email, &model.Password, &model.CreatedAt, &model.UpdatedAt); err != nil {
		return nil, errs.From(err)
	}

	return &model, nil
}

func (s *AuthPgStore) CreateUser(c context.Context, dto *model.User) (*model.User, *errs.Error) {
    query := s.psql.
		Insert(s.usersTable).
		Columns("id", "email", "password", "created_at", "updated_at").
		Values(dto.ID, dto.Email, dto.Password, dto.CreatedAt, dto.UpdatedAt).
		Suffix("RETURNING id, email, password, created_at, updated_at")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, errs.From(err)
	}

    var model model.User
	if err := s.db.QueryRowContext(c, sqlStr, args...).Scan(&model.ID, &model.Email, &model.Password, &model.CreatedAt, &model.UpdatedAt); err != nil {
		return nil, errs.From(err)
	}

	return &model, nil
}
