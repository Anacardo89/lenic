package orm

import (
	"database/sql"
	"time"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/query"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
)

func (da *DataAccess) CreatePost(p *database.Post) error {
	_, err := da.Db.Exec(query.InsertPost,
		p.GUID,
		p.AuthorId,
		p.Title,
		p.Content,
		p.Image,
		p.ImageExt,
		p.IsPublic,
		p.VoteCount,
		p.Active)
	return err
}

func (da *DataAccess) GetPosts() (*[]database.Post, error) {
	posts := []database.Post{}
	rows, err := da.Db.Query(query.SelectActivePosts)
	if err != nil {
		if err == sql.ErrNoRows {
			return &posts, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			createdAt []byte
			updatedAt []byte
		)
		p := database.Post{}
		err = rows.Scan(
			&p.Id,
			&p.GUID,
			&p.AuthorId,
			&p.Title,
			&p.Content,
			&p.Image,
			&p.ImageExt,
			&createdAt,
			&updatedAt,
			&p.IsPublic,
			&p.VoteCount,
			&p.Active,
		)
		if err != nil {
			return nil, err
		}
		p.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
		if err != nil {
			return nil, err
		}
		p.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return &posts, nil
}

func (da *DataAccess) GetPostByGUID(guid string) (*database.Post, error) {
	var (
		createdAt []byte
		updatedAt []byte
	)
	p := database.Post{}
	row := da.Db.QueryRow(query.SelectPostByGUID, guid)
	err := row.Scan(
		&p.Id,
		&p.GUID,
		&p.AuthorId,
		&p.Title,
		&p.Content,
		&p.Image,
		&p.ImageExt,
		&createdAt,
		&updatedAt,
		&p.IsPublic,
		&p.VoteCount,
		&p.Active,
	)
	if err != nil {
		return nil, err
	}
	p.CreatedAt, err = time.Parse(db.DateLayout, string(createdAt))
	if err != nil {
		return nil, err
	}
	p.UpdatedAt, err = time.Parse(db.DateLayout, string(updatedAt))
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (da *DataAccess) UpdatePostText(guid string, text string) error {
	_, err := da.Db.Exec(query.UpdatePostText, text, guid)
	if err != nil {
		return err
	}
	return nil
}

func (da *DataAccess) DisablePost(guid string) error {
	_, err := da.Db.Exec(query.SetPostAsInactive, guid)
	if err != nil {
		return err
	}
	return nil
}
