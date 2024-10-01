package x_clone_user_srv

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
}

type mongoRepository struct {
	coll *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) Repository {
	return &mongoRepository{
		coll: db.Collection("users"),
	}
}
