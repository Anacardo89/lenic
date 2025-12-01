CREATE FUNCTION increment_follow_counts()
RETURNS TRIGGER AS
$$
BEGIN
    -- Case: INSERT with accepted status
    IF TG_OP = 'INSERT' AND NEW.follow_status = 'accepted' THEN
        UPDATE users SET user_followers = user_followers + 1 WHERE id = NEW.followed_id;
        UPDATE users SET user_following = user_following + 1 WHERE id = NEW.follower_id;
    
    -- Case: UPDATE from pending/refused/blocked to accepted
    ELSIF TG_OP = 'UPDATE' AND NEW.follow_status = 'accepted' AND OLD.follow_status = 'pending' THEN
        UPDATE users SET user_followers = user_followers + 1 WHERE id = NEW.followed_id;
        UPDATE users SET user_following = user_following + 1 WHERE id = NEW.follower_id;
    END IF;

    RETURN NEW;
END;
$$
LANGUAGE plpgsql;