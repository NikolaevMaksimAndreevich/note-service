package main

import (
	"log/slog"
	"net/http"
	"note_service/internal/handlers/note/noteDeleteHandler"
	"note_service/internal/handlers/note/noteGetOneHandler"
	"note_service/internal/handlers/note/noteNewHandler"
	"note_service/internal/handlers/note/noteUpdateHandler"
	"note_service/internal/handlers/note/notesGetHandler"
	"note_service/internal/handlers/user"
	mid "note_service/internal/middleware"
	"note_service/internal/service"
	"note_service/internal/storage"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	logger := slog.Default()

	str := "user=postgres password=postgres dbname=notes sslmode=disable"
	storageDB, err := storage.New(str)
	if err != nil {
		logger.Error("failed to create storage", slog.Any("error", err))
		os.Exit(1)
		return
	}

	userServiceNew := service.NewUserService(storageDB) ////userHandler := user.New(logger, userServiceNew) можно так добавить
	noteServiceNew := &service.NoteServiceNew{Store: storageDB}
	noteServiceDelete := &service.NoteServiceDel{Store: storageDB}
	noteServiceUpdate := &service.NoteServiceUpd{Store: storageDB}
	notesServiceGet := &service.NotesServiceGet{Store: storageDB}
	noteGetOneService := &service.NoteServiceGetOne{Store: storageDB}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Post("/users", user.New(logger, userServiceNew))

	r.Group(func(r chi.Router) {
		r.Use(mid.JWTMiddleware)
		r.Post("/users/{id}/notes", noteNewHandler.New(logger, noteServiceNew))
		r.Get("/users/{id}/notes", notesGetHandler.New(logger, notesServiceGet))
		r.Get("/users/{id}/notes/{note_id}", noteGetOneHandler.New(logger, noteGetOneService))
		r.Put("/users/{id}/notes/{note_id}", noteUpdateHandler.New(logger, noteServiceUpdate))
		r.Delete("/users/{id}/notes/{note_id}", noteDeleteHandler.New(logger, noteServiceDelete))
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Error("server failed", slog.Any("error", err))
	}
	logger.Info("server started on :8080")
}
