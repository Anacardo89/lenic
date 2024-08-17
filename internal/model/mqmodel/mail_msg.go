package mqmodel

type Register struct {
	Email string `json:"email"`
	User  string `json:"user"`
	Link  string `json:"link"`
}

type PasswordRecover struct {
	Email string `json:"email"`
	User  string `json:"user"`
	Link  string `json:"link"`
}
