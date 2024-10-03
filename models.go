package x_clone_user_srv

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateRequest struct {
	User User
}

type CreateResponse struct {
	User UserResponse `json:"user"`
	Err  error        `json:"err"`
}

type User struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Username  string              `bson:"username" json:"username"`
	Password  string              `bson:"password" json:"password"`
	CreatedAt primitive.Timestamp `bson:"created_at" json:"created_at"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	CreatedAt uint32 `json:"created_at"`
}
