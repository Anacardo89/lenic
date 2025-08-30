CREATE FUNCTION decrement_follow_counts()
RETURNS TRIGGER AS
$$
BEGIN
    -- Case: DELETE and relationship was accepted
    IF TG_OP = 'DELETE' AND OLD.follow_status = 'accepted' THEN
        UPDATE users SET user_followers = user_followers - 1 WHERE id = OLD.followed_id;
        UPDATE users SET user_following = user_following - 1 WHERE id = OLD.follower_id;

    -- Case: UPDATE from accepted to anything else
    ELSIF TG_OP = 'UPDATE' AND OLD.follow_status = 'accepted' AND NEW.follow_status = 'blocked' THEN
        UPDATE users SET user_followers = user_followers - 1 WHERE id = OLD.followed_id;
        UPDATE users SET user_following = user_following - 1 WHERE id = OLD.follower_id;
    END IF;

    RETURN COALESCE(NEW, OLD);
END;
$$
LANGUAGE plpgsql;