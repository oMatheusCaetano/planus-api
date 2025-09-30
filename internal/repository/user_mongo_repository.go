package repository

import (
	"github.com/omatheuscaetano/planus-api/internal/model"
	mongo "github.com/omatheuscaetano/planus-api/pkg/db"
)

type UserMongoRepository struct {
	mongoRepository[model.User]
}

func NewUserMongoRepository(db *mongo.MongoDB) UserRepository {
	return &UserMongoRepository{
		mongoRepository: newMongoRepository(&mongoRepositoryConfig[model.User]{
			collectionName: "users",
			db:             db,
		}),
	}
}
