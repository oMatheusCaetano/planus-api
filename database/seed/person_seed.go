package seed

import (
	"context"
	"time"

	"github.com/jaswdr/faker/v2"
	"github.com/omatheuscaetano/planus-api/internal/person/model"
	"github.com/omatheuscaetano/planus-api/internal/person/store"
)

type PersonSeed struct {
    store store.PersonStore
}

func NewPersonSeed(store store.PersonStore) *PersonSeed {
    return &PersonSeed{store: store}
}

func (s *PersonSeed) Generate(ctx context.Context, quantity int) error {
    fake := faker.New()

    for i := 0; i < quantity; i++ {
        _, err := s.store.Create(ctx, &model.Person{
            Name:      fake.Person().Name(),
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        })
        if err != nil {
            return err
        }
    }
    return nil
}
