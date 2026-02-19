package noteGetOneHandler

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

type Request struct {
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

type NoteGetOneHandler interface {
	NoteGetOne(req Request) (Response, error)
}

func New(log *slog.Logger, NoteGetOneHandler NoteGetOneHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		const op = "handlers/note/noteGetOneHandler"
		log := log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		userID, ok := r.Context().Value(mid.UserIDKey).(int)
		if !ok {
			log.Error("user not authorized")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "user not authorized"})
			return
		}
		noteIDStr := chi.URLParam(r, "note_id") // если роут: /notes/{id}
		noteID, err := strconv.Atoi(noteIDStr)
		if err != nil {
			log.Error("invalid note id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "invalid note id"})
			return
		}
		req := Request{NoteID: noteID}

		resp, err := NoteGetOneHandler.NoteGetOne(req)
		if err != nil {
			log.Error("failed to get note", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "failed to get note"})
			return
		}

		if resp.Note.UserID != userID {
			log.Error("access forbidden")
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, map[string]string{"error": "access forbidden"})
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
		log.Info("note retrieved successfully")
	}
}
