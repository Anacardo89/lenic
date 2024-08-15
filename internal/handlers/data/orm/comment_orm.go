package orm

import (
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/query"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
)

func (da *DataAccess) CreateComment(c *database.Comment) error {
	_, err := da.db.Exec(query.InsertComment,
		c.PostGUID,
		c.CommentAuthor,
		c.CommentText,
		c.Active)
	return err
}

func (da *DataAccess) GetCommentsByPost(guid string) (*[]database.Comment, error) {
	comments := []database.Comment{}
	rows, err := da.db.Query(query.SelectCommentsByPost, guid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		c := database.Comment{}
		err = rows.Scan(
			c.Id,
			c.PostGUID,
			c.CommentAuthor,
			c.CommentText,
			c.CreatedAt,
			c.UpdatedAt,
			c.Active)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return &comments, nil
}

func (da *DataAccess) UpdateCommentText(id int, text string) error {
	_, err := da.db.Exec(query.UpdateCommentText, id, text)
	if err != nil {
		return err
	}
	return nil
}
