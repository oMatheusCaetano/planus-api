package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/omatheuscaetano/planus-api/pkg/env"
)

type ID string

func (id ID) String() string {
	return string(id)
}

func NewID() ID {
	res, _ := uuid.NewV7()
	return ID(res.String())
}

func UTCToLocal(t time.Time) time.Time {
	tz := env.Timezone()
	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic("failed to load timezone: " + tz + " error: " + err.Error())
	}

	return t.In(loc)
}
