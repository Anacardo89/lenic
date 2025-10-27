CREATE TRIGGER follows_set_non_accepted_decrement
AFTER UPDATE ON follows
FOR EACH ROW
EXECUTE FUNCTION decrement_follow_counts();