package router

import (
	"context"
	"testing"

	"github.com/omatheuscaetano/planus-api/internal/dto"
	"github.com/omatheuscaetano/planus-api/internal/model"
)

func TestFind(t *testing.T) {
	tester := NewRouterTester(t)

	err := tester.mongo.Collection("users").Drop(context.Background())
	if err != nil {
		t.Fatalf("failed to drop users collection: %v", err)
	}

	t.Run("Should return 404 if user does not exist", func(t *testing.T) {
		tester.
			Get("/user/111222333444555666777888999000").
			AssertStatus(404).
			AssertBodyValue("message", "Recurso não encontrado")
	})

	t.Run("Should return the user with a 200 status code", func(t *testing.T) {
		expectedModel, er := model.NewUser("John Doe", "johndoe", "hashedpassword")
		if er != nil {
			t.Fatalf("failed to create user model: %v", er)
		}

		_, err = tester.mongo.Collection("users").InsertOne(context.Background(), expectedModel)
		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}

		tester.
			Get("/user/" + expectedModel.ID.String()).
			AssertStatus(200).
			AssertBodyValue("id", expectedModel.ID.String()).
			AssertBodyValue("name", expectedModel.Name).
			AssertBodyValue("username", expectedModel.Username).
			AssertBodyValueIsDateTime("created_at", expectedModel.CreatedAt).
			AssertBodyValueIsDateTime("updated_at", expectedModel.UpdatedAt)
	})
}

func TestDelete(t *testing.T) {
	tester := NewRouterTester(t)

	err := tester.mongo.Collection("users").Drop(context.Background())
	if err != nil {
		t.Fatalf("failed to drop users collection: %v", err)
	}

	t.Run("Should return 204 even if not exists", func(t *testing.T) {
		tester.
			Delete("/user/111222333444555666777888999000").
			AssertStatus(204)
	})

	t.Run("Should delete and return 204", func(t *testing.T) {
		expectedModel, er := model.NewUser("John Doe", "johndoe", "hashedpassword")
		if er != nil {
			t.Fatalf("failed to create user model: %v", er)
		}

		_, err = tester.mongo.Collection("users").InsertOne(context.Background(), expectedModel)
		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}

		tester.
			Delete("/user/" + expectedModel.ID.String()).
			AssertStatus(204)
	})
}

func TestCreate(t *testing.T) {
	tester := NewRouterTester(t)

	err := tester.mongo.Collection("users").Drop(context.Background())
	if err != nil {
		t.Fatalf("failed to drop users collection: %v", err)
	}

	t.Run("Should validate request and return 400", func(t *testing.T) {
		dto := dto.CreateUserDTO{
			Name:     "",
			Username: "johndoe",
			Password: "abc",
		}

		tester.
			Post("/user", dto).
			AssertStatus(400).
			AssertBodyValue("message", "Dados inválidos").
			AssertBodyValue("data.name", "Campo obrigatório").
			AssertBodyValue("data.password", "Precisa ter no mínimo 6 caracteres")
	})

	// t.Run("Should create the user and return 204", func(t *testing.T) {
	// 	dto := dto.CreateUserDTO{
	// 		Name:     "John Doe",
	// 		Username: "johndoe",
	// 		Password: "plaintextpassword",
	// 	}

	// 	tester.
	// 		Post("/user", dto).
	// 		AssertStatus(201).
	// 		AssertBodyValue("name", dto.Name).
	// 		AssertBodyValue("username", dto.Username)
	// })
}
