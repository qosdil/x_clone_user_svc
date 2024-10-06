package model

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt uint32 `json:"created_at"`
}

type SecureUser struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	CreatedAt uint32 `json:"created_at"`
}
