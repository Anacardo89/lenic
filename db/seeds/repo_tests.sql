-- cleanup
TRUNCATE TABLE
    dmessages,
    conversations,
    user_tags,
    comment_ratings,
    post_ratings,
    comments,
    notifications,
    follows,
    posts,
    users
RESTART IDENTITY CASCADE;

-- seed
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
    ),
    (
        'd8e3f4a5-b6c7-4d8e-9f0a-1b2c3d4e5f6d',
        'inactiveuser',
        'Inactive User',
        'inactive@example.com',
        '$2a$10$XLgJUKvQmtOdcOmO2GPYw.Fj16.ls8QDEZyXyna1HPES8Ee8N.sA.',
        'This user is inactive by default.',
        FALSE,
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
        'b1d3c0f7-5a1a-4f9b-9b2a-2a8e4f8b9f01',
        'cfa53179-9085-4f33-86b3-5dc5f7a1465f',
        'Great insights, really makes me think!',
        2,
        TRUE
    ),
    (
        'f2b2c3d4-e5f6-4b7c-8d9e-2f3a4b5c6d7e',
        'b1d3c0f7-5a1a-4f9b-9b2a-2a8e4f8b9f01',
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

INSERT INTO notifications (
    id,
    user_id,
    from_user_id,
    notif_type,
    notif_text,
    resource_id,
    parent_id,
    is_read
) VALUES
    -- Follow request from moderata to soccerpunk
    (
        'c1a1b2c3-d4f5-4e6f-9a7b-1c2d3e4f5a6b',
        'f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e',
        'cfa53179-9085-4f33-86b3-5dc5f7a1465f',
        'follow_request',
        'moderata wants to follow you.',
        'cfa53179-9085-4f33-86b3-5dc5f7a1465f',
        NULL,
        FALSE
    ),

    -- Follow response from anacardo to soccerpunk (accepted)
    (
        'c2b2c3d4-e5f6-4a7b-8c9d-2d3e4f5a6b7c',
        'f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e',
        'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1',
        'follow_response',
        'anacardo accepted your follow request.',
        'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1',
        NULL,
        FALSE
    ),

    -- Comment on Post 1 by moderata, notify anacardo
    (
        'c3c3d4e5-f6a7-4b8c-9d0e-3e4f5a6b7c8d',
        'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1',
        'cfa53179-9085-4f33-86b3-5dc5f7a1465f',
        'post_comment',
        'moderata commented on your post: "Great insights, really makes me think!"',
        'b1d3c0f7-5a1a-4f9b-9a4d-2a8e4f8b9f01',
        'e1a1b2c3-d4f5-4a6b-9c7d-1e2f3a4b5c6d',
        FALSE
    ),

    -- Comment on Post 3 by anacardo, notify soccerpunk
    (
        'c4d4e5f6-a7b8-4c9d-0e1f-4f5a6b7c8d9e',
        'f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e',
        'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1',
        'post_comment',
        'anacardo commented on your post: "I missed the matches, thanks for the recap!"',
        'd3f5e2a9-7c3c-4b7d-9d4c-4c0a6b1d1e23',
        'c5e5f6a7-b8c9-4e0f-1a2b-5c6d7e8f9a0b',
        FALSE
    ),

    -- Post rating notification: moderata liked Post 2, notify author (moderata themselves, for example)
    (
        'c5e5f6a7-b8c9-4d0e-1f2a-5b6c7d8e9f0a',
        'cfa53179-9085-4f33-86b3-5dc5f7a1465f',
        'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1',
        'post_rating',
        'anacardo rated your post "Moderation Tips".',
        'c2e4d1f8-6b2b-4a8c-8c3b-3b9f5a9c0d12',
        NULL,
        FALSE
    ),

    -- Comment rating notification: soccerpunk upvoted Comment 2, notify author (moderata)
    (
        'c6f6a7b8-c9d0-4e1f-2a3b-6c7d8e9f0a1b',
        'cfa53179-9085-4f33-86b3-5dc5f7a1465f',
        'f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e',
        'comment_rating',
        'soccerpunk upvoted your comment: "I have some different thoughts on this topic."',
        'f2b2c3d4-e5f6-4b7c-8d9e-2f3a4b5c6d7e',
        NULL,
        FALSE
    );

INSERT INTO user_tags (user_id, target_id, resource_type) VALUES
    -- anacardo tagged moderata in Post 1
    ('cfa53179-9085-4f33-86b3-5dc5f7a1465f', 'b1d3c0f7-5a1a-4f9b-9a4d-2a8e4f8b9f01', 'post'),

    -- moderata tagged soccerpunk in Post 2
    ('f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e', 'c2e4d1f8-6b2b-4a8c-8c3b-3b9f5a9c0d12', 'post'),

    -- soccerpunk tagged anacardo in Post 3
    ('a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1', 'd3f5e2a9-7c3c-4b7d-9d4c-4c0a6b1d1e23', 'post'),

    -- anacardo tagged soccerpunk in Comment 1
    ('f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e', 'e1a1b2c3-d4f5-4a6b-9c7d-1e2f3a4b5c6d', 'comment'),

    -- moderata tagged anacardo in Comment 2
    ('a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1', 'f2b2c3d4-e5f6-4b7c-8d9e-2f3a4b5c6d7e', 'comment');

INSERT INTO conversations (id, user1_id, user2_id) VALUES
    (
        '70f31f9b-632b-4c4a-bb3a-1f2c6013f001',
        'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1',
        'cfa53179-9085-4f33-86b3-5dc5f7a1465f'
    ),
    (
        '70f31f9b-632b-4c4a-bb3a-1f2c6013f002',
        'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1',
        'f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e'
    ),
    (
        '70f31f9b-632b-4c4a-bb3a-1f2c6013f003',
        'cfa53179-9085-4f33-86b3-5dc5f7a1465f',
        'f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e'
    );

INSERT INTO dmessages (
    id,
    conversation_id,
    sender_id,
    content,
    is_read
) VALUES
    -- Conversation 1: anacardo <-> moderata
    (
        '91e0f1a2-b3c4-4d5e-8f6a-1b2c3d4e5f01',
        '70f31f9b-632b-4c4a-bb3a-1f2c6013f001',
        'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1',
        'Hey, wanted to thank you for the feedback on my post!',
        TRUE
    ),
    (
        '91e0f1a2-b3c4-4d5e-8f6a-1b2c3d4e5f02',
        '70f31f9b-632b-4c4a-bb3a-1f2c6013f001',
        'cfa53179-9085-4f33-86b3-5dc5f7a1465f',
        'No problem at all, it was a great read!',
        FALSE
    ),

    -- Conversation 2: anacardo <-> soccerpunk
    (
        '91e0f1a2-b3c4-4d5e-8f6a-1b2c3d4e5f03',
        '70f31f9b-632b-4c4a-bb3a-1f2c6013f002',
        'f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e',
        'You’re coming to the game this weekend, right?',
        FALSE
    ),
    (
        '91e0f1a2-b3c4-4d5e-8f6a-1b2c3d4e5f04',
        '70f31f9b-632b-4c4a-bb3a-1f2c6013f002',
        'a1f92e18-1d8f-4f0f-9a4d-3b9e3b26b7b1',
        'Wouldn’t miss it for anything!',
        TRUE
    ),

    -- Conversation 3: moderata <-> soccerpunk
    (
        '91e0f1a2-b3c4-4d5e-8f6a-1b2c3d4e5f05',
        '70f31f9b-632b-4c4a-bb3a-1f2c6013f003',
        'cfa53179-9085-4f33-86b3-5dc5f7a1465f',
        'Hey, can you moderate the forum tonight?',
        FALSE
    ),
    (
        '91e0f1a2-b3c4-4d5e-8f6a-1b2c3d4e5f06',
        '70f31f9b-632b-4c4a-bb3a-1f2c6013f003',
        'f7a92b5b-4b7e-4787-9c0b-2b0b6cb86b4e',
        'Sure thing, I’ll be online after 8.',
        TRUE
    );
