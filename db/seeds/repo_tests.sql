-- cleanup
TRUNCATE TABLE comment_ratings;
TRUNCATE TABLE comments;
TRUNCATE TABLE post_ratings;
TRUNCATE TABLE posts;
TRUNCATE TABLE follows;
TRUNCATE TABLE users;

INSERT INTO users (
    id,
    username, 
    display_name, 
    email, 
    password_hash, 
    bio, 
    is_active, 
    is_verified, 
    user_role, 
    last_login_at
) VALUES 
    (
        'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1',
        'anacardo',
        'Anacardo',
        'anacardo@example.com',
        '$2a$10$XLgJUKvQmtOdcOmO2GPYw.Fj16.ls8QDEZyXyna1HPES8Ee8N.sA.',
        'Local columnist and political agitator.',
        TRUE,
        TRUE,
        'admin',
        NOW()
    ),
    (
        'cfa53179-9085-4f33-86b3-5dc5f7a1465f',
        'moderata',
        'Moderata Silva',
        'moderata@example.com',
        '$2a$10$XLgJUKvQmtOdcOmO2GPYw.Fj16.ls8QDEZyXyna1HPES8Ee8N.sA.',
        'Keeps the peace in the comment sections.',
        TRUE,
        TRUE,
        'moderator',
        NOW()
    ),
    (
        'f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e',
        'soccerpunk',
        'Soccer Punk',
        'soccerpunk@example.com',
        '$2a$10$XLgJUKvQmtOdcOmO2GPYw.Fj16.ls8QDEZyXyna1HPES8Ee8N.sA.',
        'Football, antifascism, and cold beer.',
        TRUE,
        FALSE,
        'user',
        NOW()
    );

INSERT INTO follows (follower_id, followed_id, follow_status) VALUES
    -- anacardo follows moderata (accepted)
    ('a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1', 'cfa53179-9085-4f33-86b3-5dc5f7a1465f', 'accepted'),
    -- moderata follows soccerpunk (pending)
    ('cfa53179-9085-4f33-86b3-5dc5f7a1465f', 'f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e', 'pending'),
    -- soccerpunk follows anacardo (accepted)
    ('f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e', 'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1', 'accepted');

INSERT INTO posts (
    id,
    author_id,
    title,
    content,
    post_image,
    rating,
    is_public,
    is_active
) VALUES
    (
        'b1d3c0f7-5a1a-4f9b-9b2a-2a8e4f8b9f01',
        'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1',
        'Thoughts on Local Politics',
        'Sharing some ideas about our local elections and activism.',
        '',
        1,
        TRUE,
        TRUE
    ),
    (
        'c2e4d1f8-6b2b-4a8c-8c3b-3b9f5a9c0d12',
        'cfa53179-9085-4f33-86b3-5dc5f7a1465f',
        'Moderation Tips',
        'A few tips on keeping our forums safe and welcoming.',
        '',
        2,
        TRUE,
        TRUE
    ),
    (
        'd3f5e2a9-7c3c-4b7d-9d4c-4c0a6b1d1e23',
        'f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e',
        'Weekend Football Recap',
        'Quick recap of the weekend football games and highlights.',
        '',
        1,
        TRUE,
        TRUE
    );

INSERT INTO post_ratings (target_id, user_id, rating_value) VALUES
    -- Post 1 ratings (final rating: 1)
    ('b1d3c0f7-5a1a-4f9b-9b2a-2a8e4f8b9f01', 'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1', 1),
    ('b1d3c0f7-5a1a-4f9b-9b2a-2a8e4f8b9f01', 'cfa53179-9085-4f33-86b3-5dc5f7a1465f', 1),
    ('b1d3c0f7-5a1a-4f9b-9b2a-2a8e4f8b9f01', 'f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e', -1),

    -- Post 2 ratings (final rating: 2)
    ('c2e4d1f8-6b2b-4a8c-8c3b-3b9f5a9c0d12', 'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1', 1),
    ('c2e4d1f8-6b2b-4a8c-8c3b-3b9f5a9c0d12', 'cfa53179-9085-4f33-86b3-5dc5f7a1465f', 0),
    ('c2e4d1f8-6b2b-4a8c-8c3b-3b9f5a9c0d12', 'f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e', 1),

    -- Post 3 ratings (final rating: 1)
    ('d3f5e2a9-7c3c-4b7d-9d4c-4c0a6b1d1e23', 'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1', -1),
    ('d3f5e2a9-7c3c-4b7d-9d4c-4c0a6b1d1e23', 'cfa53179-9085-4f33-86b3-5dc5f7a1465f', 1),
    ('d3f5e2a9-7c3c-4b7d-9d4c-4c0a6b1d1e23', 'f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e', 1);


INSERT INTO comments (
    id,
    post_id,
    author_id,
    content,
    rating,
    is_active
) VALUES
    -- Comments on Post 1: Thoughts on Local Politics
    (
        'e1a1b2c3-d4f5-4a6b-9c7d-1e2f3a4b5c6d',
        'b1d3c0f7-5a1a-4f9b-9a4d-2a8e4f8b9f01',
        'cfa53179-9085-4f33-86b3-5dc5f7a1465f',
        'Great insights, really makes me think!',
        2,
        TRUE
    ),
    (
        'f2b2c3d4-e5f6-4b7c-8d9e-2f3a4b5c6d7e',
        'b1d3c0f7-5a1a-4f9b-9a4d-2a8e4f8b9f01',
        'f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e',
        'I have some different thoughts on this topic.',
        2,
        TRUE
    ),

    -- Comments on Post 2: Moderation Tips
    (
        'a3c3d4e5-f6a7-4c8d-9e0f-3a4b5c6d7e8f',
        'c2e4d1f8-6b2b-4a8c-8c3b-3b9f5a9c0d12',
        'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1',
        'These tips are really helpful, thanks for sharing!',
        2,
        TRUE
    ),

    -- Comments on Post 3: Weekend Football Recap
    (
        'b4d4e5f6-a7b8-4d9e-0f1a-4b5c6d7e8f9a',
        'd3f5e2a9-7c3c-4b7d-9d4c-4c0a6b1d1e23',
        'cfa53179-9085-4f33-86b3-5dc5f7a1465f',
        'Exciting games this weekend, can’t wait for the next!',
        1,
        TRUE
    ),
    (
        'c5e5f6a7-b8c9-4e0f-1a2b-5c6d7e8f9a0b',
        'd3f5e2a9-7c3c-4b7d-9d4c-4c0a6b1d1e23',
        'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1',
        'I missed the matches, thanks for the recap!',
        3,
        TRUE
    );

INSERT INTO comment_ratings (target_id, user_id, rating_value) VALUES
    -- Comment 1 ratings
    ('e1a1b2c3-d4f5-4a6b-9c7d-1e2f3a4b5c6d', 'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1', 1),
    ('e1a1b2c3-d4f5-4a6b-9c7d-1e2f3a4b5c6d', 'cfa53179-9085-4f33-86b3-5dc5f7a1465f', 1),
    ('e1a1b2c3-d4f5-4a6b-9c7d-1e2f3a4b5c6d', 'f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e', 0),

    -- Comment 2 ratings
    ('f2b2c3d4-e5f6-4b7c-8d9e-2f3a4b5c6d7e', 'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1', 0),
    ('f2b2c3d4-e5f6-4b7c-8d9e-2f3a4b5c6d7e', 'cfa53179-9085-4f33-86b3-5dc5f7a1465f', 1),
    ('f2b2c3d4-e5f6-4b7c-8d9e-2f3a4b5c6d7e', 'f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e', 1),

    -- Comment 3 ratings
    ('a3c3d4e5-f6a7-4c8d-9e0f-3a4b5c6d7e8f', 'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1', 1),
    ('a3c3d4e5-f6a7-4c8d-9e0f-3a4b5c6d7e8f', 'cfa53179-9085-4f33-86b3-5dc5f7a1465f', 0),
    ('a3c3d4e5-f6a7-4c8d-9e0f-3a4b5c6d7e8f', 'f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e', 1),

    -- Comment 4 ratings
    ('b4d4e5f6-a7b8-4d9e-0f1a-4b5c6d7e8f9a', 'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1', 1),
    ('b4d4e5f6-a7b8-4d9e-0f1a-4b5c6d7e8f9a', 'cfa53179-9085-4f33-86b3-5dc5f7a1465f', 1),
    ('b4d4e5f6-a7b8-4d9e-0f1a-4b5c6d7e8f9a', 'f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e', -1),

    -- Comment 5 ratings
    ('c5e5f6a7-b8c9-4e0f-1a2b-5c6d7e8f9a0b', 'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1', 1),
    ('c5e5f6a7-b8c9-4e0f-1a2b-5c6d7e8f9a0b', 'cfa53179-9085-4f33-86b3-5dc5f7a1465f', 1),
    ('c5e5f6a7-b8c9-4e0f-1a2b-5c6d7e8f9a0b', 'f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e', 1);
