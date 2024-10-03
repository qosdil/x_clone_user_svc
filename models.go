package x_clone_user_srv

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateResponse struct {
	User UserSecureResponse `json:"user"`
	Err  error              `json:"err"`
}

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Username  string              `bson:"username" json:"username"`
	Password  string              `bson:"password" json:"password"`
	CreatedAt primitive.Timestamp `bson:"created_at" json:"created_at"`
}

type UserSecureResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	CreatedAt uint32 `json:"created_at"`
}
