package main

import (
	"log/slog"
	"net/http"
	rout "note_service/internal/NewRouter"
	"note_service/internal/storage"
	"os"
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

	r := rout.NewRouter(storageDB)

	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Error("server failed", slog.Any("error", err))
	}
	logger.Info("server started on :8080")
}

/*
Этот вариант main был до создания internal/NewRouter. Для автотестов нужен доступ к роутеру
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
	r.Use(middleware.RequestID)

	authHandler := &handlers.Handler{
		Storage: storageDB,
	}
	r.Post("/login", authHandler.LoginHandler) //Используем для получения нового токена(авторизация уже существующего пользователя)

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

*/
