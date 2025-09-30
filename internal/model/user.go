package model

import (
	"slices"

	"github.com/omatheuscaetano/planus-api/pkg/errs"
	"golang.org/x/crypto/bcrypt"
)

const (
	PermissionUserCreate Permission = "user.create"
	PermissionUserRead   Permission = "user.read"
	PermissionUserUpdate Permission = "user.update"
	PermissionUserDelete Permission = "user.delete"
)

type User struct {
    Model                    `bson:",inline"`
	Permissions []Permission `bson:"permissions" json:"permissions"`
	Name        string       `bson:"name"       json:"name"`
	Username    string       `bson:"username"   json:"username"`
	Password    string       `bson:"password"   json:"-"`
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

func (u *User) Can(required ...Permission) bool {
	if len(required) == 0 {
		return true
	}

	if u.Permissions == nil {
		return false
	}

	if len(u.Permissions) == 0 {
		return false
	}

	for _, r := range required {
		if slices.Contains(u.Permissions, r) {
			return true
		}
	}

	return false
}
