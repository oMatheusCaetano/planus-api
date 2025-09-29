package service

import (
	"context"

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

func (s *UserService) Delete(c context.Context, id model.ID) *errs.Error {
	err := s.r.Delete(c, id)
	if err != nil {
		return err
	}
	return nil
}
