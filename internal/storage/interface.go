package storage

import "time"

type Storage interface {
	UserNew(username, email, passwordHash string) (int, string, time.Time, error)
	NoteNew(userID int, title, content string) (int, error)
	NotesGet(userID int) ([]ResultNote, error)
	NoteGetOne(id int) (ResultNote, error)
	NoteUpdate(id int, title, content string) (time.Time, error)
	NoteDelete(id int) error
	GetUserByEmail(email string) (ResultUser, error)
}

/*
Сделали интерфейс для работы с postgreSQL. Что бы сделать тесты
Из-за этого пришлось поменять все интерфейсы в service
*/
