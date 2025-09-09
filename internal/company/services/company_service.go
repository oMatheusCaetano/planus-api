package services

import (
	"github.com/omatheuscaetano/planus-api/internal/company/models"
	"github.com/omatheuscaetano/planus-api/internal/company/repositories"
	"github.com/omatheuscaetano/planus-api/internal/shared/errs"
)

type CompanyService interface {
	Find(id int) (*models.Company, *errs.Error)
}

type companyService struct {
	repo repositories.CompanyRepository
}

func NewCompanyService(repo repositories.CompanyRepository) CompanyService {
	return &companyService{repo: repo}
}

func (s *companyService) Find(id int) (*models.Company, *errs.Error) {
	return s.repo.Find(id)
}
