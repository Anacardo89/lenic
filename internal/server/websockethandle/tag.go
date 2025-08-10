package websockethandle

import (
	"encoding/json"
	"strconv"

	"github.com/Anacardo89/lenic/internal/handlers/data/orm"
	"github.com/Anacardo89/lenic/internal/model/database"
	"github.com/Anacardo89/lenic/internal/model/mapper"
	"github.com/Anacardo89/lenic/pkg/logger"
	"github.com/Anacardo89/lenic/pkg/wsocket"
)

func (h *WSHandler) HandlePostTag(msg wsocket.Message, tagged_user string) {

	dbuser, err := orm.Da.GetUserByName(tagged_user)
	if err != nil {
		logger.Error.Println("Could not get user: ", err)
		return
	}

	fromuser, err := orm.Da.GetUserByName(msg.FromUserName)
	if err != nil {
		logger.Error.Println("Could not get from user: ", err)
		return
	}

	if dbuser.Id == fromuser.Id {
		return
	}

	u := mapper.UserNotif(dbuser)
	from_u := mapper.UserNotif(fromuser)

	n := &database.Notification{
		UserID:     dbuser.Id,
		FromUserId: fromuser.Id,
		NotifType:  msg.Type,
		NotifMsg:   msg.Msg,
		ResourceId: msg.ResourceId,
		ParentId:   "",
	}

	res, err := orm.Da.CreateNotification(n)
	if err != nil {
		logger.Error.Println("Could not create notification: ", err)
		return
	}
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		logger.Error.Println("Could not get notification Id: ", err)
		return
	}

	dbnotif, err := orm.Da.GetNotificationById(int(lastInsertID))
	if err != nil {
		logger.Error.Println("Could not get notification: ", err)
		return
	}
	notif := mapper.Notification(dbnotif, *u, *from_u)
	notif.ParentId = ""

	data, err := json.Marshal(notif)
	if err != nil {
		logger.Error.Println("Could not marshal JSON: ", err)
		return
	}

	if wsocket.WSConnMan.IsConnected(dbuser.UserName) {
		wsocket.WSConnMan.SendMessage(u.UserName, data)
	}
}

func (h *WSHandler) HandleCommentTag(msg Message, tagged_user string) {
	comment_id, err := strconv.Atoi(msg.ResourceId)
	if err != nil {
		logger.Error.Printf("Could not convert %s to int: %s\n", msg.ResourceId, err)
		return
	}
	c, err := orm.Da.GetCommentById(comment_id)
	if err != nil {
		logger.Error.Println("Could not get comment: ", err)
		return
	}
	dbuser, err := orm.Da.GetUserByName(tagged_user)
	if err != nil {
		logger.Error.Println("Could not get user: ", err)
		return
	}

	fromuser, err := orm.Da.GetUserByName(msg.FromUserName)
	if err != nil {
		logger.Error.Println("Could not get from user: ", err)
		return
	}

	if dbuser.Id == fromuser.Id {
		return
	}

	u := mapper.UserNotif(dbuser)
	from_u := mapper.UserNotif(fromuser)

	n := &database.Notification{
		UserID:     dbuser.Id,
		FromUserId: fromuser.Id,
		NotifType:  msg.Type,
		NotifMsg:   msg.Msg,
		ResourceId: msg.ResourceId,
		ParentId:   msg.ParentId,
	}

	res, err := orm.Da.CreateNotification(n)
	if err != nil {
		logger.Error.Println("Could not create notification: ", err)
		return
	}
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		logger.Error.Println("Could not get notification Id: ", err)
		return
	}

	dbnotif, err := orm.Da.GetNotificationById(int(lastInsertID))
	if err != nil {
		logger.Error.Println("Could not get notification: ", err)
		return
	}
	notif := mapper.Notification(dbnotif, *u, *from_u)
	notif.ParentId = c.PostGUID

	data, err := json.Marshal(notif)
	if err != nil {
		logger.Error.Println("Could not marshal JSON: ", err)
		return
	}

	if wsocket.WSConnMan.IsConnected(dbuser.UserName) {
		wsocket.WSConnMan.SendMessage(u.UserName, data)
	}
}
