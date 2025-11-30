CREATE OR REPLACE FUNCTION t_short.insert_urls(json_data JSON)
RETURNS VOID AS $$
DECLARE
    item JSON;
BEGIN
    FOR item IN SELECT * FROM json_array_elements(json_data)
    LOOP
        INSERT INTO t_short.t_list (n_corr_id, s_alias, s_full)
        VALUES (
            CASE WHEN item->>'corr_id' IS NOT NULL AND item->>'corr_id' <> '' THEN item->>'corr_id' ELSE NULL END,
            item->>'alias',
            item->>'original'
        );
    END LOOP;
END;
$$ LANGUAGE plpgsql;
