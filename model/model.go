package model

import "errors"

var (
	ErrCodeUsernameNotAvailable = errors.New("username_not_available")

	// Errors map error codes to error messages
	Errors = map[string]string{
		ErrCodeUsernameNotAvailable.Error(): "username not available",
	}
)

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
