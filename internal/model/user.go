package model

import (
	"github.com/omatheuscaetano/planus-api/pkg/errs"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
    Model               `bson:",inline"`
	Name      string    `bson:"name"       json:"name"`
	Username  string    `bson:"username"   json:"username"`
	Password  string    `bson:"password"   json:"-"`
}

func NewUser(name, username, rawPassword string) (*User, *errs.Error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.New(500, "Failed to hash password")
	}

    return &User{
        Model:     newModel(),
        Name:      name,
        Username:  username,
        Password:  string(hashedPassword),
    }, nil
}
