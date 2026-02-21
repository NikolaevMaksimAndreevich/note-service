package newRouter

import (
	"log/slog"
	"net/http"
	"note_service/internal/handlers"
	"note_service/internal/handlers/note/noteDeleteHandler"
	"note_service/internal/handlers/note/noteGetOneHandler"
	"note_service/internal/handlers/note/noteNewHandler"
	"note_service/internal/handlers/note/noteUpdateHandler"
	"note_service/internal/handlers/note/notesGetHandler"
	"note_service/internal/handlers/user"
	mid "note_service/internal/middleware"
	"note_service/internal/service"
	"note_service/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(store storage.Storage) http.Handler {
	logger := slog.Default()

	userServiceNew := service.NewUserService(store) ////userHandler := user.New(logger, userServiceNew) можно так добавить
	noteServiceNew := &service.NoteServiceNew{Store: store}
	noteServiceDelete := &service.NoteServiceDel{Store: store}
	noteServiceUpdate := &service.NoteServiceUpd{Store: store}
	notesServiceGet := &service.NotesServiceGet{Store: store}
	noteGetOneService := &service.NoteServiceGetOne{Store: store}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	authHandler := &handlers.Handler{
		Storage: store,
	}
	r.Post("/login", authHandler.LoginHandler) //Используем для получения нового токена(авторизация уже существующего пользователя)

	r.Post("/users", user.New(logger, userServiceNew))

	r.Group(func(r chi.Router) {
		r.Use(mid.JWTMiddleware)
		r.Post("/users/{id}/note", noteNewHandler.New(logger, noteServiceNew))
		r.Get("/users/{id}/notes", notesGetHandler.New(logger, notesServiceGet))
		r.Get("/users/{id}/note/{note_id}", noteGetOneHandler.New(logger, noteGetOneService))
		r.Put("/users/{id}/note/{note_id}", noteUpdateHandler.New(logger, noteServiceUpdate))
		r.Delete("/users/{id}/note/{note_id}", noteDeleteHandler.New(logger, noteServiceDelete))
	})
	return r
}
