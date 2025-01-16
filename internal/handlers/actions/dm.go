package actions

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mapper"
	"github.com/Anacardo89/tpsi25_blog/internal/model/presentation"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/mux"
)

type JSON_Convo struct {
	User string `json:"to_user"`
}

// POST /action/user/{user_encoded}/conversations
func StartConversation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	logger.Info.Printf("POST /action/user/%s/conversations %s\n", encoded, r.RemoteAddr)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations - Could not decode user: %s\n", encoded, err)
		return
	}
	userName := string(bytes)
	logger.Info.Printf("POST /action/user/%s/conversations %s %s\n", encoded, r.RemoteAddr, userName)

	var msg JSON_Convo
	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations - Could not parse Json Data: %s\n", encoded, err)
		return
	}

	dbuser, err := orm.Da.GetUserByName(userName)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations - Could not get user: %s\n", encoded, err)
		return
	}
	u := mapper.UserNotif(dbuser)

	dbfromuser, err := orm.Da.GetUserByName(msg.User)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations - Could not get from user: %s\n", encoded, err)
		return
	}
	from_u := mapper.UserNotif(dbfromuser)

	var dbconvo *database.Conversation
	dbconvo, err = orm.Da.GetConversationByUserIds(u.Id, from_u.Id)
	if err == sql.ErrNoRows {
		convo := &database.Conversation{
			User1Id: u.Id,
			User2Id: from_u.Id,
		}
		res, err := orm.Da.CreateConversation(convo)
		if err != nil {
			logger.Error.Println("Could not create conversation: ", err)
			return
		}
		lastInsertID, err := res.LastInsertId()
		if err != nil {
			logger.Error.Printf("POST /action/user/%s/conversations - Could not get conversation id: %s\n", encoded, err)
			return
		}
		dbconvo, err = orm.Da.GetConversationById(int(lastInsertID))
		if err != nil {
			logger.Error.Printf("POST /action/user/%s/conversations - Could not get conversation: %s\n", encoded, err)
			return
		}
	} else if err != nil {
		logger.Error.Println("Could not get conversation: ", err)
		return
	}
	convo := mapper.Convesation(dbconvo, *u, *from_u, false)

	data, err := json.Marshal(convo)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations - Could not marshal conversations: %s\n", encoded, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// GET /action/user/{user_encoded}/conversations
func GetConversations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	logger.Info.Printf("GET /action/user/%s/conversations %s\n", encoded, r.RemoteAddr)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations - Could not decode user: %s\n", encoded, err)
		return
	}
	userName := string(bytes)
	logger.Info.Printf("GET /action/user/%s/conversations %s %s\n", encoded, r.RemoteAddr, userName)

	dbuser, err := orm.Da.GetUserByName(userName)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations - Could not get user: %s\n", encoded, err)
		return
	}
	u := mapper.UserNotif(dbuser)

	queryParams := r.URL.Query()
	offset := queryParams.Get("offset")
	offsetint, err := strconv.Atoi(offset)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations - Could not parse offset to int: %s\n", encoded, err)
		return
	}

	limit := queryParams.Get("limit")
	limitint, err := strconv.Atoi(limit)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations - Could not parse limit to int: %s\n", encoded, err)
		return
	}

	dbconvos, err := orm.Da.GetConversationsByUserId(dbuser.Id, limitint, offsetint)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations - Could not get conversations: %s\n", encoded, err)
		return
	}

	var convos []*presentation.Conversation
	for _, dbconvo := range dbconvos {
		fromuser_id := dbconvo.User1Id
		if dbconvo.User1Id == dbuser.Id {
			fromuser_id = dbconvo.User2Id
		}
		dbfromuser, err := orm.Da.GetUserByID(fromuser_id)
		if err != nil {
			logger.Error.Printf("GET /action/user/%s/conversations - Could not get user: %s\n", encoded, err)
			return
		}
		from_u := mapper.UserNotif(dbfromuser)
		dms, err := orm.Da.GetDMsByConversationId(dbconvo.Id, 1000, 0)
		if err != nil {
			logger.Error.Printf("GET /action/user/%s/conversations - Could not get dms: %s\n", encoded, err)
			return
		}
		is_read := true
		for _, dm := range dms {
			if dm.SenderId != dbuser.Id && !dm.IsRead {
				is_read = false
				break
			}
		}
		c := mapper.Convesation(dbconvo, *u, *from_u, is_read)
		convos = append(convos, c)
	}

	data, err := json.Marshal(convos)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations - Could not marshal conversations: %s\n", encoded, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// GET /action/user/{user_encoded}/conversations/{conversation_id}/dms
func GetDMs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	conversation_id := vars["conversation_id"]
	logger.Info.Printf("GET /action/user/%s/conversations/%s/dms %s\n", encoded, conversation_id, r.RemoteAddr)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations/%s/dms - Could not decode user: %s\n", encoded, conversation_id, err)
		return
	}
	userName := string(bytes)
	logger.Info.Printf("GET /action/user/%s/conversations/%s/dms %s %s\n", encoded, conversation_id, r.RemoteAddr, userName)

	convoid_int, err := strconv.Atoi(conversation_id)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations/%s/dms - Could not parse conversation_id to int: %s\n", encoded, conversation_id, err)
		return
	}

	queryParams := r.URL.Query()
	offset := queryParams.Get("offset")
	offsetint, err := strconv.Atoi(offset)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations/%s/dms - Could not parse offset to int: %s\n", encoded, conversation_id, err)
		return
	}

	limit := queryParams.Get("limit")
	limitint, err := strconv.Atoi(limit)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations/%s/dms - Could not parse limit to int: %s\n", encoded, conversation_id, err)
		return
	}

	dbconvo, err := orm.Da.GetConversationById(convoid_int)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations/%s/dms - Could not get conversation: %s\n", encoded, conversation_id, err)
		return
	}

	var dms []*presentation.DMessage
	dbdms, err := orm.Da.GetDMsByConversationId(dbconvo.Id, limitint, offsetint)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations/%s/dms - Could not get DMs: %s\n", encoded, conversation_id, err)
		return
	}
	for _, dbdm := range dbdms {
		dbdm_sender, err := orm.Da.GetUserByID(dbdm.SenderId)
		if err != nil {
			logger.Error.Printf("GET /action/user/%s/conversations/%s/dms - Could not get sender: %s\n", encoded, conversation_id, err)
			return
		}
		dm_sender := mapper.UserNotif(dbdm_sender)
		dm := mapper.DMessage(dbdm, *dm_sender)
		dms = append(dms, dm)
	}

	data, err := json.Marshal(dms)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations/%s/dms - Could not marshal dms: %s\n", encoded, conversation_id, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

type JSON_DM struct {
	Msg string `json:"text"`
}

// POST /action/user/{user_encoded}/conversations/{conversation_id}/dms
func SendDM(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	conversation_id := vars["conversation_id"]
	logger.Info.Printf("POST /action/user/%s/conversations/%s/dms %s\n", encoded, conversation_id, r.RemoteAddr)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations/%s/dms - Could not decode user: %s\n", encoded, conversation_id, err)
		return
	}
	userName := string(bytes)
	logger.Info.Printf("POST /action/user/%s/conversations/%s/dms %s %s\n", encoded, conversation_id, r.RemoteAddr, userName)

	convoid_int, err := strconv.Atoi(conversation_id)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations/%s/dms - Could not parse conversation_id to int: %s\n", encoded, conversation_id, err)
		return
	}

	var msg JSON_DM
	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations/%s/dms - Could not parse Json Data: %s\n", encoded, conversation_id, err)
		return
	}

	dbconvo, err := orm.Da.GetConversationById(convoid_int)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations/%s/dms - Could not get db conversation: %s\n", encoded, conversation_id, err)
		return
	}

	dbuser, err := orm.Da.GetUserByName(userName)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations/%s/dms - Could not get db user: %s\n", encoded, conversation_id, err)
		return
	}

	sender_id := dbconvo.User1Id
	if dbconvo.User1Id != dbuser.Id {
		sender_id = dbconvo.User2Id
	}

	m := &database.DMessage{
		ConversationId: convoid_int,
		SenderId:       sender_id,
		Content:        msg.Msg,
	}

	_, err = orm.Da.CreateDMessage(m)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations/%s/dms - Could not get db user: %s\n", encoded, conversation_id, err)
		return
	}

	logger.Info.Printf("OK - POST /action/user/%s/conversations/%s/dms\n", encoded, conversation_id)
	w.WriteHeader(http.StatusOK)
}

// PUT /action/user/{user_encoded}/conversations/{conversation_id}/read
func ReadConversation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	convo_id := vars["conversation_id"]
	logger.Info.Printf("PUT /action/user/%s/conversations/%s/read %s\n", encoded, convo_id, r.RemoteAddr)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/conversations/%s/read - Could not decode user: %s\n", encoded, convo_id, err)
		return
	}
	userName := string(bytes)
	logger.Info.Printf("PUT /action/user/%s/conversations/%s/read %s %s\n", encoded, convo_id, r.RemoteAddr, userName)

	dbuser, err := orm.Da.GetUserByName(userName)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/conversations/%s/read - Could get user: %s\n", encoded, convo_id, err)
		return
	}

	convo_id_int, err := strconv.Atoi(convo_id)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/conversations/%s/read - Could not parse convo_id to int: %s\n", encoded, convo_id, err)
		return
	}

	dms, err := orm.Da.GetDMsByConversationId(convo_id_int, 1000, 0)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/conversations/%s/read - Could get DMs: %s\n", encoded, convo_id, err)
		return
	}

	for _, dm := range dms {
		if dm.SenderId != dbuser.Id {
			err = orm.Da.UpdateDMReadById(dm.Id)
			if err != nil {
				logger.Error.Printf("PUT /action/user/%s/conversations/%s/read - Could not update notif: %s\n", encoded, convo_id, err)
				return
			}
		}
	}

	logger.Info.Printf("OK - PUT /action/user/%s/conversations/%s/read\n", encoded, convo_id)
	w.WriteHeader(http.StatusOK)
}
