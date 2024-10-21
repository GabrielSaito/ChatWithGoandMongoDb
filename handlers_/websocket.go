package handlers_

import (
	"chat-go/models"
	"chat-go/services"
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

var clients = make(map[*websocket.Conn]bool)
var rooms = make(map[string]map[*websocket.Conn]bool)
var broadcast = make(chan models.Message)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Erro ao atualizar para WebSocket: %v", err)
		return
	}
	defer ws.Close()

	clients[ws] = true
	room := r.URL.Query().Get("room")
	if room == "" {
		room = "default"
	}
	sendChatHistory(ws, room)

	if rooms[room] == nil {
		rooms[room] = make(map[*websocket.Conn]bool)
	}
	rooms[room][ws] = true

	for {
		var msg models.Message

		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Erro ao ler mensagem: %v", err)
			delete(clients, ws)
			delete(rooms[room], ws)
			break
		}

		err = services.SaveMessage(msg)
		if err != nil {
			log.Printf("Erro ao salvar a mensagem: %v", err)
		}

		broadcast <- msg
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast
		for client := range rooms[msg.Room] {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Erro ao enviar mensagem: %v", err)
				client.Close()
				delete(clients, client)
				delete(rooms[msg.Room], client)
			}
		}
	}
}

func sendChatHistory(ws *websocket.Conn, room string) {
	cursor, err := services.GetMessageCollection().Find(context.TODO(), bson.M{"room": room})
	if err != nil {
		log.Printf("Erro ao recuperar histórico de mensagens: %v", err)
		return
	}
	defer cursor.Close(context.TODO())

	var messages []models.Message
	for cursor.Next(context.TODO()) {
		var msg models.Message
		if err := cursor.Decode(&msg); err != nil {
			log.Printf("Erro ao decodificar mensagem: %v", err)
			continue
		}
		messages = append(messages, msg)
	}

	for _, msg := range messages {
		ws.WriteJSON(msg)
	}
}
