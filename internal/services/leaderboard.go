package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"

	"github.com/phongln/elsa-real-time-quiz/internal/config"
	"github.com/phongln/elsa-real-time-quiz/internal/dao"
	"github.com/phongln/elsa-real-time-quiz/internal/models"
	"github.com/phongln/elsa-real-time-quiz/internal/repo"
	"github.com/phongln/elsa-real-time-quiz/pkg/kafka"
)

type (
	LeaderBoardUpdater interface {
		Update(ctx context.Context, quizID, userID, score int) (scores models.ScoreSlice, err error)
	}
	LeaderBoardUpdaterService struct {
		kafkaWriter     kafka.Writer
		redisClient     *redis.Client
		scoreSetterRepo repo.ScoreSetter
		scoreGetterRepo repo.ScoreGetter
	}
)

func GetLeaderBoardUpdaterService(kafkaWriter kafka.Writer, redisClient *redis.Client, scoreSetterRepo repo.ScoreSetter, scoreGetterRepo repo.ScoreGetter) *LeaderBoardUpdaterService {
	return &LeaderBoardUpdaterService{
		kafkaWriter:     kafkaWriter,
		redisClient:     redisClient,
		scoreSetterRepo: scoreSetterRepo,
		scoreGetterRepo: scoreGetterRepo,
	}
}

func (s *LeaderBoardUpdaterService) Update(ctx context.Context, quizID, userID, score int) (scores models.ScoreSlice, err error) {
	leaderBoardRedisKey := config.GetLeaderBoardRedisKeyByQuizID(quizID)
	leaderBoard := s.redisClient.Get(ctx, leaderBoardRedisKey)
	if leaderBoard.Err() != nil {
		scores, err = s.scoreGetterRepo.GetByQuizID(ctx, quizID)
		if err != nil {
			logrus.Error("Error getting scores:", err)
			return
		}
	} else {
		err = json.Unmarshal([]byte(leaderBoard.Val()), &scores)
		if err != nil {
			logrus.Error("Error unmarshal scores:", err)
			return
		}
	}

	if len(scores) == 0 {
		scores = models.ScoreSlice{
			&models.Score{
				QuizID: quizID,
				UserID: userID,
				Score:  score,
			},
		}
	}

	sc := dao.ScoreUpdate{
		QuizID: quizID,
		UserID: userID,
		Score:  score,
	}
	isExistScore := false
	for _, s := range scores {
		if s.QuizID == sc.QuizID && s.UserID == sc.UserID {
			s.Score += sc.Score
			sc.Score = s.Score
			isExistScore = true
			break
		}
	}

	if !isExistScore {
		scores = append(scores, &models.Score{
			QuizID: quizID,
			UserID: userID,
			Score:  score,
		})
	}

	msg, err := json.Marshal(scores)
	if err != nil {
		logrus.Error("Error marshalling scores:", err)
		return
	}
	if sttCmd := s.redisClient.Set(ctx, leaderBoardRedisKey, msg, redis.KeepTTL); sttCmd.Err() != nil {
		logrus.Error("Error setting scores:", sttCmd.Err())
		return
	}

	// s.redisClient.ZAdd(context.Background(), leaderBoardRedisKey, &redis.Z{
	// 	Score:  float64(score),
	// 	Member: userID,
	// })

	message, err := json.Marshal(sc)
	if err != nil {
		logrus.Error("Error marshalling score update:", err)
		return
	}

	key, err := json.Marshal(fmt.Sprintf("quiz_%d", quizID))
	if err != nil {
		logrus.Error("Error marshalling score key:", err)
		return
	}

	go s.kafkaWriter.ProduceMessageWithKey(ctx, config.KafkaScoreUpdatesTopic, key, message)

	return
}
