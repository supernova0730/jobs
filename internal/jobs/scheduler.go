package jobs

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/supernova0730/job/internal/models"
	"github.com/supernova0730/job/internal/repository"
	"github.com/supernova0730/job/internal/task"
	"github.com/supernova0730/job/pkg/logger"
	"github.com/supernova0730/job/pkg/uuid"
	"go.uber.org/zap"
	"time"
)

const (
	ResultSuccess = "SUCCESS"
	ResultError   = "ERROR"
	ResultPanic   = "PANIC"
)

type Scheduler struct {
	s              *gocron.Scheduler
	jobRepo        *repository.JobRepository
	jobHistoryRepo *repository.JobHistoryRepository
}

func NewScheduler(
	jobRepo *repository.JobRepository,
	jobHistoryRepo *repository.JobHistoryRepository,
) *Scheduler {
	return &Scheduler{
		s:              gocron.NewScheduler(time.UTC),
		jobRepo:        jobRepo,
		jobHistoryRepo: jobHistoryRepo,
	}
}

func (s *Scheduler) RegisterTasks(ctx context.Context) (err error) {
	jobs, err := s.jobRepo.ListActive(ctx)
	if err != nil {
		return
	}

	for _, job := range jobs {
		registeredTask := task.FindByCode(job.Code)
		if registeredTask == nil {
			return fmt.Errorf("not registered task: %s", job.Code)
		}

		err = s.registerTask(ctx, job, registeredTask)
		if err != nil {
			return
		}
	}

	return nil
}

func (s *Scheduler) registerTask(ctx context.Context, job models.Job, task task.Task) error {
	_, err := s.s.CronWithSeconds(job.Schedule).Tag(job.Code).Do(func() {
		id := uuid.Generate()
		log := logger.Log.
			Named(fmt.Sprintf("[%s]", job.Code)).
			With(zap.String("id", id))

		log.Info("started")
		started := time.Now()

		jobHistoryID, err := s.setJobStarting(ctx, job.Code)
		if err != nil {
			log.Error("set job starting failed", zap.Error(err))
		}

		defer func() {
			if r := recover(); r != nil {
				log.Error("panic", zap.Any("recover", r))
				if err := s.setJobStopping(ctx, job.Code, jobHistoryID, ResultPanic, r.(error).Error()); err != nil {
					log.Error("set job stopping failed", zap.Error(err))
				}
			}
		}()

		result := ResultSuccess
		resultMessage := ""

		err = task.Do()
		if err != nil {
			log.Error("Do() failed", zap.Error(err))
			result = ResultError
			resultMessage = err.Error()
		}

		err = s.setJobStopping(ctx, job.Code, jobHistoryID, result, resultMessage)
		if err != nil {
			log.Error("set job stopping failed", zap.Error(err))
		}

		log.Info(
			"finished",
			zap.String("result", result),
			zap.String("resultMessage", resultMessage),
			zap.Duration("elapsed", time.Since(started)),
		)
	})
	return err
}

func (s *Scheduler) setJobStarting(ctx context.Context, code string) (jobHistoryID int64, err error) {
	err = s.jobRepo.SetRunning(ctx, code, true)
	if err != nil {
		return
	}

	return s.jobHistoryRepo.Insert(ctx, models.JobHistory{
		JobCode: code,
		Started: time.Now(),
	})
}

func (s *Scheduler) setJobStopping(
	ctx context.Context,
	code string,
	jobHistoryID int64,
	result string,
	resultMessage string,
) (err error) {
	err = s.jobRepo.SetRunning(ctx, code, false)
	if err != nil {
		return
	}

	return s.jobHistoryRepo.Update(ctx, jobHistoryID, models.JobHistory{
		Finished:      time.Now(),
		Result:        result,
		ResultMessage: resultMessage,
	})
}

func (s *Scheduler) Start() {
	s.s.StartBlocking()
}

func (s *Scheduler) Stop() {
	s.s.Stop()
}
