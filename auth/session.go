package auth

type User struct {
	Id   int
	Name string
}

type Session struct {
	Id            string
	Authenticated bool
	User          User
}
