package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"github.com/tartushkin/TSHORT.git/internal/model"
)

// InsertURLJson - вставка списком.
func (r *Repo) InsertURLJson(ctx context.Context, list []byte) error {
	_, err := r.conn.ExecContext(ctx, `SELECT t_short.insert_urls($1)`, list)
	if err != nil {
		return err
	}
	return nil
}

// InsertURL - вставка новой записи.
func (r *Repo) InsertURL(ctx context.Context, couple *model.AliasFullCore) error {
	_, err := r.conn.ExecContext(ctx, `
	INSERT INTO t_short.t_list(s_alias, s_full, n_corr_id, s_user_id)
	VALUES ($1,$2,$3,$4);
	`, couple.Alias, couple.OriginalURL, "-", couple.UserID)
	if err != nil {
		// Проверяем, является ли ошибка ошибкой уникальности
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("ERRCONFLICT")
			}
		}
		return fmt.Errorf("ошибка при вставке URL: %w", err)
	}
	return nil
}

// LoadCache - подгрузка всех записей с базы.
func (r *Repo) LoadCache(ctx context.Context) ([]*model.AliasFullCore, error) {
	rows, err := r.conn.QueryContext(ctx, `SELECT s_alias, s_full, s_user_id, b_is_deleted FROM t_short.t_list`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []*model.AliasFullCore{}
	for rows.Next() {
		url := &model.AliasFullCore{}
		if err := rows.Scan(&url.Alias, &url.OriginalURL, &url.UserID, &url.DeletedFlag); err != nil {
			return nil, err
		}
		list = append(list, url)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

// GetAlias - получение алиса по полному url.
func (r *Repo) GetAlias(ctx context.Context, orig string) (string, error) {
	var alias string
	err := r.conn.QueryRowContext(ctx, `SELECT s_alias FROM t_short.t_list WHERE s_full = $1`, orig).Scan(&alias)
	if err != nil {
		return "", err
	}
	return alias, nil
}

// DeleteURL - удаление списка записей по их алиасу.
func (r *Repo) DeleteURL(ctx context.Context, aliasList []string) error {
	query := "UPDATE t_short.t_list SET b_is_deleted = true WHERE s_alias = ANY($1)"
	_, err := r.conn.ExecContext(ctx, query, pq.Array(aliasList))
	if err != nil {
		return err
	}

	return nil
}

// GetOriginalURL - получить оригинальный url по алиасу.
func (r *Repo) GetOriginalURL(ctx context.Context, alias string) (*model.AliasFullCore, error) {
	var coupe model.AliasFullCore
	err := r.conn.QueryRowContext(ctx, `SELECT s_full, b_is_deleted FROM t_short.t_list WHERE s_alias = $1`, alias).Scan(&coupe.OriginalURL, &coupe.DeletedFlag)
	if err != nil {
		return nil, err
	}
	return &coupe, nil
}

// GetUserURL - получить список записей по id пользователя.
func (r *Repo) GetUserURL(ctx context.Context, userID string) ([]*model.AliasFullCore, error) {
	if r.conn == nil {
		return nil, nil
	}
	rows, err := r.conn.QueryContext(ctx, `SELECT s_alias, s_full, b_is_deleted FROM t_short.t_list WHERE s_user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []*model.AliasFullCore{}
	for rows.Next() {
		url := &model.AliasFullCore{}
		if err := rows.Scan(&url.Alias, &url.OriginalURL, &url.DeletedFlag); err != nil {
			return nil, err
		}
		list = append(list, url)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

// DeleteMarkedURLs - удаление записей.
func (r *Repo) DeleteMarkedURLs(ctx context.Context) error {
	_, err := r.conn.ExecContext(ctx, `delete from t_short.t_list where b_is_deleted = true`)
	if err != nil {
		return err
	}

	return nil
}
