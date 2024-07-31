package config

import "fmt"

func GetLeaderBoardRedisKeyByQuizID(quizID int) string {
	return fmt.Sprintf("leaderboard:quiz_id:%d", quizID)
}
