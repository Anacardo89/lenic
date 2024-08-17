package orm

import (
	"database/sql"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/query"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
)

func (da *DataAccess) CreatePost(p *database.Post) error {
	_, err := da.Db.Exec(query.InsertPost,
		p.GUID,
		p.Title,
		p.User,
		p.Content,
		p.Image,
		p.ImageExtention,
		p.Active)
	return err
}

func (da *DataAccess) GetPosts() (*[]database.Post, error) {
	posts := []database.Post{}
	rows, err := da.Db.Query(query.SelectPosts)
	if err != nil {
		if err == sql.ErrNoRows {
			return &posts, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		p := database.Post{}
		err = rows.Scan(
			&p.Id, &p.GUID, &p.Title, &p.User, &p.Content, &p.Image,
			&p.ImageExtention, &p.CreatedAt, &p.UpdatedAt, &p.Active)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return &posts, nil
}

func (da *DataAccess) GetPostByGUID(guid string) (*database.Post, error) {
	p := database.Post{}
	row := da.Db.QueryRow(query.SelectPostByGUID, guid)
	err := row.Scan(
		&p.Id, &p.GUID, &p.Title, &p.User, &p.Content, &p.Image,
		&p.ImageExtention, &p.CreatedAt, &p.UpdatedAt, &p.Active)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
