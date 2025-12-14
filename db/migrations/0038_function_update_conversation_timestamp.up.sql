CREATE FUNCTION update_conversation_timestamp()
RETURNS TRIGGER AS
$$
BEGIN
    UPDATE conversations
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.conversation_id;
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;