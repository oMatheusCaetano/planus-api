package service

import (
	"context"
	"time"

	"github.com/omatheuscaetano/planus-api/internal/model"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

type UserService struct {}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Find(c context.Context, id string) (*model.User, *errs.Error) {
	return &model.User{
		ID:        "1",
		Name:      "John Doe",
		Username:  "johndoe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
