package handlers_

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"chat-go/models"
	"chat-go/services"
)

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	var room models.Room
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	if room.Password == "" {
		http.Error(w, "Senha obrigatória!", http.StatusBadRequest)
		return
	}

	room.ID = primitive.NewObjectID().Hex()
	room.Created = time.Now()

	collection := services.GetRoomCollection()
	_, err := collection.InsertOne(context.Background(), room)
	if err != nil {
		http.Error(w, "Não foi possível criar sala!", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(room)
}

func Rooms(w http.ResponseWriter, r *http.Request) {
	collection := services.GetRoomCollection()
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		http.Error(w, "Nçao foi possivel achar a sala", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var rooms []models.Room
	if err := cursor.All(context.Background(), &rooms); err != nil {
		http.Error(w, "Não foi possivel entrar na sala", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rooms)
}

func JoinRoom(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RoomID   string `json:"room_id"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	collection := services.GetRoomCollection()
	var room models.Room
	err := collection.FindOne(context.Background(), bson.M{"_id": request.RoomID}).Decode(&room)
	if err != nil {
		http.Error(w, "Sala não encontrada!", http.StatusNotFound)
		return
	}

	if room.Password != request.Password {
		http.Error(w, "Senha inválida!", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Joined room successfully"})
}

func UploadProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Erro ao obter arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	userID := r.FormValue("user_id")
	if userID == "" {
		http.Error(w, "ID do usuário não fornecido", http.StatusBadRequest)
		return
	}

	filename := userID + "_profile" + filepath.Ext("")

	filePath, err := services.UploadProfilePic(file, filename, userID)
	if err != nil {
		http.Error(w, "Erro ao fazer upload da imagem: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"profile_pic": filePath})
}

func UploadImage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Erro ao fazer o parse do form: %v", err)
		http.Error(w, "Erro ao fazer o parse do form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("Erro ao obter o arquivo: %v", err)
		http.Error(w, "Erro ao obter o arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		log.Printf("Erro ao criar diretório de uploads: %v", err)
		http.Error(w, "Erro ao criar diretório de uploads", http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join("uploads", header.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		log.Printf("Erro ao criar arquivo: %v", err)
		http.Error(w, "Erro ao criar arquivo", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		log.Printf("Erro ao salvar arquivo: %v", err)
		http.Error(w, "Erro ao salvar arquivo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Arquivo enviado com sucesso!"))
}
