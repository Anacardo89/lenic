package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

// Posts
func (c *dbClient) CreatePost(ctx context.Context, p *Post) (uuid.UUID, error) {

	query := `
	INSERT INTO posts (
		id,
		author_id,
		title,
		content,
		post_image,
		is_public
	)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	;`

	var ID uuid.UUID
	err := c.Pool().QueryRow(ctx, query,
		p.ID,
		p.AuthorID,
		p.Title,
		p.Content,
		p.PostImage,
		p.IsPublic,
	).Scan(&ID)
	return ID, err
}

func (c *dbClient) GetFeed(ctx context.Context, userID uuid.UUID) ([]*Post, error) {
	query := `
	SELECT p.* 
	FROM posts p
	LEFT JOIN follows f 
		ON p.author_id = f.followed_id AND f.follower_id = $1
	WHERE
		p.is_active = TRUE AND (
			p.author_id = $1 OR
			p.is_public = TRUE OR
			(f.follower_id = $1 AND f.follow_status = 'accepted')
		)
	ORDER BY 
		CASE 
			WHEN p.created_at >= NOW() - INTERVAL '24 hours' 
			THEN 1 
			ELSE 2 
		END
		ASC,
		p.rating DESC,
		p.created_at DESC
	;`

	posts := []*Post{}
	rows, err := c.Pool().Query(ctx, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return posts, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		p := Post{}
		err = rows.Scan(
			&p.ID,
			&p.AuthorID,
			&p.Title,
			&p.Content,
			&p.PostImage,
			&p.Rating,
			&p.IsPublic,
			&p.IsActive,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}
	return posts, nil
}

func (c *dbClient) GetUserPosts(ctx context.Context, userID uuid.UUID) ([]*Post, error) {

	query := `
	SELECT *
	FROM posts
	WHERE author_id = $1 AND is_active = TRUE
	ORDER BY created_at DESC
	;`

	posts := []*Post{}
	rows, err := c.Pool().Query(ctx, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return posts, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := Post{}
		err = rows.Scan(
			&p.ID,
			&p.AuthorID,
			&p.Title,
			&p.Content,
			&p.PostImage,
			&p.Rating,
			&p.IsPublic,
			&p.IsActive,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}
	return posts, nil
}

func (c *dbClient) GetUserPublicPosts(ctx context.Context, userID uuid.UUID) ([]*Post, error) {

	query := `
	SELECT * 
	FROM posts
	WHERE author_id = $1 AND is_public = TRUE AND is_active = TRUE
	ORDER BY created_at DESC
	;`

	posts := []*Post{}
	rows, err := c.Pool().Query(ctx, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return posts, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := Post{}
		err = rows.Scan(
			&p.ID,
			&p.AuthorID,
			&p.Title,
			&p.Content,
			&p.PostImage,
			&p.Rating,
			&p.IsPublic,
			&p.IsActive,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}
	return posts, nil
}

func (c *dbClient) GetPost(ctx context.Context, ID uuid.UUID) (*Post, error) {

	query := `
	SELECT * 
	FROM posts
	WHERE id = $1
	;`

	p := Post{}
	err := c.Pool().QueryRow(ctx, query, ID).
		Scan(
			&p.ID,
			&p.AuthorID,
			&p.Title,
			&p.Content,
			&p.PostImage,
			&p.Rating,
			&p.IsPublic,
			&p.IsActive,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.DeletedAt,
		)
	return &p, err
}

func (c *dbClient) UpdatePost(ctx context.Context, post *Post) error {

	query := `
	UPDATE posts
	SET title = $2,
		content = $3,
		is_public = $4
	WHERE id = $1
	;`

	_, err := c.Pool().Exec(ctx, query,
		post.ID,
		post.Title,
		post.Content,
		post.IsPublic,
	)
	return err
}

func (c *dbClient) DisablePost(ctx context.Context, ID uuid.UUID) error {

	query := `
	UPDATE posts
	SET is_active = FALSE,
		deleted_at = NOW()
	WHERE id = $1
	;`

	_, err := c.Pool().Exec(ctx, query, ID)
	return err
}

// Post Ratings
func (c *dbClient) RatePostUp(ctx context.Context, targetID, userID uuid.UUID) error {

	query := `
	INSERT INTO post_ratings (
		target_id,
		user_id,
		rating_value
	)
	VALUES ($1, $2, 1)
	ON CONFLICT (target_id, user_id) DO UPDATE
	SET rating_value = 
		CASE
			WHEN post_ratings.rating_value = 1
			THEN 0
			ELSE 1
		END;
	;`
	_, err := c.Pool().Exec(ctx, query, targetID, userID)
	return err
}

func (c *dbClient) RatePostDown(ctx context.Context, targetID, userID uuid.UUID) error {

	query := `
	INSERT INTO post_ratings (
		target_id,
		user_id,
		rating_value
	)
	VALUES ($1, $2, -1)
	ON CONFLICT (target_id, user_id) DO UPDATE
	SET rating_value = 
		CASE
			WHEN post_ratings.rating_value = -1
			THEN 0
			ELSE -1
		END;
	;`

	_, err := c.Pool().Exec(ctx, query, targetID, userID)
	return err
}

func (c *dbClient) GetPostUserRating(ctx context.Context, targetID, userID uuid.UUID) (*PostRatings, error) {
	query := `
	SELECT *
	FROM post_ratings
	WHERE target_id = $1 AND user_id = $2
	;`
	pr := PostRatings{}
	err := c.Pool().QueryRow(ctx, query, targetID, userID).
		Scan(
			&pr.TargetID,
			&pr.UserID,
			&pr.RatingValue,
			&pr.CreatedAt,
			&pr.UpdatedAt,
		)
	return &pr, err
}
