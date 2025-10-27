CREATE FUNCTION enforce_user_order()
RETURNS TRIGGER AS
$$
DECLARE
    temp UUID;
BEGIN
    IF NEW.user1_id > NEW.user2_id THEN
        temp := NEW.user1_id;
        NEW.user1_id := NEW.user2_id;
        NEW.user2_id := temp;
    END IF;
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;