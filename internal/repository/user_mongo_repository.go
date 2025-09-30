package repository

import (
	"context"
	"strings"

	"github.com/omatheuscaetano/planus-api/internal/model"
	mongo "github.com/omatheuscaetano/planus-api/pkg/db"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
	"go.mongodb.org/mongo-driver/bson"
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

func (r *UserMongoRepository) FindByUsername(ctx context.Context, username string) (*model.User, *errs.Error) {
	return r.findWhere(ctx, bson.M{ "username": strings.TrimSpace(strings.ToLower(username)) })
}
