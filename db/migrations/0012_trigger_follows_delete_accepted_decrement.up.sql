CREATE TRIGGER follows_delete_accepted_decrement
AFTER DELETE ON follows
FOR EACH ROW
EXECUTE FUNCTION decrement_follow_counts();