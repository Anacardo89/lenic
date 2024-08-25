package orm

import (
	"database/sql"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/query"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
)

func (da *DataAccess) GetPostUserRating(post_id int, user_id int) (*database.PostRatings, error) {
	pr := database.PostRatings{}
	row := da.Db.QueryRow(query.SelectPostUserRating, post_id, user_id)
	err := row.Scan(
		&pr.PostId,
		&pr.UserId,
		&pr.RatingValue,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pr, err
		} else {
			return nil, err
		}
	}
	return &pr, nil
}

func (da *DataAccess) GetCommentUserRating(comment_id int, user_id int) (*database.CommentRatings, error) {
	cr := database.CommentRatings{}
	row := da.Db.QueryRow(query.SelectCommentUserRating, comment_id, user_id)
	err := row.Scan(
		&cr.CommentId,
		&cr.UserId,
		&cr.RatingValue,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return &cr, err
		} else {
			return nil, err
		}
	}
	return &cr, nil
}
