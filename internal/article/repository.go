package article

import (
	"database/sql"
	"errors"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db} }

func (r *Repository) Create(a Article) error {
	_, err := r.db.Exec("INSERT INTO posts (title, content, category, status) VALUES (?,?,?,?)",
		a.Title, a.Content, a.Category, a.Status)
	return err
}

func (r *Repository) List(limit, offset int) ([]Article, error) {
	rows, err := r.db.Query("SELECT id, title, content, category, status FROM posts ORDER BY id LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Article{}
	for rows.Next() {
		var a Article
		if err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.Category, &a.Status); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *Repository) Get(id int) (Article, error) {
	var a Article
	err := r.db.QueryRow("SELECT id, title, content, category, status FROM posts WHERE id=?", id).
		Scan(&a.ID, &a.Title,&a.Content, &a.Category, &a.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return a, ErrNotFound
	}
	return a, err
}

func (r *Repository) Update(id int, a Article) error {
	res, err := r.db.Exec("UPDATE posts SET title=?, content=?, category=?, status=? WHERE id=?",
		a.Title, a.Content, a.Category, a.Status, id)
	if err != nil {
		return err
	}
	// MySQL melaporkan 0 affected rows jika data tidak berubah, jadi cek keberadaan baris terpisah.
	if n, _ := res.RowsAffected(); n == 0 {
		_, err := r.Get(id)
		return err
	}
	return nil
}

func (r *Repository) Delete(id int) error {
	res, err := r.db.Exec("DELETE FROM posts WHERE id=?", id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}
