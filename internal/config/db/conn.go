package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewConnection(ps string) (*sql.DB, error) {

	db, err := sql.Open("pgx", ps)
	if err != nil {
		return nil, err
	}

	dbURL, err := convertDSNToURL(ps)
	if err != nil {
		return nil, err
	}
	err = migration(dbURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}
func migration(ps string) error {
	// Запуск миграций при старте приложения
	m, err := migrate.New(
		"file://./migrations",
		ps)
	if err != nil {
		//log.Fatalf("Ошибка создания объекта миграции: %v", err)
		log.Println("Папка migrations не найдена, пропускаем миграции, ошибка - " + err.Error())
		return err
	}

	// Применение миграций
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// convertDSNToURL преобразует строку подключения в формате key=value в URL-формат
func convertDSNToURL(dsn string) (string, error) {
	// Разбиваем строку на пары key=value
	parts := strings.Fields(dsn)
	params := make(map[string]string)
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			params[kv[0]] = kv[1]
		}
	}

	// Формируем URL
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(params["user"], params["password"]),
		Host:   fmt.Sprintf("%s:%s", params["host"], params["port"]),
		Path:   params["dbname"],
	}
	// Добавляем параметры запроса
	query := url.Values{}
	if sslmode, ok := params["sslmode"]; ok {
		query.Add("sslmode", sslmode)
	}
	u.RawQuery = query.Encode()

	return u.String(), nil
}
