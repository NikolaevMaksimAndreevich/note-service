package service

import (
	"note_service/internal/storage"
	"time"
)

type NoteServiceUpd struct {
	Store storage.Storage
}

func (s *NoteServiceUpd) NoteUpdate(id int, title, content string) (time.Time, error) {
	return s.Store.NoteUpdate(id, title, content)
}

func (s *NoteServiceUpd) NoteGetOne(id int) (storage.ResultNote, error) {
	return s.Store.NoteGetOne(id)
}
