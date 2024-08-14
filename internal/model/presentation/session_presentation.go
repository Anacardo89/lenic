package presentation

type Session struct {
	Id            int
	SessionId     string
	Authenticated bool
	User          User
}
