package services

import (
	"chat-go/models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// jogar pro contexto depois
var mongoClient *mongo.Client
var messageCollection *mongo.Collection
var roomCollection *mongo.Collection

func ConnectMongoDB() {
	var err error
	mongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Erro ao conectar ao MongoDB: ", err)
	}
	messageCollection = mongoClient.Database("chatDB").Collection("messages")
	roomCollection = mongoClient.Database("chatDB").Collection("rooms")
}

func GetRoomCollection() *mongo.Collection {
	return roomCollection
}

func GetMessageCollection() *mongo.Collection {
	return messageCollection
}

func SaveMessage(msg models.Message) error {
	_, err := messageCollection.InsertOne(context.TODO(), msg)
	return err
}
