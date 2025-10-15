package repo

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

// Comments
func (db *dbHandler) CreateComment(ctx context.Context, c *Comment) error {
	query := `
	INSERT INTO comments (
		post_id,
		author_id,
		content
	)
	VALUES ($1, $2, $3)
	RETURNING
		id,
		post_id,
		author_id,
		content,
		rating,
		is_active,
		created_at
	;`

	err := db.pool.QueryRow(ctx, query,
		c.PostID,
		c.AuthorID,
		c.Content,
	).Scan(
		&c.ID,
		&c.PostID,
		&c.AuthorID,
		&c.Content,
		&c.Rating,
		&c.IsActive,
		&c.CreatedAt,
	)
	return err
}

func (db *dbHandler) GetComment(ctx context.Context, ID uuid.UUID) (*Comment, error) {

	query := `
	SELECT *
	FROM comments
	WHERE id = $1
	;`

	comment := Comment{}
	err := db.pool.QueryRow(ctx, query, ID).
		Scan(
			&comment.ID,
			&comment.PostID,
			&comment.AuthorID,
			&comment.Content,
			&comment.Rating,
			&comment.IsActive,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.DeletedAt,
		)
	return &comment, err
}

func (db *dbHandler) GetCommentsByPost(ctx context.Context, postID uuid.UUID) ([]*Comment, error) {
	query := `
	SELECT *
	FROM comments
	WHERE post_id = $1 AND is_active = 1
	ORDER BY rating DESC
	;`

	comments := []*Comment{}
	rows, err := db.pool.Query(ctx, query, postID)
	if err != nil {
		if err == sql.ErrNoRows {
			return comments, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		comment := Comment{}
		err = rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.AuthorID,
			&comment.Content,
			&comment.Rating,
			&comment.IsActive,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

func (db *dbHandler) UpdateComment(ctx context.Context, c *Comment) error {
	query := `
	UPDATE comments
	SET content = $2
	WHERE id = $1
	RETURNING
		id,
		post_id,
		author_id,
		content,
		rating,
		is_active,
		created_at
	;`

	err := db.pool.QueryRow(ctx, query,
		c.ID,
		c.Content,
	).Scan(
		&c.ID,
		&c.PostID,
		&c.AuthorID,
		&c.Content,
		&c.Rating,
		&c.IsActive,
		&c.CreatedAt,
	)
	return err
}

func (db *dbHandler) DisableComment(ctx context.Context, ID uuid.UUID) (*Comment, error) {

	query := `
	UPDATE comments
	SET active = FALSE,
		deleted_at = NOW()
	WHERE id = $1
	RETURNING
		id,
		content
	;`

	var c Comment
	err := db.pool.QueryRow(ctx, query, ID).Scan(
		&c.ID,
		&c.Content,
	)
	return &c, err
}

// Comment Ratings
func (db *dbHandler) RateCommentUp(ctx context.Context, targetID, userID uuid.UUID) error {

	query := `
	INSERT INTO comment_ratings (
		target_id,
		user_id,
		rating_value
	)
	VALUES ($1, $2, 1)
	ON CONFLICT (target_id, user_id) DO UPDATE
	SET rating_value = 
		CASE
			WHEN comment_ratings.rating_value = 1
			THEN 0
			ELSE 1
		END
	;`

	_, err := db.pool.Exec(ctx, query, targetID, userID)
	return err
}

func (db *dbHandler) RateCommentDown(ctx context.Context, targetID, userID uuid.UUID) error {

	query := `
	INSERT INTO comment_ratings (
		target_id,
		user_id,
		rating_value
	)
	VALUES ($1, $2, -1)
	ON CONFLICT (target_id, user_id) DO UPDATE
	SET rating_value = 
		CASE
			WHEN comment_ratings.rating_value = -1
			THEN 0
			ELSE -1
		END
	;`

	_, err := db.pool.Exec(ctx, query, targetID, userID)
	return err
}

func (db *dbHandler) GetCommentUserRating(ctx context.Context, targetID, userID uuid.UUID) (*CommentRatings, error) {

	query := `
	SELECT *
	FROM comment_ratings
	WHERE target_id = $1 AND user_id = $2
	;`

	cr := CommentRatings{}
	err := db.pool.QueryRow(ctx, query, targetID, userID).
		Scan(
			&cr.TargetID,
			&cr.UserID,
			&cr.RatingValue,
			&cr.CreatedAt,
			&cr.UpdatedAt,
		)
	return &cr, err
}
