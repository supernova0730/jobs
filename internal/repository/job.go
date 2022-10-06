package repository

import (
	"context"
	"github.com/supernova0730/job/internal/models"
	"github.com/supernova0730/job/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type JobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (repo *JobRepository) ListActive(ctx context.Context) (result []models.Job, err error) {
	defer func() {
		if err != nil {
			logger.Log.Error(
				"JobRepository.ListActive failed",
				zap.Error(err),
			)
		}
	}()

	err = repo.db.
		Model(&models.Job{}).
		Where("is_active = ?", true).
		Find(&result).Error
	return
}

func (repo *JobRepository) ListCodes(ctx context.Context) (result []string, err error) {
	defer func() {
		if err != nil {
			logger.Log.Error(
				"JobRepository.ListActiveCodes failed",
				zap.Error(err),
			)
		}
	}()

	err = repo.db.
		Model(&models.Job{}).
		Select("code").
		Find(&result).Error
	return
}

func (repo *JobRepository) SetRunning(ctx context.Context, code string, isRunning bool) (err error) {
	defer func() {
		if err != nil {
			logger.Log.Error(
				"JobRepository.SetRunning failed",
				zap.Error(err),
				zap.String("code", code),
				zap.Bool("isRunning", isRunning),
			)
		}
	}()

	err = repo.db.
		Model(&models.Job{}).
		Where("code = ?", code).
		Update("is_running", isRunning).Error
	return
}

func (repo *JobRepository) SetActive(ctx context.Context, code string, isActive bool) (err error) {
	defer func() {
		if err != nil {
			logger.Log.Error(
				"JobRepository.SetActive failed",
				zap.Error(err),
				zap.String("code", code),
				zap.Bool("isActive", isActive),
			)
		}
	}()

	err = repo.db.
		Model(&models.Job{}).
		Where("code = ?", code).
		Update("is_active", isActive).Error
	return
}

func (repo *JobRepository) SetSchedule(ctx context.Context, code string, schedule string) (err error) {
	defer func() {
		if err != nil {
			logger.Log.Error(
				"JobRepository.SetSchedule failed",
				zap.Error(err),
				zap.String("code", code),
				zap.String("schedule", schedule),
			)
		}
	}()

	err = repo.db.
		Model(&models.Job{}).
		Where("code = ?", code).
		Update("schedule", schedule).Error
	return
}
