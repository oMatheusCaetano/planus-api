package model

import (
	"time"
)

type Company struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    CNPJ      *string   `json:"cnpj"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
