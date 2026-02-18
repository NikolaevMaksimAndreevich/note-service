package user

import (
	"log/slog"
	"net/http"
	"note_service/internal/authorization"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type Request struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	Created_at string `json:"created_at"`
}

type NewUserHandler interface {
	UserNew(req Request) (Response, error)
}

func New(log *slog.Logger, NewUserHandler NewUserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		const op = "handlers/user/handler_user_registration.New"
		log = log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))
		//request_id — уникальный идентификатор запроса, который генерирует chi middleware. Полезно для связывания всех логов одного HTTP-запроса.

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "failed to decode request body"})
			return
		}
		log.Info("received request", slog.Any("request", req))

		if req.Username == "" || req.Password == "" || req.Email == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "username cannot be empty"})
			return
		}

		//Хешируем пароль
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to hash password")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "failed to hash password"})
			return
		}
		req.Password = string(hash)

		// Сохраняем пользователя в базе
		resp, err := NewUserHandler.UserNew(req)
		if err != nil {
			log.Error("failed to create user", slog.Any("error", err))
			render.Status(r, http.StatusInternalServerError) //Статус 500 - если ошибка с бд
			render.JSON(w, r, map[string]string{"error": err.Error()})
			return
		}

		//  Генерируем JWT
		token, err := authorization.GenerateToken(int(resp.ID))
		if err != nil {
			log.Error("failed to generate token", slog.Any("error", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "failed to generate token"})
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, map[string]interface{}{
			"user":  resp,
			"token": token,
		}) //Клиент получает ID, Username и Created_at и токен
		log.Info("user created", slog.Any("response", resp))

	}
}
