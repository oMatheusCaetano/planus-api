package service

import (
	"context"
	"strings"
	"time"

	"github.com/omatheuscaetano/planus-api/internal/person/dto"
	"github.com/omatheuscaetano/planus-api/internal/person/model"
	"github.com/omatheuscaetano/planus-api/internal/person/store"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

type PersonService struct {
	store store.PersonStore
}

func NewPersonService(store store.PersonStore) *PersonService {
	return &PersonService{store: store}
}

func (s *PersonService) All(c context.Context, dto *dto.ListPerson) ([]*model.Person, *errs.Error) {
	return s.store.All(c, dto)
}

func (s *PersonService) Find(c context.Context, id int) (*model.Person, *errs.Error) {
	return s.store.Find(c, id)
}

func (s *PersonService) Create(c context.Context, dto *dto.CreatePerson) (*model.Person, *errs.Error) {
    return s.store.Create(c, &model.Person{
        Name:      strings.TrimSpace(dto.Name),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    })
}

func (s *PersonService) Update(c context.Context, id int, dto *dto.UpdatePerson) (*model.Person, *errs.Error) {
	model := &model.Person{
		Name:      strings.TrimSpace(dto.Name),
		UpdatedAt: time.Now(),
	}

	validToUpdate := false

	if (model.Name != "") {
		validToUpdate = true
}

	if !validToUpdate {
		return nil, errs.BadRequest("Sem campos para atualizar")
	}

    return s.store.Update(c, id, model)
}

func (s *PersonService) Delete(c context.Context, id int) *errs.Error {
	return s.store.Delete(c, id)
}
