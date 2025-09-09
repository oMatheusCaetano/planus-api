package models

import "time"

type Company struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name" form:"name" binding:"required,min=3,max=255,unique=companies.name"`
    CNPJ      *string   `json:"cnpj" form:"cnpj" binding:"omitempty,len=14,cnpj,unique=companies.cnpj"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
