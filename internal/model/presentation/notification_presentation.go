package presentation

type Notification struct {
	Id         int
	User       User
	FromUser   User
	NotifType  string
	NotifMsg   string
	ResourceId string
	IsRead     bool
}
