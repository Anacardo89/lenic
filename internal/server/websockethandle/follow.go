package websockethandle

import (
	"encoding/base64"
	"encoding/json"

	"github.com/Anacardo89/lenic/internal/handlers/data/orm"
	"github.com/Anacardo89/lenic/internal/model/database"
	"github.com/Anacardo89/lenic/internal/model/mapper"
	"github.com/Anacardo89/lenic/pkg/logger"
	"github.com/Anacardo89/lenic/pkg/wsocket"
)

func (h *WSHandler) handleFollowRequest(msg Message) {

	bytes, err := base64.URLEncoding.DecodeString(msg.ResourceId)
	if err != nil {
		logger.Error.Printf("Could not decode user %s: %s\n", msg.ResourceId, err)
		return
	}
	userName := string(bytes)

	dbuser, err := orm.Da.GetUserByName(userName)
	if err != nil {
		logger.Error.Println("Could not get post: ", err)
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
		UserID:     u.Id,
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

	wsocket.WSConnMan.SendMessage(u.UserName, data)
}

func (h *WSHandler) handleFollowAccept(msg Message) {

	bytes, err := base64.URLEncoding.DecodeString(msg.ResourceId)
	if err != nil {
		logger.Error.Printf("Could not decode user %s: %s\n", msg.ResourceId, err)
		return
	}
	userName := string(bytes)

	dbuser, err := orm.Da.GetUserByName(userName)
	if err != nil {
		logger.Error.Println("Could not get post: ", err)
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
		UserID:     u.Id,
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

	wsocket.WSConnMan.SendMessage(u.UserName, data)
}
