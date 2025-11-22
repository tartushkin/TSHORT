package repository

import (
	"context"

	"github.com/tartushkin/TSHORT.git/internal/model"
)

// вставка новых ссылок
func (r *Repo) InsertURL(ctx context.Context, fullURL, alias string) error {

	_, err := r.conn.ExecContext(ctx, `INSERT INTO t_short.t_list (s_alias, s_full) VALUES ($1, $2)`, alias, fullURL)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repo) GetURLList(ctx context.Context) ([]*model.StorageURL, error) {
	rows, err := r.conn.QueryContext(ctx, `SELECT s_alias, s_full FROM t_short.t_list`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []*model.StorageURL{}
	for rows.Next() {
		url := &model.StorageURL{}
		if err := rows.Scan(&url.Alias, &url.Original); err != nil {
			return nil, err
		}
		list = append(list, url)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
