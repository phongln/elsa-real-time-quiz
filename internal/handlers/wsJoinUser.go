package handlers

import (
	"context"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/phongln/elsa-real-time-quiz/internal/dao"
	"github.com/phongln/elsa-real-time-quiz/internal/services"
)

type UserJoinQuizRequest struct {
	UserID int `json:"user_id"`
	QuizID int `json:"quiz_id"`
}

func joiningUserToQuiz(ctx context.Context, ws *websocket.Conn, quizGetterService services.QuizGetter) (quizID, userID int, err error) {
	var userQuiz UserJoinQuizRequest
	if err = ws.ReadJSON(&userQuiz); err != nil {
		logrus.Error("Error reading JSON:", err)
		return
	}

	quizID = userQuiz.QuizID
	userID = userQuiz.UserID

	quizzesLock.Lock()

	if _, err = quizGetterService.GetActiveQuizByID(ctx, quizID); err != nil {
		logrus.Error(err)
		if wsErr := ws.WriteJSON(err.Error()); wsErr != nil {
			logrus.Error("write:", wsErr)
		}
		return
	}

	if _, ok := quizzes[quizID]; !ok {
		quizzes[quizID] = &dao.Quiz{ID: quizID, Users: make(map[int]*dao.User)}
	}
	if _, ok := quizzes[quizID].Users[userID]; !ok {
		quizzes[quizID].Users[userID] = &dao.User{ID: userID, Score: 0, Conn: ws}
	} else {
		quizzes[quizID].Users[userID].Conn = ws
	}
	quizzesLock.Unlock()

	for _, user := range quizzes[quizID].Users {
		err = user.Conn.WriteJSON(fmt.Sprintf("User %d joined quiz %d successfully", userID, quizID))
		if err != nil {
			logrus.Error("write:", err)
		}
	}

	return
}
