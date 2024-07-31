package handlers

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/gorilla/websocket"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	"github.com/phongln/elsa-real-time-quiz/internal/services"
)

type AnswerRequest struct {
	Answer string `json:"answer"`
	Score  int    `json:"score"`
}

func receivedAnswers(ctx context.Context, ws *websocket.Conn, leaderBoardUpdaterService services.LeaderBoardUpdater, quizID, userID int) (err error) {
	for {
		var answer AnswerRequest
		if err = ws.ReadJSON(&answer); err != nil {
			logrus.Error("Error reading JSON:", err)
			return
		}

		if lo.IsEmpty(answer.Answer) {
			err = errors.New("Answer is required and should not empty")
			if wsErr := ws.WriteJSON(err.Error()); wsErr != nil {
				logrus.Error("write:", wsErr)
				return
			}
			continue
		}

		scores, er := leaderBoardUpdaterService.Update(ctx, quizID, userID, answer.Score)
		if er != nil {
			logrus.Error("Error updating score: ", er)
			continue
		}

		sort.Slice(scores, func(i, j int) bool {
			return scores[i].Score >= scores[j].Score
		})

		for _, user := range quizzes[quizID].Users {
			er = user.Conn.WriteJSON(scores)
			if websocket.IsCloseError(er, []int{
				websocket.CloseMessage,
				websocket.CloseNormalClosure,
			}...) {
				delete(quizzes[quizID].Users, userID)
			} else if er != nil {
				logrus.Error(fmt.Sprintf("write: %s", er.Error()))
			}
		}
	}
}
