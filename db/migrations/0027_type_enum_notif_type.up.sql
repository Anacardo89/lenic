CREATE TYPE notif_type AS ENUM (
    'follow_request',
    'follow_response',
    'post_comment',
    'post_mention',
    'comment_mention',
    'post_rating',
    'comment_rating',
);