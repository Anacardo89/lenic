package presentation

type User struct {
	Id         int
	UserName   string
	UserEmail  string
	UserPass   string
	HashedPass string
	Active     int
}
