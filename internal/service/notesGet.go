package service

import (
	"note_service/internal/handlers/note/notesGetHandler"
	"note_service/internal/storage"
)

type NotesServiceGet struct {
	Store *storage.PostgreSQL
}

func (s *NotesServiceGet) NotesGet(req notesGetHandler.Request) ([]storage.ResultNote, error) {
	return s.Store.NotesGet(req.UserID)
}
