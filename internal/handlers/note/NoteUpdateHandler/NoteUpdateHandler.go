package NoteUpdateHandler

import (
	"encoding/json"
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
	UserID  int    `json:"user_id"`
	NoteID  int    `json:"note_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Response struct {
	UpdateAT time.Time `json:"update_at"`
}

type NoteUpdateHandler interface {
	NoteUpdate(noteID int, title, content string) (time.Time, error)
}

func New(log *slog.Logger, NoteUpdateHandler NoteUpdateHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers/note/NoteUpdateHandler"
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
			log.Error("invalid note id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "invalid note id"})
			return
		}
		if noteID <= 0 {
			log.Error("invalid note id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "invalid note id"})
			return
		}

		var body struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}

		err = json.NewDecoder(r.Body).Decode(&body) //Читаем JSON напрямую. После этого body.Title и body.Content содержат нужные значения.
		if err != nil {
			log.Error("invalid JSON")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "invalid JSON"})
			return
		}

		req := Request{
			UserID:  userID,
			NoteID:  noteID,
			Title:   body.Title,
			Content: body.Content,
		}
		if req.Title == "" || req.Content == "" {
			log.Error("title or content cannot be empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "title or content cannot be empty"})
			return
		}

		updateAt, err := NoteUpdateHandler.NoteUpdate(req.NoteID, req.Title, req.Content)
		if err != nil {
			log.Error("failed to update note")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "failed to update note"})
			return
		}

		resp := Response{UpdateAT: updateAt}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
		log.Info("note updated successfully")
	}
}
