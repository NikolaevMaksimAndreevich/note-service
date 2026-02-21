package internal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"note_service/internal/mocks"
	rout "note_service/internal/newRouter"
)

// Объявление тестовой функции в Go
func TestCreateUser(t *testing.T) {
	mockStore := &mocks.MockStorage{}   //Создаём мок хранилища. MockStorage реализует интерфейс storage.Storage, но не ходит в реальную БД, возвращает заранее заданные значения
	router := rout.NewRouter(mockStore) // Вызываем функцию, которая создаёт наш HTTP-роутер (chi.Router) с нужными эндпоинтами.  Передаём моковое хранилище вместо реальной БД.Теперь все обработчики (/users, /login, /notes) внутри роутера будут работать через mockStore.

	req := httptest.NewRequest( //Создаём искусственный HTTP-запрос
		http.MethodPost,
		"/users",
		strings.NewReader(`{"username":"u","email":"u@mail.com","password":"123"}`),
	)
	req.Header.Set("Content-Type", "application/json") //Устанавливаем заголовок Content-Type, чтобы сервер понимал, что тело запроса в формате JSON.  Важно, иначе хендлер может возвращать ошибку из-за неподходящего типа данных.

	w := httptest.NewRecorder() //Создаём искусственный ResponseWriter. w будет собирать всё, что сервер “отправляет” в ответ (статус-код, тело ответа, заголовки). Фактически это аналог того, что браузер или Postman получает при настоящем HTTP-запросе.
	router.ServeHTTP(w, req)    //Отправляем наш искусственный запрос в роутер.

	if w.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, w.Code)
	} //Проверяем HTTP-код ответа.
}

func TestNoteNewHandler(t *testing.T) {
	mockStore := &mocks.MockStorage{}
	router := rout.NewRouter(mockStore)

	req := httptest.NewRequest(http.MethodPost, "/users/4/note",
		strings.NewReader(`{"title":"TEST","content":"test test test"}`))
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJpc3MiOiJub3RlX3NlcnZpY2UiLCJzdWIiOiI0IiwiZXhwIjoxNzcxNzk2NjkyLCJpYXQiOjE3NzE3MTAyOTJ9.YcUeLeDZBzMx0pxmq-JUyNRmFn6KdLI6puWaJ6n4JT8")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestNoteGetOneHandler(t *testing.T) {
	mockStore := &mocks.MockStorage{}
	router := rout.NewRouter(mockStore)

	req := httptest.NewRequest(http.MethodGet, "/users/4/note/1",
		nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJpc3MiOiJub3RlX3NlcnZpY2UiLCJzdWIiOiI0IiwiZXhwIjoxNzcxNzk2NjkyLCJpYXQiOjE3NzE3MTAyOTJ9.YcUeLeDZBzMx0pxmq-JUyNRmFn6KdLI6puWaJ6n4JT8")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestNotesGetHandler(t *testing.T) {
	mockStore := &mocks.MockStorage{}
	router := rout.NewRouter(mockStore)

	req := httptest.NewRequest(http.MethodGet, "/users/4/notes",
		nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJpc3MiOiJub3RlX3NlcnZpY2UiLCJzdWIiOiI0IiwiZXhwIjoxNzcxNzk2NjkyLCJpYXQiOjE3NzE3MTAyOTJ9.YcUeLeDZBzMx0pxmq-JUyNRmFn6KdLI6puWaJ6n4JT8")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestNoteUpdateHandler(t *testing.T) {
	mockStore := &mocks.MockStorage{}
	router := rout.NewRouter(mockStore)

	req := httptest.NewRequest(http.MethodPut, "/users/4/note/1",
		strings.NewReader(`{"title":"TEST","content":"test test test"}`))
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJpc3MiOiJub3RlX3NlcnZpY2UiLCJzdWIiOiI0IiwiZXhwIjoxNzcxNzk2NjkyLCJpYXQiOjE3NzE3MTAyOTJ9.YcUeLeDZBzMx0pxmq-JUyNRmFn6KdLI6puWaJ6n4JT8")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestNoteDeleteHandler(t *testing.T) {
	mockStore := &mocks.MockStorage{}
	router := rout.NewRouter(mockStore)

	req := httptest.NewRequest(http.MethodDelete, "/users/4/note/1",
		nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJpc3MiOiJub3RlX3NlcnZpY2UiLCJzdWIiOiI0IiwiZXhwIjoxNzcxNzk2NjkyLCJpYXQiOjE3NzE3MTAyOTJ9.YcUeLeDZBzMx0pxmq-JUyNRmFn6KdLI6puWaJ6n4JT8")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}
}
