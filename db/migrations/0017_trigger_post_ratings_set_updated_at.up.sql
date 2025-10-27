CREATE TRIGGER post_ratings_set_updated_at
BEFORE UPDATE ON post_ratings
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();