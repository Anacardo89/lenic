package repo

import (
	"context"

	"github.com/google/uuid"
)

// UserTags
func (db *dbHandler) CreateUserTag(ctx context.Context, t *UserTag) error {

	query := `
	INSERT INTO user_tags (
		user_id,
		target_id,
		resource_type
	)
	VALUES ($1, $2, $3)
	;`

	_, err := db.pool.Exec(ctx, query,
		t.UserID,
		t.TargetID,
		t.ResourceTpe,
	)
	return err
}

func (db *dbHandler) GetUserTagByTarget(ctx context.Context, userID, targetID uuid.UUID) (*UserTag, error) {
	query := `
	SELECT *
	FROM user_tags
	WHERE user_id = $1 AND target_id = $2
	;`

	t := UserTag{}
	err := db.pool.QueryRow(ctx, query, userID, targetID).
		Scan(
			&t.UserID,
			&t.TargetID,
			&t.ResourceTpe,
		)
	return &t, err
}

func (db *dbHandler) DeleteUserTag(ctx context.Context, userID uuid.UUID, targetID uuid.UUID) error {

	query := `
	DELETE FROM user_tags
	WHERE user_id = $1 AND target_id = $2
	;`

	_, err := db.pool.Exec(ctx, query, userID, targetID)
	return err
}

// HashTags
func (db *dbHandler) CreateHashtag(ctx context.Context, t *HashTag) (uuid.UUID, error) {
	query := `
	INSERT INTO hashtags (tag_name)
	VALUES ($1)
	ON CONFLICT (tag_name) DO NOTHING
	RETURNING id
	;`

	var ID uuid.UUID
	err := db.pool.QueryRow(ctx, query,
		t.TagName,
	).Scan(&ID)
	return ID, err
}

func (db *dbHandler) GetHashTagByName(ctx context.Context, tagName string) (*HashTag, error) {

	query := `
	SELECT *
	FROM hashtags
	WHERE tag_name = $1
	;`

	t := HashTag{}
	err := db.pool.QueryRow(ctx, query, tagName).
		Scan(
			&t.ID,
			&t.TagName,
			&t.CreatedAt,
		)
	return &t, err
}

// HashTag Resources
func (db *dbHandler) CreateHashTagResource(ctx context.Context, t *HashTagResource) error {

	query := `
	INSERT INTO hashtag_resources (
		tag_id,
		target_id,
		resource_type
	)
	VALUES ($1, $2, $3)
	;`

	_, err := db.pool.Exec(ctx, query,
		t.TagID,
		t.TargetID,
		t.ResourceTpe,
	)
	return err
}

func (db *dbHandler) GetHashTagResourceByTarget(ctx context.Context, tagID, targetID uuid.UUID) (*HashTagResource, error) {

	query := `
	SELECT *
	FROM hashtag_resources
	WHERE tag_id = $1 AND target_id = $2
	;`

	t := HashTagResource{}
	err := db.pool.QueryRow(ctx, query, tagID, targetID).
		Scan(
			&t.TagID,
			&t.TargetID,
			&t.ResourceTpe,
		)
	return &t, err
}
