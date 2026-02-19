package NoteDeleteHandler

import (
	"log/slog"
	"net/http"
	mid "note_service/internal/middleware"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	NoteID int `json:"note_id"`
}

type NoteDeleteHandler interface {
	NoteDelete(noteID int) error
}

func New(log *slog.Logger, NoteDeleteHandler NoteDeleteHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers/note/NoteDeleteHandler"
		log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		_, ok := r.Context().Value(mid.UserIDKey).(int)
		if !ok {
			log.Error("UserID not found in context")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "UserID not found in context"})
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
