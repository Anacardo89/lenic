CREATE TABLE user_tags (
    user_id UUID NOT NULL REFERENCES users(id),
    tagged_resource_id UUID NOT NULL,
    resource_type resource_type NOT NULL,
    PRIMARY KEY (user_id, tagged_resource_id)
);