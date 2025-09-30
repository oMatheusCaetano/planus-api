package router

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/omatheuscaetano/planus-api/internal/dto"
	"github.com/omatheuscaetano/planus-api/internal/model"
	"github.com/omatheuscaetano/planus-api/pkg/env"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestLogin(t *testing.T) {
	tester := NewRouterTester(t)

	tester.DropTable("users")

	t.Run("Should validate request and return 400", func(t *testing.T) {
		dto := dto.LoginData{ Username: "", Password: "" }

		tester.
			Post("/auth/login", dto).
			AssertStatus(400).
			AssertBodyValue("message", "Dados inválidos").
			AssertBodyValue("data.username", "Campo obrigatório").
			AssertBodyValue("data.password", "Campo obrigatório")
	})

	t.Run("Should create the user and return 201", func(t *testing.T) {
		rawPassword := "123456789"
		user, err := model.NewUser(
			"John Doe",
			"johndoe@email.com",
			rawPassword,
		)
		if err != nil {
			t.Fatalf("failed to create user model: %v", err)
		}

		insertedID := tester.InsertToTable("users", user)
		createdUser := tester.FindFromTable("users", map[string]any{
			"_id": insertedID,
		})

		tester.Post("/auth/login", dto.LoginData{
			Username: user.Username,
			Password: rawPassword,
		}).
			AssertStatus(200).
			AssertBodyKeyIsJWT("token").
			AssertBodyKeyIsNumeric("expires_in").
			AssertBodyValue("user.id", createdUser["_id"]).
			AssertBodyValue("user.name", createdUser["name"]).
			AssertBodyValue("user.username", createdUser["username"]).
			AssertBodyValueIsDateTime("user.created_at", createdUser["created_at"].(primitive.DateTime).Time()).
			AssertBodyValueIsDateTime("user.updated_at", createdUser["updated_at"].(primitive.DateTime).Time())

		token, er := jwt.Parse((*tester.Body())["token"].(string), func(token *jwt.Token) (interface{}, error) {
			return []byte(env.JWTSecret()), nil
		})

		if er != nil {
			t.Fatalf("failed to parse jwt token: %v", er)
		}

		if !token.Valid {
			t.Fatal("token is not valid")
		}

		claims := token.Claims.(jwt.MapClaims)
		tester.AssertEqual(claims["sub"], createdUser["_id"])
	})
}
