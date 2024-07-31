package dao

import (
	"github.com/gorilla/websocket"
)

type User struct {
	ID    int
	Score int
	Conn  *websocket.Conn
}

type Quiz struct {
	ID    int
	Users map[int]*User
}

type ScoreUpdate struct {
	QuizID int `json:"quiz_id"`
	UserID int `json:"user_id"`
	Score  int `json:"score"`
}
