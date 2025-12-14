CREATE TRIGGER update_conversation_timestamp_on_new_message
AFTER INSERT ON dmessages
FOR EACH ROW
EXECUTE FUNCTION update_conversation_timestamp();