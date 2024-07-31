package repo

import (
	"context"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/phongln/elsa-real-time-quiz/internal/models"
)

type (
	ScoreSetter interface {
		UpdateScore(ctx context.Context, score *models.Score, cols boil.Columns) error
		InsertScore(ctx context.Context, score *models.Score) error
		Upsert(ctx context.Context, score *models.Score, updatedCols boil.Columns) error
	}
	ScoreGetter interface {
		GetByUserIDAndQuizID(ctx context.Context, userID, quizID int) (*models.Score, error)
		GetByQuizID(ctx context.Context, quizID int) ([]*models.Score, error)
	}

	ScoreSetterRepo struct {
		dbConn boil.ContextExecutor
	}
	ScoreGetterRepo struct {
		dbConn boil.ContextExecutor
	}
)

func GetScoreSetterRepo(dbConn boil.ContextExecutor) *ScoreSetterRepo {
	return &ScoreSetterRepo{
		dbConn: dbConn,
	}
}

func GetScoreGetterRepo(dbConn boil.ContextExecutor) *ScoreGetterRepo {
	return &ScoreGetterRepo{
		dbConn: dbConn,
	}
}

func (repo *ScoreSetterRepo) UpdateScore(ctx context.Context, score *models.Score, cols boil.Columns) error {
	_, err := score.Update(ctx, repo.dbConn, cols)

	return err
}

func (repo *ScoreSetterRepo) InsertScore(ctx context.Context, score *models.Score) error {
	return score.Insert(ctx, repo.dbConn, boil.Infer())
}

func (repo *ScoreSetterRepo) Upsert(ctx context.Context, score *models.Score, updatedCols boil.Columns) error {
	return score.Upsert(ctx, repo.dbConn, updatedCols, boil.Infer())
}

func (repo *ScoreGetterRepo) GetByUserIDAndQuizID(ctx context.Context, userID, quizID int) (*models.Score, error) {
	return models.Scores(
		models.ScoreWhere.UserID.EQ(userID),
		models.ScoreWhere.QuizID.EQ(quizID),
	).One(ctx, repo.dbConn)
}

func (repo *ScoreGetterRepo) GetByQuizID(ctx context.Context, quizID int) ([]*models.Score, error) {
	return models.Scores(
		models.ScoreWhere.QuizID.EQ(quizID),
	).All(ctx, repo.dbConn)
}
