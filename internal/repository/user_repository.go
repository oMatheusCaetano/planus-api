package repository

import "github.com/omatheuscaetano/planus-api/internal/model"

type UserRepository interface {
	Repository[model.User]
}
