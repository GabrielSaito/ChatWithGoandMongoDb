package handlers_

import (
	"chat-go/services"
	"encoding/json"
	"net/http"
)

func GenerateToken(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "O nome de usuário é obrigatório", http.StatusBadRequest)
		return
	}

	token := services.GenerateToken(username)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
