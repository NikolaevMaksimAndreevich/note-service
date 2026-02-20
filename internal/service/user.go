package service

import (
	"note_service/internal/handlers/user"
	"time"
)

// Интерфейс для storage
type UserStorage interface {
	UserNew(username, email, passwordHash string) (int, string, time.Time, error)
}

// Создаём сервис
type UserService struct {
	storage UserStorage
}

func NewUserService(storage UserStorage) *UserService {
	return &UserService{
		storage: storage,
	}
}

func (s *UserService) UserNew(req user.Request) (user.Response, error) {
	id, username, createdAt, err :=
		s.storage.UserNew(req.Username, req.Email, req.Password)

	if err != nil {
		return user.Response{}, err
	}

	return user.Response{
		ID:         int64(id),
		Username:   username,
		Created_at: createdAt.Format(time.RFC3339),
	}, nil
} //HTTP → handler → service → storage → DB
