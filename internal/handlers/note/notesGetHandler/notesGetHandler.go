package notesGetHandler

import (
	mid "note_service/internal/middleware"
	"note_service/internal/storage"

	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	UserID int `json:"user_id"`
}

/*
type Note struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Response struct {
	Notes []storage.ResultNote `json:"notes"`
}*/

type NotesGetHandler interface {
	NotesGet(req Request) ([]storage.ResultNote, error)
}

func New(log *slog.Logger, NotGetHandler NotesGetHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		const op = "internal/handlers/note/notesGetHandler/NoteGet"
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
		req := Request{UserID: userID}
		resp, err := NotGetHandler.NotesGet(req)
		if err != nil {
			log.Error("failed to get notes", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "failed to get notes"})
			return
		}

		if req.UserID != userID {
			log.Error("access forbidden")
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, map[string]string{"error": "access forbidden"})
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
		log.Info("notes retrieved successfully")
	}
}
