

/*  Create DB */
CREATE DATABASE tpsi25_blog;
USE tpsi25_blog;

CREATE TABLE users (
	id int unsigned NOT NULL AUTO_INCREMENT,
	user_name varchar(32) NOT NULL DEFAULT '',
	user_email varchar(128) NOT NULL DEFAULT '',
	user_password varchar(128) NOT NULL DEFAULT '',
	user_salt varchar(128) NOT NULL DEFAULT '',
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	user_active tinyint NOT NULL,
	PRIMARY KEY(id),
	UNIQUE KEY user_name (user_name),
	UNIQUE KEY user_email (user_email)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE sessions (
	id int unsigned NOT NULL AUTO_INCREMENT,
	session_id varchar(256) NOT NULL DEFAULT '',
	user_id int DEFAULT NULL,
	session_start timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	session_update timestamp NOT NULL DEFAULT '2023-01-01 23:59:59',
	session_active tinyint NOT NULL,
	PRIMARY KEY (id),
	UNIQUE KEY session_id (session_id)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE posts (
	id int unsigned NOT NULL AUTO_INCREMENT,
	page_guid varchar(256) NOT NULL DEFAULT '',
	page_title varchar(256) DEFAULT NULL,
	page_content mediumtext,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	UNIQUE KEY page_guid (page_guid)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE comments (
	id int unsigned NOT NULL AUTO_INCREMENT,
	page_guid varchar(256) NOT NULL,
	comment_guid varchar(256) DEFAULT NULL,
	comment_name varchar(64) DEFAULT NULL,
	comment_email varchar(128) DEFAULT NULL,
	comment_text mediumtext,
	comment_date timestamp NULL DEFAULT NULL,
	PRIMARY KEY (id),
	KEY page_guiid (page_guid)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



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