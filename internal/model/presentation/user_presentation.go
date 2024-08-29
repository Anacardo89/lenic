package presentation

type User struct {
	Id          int    `json:"id"`
	UserName    string `json:"username"`
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
