CREATE TYPE notif_type AS ENUM (
    'follow_request',
    'follow_response',
    'post_comment',
    'post_tag',
    'comment_tag',
    'post_rating',
    'comment_rating'
);