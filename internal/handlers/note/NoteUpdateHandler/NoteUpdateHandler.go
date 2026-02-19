package NoteUpdateHandler

import "time"

type Request struct {
	UserID int `json:"user_id"`
	NoteID int `json:"note_id"`
}

type Note struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Response struct {
	Note Note `json:"note"`
}

type NoteUpdateHandler interface {
	NoteUpdate(req Request) (Response, error)
}
