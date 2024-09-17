package orm

import (
	"database/sql"
	"time"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/query"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
)

func (da *DataAccess) CreateConversation(c *database.Conversation) (sql.Result, error) {
	result, err := da.Db.Exec(query.InsertConversation,
		c.User1Id,
		c.User2Id,
	)
	return result, err
}

func (da *DataAccess) CreateDMessage(d *database.DMessage) (sql.Result, error) {
	result, err := da.Db.Exec(query.InsertDMessage,
		d.ConversationId,
		d.SenderId,
		d.Content,
	)
	return result, err
}

func (da *DataAccess) GetConversationsByUserId(user_id int, limit int, offset int) ([]*database.Conversation, error) {
	conversations := []*database.Conversation{}
	rows, err := da.Db.Query(query.SelectConversationsByUserId, user_id, user_id, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return conversations, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			createdAt []byte
			updatedAt []byte
		)
		c := database.Conversation{}
		err = rows.Scan(
			&c.Id,
			&c.User1Id,
			&c.User2Id,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		c.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
		if err != nil {
			return nil, err
		}
		c.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, &c)
	}
	return conversations, nil
}

func (da *DataAccess) GetDMsByConversationId(conversation_id int, limit int, offset int) ([]*database.DMessage, error) {
	dms := []*database.DMessage{}
	rows, err := da.Db.Query(query.SelectDMsByConversationId, conversation_id, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return dms, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var createdAt []byte
		m := database.DMessage{}
		err = rows.Scan(
			&m.Id,
			&m.ConversationId,
			&m.SenderId,
			&m.Content,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}
		m.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
		if err != nil {
			return nil, err
		}
		dms = append(dms, &m)
	}
	return dms, nil
}

func (da *DataAccess) UpdateConversationById(id int) error {
	_, err := da.Db.Exec(query.UpdateConversationById, id)
	if err != nil {
		return err
	}
	return nil
}
