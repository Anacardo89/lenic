CREATE TRIGGER post_ratings_apply_rating
AFTER INSERT OR UPDATE ON post_ratings
FOR EACH ROW
EXECUTE FUNCTION apply_rating('posts');