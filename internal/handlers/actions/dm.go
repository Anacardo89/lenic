package actions

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mapper"
	"github.com/Anacardo89/tpsi25_blog/internal/model/presentation"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/mux"
)

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
	u := mapper.User(dbuser)

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
		from_u := mapper.User(dbfromuser)
		c := mapper.Convesation(dbconvo, *u, *from_u)
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
		dm_sender := mapper.User(dbdm_sender)
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
