package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/phongln/elsa-real-time-quiz/internal/config"
	"github.com/phongln/elsa-real-time-quiz/internal/dao"
	"github.com/phongln/elsa-real-time-quiz/internal/services"
	"github.com/phongln/elsa-real-time-quiz/pkg/kafka"
)

var (
	quizzes     = make(map[int]*dao.Quiz)
	quizzesLock sync.Mutex
)

func HandleConnections(leaderBoardUpdaterService services.LeaderBoardUpdater, quizGetterService services.QuizGetter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			logrus.Error("Error upgrading to websocket:", err)
			return
		}
		defer ws.Close()

		var quizID, userID int
		for {
			quizID, userID, err = joiningUserToQuiz(ctx.Request.Context(), ws, quizGetterService)
			if err != nil {
				continue
			}
			break
		}

		if err = receivedAnswers(ctx.Request.Context(), ws, leaderBoardUpdaterService, quizID, userID); err != nil {
			return
		}
	}
}

func ConsumeScoreUpdate(scoreSetterService services.ScoreSetter) {
	scoreUpdateKafkaReader := kafka.GetKafkaReader(
		[]string{os.Getenv("KAFKA_BROKER")},
		config.KafkaScoreUpdatesTopic,
		config.KafkaConsumerGroup,
	)

	go func() {
		ctx := context.Background()
		for {
			m, err := scoreUpdateKafkaReader.ReadMessage(ctx)
			if err != nil {
				logrus.Fatal("Error reading message from Kafka: ", err)
				continue
			}

			var scoreUpdate dao.ScoreUpdate
			if err := json.Unmarshal(m.Value, &scoreUpdate); err != nil {
				logrus.Fatal("Error unmarshalling JSON: ", err)
				continue
			}

			if err = scoreSetterService.Update(ctx, scoreUpdate.QuizID, scoreUpdate.UserID, scoreUpdate.Score); err != nil {
				logrus.Fatal("ConsumeKafkaMessages:leaderBoardService:Update: ", err)
			}
		}
	}()
}
