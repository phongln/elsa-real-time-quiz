package services

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/phongln/elsa-real-time-quiz/internal/models"
	"github.com/phongln/elsa-real-time-quiz/internal/repo"
	"github.com/phongln/elsa-real-time-quiz/pkg/kafka"
)

type (
	ScoreSetter interface {
		Update(ctx context.Context, quizID, userID int, score int) (err error)
	}
	ScoreSetterService struct {
		kafkaWriter     kafka.Writer
		scoreSetterRepo repo.ScoreSetter
		scoreGetterRepo repo.ScoreGetter
	}
)

func GetScoreSetterService(kafkaWriter kafka.Writer, scoreSetterRepo repo.ScoreSetter, scoreGetterRepo repo.ScoreGetter) *ScoreSetterService {
	return &ScoreSetterService{
		kafkaWriter:     kafkaWriter,
		scoreSetterRepo: scoreSetterRepo,
		scoreGetterRepo: scoreGetterRepo,
	}
}

func (s *ScoreSetterService) Update(ctx context.Context, quizID, userID, score int) (err error) {
	sc, err := s.scoreGetterRepo.GetByUserIDAndQuizID(ctx, userID, quizID)
	if err == sql.ErrNoRows {
		sc = &models.Score{
			QuizID: quizID,
			UserID: userID,
			Score:  score,
		}

		return s.scoreSetterRepo.InsertScore(ctx, sc)
	} else if err != nil {
		return err
	}

	sc.Score = score
	return s.scoreSetterRepo.UpdateScore(ctx, sc, boil.Whitelist(
		models.ScoreColumns.Score,
	))
}
