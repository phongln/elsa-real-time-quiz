package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/phongln/elsa-real-time-quiz/internal/handlers"
	"github.com/phongln/elsa-real-time-quiz/internal/infra/db"
	"github.com/phongln/elsa-real-time-quiz/internal/repo"
	"github.com/phongln/elsa-real-time-quiz/internal/services"
	"github.com/phongln/elsa-real-time-quiz/pkg/kafka"
	"github.com/phongln/elsa-real-time-quiz/pkg/logger"
	"github.com/phongln/elsa-real-time-quiz/pkg/middleware"
	"github.com/phongln/elsa-real-time-quiz/pkg/redis"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logger.InitLogger()
	db.InitMySQL()

	redis := redis.GetRedisClient()
	kafkaWriter := kafka.GetKafkaWriter([]string{os.Getenv("KAFKA_BROKER")})

	scoreSetterRepo := repo.GetScoreSetterRepo(boil.GetContextDB())
	scoreGetterRepo := repo.GetScoreGetterRepo(boil.GetContextDB())
	scoreSetterService := services.GetScoreSetterService(kafkaWriter, scoreSetterRepo, scoreGetterRepo)

	quizGetterRepo := repo.GetQuizGetterRepo(boil.GetContextDB())
	quizGetterService := services.GetQuizGetterService(redis, quizGetterRepo)

	leaderBoardUpdaterService := services.GetLeaderBoardUpdaterService(kafkaWriter, redis, scoreSetterRepo, scoreGetterRepo)

	handlers.ConsumeScoreUpdate(scoreSetterService)

	r := gin.Default()
	r.Use(middleware.LoggerMiddleware())

	r.GET("/ws", handlers.HandleConnections(leaderBoardUpdaterService, quizGetterService))
	r.GET("/metrics", gin.WrapH(middleware.MetricsHandler()))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
