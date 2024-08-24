package orm

import (
	"database/sql"
	"time"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/query"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
)

func (da *DataAccess) CreateComment(c *database.Comment) error {
	_, err := da.Db.Exec(query.InsertComment,
		c.PostGUID,
		c.AuthorId,
		c.Content,
		c.VoteCount,
		c.Active)
	return err
}

func (da *DataAccess) GetCommentsByPost(guid string) (*[]database.Comment, error) {
	comments := []database.Comment{}
	rows, err := da.Db.Query(query.SelectActiveCommentsByPost, guid)
	if err != nil {
		if err == sql.ErrNoRows {
			return &comments, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			createdAt []byte
			updatedAt []byte
		)
		c := database.Comment{}
		err = rows.Scan(
			&c.Id,
			&c.PostGUID,
			&c.AuthorId,
			&c.Content,
			&createdAt,
			&updatedAt,
			&c.VoteCount,
			&c.Active)
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
		comments = append(comments, c)
	}
	return &comments, nil
}

func (da *DataAccess) UpdateCommentText(id int, text string) error {
	_, err := da.Db.Exec(query.UpdateCommentText, text, id)
	if err != nil {
		return err
	}
	return nil
}

func (da *DataAccess) DisableComment(id int) error {
	_, err := da.Db.Exec(query.SetCommentAsInactive, id)
	if err != nil {
		return err
	}
	return nil
}
