package models

import (
	"time"
)

type Room struct {
	ID       string    `json:"id" bson:"_id"`
	Name     string    `json:"name" bson:"name"`
	Password string    `json:"password,omitempty" bson:"password"`
	Created  time.Time `json:"created" bson:"created"`
}

type Message struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Text     string `json:"text"`
	Room     string `json:"room"`
	ImageURL string `json:"image_url,omitempty"`
	UserID   string `json:"user_id"`
	FileURL  string `json:"file_url,omitempty"`
}

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ProfilePic string `json:"profile_pic"`
}
