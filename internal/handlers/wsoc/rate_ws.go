package wsoc

import (
	"encoding/json"
	"strconv"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mapper"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/Anacardo89/tpsi25_blog/pkg/wsocket"
)

func handleCommentRate(msg wsocket.Message) {
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
	dbuser, err := orm.Da.GetUserByID(c.AuthorId)
	if err != nil {
		logger.Error.Println("Could not get user: ", err)
	}
	u := mapper.UserNotif(dbuser)

	fromuser, err := orm.Da.GetUserByName(msg.FromUserName)
	if err != nil {
		logger.Error.Println("Could not get from user: ", err)
		return
	}
	from_u := mapper.UserNotif(fromuser)

	n := &database.Notification{
		UserID:     c.AuthorId,
		FromUserId: fromuser.Id,
		NotifType:  msg.Type,
		NotifMsg:   msg.Msg,
		ResourceId: msg.ResourceId,
	}

	res, err := orm.Da.CreateNotification(n)
	if err != nil {
		logger.Error.Println("Could not create notification: ", err)
		return
	}
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		logger.Error.Println("Could not get notification Id: ", err)
	}

	dbnotif, err := orm.Da.GetNotificationById(int(lastInsertID))
	if err != nil {
		logger.Error.Println("Could not get notification: ", err)
	}
	notif := mapper.Notification(dbnotif, *u, *from_u)
	notif.ParentId = c.PostGUID

	data, err := json.Marshal(notif)
	if err != nil {
		logger.Error.Println("Could not marshal JSON: ", err)
	}

	wsocket.WSConnMan.SendMessage(u.UserName, data)
}

func handlePostRate(msg wsocket.Message) {

	p, err := orm.Da.GetPostByGUID(msg.ResourceId)
	if err != nil {
		logger.Error.Println("Could not get post: ", err)
		return
	}
	dbuser, err := orm.Da.GetUserByID(p.AuthorId)
	if err != nil {
		logger.Error.Println("Could not get user: ", err)
	}
	u := mapper.UserNotif(dbuser)

	fromuser, err := orm.Da.GetUserByName(msg.FromUserName)
	if err != nil {
		logger.Error.Println("Could not get from user: ", err)
		return
	}
	from_u := mapper.UserNotif(fromuser)

	n := &database.Notification{
		UserID:     p.AuthorId,
		FromUserId: fromuser.Id,
		NotifType:  msg.Type,
		NotifMsg:   msg.Msg,
		ResourceId: msg.ResourceId,
	}

	res, err := orm.Da.CreateNotification(n)
	if err != nil {
		logger.Error.Println("Could not create notification: ", err)
		return
	}
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		logger.Error.Println("Could not get notification Id: ", err)
	}

	dbnotif, err := orm.Da.GetNotificationById(int(lastInsertID))
	if err != nil {
		logger.Error.Println("Could not get notification: ", err)
	}
	notif := mapper.Notification(dbnotif, *u, *from_u)
	notif.ParentId = ""

	data, err := json.Marshal(notif)
	if err != nil {
		logger.Error.Println("Could not marshal JSON: ", err)
	}

	wsocket.WSConnMan.SendMessage(u.UserName, data)
}
