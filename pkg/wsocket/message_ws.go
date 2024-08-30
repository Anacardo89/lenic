package wsocket

type Message struct {
	FromUserName string `json:"from_username"`
	Type         string `json:"type"`
	Msg          string `json:"msg"`
	ResourceId   string `json:"resource_id"`
	ParentId     string `json:"parent_id"`
}
