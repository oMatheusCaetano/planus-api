package dto

import (
	"github.com/omatheuscaetano/planus-api/internal/person/model"
	"github.com/omatheuscaetano/planus-api/pkg/db/dto"
)

type PaginatePerson struct {
    Page    int                    `json:"page"    binding:"omitempty,min=1"`
    PerPage int                    `json:"per_page" binding:"omitempty,min=1"`
    SortBy  []*dto.SortBy          `json:"sort_by" binding:"omitempty,dive"`
    Where   []*dto.WhereLogicBlock `json:"where"   binding:"omitempty,dive"`
}

type PaginatedPersonMeta struct {
    Page      int                    `json:"page"`
    PerPage   int                    `json:"per_page"`
    LastPage  int                    `json:"last_page"`
    Total     int                    `json:"total"`
    SortBy    []*dto.SortBy          `json:"sort_by"`
    Where     []*dto.WhereLogicBlock `json:"where"`
}

type PaginatedPerson struct {
    Data  []*model.Person `json:"data"`
    Meta  *PaginatedPersonMeta `json:"meta"`
}
    
type ListPerson struct {
    SortBy []*dto.SortBy          `json:"sort_by" binding:"omitempty,dive"`
    Where  []*dto.WhereLogicBlock `json:"where"   binding:"omitempty,dive"`
}

type CreatePerson struct {
    Name  string `json:"name" binding:"required,filled,min=2,max=255"`
}

type UpdatePerson struct {
    Name  string `json:"name" binding:"omitempty,filled,min=2,max=255"`
}
