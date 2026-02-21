package service

import "note_service/internal/storage"

// Обёртка для интерфейса NoteDeleteHandler
type NoteServiceDel struct {
	Store storage.Storage
}

func (s *NoteServiceDel) NoteDelete(ID int) error {
	return s.Store.NoteDelete(ID)
}

func (s *NoteServiceDel) NoteGetOne(ID int) (storage.ResultNote, error) {
	return s.Store.NoteGetOne(ID)
}
