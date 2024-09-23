package wsoc

import (
	"encoding/json"
	"strconv"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mapper"
	"github.com/Anacardo89/tpsi25_blog/internal/model/presentation"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/Anacardo89/tpsi25_blog/pkg/wsocket"
)

func handleDM(msg wsocket.Message) {
	logger.Info.Println("/ws handling DM")
	logger.Debug.Println(msg)

	dbuser, err := orm.Da.GetUserByName(msg.ResourceId)
	if err != nil {
		logger.Error.Println("Could not get post: ", err)
		return
	}
	u := mapper.UserNotif(dbuser)

	fromuser, err := orm.Da.GetUserByName(msg.FromUserName)
	if err != nil {
		logger.Error.Println("Could not get from user: ", err)
		return
	}
	from_u := mapper.UserNotif(fromuser)

	dbConvo, err := orm.Da.GetConversationByUserIds(u.Id, from_u.Id)
	if err != nil {
		logger.Error.Println("Could not get conversation: ", err)
		return
	}

	err = orm.Da.UpdateConversationById(dbConvo.Id)
	if err != nil {
		logger.Error.Println("Could not update conversation: ", err)
		return
	}

	convo_id := strconv.Itoa(dbConvo.Id)

	n := &presentation.Notification{
		User:       *u,
		FromUser:   *from_u,
		NotifType:  msg.Type,
		NotifMsg:   msg.Msg,
		ResourceId: convo_id,
		ParentId:   "",
		IsRead:     false,
	}

	data, err := json.Marshal(n)
	if err != nil {
		logger.Error.Println("Could not marshal JSON: ", err)
		return
	}

	wsocket.WSConnMan.SendMessage(u.UserName, data)
}
