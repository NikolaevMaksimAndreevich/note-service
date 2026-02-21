package service

import (
	"note_service/internal/handlers/note/noteNewHandler"
	"note_service/internal/storage"
	"time"
)

type NoteServiceNew struct {
	Store storage.Storage
}

func (s *NoteServiceNew) NoteNew(req noteNewHandler.Request) (noteNewHandler.Response, error) {
	// Создаём заметку и получаем ID
	id, err := s.Store.NoteNew(req.ID_user, req.Title, req.Content)
	if err != nil {
		return noteNewHandler.Response{}, err
	}

	resp := noteNewHandler.Response{
		ID:         id,          // новый ID из БД
		ID_user:    req.ID_user, // берём из запроса / JWT
		Title:      req.Title,   // берём из запроса
		Content:    req.Content, // берём из запроса
		Created_at: time.Now(),  // текущее время
	}

	return resp, nil
}
