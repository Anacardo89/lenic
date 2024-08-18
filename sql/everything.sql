

/*  Create DB */
CREATE DATABASE tpsi25_blog;
USE tpsi25_blog;

CREATE TABLE users (
	id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	user_name VARCHAR(32) NOT NULL DEFAULT '',
	user_email VARCHAR(128) NOT NULL DEFAULT '',
	user_password VARCHAR(128) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	user_active TINYINT NOT NULL,
	PRIMARY KEY(id),
	UNIQUE KEY user_name (user_name),
	UNIQUE KEY user_email (user_email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE sessions (
	id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	session_id VARCHAR(256) NOT NULL DEFAULT '',
	user_id INT NOT NULL REFERENCES users(id),
	session_start TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	session_update TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	session_active TINYINT NOT NULL,
	PRIMARY KEY (id),
	UNIQUE KEY session_id (session_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE posts (
	id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	post_guid VARCHAR(256) NOT NULL DEFAULT '',
	post_title VARCHAR(256) DEFAULT NULL,
	post_user VARCHAR(32) NOT NULL REFERENCES users(user_name),
	post_content MEDIUMTEXT,
	post_image VARCHAR(64),
	post_image_ext VARCHAR(10),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	post_active TINYINT NOT NULL,
	PRIMARY KEY (id),
	UNIQUE KEY post_guid (post_guid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE comments (
	id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	post_guid VARCHAR(256) NOT NULL REFERENCES posts(post_guid),
	comment_author VARCHAR(64) NOT NULL REFERENCES users(user_name),
	comment_text MEDIUMTEXT,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	comment_active TINYINT NOT NULL,
	PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



/*  Mock Data */
INSERT INTO pages (id, page_guid, page_title, page_content, page_date)
VALUES (
	1,
	"hello-world",
	"Hello, World",
	"I'm so glad you found this page! It's been sitting patiently on the Internet for some time, just waiting for a visitor.",
	CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);

INSERT INTO pages (id, page_guid, page_title, page_content, page_date)
VALUES (
	2,
	"a-new-blog",
	"A New Blog",
	"I hope you enjoyed the last blog! Well brace yourself, because my latest blog is even <i>better</i> than the last!",
	CURRENT_TIMESTAMP
);

INSERT INTO pages (id, page_guid, page_title, page_content, page_date)
VALUES (
	3,
	"lorem-ipsum",
	"Lorem Ipsum",
	"'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas sem tortor, lobortis in posuere sit amet, ornare non eros. Pellentesque vel lorem sed nisl dapibus fringilla. In pretium...",
	CURRENT_TIMESTAMP
);