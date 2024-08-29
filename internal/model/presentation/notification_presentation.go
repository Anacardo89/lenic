package presentation

type Notification struct {
	Id         int    `json:"id"`
	User       string `json:"user"`
	FromUser   string `json:"fromuser"`
	NotifType  string `json:"type"`
	NotifMsg   string `json:"msg"`
	ResourceId string `json:"resouce_id"`
	ParentId   string `json:"parent_id"`
	IsRead     bool   `json:"is_read"`
}
