package service

import (
	"note_service/internal/handlers/note/noteGetOneHandler"
	"note_service/internal/storage"
)

type NoteServiceGetOne struct {
	Store storage.Storage
}

func (s *NoteServiceGetOne) NoteGetOne(req noteGetOneHandler.Request) (noteGetOneHandler.Response, error) {
	n, err := s.Store.NoteGetOne(req.NoteID)
	if err != nil {
		return noteGetOneHandler.Response{}, err
	}

	resp := noteGetOneHandler.Response{
		Note: noteGetOneHandler.Note{
			ID:        n.Id,
			UserID:    n.UserId,
			Title:     n.Title,
			Content:   n.Content,
			CreatedAt: n.Created_at,
			UpdatedAt: n.Updated_at,
		},
	}

	return resp, nil
}
