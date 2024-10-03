package x_clone_user_svc

import "go.mongodb.org/mongo-driver/bson/primitive"

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	User UserSecureResponse `json:"user"`
	Err  error              `json:"err"`
}

type User struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Username  string              `bson:"username" json:"username"`
	Password  string              `bson:"password" json:"password"`
	CreatedAt primitive.Timestamp `bson:"created_at" json:"created_at"`
}

type UserNotSecureResponse struct {
	User User  `json:"user"`
	Err  error `json:"err"`
}

type UserSecureResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	CreatedAt uint32 `json:"created_at"`
}
