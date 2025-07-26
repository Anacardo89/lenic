CREATE TRIGGER comment_ratings_set_updated_at
BEFORE UPDATE ON comment_ratings
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();