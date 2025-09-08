package services

import (
	"github.com/omatheuscaetano/planus-api/internal/company/models"
	"github.com/omatheuscaetano/planus-api/internal/company/repositories"
)

type CompanyService interface {
	Find(id int64) (*models.Company, error)
}

type companyService struct {
	repo repositories.CompanyRepository
}

func NewCompanyService(repo repositories.CompanyRepository) CompanyService {
	return &companyService{repo: repo}
}

func (s *companyService) Find(id int64) (*models.Company, error) {
	return s.repo.Find(id)
}
