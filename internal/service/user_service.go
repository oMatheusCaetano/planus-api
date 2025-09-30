package service

import (
	"context"

	"github.com/omatheuscaetano/planus-api/internal/dto"
	"github.com/omatheuscaetano/planus-api/internal/model"
	"github.com/omatheuscaetano/planus-api/internal/repository"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

type UserService struct {
	r repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{r: r}
}

func (s *UserService) Find(c context.Context, id model.ID) (*model.User, *errs.Error) {
	return s.r.Find(c, id)
}

func (s *UserService) FindByUsername(c context.Context, username string) (*model.User, *errs.Error) {
	return s.r.FindByUsername(c, username)
}

func (s *UserService) Delete(c context.Context, id model.ID) *errs.Error {
	return s.r.Delete(c, id)
}

func (s *UserService) Create(c context.Context, entity *dto.CreateUser) (*model.User, *errs.Error) {
	user, err := model.NewUser(
		entity.Name,
		entity.Username,
		entity.Password,
	)
	if err != nil {
		return nil, errs.From(err)
	}

	err = s.r.Create(c, user)
	if err != nil {
		return nil, errs.From(err)
	}
	return user, nil
}

