-- Создание таблицы
CREATE table IF NOT EXISTS t_short.t_url (
    n_id SERIAL PRIMARY KEY,
    s_alias VARCHAR(50),
    s_full_url VARCHAR(100)
);

COMMENT ON TABLE t_short.t_url IS 'Таблица ссылок';
COMMENT ON COLUMN t_short.t_url.n_id IS 'Идентификатор URL';
COMMENT ON COLUMN t_short.t_url.s_alias IS 'Псевдоним URL';
COMMENT ON COLUMN t_short.t_url.s_full_url IS 'Полный URL';
