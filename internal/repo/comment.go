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
	VALUES ($1, $2, $3, $4)
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
	if err := db.pool.QueryRow(ctx, query,
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
	); err != nil {
		return err
	}
	return nil
}

// Endpoints:
//
// ws - comment_rating
// ws - comment_tag
func (db *dbHandler) GetComment(ctx context.Context, ID uuid.UUID) (*Comment, error) {
	query := `
	SELECT
		id,
		post_id,
		author_id,
		content,
		rating,
		is_active,
		created_at,
		updated_at
	FROM comments
	WHERE id = $1
		AND is_active = TRUE
	;`
	comment := Comment{}
	if err := db.pool.QueryRow(ctx, query, ID).
		Scan(
			&comment.ID,
			&comment.PostID,
			&comment.AuthorID,
			&comment.Content,
			&comment.Rating,
			&comment.IsActive,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		); err != nil {
		return nil, err
	}
	return &comment, nil
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
		created_at,
		updated_at
	;`
	if err := db.pool.QueryRow(ctx, query,
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
		&c.UpdatedAt,
	); err != nil {
		return err
	}
	return nil
}

// Endpoints:
//
// DELETE /action/post/{post_id}/comment/{comment_id}
func (db *dbHandler) DisableComment(ctx context.Context, ID uuid.UUID) (*Comment, error) {
	query := `
	UPDATE comments
	SET is_active = FALSE,
		deleted_at = NOW()
	WHERE id = $1
		AND is_active = TRUE
	RETURNING
		id,
		content
	;`
	var c Comment
	if err := db.pool.QueryRow(ctx, query, ID).Scan(
		&c.ID,
		&c.Content,
	); err != nil {
		return nil, err
	}
	return &c, nil
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
	if _, err := db.pool.Exec(ctx, query, targetID, userID); err != nil {
		return err
	}
	return nil
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
	if _, err := db.pool.Exec(ctx, query, targetID, userID); err != nil {
		return err
	}
	return nil
}
