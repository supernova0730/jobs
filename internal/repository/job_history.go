package repository

import (
	"context"
	"github.com/supernova0730/job/internal/models"
	"github.com/supernova0730/job/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type JobHistoryRepository struct {
	db *gorm.DB
}

func NewJobHistoryRepository(db *gorm.DB) *JobHistoryRepository {
	return &JobHistoryRepository{db: db}
}

func (repo *JobHistoryRepository) Insert(ctx context.Context, jobHistory models.JobHistory) (id int64, err error) {
	defer func() {
		if err != nil {
			logger.Log.Error(
				"JobHistoryRepository.Insert failed",
				zap.Error(err),
				zap.Any("jobHistory", jobHistory),
			)
		}
	}()

	err = repo.db.
		Model(&models.JobHistory{}).
		Save(&jobHistory).Error
	if err != nil {
		return
	}

	return jobHistory.ID, nil
}

func (repo *JobHistoryRepository) Update(ctx context.Context, id int64, jobHistory models.JobHistory) (err error) {
	defer func() {
		if err != nil {
			logger.Log.Error(
				"JobHistoryRepository.Update failed",
				zap.Error(err),
				zap.Int64("id", id),
				zap.Any("jobHistory", jobHistory),
			)
		}
	}()

	err = repo.db.
		Model(&models.JobHistory{}).
		Where("id = ?", id).
		Updates(&jobHistory).Error
	return
}
