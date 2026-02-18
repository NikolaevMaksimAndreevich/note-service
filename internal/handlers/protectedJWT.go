package handlers

import (
	"encoding/json"
	"net/http"
	"note_service/internal/middleware"
)

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middleware.UserIDKey).(int)

	response := map[string]interface{}{
		"message": "You are authorized",
		"user_id": userID,
	}

	json.NewEncoder(w).Encode(response)
}
