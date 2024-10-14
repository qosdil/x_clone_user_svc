package databases

import (
	"context"
	"errors"
	"time"

	"github.com/qosdil/x_clone_user_svc/model"
	"github.com/qosdil/x_clone_user_svc/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	coll *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) (repository.Repository, error) {
	coll := db.Collection("users")

	// Create the unique index on "username" field
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	return &mongoRepository{
		coll: coll,
	}, nil
}

type mongoRepoUser struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty"`
	Username  string              `bson:"username"`
	Password  string              `bson:"password"`
	CreatedAt primitive.Timestamp `bson:"created_at"`
}

func (r *mongoRepository) Create(ctx context.Context, user model.User) (model.User, error) {
	user.CreatedAt = uint32(time.Now().Unix())
	repoUser := mongoRepoUser{
		Username:  user.Username,
		Password:  user.Password,
		CreatedAt: primitive.Timestamp{T: user.CreatedAt},
	}
	result, err := r.coll.InsertOne(ctx, repoUser)
	if mongo.IsDuplicateKeyError(err) {
		return model.User{}, model.ErrCodeUsernameNotAvailable
	}
	if err != nil {
		return model.User{}, err
	}
	insertedID, _ := result.InsertedID.(primitive.ObjectID)
	user.ID = insertedID.Hex()
	return user, nil
}

func (r *mongoRepository) Find(ctx context.Context) (resp []model.SecureUser, err error) {
	projection := bson.D{
		{Key: "username", Value: 1},
		{Key: "created_at", Value: 1},
	}
	cursor, err := r.coll.Find(ctx, bson.M{}, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var repoUser mongoRepoUser
	for cursor.Next(ctx) {
		if err := cursor.Decode(&repoUser); err != nil {
			return nil, err
		}
		resp = append(resp, model.SecureUser{
			ID:        repoUser.ID.Hex(),
			Username:  repoUser.Username,
			CreatedAt: repoUser.CreatedAt.T,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *mongoRepository) FirstByUsername(ctx context.Context, username string) (user model.User, err error) {
	var repoUser mongoRepoUser
	err = r.coll.FindOne(ctx, bson.M{"username": username}).Decode(&repoUser)
	if err == mongo.ErrNoDocuments {
		return model.User{}, errors.New("not found")
	}
	if err != nil {
		return model.User{}, err
	}
	user.ID = repoUser.ID.Hex()
	user.Username = username
	user.Password = repoUser.Password
	user.CreatedAt = repoUser.CreatedAt.T
	return user, nil
}
