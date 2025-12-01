CREATE TRIGGER follows_set_accepted_increment
AFTER INSERT OR UPDATE ON follows
FOR EACH ROW
EXECUTE FUNCTION increment_follow_counts();