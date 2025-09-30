package repository

import (
	"context"

	"github.com/omatheuscaetano/planus-api/internal/model"
	mongo "github.com/omatheuscaetano/planus-api/pkg/db"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
	"go.mongodb.org/mongo-driver/bson"
	driver "go.mongodb.org/mongo-driver/mongo"
)

type mongoRepositoryConfig[T any] struct {
	collectionName string
	db             *mongo.MongoDB
}

type mongoRepository[T any] struct {
	c *mongoRepositoryConfig[T]
}

func newMongoRepository[T any](c *mongoRepositoryConfig[T]) mongoRepository[T] {
	return mongoRepository[T]{c: c}
}

func (r *mongoRepository[T]) collection() *driver.Collection {
	return r.c.db.Collection(r.c.collectionName)
}

func (r *mongoRepository[T]) Find(ctx context.Context, id model.ID) (*T, *errs.Error) {
	result := new(T)
	err := r.collection().FindOne(ctx, bson.M{"_id": id.String()}).Decode(result)
	if err != nil {
		return nil, errs.From(err)
	}

	if hook, ok := any(result).(interface{ OnRead() *errs.Error }); ok {
		if err := hook.OnRead(); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (r *mongoRepository[T]) Delete(ctx context.Context, id model.ID) *errs.Error {
	_, err := r.collection().DeleteOne(ctx, bson.M{ "_id": id })
	if err != nil {
		return errs.From(err)
	}
	return nil
}

func (r *mongoRepository[T]) Create(ctx context.Context, entity *T) *errs.Error {
	_, err := r.collection().InsertOne(ctx, entity)
	if err != nil {
		return errs.From(err)
	}


	if hook, ok := any(entity).(interface{ OnCreate() *errs.Error }); ok {
		if err := hook.OnCreate(); err != nil {
			return err
		}
	}

	return nil
}
