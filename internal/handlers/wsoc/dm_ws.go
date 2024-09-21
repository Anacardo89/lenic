package wsoc

import (
	"encoding/json"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mapper"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/Anacardo89/tpsi25_blog/pkg/wsocket"
)

func handleDM(msg wsocket.Message) {

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

	m := &database.DMessage{
		ConversationId: dbConvo.Id,
		SenderId:       from_u.Id,
		Content:        msg.Msg,
	}
	res, err := orm.Da.CreateDMessage(m)
	if err != nil {
		logger.Error.Println("Could not create message: ", err)
		return
	}
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		logger.Error.Println("Could not get dm Id: ", err)
		return
	}

	dbM, err := orm.Da.GetDMById(int(lastInsertID))
	if err != nil {
		logger.Error.Println("Could not get dm by Id: ", err)
		return
	}
	dm := mapper.DMessage(dbM, *from_u)

	err = orm.Da.UpdateConversationById(dbConvo.Id)
	if err != nil {
		logger.Error.Println("Could not update conversation: ", err)
		return
	}

	data, err := json.Marshal(dm)
	if err != nil {
		logger.Error.Println("Could not marshal JSON: ", err)
		return
	}

	wsocket.WSConnMan.SendMessage(u.UserName, data)
}
