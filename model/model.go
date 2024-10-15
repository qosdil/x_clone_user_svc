package model

const (
	ErrCodeUsernameNotAvailable = "username_not_available"
)

type Error struct {
	Code string
}

func (e Error) Error() string {
	switch e.Code {
	case ErrCodeUsernameNotAvailable:
		return "username not available"
	default:
		return "unknown user error"
	}
}

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
