package auth

import "github.com/gorilla/sessions"

type SessionConfig struct {
	Pass string `yaml:"pass"`
}

type User struct {
	Id         int
	UserName   string
	UserEmail  string
	UserPass   string
	HashedPass string
}

type Session struct {
	Id            string
	Authenticated bool
	User          User
}

var (
	SessionStore *sessions.CookieStore
)
