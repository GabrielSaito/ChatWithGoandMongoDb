package main

import (
	"chat-go/handlers_"
	"chat-go/services"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	log.SetFormatter(&logrus.JSONFormatter{})

	services.ConnectMongoDB()

	http.HandleFunc("/create-room", handlers_.CreateRoom)
	http.HandleFunc("/rooms", handlers_.Rooms)
	http.HandleFunc("/join-room", handlers_.JoinRoom)
	http.HandleFunc("/generate-token", handlers_.GenerateToken)
	http.HandleFunc("/upload-image", handlers_.UploadImage)

	http.HandleFunc("/ws", handlers_.HandleConnections)

	go handlers_.HandleMessages()

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(http.DefaultServeMux)

	log.Info("Servidor iniciado na porta :7120")
	err := http.ListenAndServe(":7120", corsHandler)
	if err != nil {
		log.Fatal("Erro ao iniciar servidor: ", err)
	}
}
