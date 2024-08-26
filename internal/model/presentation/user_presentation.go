package presentation

type User struct {
	Id          int
	UserName    string
	EncodedName string
	Email       string
	Pass        string
	ProfilePic  string
	Followers   int
	Following   int
	HashPass    string
	Active      int
}

type Follows struct {
	FollowerId int
	FollowedId int
}
