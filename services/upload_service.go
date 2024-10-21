package services

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UploadProfilePic(file io.Reader, filename, userID string) (string, error) {
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		return "", err
	}

	filePath := filepath.Join("uploads", filename)
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", err
	}

	collection := GetUserCollection()

	filter := bson.M{"_id": userID}

	update := bson.M{"$set": bson.M{"profile_pic": filePath}}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func GetProfilePic(w http.ResponseWriter, r *http.Request) {
	imageName := r.URL.Query().Get("filename")
	if imageName == "" {
		http.Error(w, "Nome da imagem não fornecido", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("uploads", imageName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Imagem não encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")

	imageFile, err := os.Open(filePath)
	if err != nil {
		log.Printf("Erro ao abrir o arquivo de imagem: %v", err)
		http.Error(w, "Erro interno ao abrir a imagem", http.StatusInternalServerError)
		return
	}
	defer imageFile.Close()

	_, err = io.Copy(w, imageFile)
	if err != nil {
		log.Printf("Erro ao enviar o arquivo de imagem: %v", err)
		http.Error(w, "Erro ao enviar a imagem", http.StatusInternalServerError)
	}
}

func GetUserCollection() *mongo.Collection {
	return mongoClient.Database("chatDB").Collection("users")
}
