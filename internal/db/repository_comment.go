package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

// Comments
func (c *dbClient) CreateComment(ctx context.Context, comment *Comment) (uuid.UUID, error) {
	query := `
	INSERT INTO comments (
		post_id,
		author_id,
		content,
		is_active
	)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	;`

	var ID uuid.UUID
	err := c.Pool().QueryRow(ctx, query,
		comment.PostID,
		comment.AuthorID,
		comment.Content,
		comment.IsActive,
	).Scan(&ID)
	return ID, err
}

func (c *dbClient) GetComment(ctx context.Context, ID uuid.UUID) (*Comment, error) {

	query := `
	SELECT *
	FROM comments
	WHERE id = $1
	;`

	comment := Comment{}
	err := c.Pool().QueryRow(ctx, query, ID).
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

func (c *dbClient) GetCommentsByPost(ctx context.Context, postID uuid.UUID) ([]*Comment, error) {

	query := `
	SELECT *
	FROM comments
	WHERE post_id = $1 AND is_active = 1
	ORDER BY rating DESC
	;`

	comments := []*Comment{}
	rows, err := c.Pool().Query(ctx, query, postID)
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

func (c *dbClient) UpdateCommentContent(ctx context.Context, ID uuid.UUID, content string) error {

	query := `
	UPDATE comments
	SET content = $2
	WHERE id = $1
	;`

	_, err := c.Pool().Exec(ctx, query, ID, content)
	return err
}

func (c *dbClient) DisableComment(ctx context.Context, ID uuid.UUID) error {

	query := `
	UPDATE comments
	SET active = FALSE,
		deleted_at = NOW()
	WHERE id = $1
	;`

	_, err := c.Pool().Exec(ctx, query, ID)
	return err
}

// Comment Ratings
func (c *dbClient) RateCommentUp(ctx context.Context, targetID, userID uuid.UUID) error {

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
		END;
	;`

	_, err := c.Pool().Exec(ctx, query, targetID, userID)
	return err
}

func (c *dbClient) RateCommentDown(ctx context.Context, targetID, userID uuid.UUID) error {

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
		END;
	;`

	_, err := c.Pool().Exec(ctx, query, targetID, userID)
	return err
}

func (c *dbClient) GetCommentUserRating(ctx context.Context, targetID, userID uuid.UUID) (*CommentRatings, error) {

	query := `
	SELECT *
	FROM comment_ratings
	WHERE target_id = $1 AND user_id = $2
	;`

	cr := CommentRatings{}
	err := c.Pool().QueryRow(ctx, query, targetID, userID).
		Scan(
			&cr.TargetID,
			&cr.UserID,
			&cr.RatingValue,
			&cr.CreatedAt,
			&cr.UpdatedAt,
		)
	return &cr, err
}
