CREATE OR REPLACE FUNCTION t_short.insert_urls(json_data JSON)
RETURNS VOID AS $$
DECLARE
    item JSON;
BEGIN
    FOR item IN SELECT * FROM json_array_elements(json_data)
    LOOP
        INSERT INTO t_short.t_list (n_corr_id, s_alias, s_full, s_user_id)
        VALUES (
            item->>'correlation_id',
            item->>'alias',
            item->>'original_url',
            item->>'user_id'
        );
    END LOOP;
END;
$$ LANGUAGE plpgsql;
