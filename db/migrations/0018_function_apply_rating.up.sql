CREATE FUNCTION apply_rating(table_name TEXT)
RETURNS TRIGGER AS 
$$
DECLARE
    delta INTEGER;
BEGIN
    -- Delta relates to how the final value changed in relation to the old onde
    -- 1 → -1 results in a delta of -2 (remove +1, apply -1)   
    IF TG_OP = 'INSERT' THEN
        delta := NEW.rating_value;
    ELSIF TG_OP = 'UPDATE' THEN
        delta := NEW.rating_value - OLD.rating_value;
    ELSE
        RETURN NEW;
    END IF;

    EXECUTE format(
        'UPDATE %I SET rating = rating + $1 WHERE id = $2',
        table_name
    )
    USING delta, NEW.target_id;
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;