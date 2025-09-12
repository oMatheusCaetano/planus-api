package model

import (
	"time"

	"github.com/omatheuscaetano/planus-api/internal/auth/model"
)

type Person struct {
    ID        int         `json:"id"`
    Name      string      `json:"name"`
    CreatedAt time.Time   `json:"created_at"`
    UpdatedAt time.Time   `json:"updated_at"`

    User      *model.User `json:"user"`
}
