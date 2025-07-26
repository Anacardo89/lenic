CREATE TRIGGER conversations_enforce_user_order
BEFORE INSERT ON conversations
FOR EACH ROW
EXECUTE FUNCTION enforce_user_order();