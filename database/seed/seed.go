package seed

import (
	"context"
	"log"
	"time"

	"github.com/jaswdr/faker/v2"
	authModel "github.com/omatheuscaetano/planus-api/internal/auth/model"
	authStore "github.com/omatheuscaetano/planus-api/internal/auth/store"
	personModel "github.com/omatheuscaetano/planus-api/internal/person/model"
	personStore "github.com/omatheuscaetano/planus-api/internal/person/store"
	"golang.org/x/crypto/bcrypt"
)

type Seeder struct {
    personStore personStore.PersonStore
    authStore   authStore.AuthStore
}

func NewSeeder(personStore personStore.PersonStore, authStore authStore.AuthStore) *Seeder {
    return &Seeder{personStore: personStore, authStore: authStore}
}

func (s *Seeder) Generate(ctx context.Context, quantity int, withUser bool) {
    for i := 0; i < quantity; i++ {
        s.CreatePerson(ctx, withUser)
    }
}

func (s *Seeder) CreatePerson(ctx context.Context, withUser bool) {
    fake := faker.New()

    person, err := s.personStore.Create(ctx, &personModel.Person{
        Name:      fake.Person().Name(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    })
    if err != nil {
        log.Printf("failed to create person: %v", err)
        return
    }

    if !withUser { return }

    hashedPassword, e := bcrypt.GenerateFromPassword([]byte("123456789"), bcrypt.DefaultCost)
    if e != nil {
        log.Printf("failed to hash password: %v", e)
        return
    }

    _, erru := s.authStore.CreateUser(ctx, &authModel.User{
        ID:        person.ID,
        Email:     fake.Internet().Email(),
        Password:  string(hashedPassword),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    })

    if erru != nil {
        log.Printf("failed to create user: %v", erru)
        return
    }
}
