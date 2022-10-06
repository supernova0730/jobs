package main

import (
	"context"
	"github.com/supernova0730/job/config"
	"github.com/supernova0730/job/internal/repository"
	"github.com/supernova0730/job/internal/scheduler"
	"github.com/supernova0730/job/pkg/logger"
	"github.com/supernova0730/job/pkg/postgres"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()

	logger.Init()

	conf, err := config.Load("")
	if err != nil {
		logger.Log.Fatal("failed to load config", zap.Error(err))
	}

	logger.Log.Info("config successfully loaded", zap.Any("config", conf))

	dbConfig := postgres.Config{
		Host:     conf.DBHost,
		Port:     conf.DBPort,
		Username: conf.DBUser,
		Password: conf.DBPassword,
		DBName:   conf.DBName,
		SSLMode:  conf.DBSSLMode,
	}
	db, err := postgres.Connect(ctx, dbConfig)
	if err != nil {
		logger.Log.Fatal("failed to connect to database", zap.Error(err))
	}

	logger.Log.Info("connected to database...")

	jobRepo := repository.NewJobRepository(db)
	jobHistoryRepo := repository.NewJobHistoryRepository(db)

	schd := scheduler.New(jobRepo, jobHistoryRepo)
	go schd.Start(ctx, conf.SchedulerRefreshRate)
	logger.Log.Info("scheduler started...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	logger.Log.Info("shutting down...")

	schd.Stop()
	logger.Log.Info("scheduler stopped")
}
