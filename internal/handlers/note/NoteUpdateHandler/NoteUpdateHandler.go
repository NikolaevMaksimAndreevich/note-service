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
type resultNote struct {
	Id         int
	UserId     int
	Title      string
	Content    string
	Created_at time.Time
	Updated_at time.Time
}
type Response struct {
	UpdateAT time.Time `json:"update_at"`
}

type NoteUpdateHandler interface {
	NoteUpdate(noteID int, title, content string) (time.Time, error)
	NoteGetOne(noteID int) (resultNote, error)
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

		req := Request{ //отсюда вытащем нужные нам поля
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
		//Проверка доступа
		note, err := NoteUpdateHandler.NoteGetOne(req.NoteID)
		if err != nil {
			log.Error("note not found")
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "note not found"})
			return
		}
		if note.UserId != userID {
			log.Error("access forbidden")
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, map[string]string{"error": "access forbidden"})
			return
		}
		//конец проверки доступа
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
