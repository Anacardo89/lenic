CREATE TABLE hashtag_resources (
    tag_id UUID NOT NULL REFERENCES hashtags(id),
    target_id UUID NOT NULL,
    resource_type resource_type NOT NULL,
    PRIMARY KEY (tag_id, target_id)
);