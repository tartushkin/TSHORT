-- Создание таблицы
CREATE table IF NOT EXISTS t_short.t_list (
    n_id SERIAL PRIMARY KEY,
    s_alias VARCHAR(50),
    s_full VARCHAR(100)
);

COMMENT ON TABLE t_short.t_list IS 'Таблица ссылок';
COMMENT ON COLUMN t_short.t_list.n_id IS 'Идентификатор URL';
COMMENT ON COLUMN t_short.t_list.s_alias IS 'Псевдоним URL';
COMMENT ON COLUMN t_short.t_list.s_full IS 'Полный URL';
