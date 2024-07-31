package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/phongln/elsa-real-time-quiz/internal/models"
	"github.com/phongln/elsa-real-time-quiz/internal/repo"
)

type (
	QuizGetter interface {
		GetActiveQuizByID(ctx context.Context, id int) (*models.Quiz, error)
	}

	QuizGetterService struct {
		redisClient    *redis.Client
		quizGetterRepo repo.QuizGetter
	}
)

func GetQuizGetterService(redisClient *redis.Client, quizGetterRepo repo.QuizGetter) *QuizGetterService {
	return &QuizGetterService{
		redisClient:    redisClient,
		quizGetterRepo: quizGetterRepo,
	}
}

func (q *QuizGetterService) GetActiveQuizByID(ctx context.Context, id int) (quiz *models.Quiz, err error) {
	quiz, err = q.quizGetterRepo.FindByID(ctx, id)
	if err == sql.ErrNoRows {
		return nil, errors.New(fmt.Sprintf("quiz %d not found", id))
	} else if err != nil {
		return nil, err
	}

	if quiz.IsActivated.Int8 < 1 {
		return nil, errors.New(fmt.Sprintf("quiz %d is deactivated", id))
	}

	return
}
