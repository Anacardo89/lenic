package presentation

type Session struct {
	Authenticated bool
	User          User
	Notifs        []Notification
	DMs           []Conversation
}
