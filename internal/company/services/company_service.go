package services

import (
	"time"

	"github.com/omatheuscaetano/planus-api/internal/company/models"
	"github.com/omatheuscaetano/planus-api/internal/company/repositories"
	"github.com/omatheuscaetano/planus-api/internal/shared/dto"
	"github.com/omatheuscaetano/planus-api/internal/shared/errs"
)

type CompanyService interface {
	Paginate(offset, limit int) (dto.Paginated[models.Company], *errs.Error)
	All() ([]models.Company, *errs.Error)
	Find(id int) (*models.Company, *errs.Error)
	Create(company *models.Company) *errs.Error
}

type companyService struct {
	repo repositories.CompanyRepository
}

func NewCompanyService(repo repositories.CompanyRepository) CompanyService {
	return &companyService{repo: repo}
}

func (s *companyService) Paginate(offset, limit int) (dto.Paginated[models.Company], *errs.Error) {
	return s.repo.Paginate(offset, limit)
}

func (s *companyService) All() ([]models.Company, *errs.Error) {
	return s.repo.All()
}

func (s *companyService) Find(id int) (*models.Company, *errs.Error) {
	return s.repo.Find(id)
}

func (s *companyService) Create(company *models.Company) *errs.Error {
	company.CreatedAt = time.Now()
	company.UpdatedAt = company.CreatedAt
	return s.repo.Create(company)
}