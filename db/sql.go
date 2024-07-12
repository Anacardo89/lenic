package db

const (
	SelectUserByName = `
	SELECT user_name FROM users
		WHERE user_name = ?;
	`

	SelectUserByEmail = `
	SELECT user_email FROM users
		WHERE user_email = ?;
	`

	CreateTableProjects = `
	CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL
	);`
	CreateTableBoards = `
	CREATE TABLE IF NOT EXISTS boards (
		id INTEGER PRIMARY KEY,
		position INTEGER,
		title TEXT NOT NULL,
		project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE
	);`
	CreateTableLabels = `
	CREATE TABLE IF NOT EXISTS labels (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		color TEXT NOT NULL,
		project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE
	);`
	CreateTableCards = `
	CREATE TABLE IF NOT EXISTS cards (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		card_desc TEXT,
		board_id INTEGER NOT NULL REFERENCES boards(id) ON DELETE CASCADE
	);`
	CreateTableCardLabels = `
	CREATE TABLE IF NOT EXISTS card_labels (
		card_id INTEGER NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
		label_id INTEGER NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
		UNIQUE(card_id, label_id)
	);`
	CreateTableCheckItems = `
	CREATE TABLE IF NOT EXISTS check_items (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		done INTEGER NOT NULL,
		card_id INTEGER NOT NULL REFERENCES cards(id) ON DELETE CASCADE
	);`

	// queries
	// projects
	SelectAllProjectsSql = `
	SELECT * FROM projects;`
	CreateProjectSql = `
	INSERT INTO projects (title)
		VALUES ($1)
		RETURNING *;`
	UpdateProjectSql = `
	UPDATE projects
		SET title = ?2
		WHERE id = ?1
		RETURNING *;`
	DeleteProjectSql = `
	DELETE FROM projects
		WHERE id = ?1;`

	// boards
	SelectAllBoardsSql = `
	SELECT * FROM boards;`
	SelectBoardsWithParentOrderedSql = `
	SELECT * FROM boards
		WHERE project_id = ?1
		ORDER BY position ASC;`
	CreateBoardSql = `
	INSERT INTO boards (title, project_id)
		VALUES (?1, ?2)
		RETURNING *;`
	UpdateBoardTitleSql = `
	UPDATE boards
		SET title = ?2
		WHERE id = ?1
		RETURNING *;`
	UpdateBoardPositionSql = `
	UPDATE boards
		SET position = ?2
		WHERE id = ?1
		RETURNING *;`
	DeleteBoardSql = `
	DELETE FROM boards
		WHERE id = ?1;`

	// labels
	SelectAllLabelsSql = `
	SELECT * FROM labels;`
	SelectLabelsWithParentSql = `
	SELECT * FROM labels
		WHERE project_id = ?1;`
	CreateLabelSql = `
	INSERT INTO labels (title, color, project_id)
		VALUES ($1, $2, $3)
		RETURNING *;`
	UpdateLabelTitleSql = `
	UPDATE labels
		SET title = ?2
		WHERE id = ?1
		RETURNING *;`
	UpdateLabelColorSql = `
	UPDATE labels
		SET color = ?2
		WHERE id = ?1
		RETURNING *;`
	DeleteLabelSql = `
	DELETE FROM labels
		WHERE id = ?1;`

	// cards
	SelectAllCardsSql = `
	SELECT * FROM cards;`
	SelectCardsWithParentSql = `
	SELECT * FROM cards
		WHERE board_id = ?1;`
	CreateCardSql = `
	INSERT INTO cards (title, board_id)
		VALUES ($1, $2)
		RETURNING *;`
	UpdateCardTitleSql = `
	UPDATE cards
		SET title = ?2
		WHERE id = ?1
		RETURNING *;`
	UpdateCardDescSql = `
	UPDATE cards
		SET card_desc = ?2
		WHERE id = ?1
		RETURNING *;`
	UpdateCardParentSql = `
	UPDATE cards
		SET board_id = ?2
		WHERE id = ?1
		RETURNING *;`
	DeleteCardSql = `
	DELETE FROM cards
		WHERE id = ?1;`

	// card_labels
	SelectAllCardLabelsSql = `
	SELECT * FROM card_labels;`
	SelectLabelsInCardSql = `
	SELECT * FROM card_labels
		WHERE card_id = ?1;`
	CreateCardLabelSql = `
	INSERT INTO card_labels (card_id, label_id)
		VALUES ($1, $2)
		RETURNING *;`
	DeleteCardLabelSql = `
	DELETE FROM card_labels
		WHERE label_id = ?1;`

	// check_items
	SelectAllCheckItemsSql = `
	SELECT * FROM check_items;`
	SelectCheckItemsWithParentSql = `
	SELECT * FROM check_items
		WHERE card_id = ?1;`
	CreateCheckItemSql = `
	INSERT INTO check_items (title, done, card_id)
		VALUES ($1, $2, $3)
		RETURNING *;`
	UpdateCheckItemTitleSql = `
	UPDATE check_items
		SET title = ?2
		WHERE id = ?1
		RETURNING *;`
	UpdateCheckItemDoneSql = `
	UPDATE check_items
		SET done = ?2
		WHERE id = ?1
		RETURNING *;`
	DeleteCheckItemSql = `
	DELETE FROM check_items
		WHERE id = ?1;`
)
