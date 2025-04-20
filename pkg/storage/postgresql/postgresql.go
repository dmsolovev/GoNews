package postgres

import (
	"GoNews/pkg/storage"
	"database/sql"

	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func New(connstr string) (*Store, error) {
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(`
		SELECT p.id, p.title, p.content, p.author_id, a.name, p.created_at, p.published_at
		FROM posts p
		JOIN authors a ON p.author_id = a.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []storage.Post
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.AuthorID, &p.AuthorName, &p.CreatedAt, &p.PublishedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (s *Store) AddPost(p storage.Post) error {
	_, err := s.db.Exec(`
		INSERT INTO posts (title, content, author_id, created_at, published_at)
		VALUES ($1, $2, $3, $4, $5)
	`, p.Title, p.Content, p.AuthorID, p.CreatedAt, p.PublishedAt)
	return err
}

func (s *Store) UpdatePost(p storage.Post) error {
	_, err := s.db.Exec(`
		UPDATE posts
		SET title = $1, content = $2, author_id = $3, created_at = $4, published_at = $5
		WHERE id = $6
	`, p.Title, p.Content, p.AuthorID, p.CreatedAt, p.PublishedAt, p.ID)
	return err
}

func (s *Store) DeletePost(p storage.Post) error {
	_, err := s.db.Exec("DELETE FROM posts WHERE id = $1", p.ID)
	return err
}
