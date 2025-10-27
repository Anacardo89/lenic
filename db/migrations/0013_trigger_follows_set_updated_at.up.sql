CREATE TRIGGER follows_set_updated_at
BEFORE UPDATE ON follows
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();