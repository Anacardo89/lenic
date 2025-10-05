package repo

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
)

// Posts
func (db *dbHandler) CreatePost(ctx context.Context, p *Post) (uuid.UUID, error) {

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
	err := db.pool.QueryRow(ctx, query,
		p.ID,
		p.AuthorID,
		p.Title,
		p.Content,
		p.PostImage,
		p.IsPublic,
	).Scan(&ID)
	return ID, err
}

func (db *dbHandler) GetFeed(ctx context.Context, username string) ([]*Post, error) {
	query := `
	WITH active_user AS (
		SELECT id AS user_id 
		FROM users 
		WHERE username = $1
	)
	SELECT 
		p.id,
		p.author_id,
		p.title,
		p.content,
		p.post_image,
		p.rating,
		p.is_public,
		p.is_active,
		p.created_at,
		p.updated_at
	FROM posts p
	LEFT JOIN follows f
		ON p.author_id = f.followed_id
		AND f.follower_id = (SELECT user_id FROM active_user)
	WHERE
		p.is_active = TRUE
		AND (
			p.author_id = (SELECT user_id FROM active_user)
			OR p.is_public = TRUE
			OR (f.follower_id = (SELECT user_id FROM active_user) AND f.follow_status = 'accepted')
		)
	ORDER BY 
		CASE 
			WHEN p.created_at >= NOW() - INTERVAL '24 hours' THEN 1
			ELSE 2
		END ASC,
		p.rating DESC,
		p.created_at DESC
	;`

	posts := []*Post{}
	rows, err := db.pool.Query(ctx, query, username)
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
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}
	return posts, nil
}

func (db *dbHandler) GetPostAuthorFromComment(ctx context.Context, commentID uuid.UUID) (*User, error) {
	query := `
	SELECT 
		u.id AS user_id,
		u.username,
		u.profile_pic
	FROM comments c
	JOIN posts p ON c.post_id = p.id
	JOIN users u ON p.author_id = u.id
	WHERE c.id = $1
	;`

	var u User
	err := db.pool.QueryRow(ctx, query, commentID).Scan(
		&u.ID,
		&u.Username,
		&u.ProfilePic,
	)
	return &u, err
}

func (db *dbHandler) GetUserPosts(ctx context.Context, userID uuid.UUID) ([]*Post, error) {

	query := `
	SELECT *
	FROM posts
	WHERE author_id = $1 AND is_active = TRUE
	ORDER BY created_at DESC
	;`

	posts := []*Post{}
	rows, err := db.pool.Query(ctx, query, userID)
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

func (db *dbHandler) GetUserPublicPosts(ctx context.Context, userID uuid.UUID) ([]*Post, error) {

	query := `
	SELECT * 
	FROM posts
	WHERE author_id = $1 AND is_public = TRUE AND is_active = TRUE
	ORDER BY created_at DESC
	;`

	posts := []*Post{}
	rows, err := db.pool.Query(ctx, query, userID)
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

func (db *dbHandler) GetPost(ctx context.Context, ID uuid.UUID) (*Post, error) {

	query := `
	SELECT * 
	FROM posts
	WHERE id = $1
	;`

	p := Post{}
	err := db.pool.QueryRow(ctx, query, ID).
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

func (db *dbHandler) GetPostForPage(ctx context.Context, ID, userID uuid.UUID) (*PostWithComments, error) {
	query := `
	SELECT 
		p.id AS post_id,
		p.author_id,
		p.title,
		p.content,
		p.post_image,
		p.rating,
		p.is_public,
		p.is_active,
		p.created_at,
		p.updated_at,
		-- Post rating for the specific user
		COALESCE(MAX(pr.rating_value), 0) AS user_post_rating,
		-- Post author
		json_build_object(
			'id', u.id,
			'username', u.username,
			'profile_pic', u.profile_pic
		) AS author,
		-- Comments array
		COALESCE(
			json_agg(cdata) FILTER (WHERE cdata IS NOT NULL),
			'[]'
		) AS comments
	FROM posts p
	JOIN users u
		ON u.id = p.author_id
	LEFT JOIN post_ratings pr 
		ON pr.target_id = p.id AND pr.user_id = $2
	-- comment + author + rating packed into subquery
	LEFT JOIN (
		SELECT 
			c.id,
			c.post_id,
			c.author_id,
			c.content,
			c.rating,
			c.is_active,
			c.created_at,
			c.updated_at,
			COALESCE(MAX(cr.rating_value), 0) AS user_rating,
			json_build_object(
				'id', cu.id,
				'username', cu.username,
				'profile_pic', cu.profile_pic
			) AS author
		FROM comments c
		JOIN users cu ON cu.id = c.author_id
		LEFT JOIN comment_ratings cr 
			ON cr.target_id = c.id AND cr.user_id = $2
		WHERE c.is_active = TRUE
		GROUP BY c.id, cu.id
	) AS cdata ON cdata.post_id = p.id
	WHERE p.id = $1 AND p.is_active = TRUE
	GROUP BY p.id, u.id
	;`
	var (
		p     PostWithComments
		uJSON []byte
		cJSON []byte
	)
	if err := db.pool.QueryRow(ctx, query, ID, userID).
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
			&p.UserRating,
			&uJSON,
			&cJSON,
		); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(uJSON, &p.Author); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(cJSON, &p.Comments); err != nil {
		return nil, err
	}
	return &p, nil
}

func (db *dbHandler) UpdatePost(ctx context.Context, post *Post) error {

	query := `
	UPDATE posts
	SET title = $2,
		content = $3,
		is_public = $4
	WHERE id = $1
	;`

	_, err := db.pool.Exec(ctx, query,
		post.ID,
		post.Title,
		post.Content,
		post.IsPublic,
	)
	return err
}

func (db *dbHandler) DisablePost(ctx context.Context, ID uuid.UUID) (*Post, error) {
	query := `
	UPDATE posts
	SET active = FALSE,
		deleted_at = NOW()
	WHERE id = $1
	RETURNING
		id,
		title,
		content
	;`
	var p Post
	err := db.pool.QueryRow(ctx, query, ID).Scan(
		&p.ID,
		&p.Title,
		&p.Content,
	)
	return &p, err
}

// Post Ratings
func (db *dbHandler) RatePostUp(ctx context.Context, targetID, userID uuid.UUID) error {

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
	_, err := db.pool.Exec(ctx, query, targetID, userID)
	return err
}

func (db *dbHandler) RatePostDown(ctx context.Context, targetID, userID uuid.UUID) error {

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

	_, err := db.pool.Exec(ctx, query, targetID, userID)
	return err
}

func (db *dbHandler) GetPostUserRating(ctx context.Context, targetID, userID uuid.UUID) (*PostRatings, error) {
	query := `
	SELECT *
	FROM post_ratings
	WHERE target_id = $1 AND user_id = $2
	;`
	pr := PostRatings{}
	err := db.pool.QueryRow(ctx, query, targetID, userID).
		Scan(
			&pr.TargetID,
			&pr.UserID,
			&pr.RatingValue,
			&pr.CreatedAt,
			&pr.UpdatedAt,
		)
	return &pr, err
}
