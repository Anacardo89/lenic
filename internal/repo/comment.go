package repo

import (
	"context"

	"github.com/google/uuid"
)

// Comments

// Endpoints:
//
// POST /action/post/{post_id}/comment
func (db *dbHandler) CreateComment(ctx context.Context, c *Comment) error {
	query := `
	INSERT INTO comments (
		id,
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

	ID := uuid.New()
	err := db.pool.QueryRow(ctx, query,
		ID,
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

// Endpoints:
//
// ws - comment_rating
// ws - comment_tag
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

// Endpoints:
//
// PUT /action/post/{post_id}/comment/{comment_id}
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

// Endpoints:
//
// DELETE /action/post/{post_id}/comment/{comment_id}
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

// Endpoints:
//
// POST /action/post/{post_id}/comment/{comment_id}/up
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

// Endpoints:
//
// POST /action/post/{post_id}/comment/{comment_id}/down
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
