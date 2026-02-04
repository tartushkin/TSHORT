package db

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewConnection создаёт и инициализирует подключение к PostgreSQL.
//
// Выполняет:
//   - Открытие соединения с использованием драйвера pgx
//   - Запуск миграций из папки ./migrations
//
// Возвращает:
//   - *sql.DB: готовое соединение
//   - error: если не удалось подключиться или применить миграции
//
// Пример:
//
//	db, err := NewConnection("postgres://user:pass@localhost/db")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer db.Close()
func NewConnection(ps string) (*sql.DB, error) {

	db, err := sql.Open("pgx", ps)
	if err != nil {
		return nil, err
	}
	err = migration(ps)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// migration применяет SQL-миграции из папки ./migrations к базе данных.
//
// Использует библиотеку github.com/golang-migrate/migrate.
// Игнорирует ошибку migrate.ErrNoChange (нет новых миграций).
//
// Параметры:
//   - connectionString: строка подключения к PostgreSQL
//
// Возвращает:
//   - error: если миграции не удалось применить (кроме случая отсутствия изменений)
func migration(ps string) error {
	// Запуск миграций при старте приложения
	m, err := migrate.New(
		"file://./migrations",
		ps)
	if err != nil {
		log.Println("Ошибка создания объекта миграции:  " + err.Error() + ". Строка подключения - " + ps)
		return err
	}

	// Применение миграций
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
