package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgreSQL struct {
	Pool *pgxpool.Pool
}

// При запуске создаём таблицу, если она не была создана
func New(storagePath string) (*PostgreSQL, error) {
	ctx := context.Background()
	const op = "internal/storage/postgreSQL.NewUsers"

	Pool, err := pgxpool.New(ctx, storagePath)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to database: %w, %s", err, op)
	}

	_, err = Pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, username VARCHAR(50) UNIQUE NOT NULL, email VARCHAR(100) UNIQUE NOT NULL, password_hash VARCHAR(100) NOT NULL, created_at timestamp)")
	if err != nil {
		return nil, fmt.Errorf("cannot create users table: %w, %s", err, op)
	}

	_, err = Pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS notes (id SERIAL PRIMARY KEY, user_id integer UNIQUE NOT NULL,title text UNIQUE NOT NULL, content text, created_at timestamp, updated_at timestamp DEFAULT NOW()")
	if err != nil {
		return nil, fmt.Errorf("cannot create notes table: %w, %s", err, op)
	}

	return &PostgreSQL{Pool: Pool}, nil
}

// Создание нового пользователя
func (p *PostgreSQL) UserNew(username, email, passwordHash string) (int, string, time.Time, error) {
	const op = "internal/storage/postgreSQL.UserNew"
	if username == "" || email == "" || passwordHash == "" {
		return 0, "", time.Time{}, fmt.Errorf("username, email, and password cannot be empty")
	}
	ctx := context.Background()
	var id int
	var createdAt time.Time
	err := p.Pool.QueryRow(ctx,
		"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at",
		username, email, passwordHash).Scan(&id, &createdAt)
	if err != nil {
		return 0, "", time.Time{}, fmt.Errorf("cannot insert user: %w, %s", err, op)
	}
	return id, username, createdAt, nil
}

// Создание новой заметки
func (p *PostgreSQL) NoteNew(userID int, title, content string) (int, error) {
	const op = "internal/storage/postgreSQL.NoteNew"
	ctx := context.Background()
	var id int
	err := p.Pool.QueryRow(ctx, "INSERT INTO notes (user_id, title, content) VALUES ($1, $2, $3) RETURNING id", userID, title, content).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("cannot insert note: %w, %s", err, op)
	}
	return id, nil
}

// Получение всех заметок пользователя
type result struct {
	Id         int
	UserId     int
	Title      string
	Content    string
	Created_at string
	Updated_at string
}

func (p *PostgreSQL) NoteGet(user_id int) ([]result, error) {
	const op = "internal/storage/postgreSQL.NoteGet"
	ctx := context.Background()
	rows, err := p.Pool.Query(ctx, "SELECT * FROM notes WHERE user_id = $1", user_id)
	if err != nil {
		return nil, fmt.Errorf("cannot get notes: %w, %s", err, op)
	}
	results := []result{}
	for rows.Next() {
		var note result
		err := rows.Scan(&note.Id, &note.UserId, &note.Title, &note.Content, &note.Created_at, &note.Updated_at)
		if err != nil {
			return nil, fmt.Errorf("cannot scan note: %w, %s", err, op)
		}
		results = append(results, note)
	}
	return results, nil
}

// Получаем одну заметку
func (p *PostgreSQL) NoteGetOne(id int) (result, error) {
	const op = "internal/storage/postgreSQL.NoteGetOne"
	ctx := context.Background()
	var note result
	err := p.Pool.QueryRow(ctx, "SELECT * FROM notes WHERE id = $1", id).
		Scan(&note.Id, &note.UserId, &note.Title, &note.Content, &note.Created_at, &note.Updated_at)
	if err != nil {
		return note, fmt.Errorf("cannot get note: %w, %s", err, op)
	}
	return note, nil
}

// Обновляем заметку
func (p *PostgreSQL) NoteUpdate(id int, title, content string) error {
	const op = "internal/storage/postgreSQL.NoteUpdate"
	ctx := context.Background()
	_, err := p.Pool.Exec(ctx, "UPDATE notes SET title = $1, content = $2, updated_at = NOW() WHERE id = $3", title, content, id)
	if err != nil {
		return fmt.Errorf("cannot update note: %w, %s", err, op)
	}
	return nil
}

// Удаляем заметку
func (p *PostgreSQL) NoteDelete(id int) error {
	const op = "internal/storage/postgreSQL.NoteDelete"
	ctx := context.Background()
	_, err := p.Pool.Exec(ctx, "DELETE FROM notes WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("cannot delete note: %w, %s", err, op)
	}
	return nil
}

// Закрываем соединение с бд
func (p *PostgreSQL) Close() {
	p.Pool.Close()
}
