package dto

import "github.com/omatheuscaetano/planus-api/pkg/db/dto"

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
