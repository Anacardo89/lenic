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

type UserNotif struct {
	Id          int    `json:"id"`
	UserName    string `json:"username"`
	EncodedName string `json:"encoded"`
	ProfilePic  string `json:"profile_pic"`
}

type Follows struct {
	FollowerId int
	FollowedId int
	Status     int
}
