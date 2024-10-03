package x_clone_user_srv

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Create(ctx context.Context, user User) (User, error)
	Find(ctx context.Context) (users []UserResponse, err error)
}

type mongoRepository struct {
	coll *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) Repository {
	return &mongoRepository{
		coll: db.Collection("users"),
	}
}

func (r *mongoRepository) Create(ctx context.Context, user User) (User, error) {
	user.CreatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}
	result, err := r.coll.InsertOne(ctx, user)
	if err != nil {
		return user, err
	}
	insertedID, _ := result.InsertedID.(primitive.ObjectID)
	user.ID = insertedID
	return user, nil
}

func (r *mongoRepository) Find(ctx context.Context) (userResponses []UserResponse, err error) {
	projection := bson.D{
		{Key: "username", Value: 1},
		{Key: "created_at", Value: 1},
	}
	cursor, err := r.coll.Find(ctx, bson.M{}, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var user User
	for cursor.Next(ctx) {
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		userResponses = append(userResponses, UserResponse{
			ID:        user.ID.Hex(),
			Username:  user.Username,
			CreatedAt: user.CreatedAt.T,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return userResponses, nil
}
