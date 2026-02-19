package noteNewHandler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	ID_user int    `json:"id_user"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Response struct {
	ID         int       `json:"id"`
	ID_user    int       `json:"id_user"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Created_at time.Time `json:"created_at"`
	//Updated_at time.Time `json:"updated_at"`
}

type NoteNewHandler interface {
	NoteNew(req Request) (Response, error)
}

func New(log *slog.Logger, NoteNewHandler NoteNewHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		const op = "handlers/note/noteNewHandler.New"
		log := log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		userId, ok := r.Context().Value("user_id").(int) //По переданному токену получаем userID
		if !ok {
			log.Error("user not authorized")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "user not authorized"})
			return
		}

		var req Request
		req.ID_user = userId //Передали полученный userID
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "failed to decode request body"})
			return
		}
		log.Info("received request", slog.Any("request", req))

		idUser := req.ID_user
		title := req.Title
		content := req.Content

		if idUser == 0 || title == "" || content == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "id_user, title or content cannot be empty"})
			return
		}
		resp, err := NoteNewHandler.NoteNew(req)
		if err != nil {
			log.Error("failed to create note", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "failed to create note"})
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, resp)
		log.Info("note created", slog.Int("id", resp.ID))
	}
}
