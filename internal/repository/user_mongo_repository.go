package repository

import (
	"context"

	"github.com/omatheuscaetano/planus-api/internal/model"
	mongo "github.com/omatheuscaetano/planus-api/pkg/db"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
	"go.mongodb.org/mongo-driver/bson"
)


type UserMongoRepository struct {
	collectionName string
	db             *mongo.MongoDB
}

func NewUserMongoRepository(db *mongo.MongoDB) UserRepository {
	return &UserMongoRepository{collectionName: "users", db: db}
}

func (r *UserMongoRepository) Find(ctx context.Context, id model.ID) (*model.User, *errs.Error) {
	collection := r.db.Collection(r.collectionName)

	var user model.User
	result := collection.FindOne(ctx, bson.M{"_id": id})

	if result.Err() != nil {
		return nil, errs.ResourceNotFound()
	}

	err := result.Decode(&user)
	if err != nil {
		return nil, errs.From(err)
	}

	user.CreatedAt = model.UTCToLocal(user.CreatedAt)
	user.UpdatedAt = model.UTCToLocal(user.UpdatedAt)
	return &user, nil
}
