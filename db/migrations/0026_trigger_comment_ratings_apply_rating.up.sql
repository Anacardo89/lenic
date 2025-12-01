CREATE TRIGGER comment_ratings_apply_rating
AFTER INSERT OR UPDATE ON comment_ratings
FOR EACH ROW
EXECUTE FUNCTION apply_rating('comments');