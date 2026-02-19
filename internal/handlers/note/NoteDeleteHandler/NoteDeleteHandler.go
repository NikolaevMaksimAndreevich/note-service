package NoteDeleteHandler

import (
	"log/slog"
	"net/http"
	mid "note_service/internal/middleware"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type resultNote struct {
	Id         int
	UserId     int
	Title      string
	Content    string
	Created_at time.Time
	Updated_at time.Time
}
type Request struct {
	NoteID int `json:"note_id"`
}

type NoteDeleteHandler interface {
	NoteDelete(noteID int) error
	NoteGetOne(noteID int) (resultNote, error)
}

func New(log *slog.Logger, NoteDeleteHandler NoteDeleteHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers/note/NoteDeleteHandler"
		log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		userID, ok := r.Context().Value(mid.UserIDKey).(int)
		if !ok {
			log.Error("user not authorized")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "user not authorized"})
			return
		}

		noteIDStr := chi.URLParam(r, "note_id")
		noteID, err := strconv.Atoi(noteIDStr)
		if err != nil {
			log.Error("Invalid note ID", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid note ID"})
			return
		}
		if noteID <= 0 {
			log.Error("note_id must be greater than 0")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "note_id must be greater than 0"})
			return
		}
		//Отсюда начали проверять, принадлежит ли заметка данному пользователю
		note, err := NoteDeleteHandler.NoteGetOne(noteID)
		if err != nil {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "note not found"})
			return
		}

		if note.UserId != userID {
			render.Status(r, http.StatusForbidden) // 403
			render.JSON(w, r, map[string]string{"error": "access forbidden"})
			return
		}
		//Здесь закончили проверять пользователя и заметку, для этого создали структуру resultNote

		err = NoteDeleteHandler.NoteDelete(noteID)
		if err != nil {
			log.Error("Failed to delete note", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to delete note"})
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]string{"message": "Note deleted successfully"})
		log.Info("Note deleted successfully")
	}
}
