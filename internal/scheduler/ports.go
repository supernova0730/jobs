package scheduler

import (
	"context"
	"github.com/supernova0730/job/internal/models"
)

type JobRepository interface {
	ListActive(ctx context.Context) (result []models.Job, err error)
	ListCodes(ctx context.Context) (result []string, err error)
	SetRunning(ctx context.Context, code string, isRunning bool) (err error)
}

type JobHistoryRepository interface {
	Insert(ctx context.Context, jobHistory models.JobHistory) (id int64, err error)
	Update(ctx context.Context, id int64, jobHistory models.JobHistory) (err error)
}
