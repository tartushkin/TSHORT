package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/tartushkin/TSHORT.git/internal/model"
)

// вставка новых ссылок
func (r *Repo) InsertURLJson(ctx context.Context, list []byte) error {
	_, err := r.conn.ExecContext(ctx, `SELECT t_short.insert_urls($1)`, list)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repo) InsertURL(ctx context.Context, couple *model.AliasFullCore) error {
	_, err := r.conn.ExecContext(ctx, `
	INSERT INTO t_short.t_list(s_alias, s_full, n_corr_id, s_user_id)
	VALUES ($1,$2,$3,$4);
	`, couple.Alias, couple.OriginalURL, "-", "undefined")
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

func (r *Repo) GetAlias(ctx context.Context, orig string) (string, error) {
	var alias string
	err := r.conn.QueryRowContext(ctx, `SELECT s_alias FROM t_short.t_list WHERE s_full = $1`, orig).Scan(&alias)
	if err != nil {
		return "", err
	}
	return alias, nil
}

func (r *Repo) DeleteURL(ctx context.Context, delStr string) error {
	query := "UPDATE t_short.t_list SET b_is_deleted = true WHERE s_alias IN (?)"
	query = strings.Replace(query, "?", delStr, 1)
	_, err := r.conn.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetOriginalURL(ctx context.Context, alias string) (string, error) {
	var original string
	err := r.conn.QueryRowContext(ctx, `SELECT s_full FROM t_short.t_list WHERE s_alias = $1`, alias).Scan(&original)
	if err != nil {
		return "", err
	}
	return original, nil
}
