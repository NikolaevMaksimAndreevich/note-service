package mocks

import (
	"errors"
	"note_service/internal/storage"
	"time"
)

//Создаём мок для тестов. Полностью написан ИИ
//Теперь мок полностью реализует интерфейс.

type MockStorage struct{}

func (m *MockStorage) UserNew(username, email, passwordHash string) (int, string, time.Time, error) {
	return 1, username, time.Now(), nil
}

func (m *MockStorage) NoteNew(userID int, title, content string) (int, error) {
	return 1, nil
}

func (m *MockStorage) NotesGet(userID int) ([]storage.ResultNote, error) {
	return []storage.ResultNote{}, nil
}

func (m *MockStorage) NoteGetOne(id int) (storage.ResultNote, error) {
	return storage.ResultNote{}, errors.New("not found")
}

func (m *MockStorage) NoteUpdate(id int, title, content string) (time.Time, error) {
	return time.Now(), nil
}

func (m *MockStorage) NoteDelete(id int) error {
	return nil
}

func (m *MockStorage) GetUserByEmail(email string) (storage.ResultUser, error) {
	return storage.ResultUser{
		Id:           1,
		Email:        email,
		PasswordHash: "$2a$10$hash", // если используешь bcrypt
	}, nil
}
