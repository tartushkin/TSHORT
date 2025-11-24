package repository

import (
	"context"

	"github.com/tartushkin/TSHORT.git/internal/model"
)

// вставка новых ссылок
func (r *Repo) InsertURL(ctx context.Context, list []byte) error {
	_, err := r.conn.ExecContext(ctx, `SELECT t_short.insert_urls($1)`, list)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repo) GetURLList(ctx context.Context) ([]*model.AliasFullCore, error) {
	rows, err := r.conn.QueryContext(ctx, `SELECT s_alias, s_full FROM t_short.t_list`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []*model.AliasFullCore{}
	for rows.Next() {
		url := &model.AliasFullCore{}
		if err := rows.Scan(&url.Alias, &url.OriginalURL); err != nil {
			return nil, err
		}
		list = append(list, url)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
