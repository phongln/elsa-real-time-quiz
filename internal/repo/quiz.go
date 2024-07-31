package repo

import (
	"context"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/phongln/elsa-real-time-quiz/internal/models"
)

type (
	QuizGetter interface {
		FindByID(ctx context.Context, id int) (*models.Quiz, error)
	}

	QuizGetterRepo struct {
		dbConn boil.ContextExecutor
	}
)

func GetQuizGetterRepo(dbConn boil.ContextExecutor) *QuizGetterRepo {
	return &QuizGetterRepo{
		dbConn: dbConn,
	}
}

func (repo *QuizGetterRepo) FindByID(ctx context.Context, id int) (*models.Quiz, error) {
	return models.Quizzes(
		models.QuizWhere.ID.EQ(id),
		models.QuizWhere.IsActivated.EQ(null.Int8From(1)),
	).One(ctx, repo.dbConn)
}
