package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/omatheuscaetano/planus-api/pkg/env"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

type ID string
type Permission string


type Model struct {
	ID          ID           `bson:"_id"         json:"id"`
	CreatedAt   time.Time    `bson:"created_at"  json:"created_at"`
	UpdatedAt   time.Time    `bson:"updated_at"  json:"updated_at"`
}

func newModel() Model {
	now := time.Now()

	return Model{
		ID:        NewID(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (m *Model) SetID(id ID) {
	m.ID = id
}

func (m *Model) OnRead() *errs.Error {
	m.CreatedAt = ToLocalTimezone(m.CreatedAt)
	m.UpdatedAt = ToLocalTimezone(m.UpdatedAt)
	return nil
}


func (id ID) String() string {
	return string(id)
}

func NewID() ID {
	res, _ := uuid.NewV7()
	return ID(res.String())
}

func ToLocalTimezone(t time.Time) time.Time {
	tz := env.Timezone()
	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic("failed to load timezone: " + tz + " error: " + err.Error())
	}

	return t.In(loc)
}
