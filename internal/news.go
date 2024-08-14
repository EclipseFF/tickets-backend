package internal

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type NewsRepo struct {
	DB *pgxpool.Pool
}

type News struct {
	Id          *int       `json:"id" form:"id"`
	Name        *string    `json:"name" form:"name"`
	Images      []*string  `json:"images" form:"images"`
	Description *string    `json:"description" form:"description"`
	CreatedAt   *time.Time `json:"created_at" form:"created_at"`
}

func (m *NewsRepo) GetAllNews() ([]*News, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	rows, err := tx.Query(context.Background(), `SELECT id, name, images, description, created_at FROM news`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	news := make([]*News, 0)
	for rows.Next() {
		var n News
		err = rows.Scan(&n.Id, &n.Name, &n.Images, &n.Description, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		news = append(news, &n)
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return news, nil
}

func (m *NewsRepo) GetNewsById(id *int) (*News, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	row := tx.QueryRow(context.Background(), `SELECT id, name, images, description, created_at FROM news WHERE id = $1`, id)
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	var n News
	err = row.Scan(&n.Id, &n.Name, &n.Images, &n.Description, &n.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (m *NewsRepo) CreateNews(n *News) (*News, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	err = tx.QueryRow(context.Background(), `INSERT INTO news (id, name, description, created_at) VALUES (default, $1, $2, now()) RETURNING id, name, images, description, created_at`, n.Name, n.Description).Scan(&n.Id, &n.Name, &n.Images, &n.Description, &n.CreatedAt)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (m *NewsRepo) UpdateNews(n *News) error {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(context.Background(), `UPDATE news SET name = $1, images = $2, description = $3, created_at = $4 WHERE id = $5`, n.Name, n.Images, n.Description, n.CreatedAt, n.Id)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (m *NewsRepo) GetPaginatedNews(limit int, offset int) ([]*News, *int, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(context.Background())
	rows, err := tx.Query(context.Background(), `SELECT id, name, images, description, created_at FROM news ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	err = tx.Commit(context.Background())

	news := make([]*News, 0)
	for rows.Next() {
		var n News
		err = rows.Scan(&n.Id, &n.Name, &n.Images, &n.Description, &n.CreatedAt)
		if err != nil {
			return nil, nil, err
		}
		news = append(news, &n)
	}

	if err != nil {
		return nil, nil, err
	}
	var totalRows int

	tx, err = m.DB.Begin(context.Background())
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(context.Background())

	err = tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM news`).Scan(&totalRows)
	if err != nil {
		return nil, nil, err
	}
	totalPages := totalRows / limit
	if totalRows%limit != 0 {
		totalPages++
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, nil, err
	}
	return news, &totalPages, nil
}

func (m *NewsRepo) DeleteNews(id *int) error {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(context.Background(), `DELETE FROM news WHERE id = $1`, id)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (m *NewsRepo) GetLatestNews(limit int) ([]*News, error) {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	rows, err := tx.Query(context.Background(), `SELECT id, name, images, description, created_at FROM news ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	news := make([]*News, 0)
	for rows.Next() {
		var n News
		err = rows.Scan(&n.Id, &n.Name, &n.Images, &n.Description, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		news = append(news, &n)
	}
	return news, nil
}

func (m *NewsRepo) SetNewsImages(images []*string, id *int) error {
	tx, err := m.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `UPDATE news SET images = $1 WHERE id = $2`, images, id)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil

}
