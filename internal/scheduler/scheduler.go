package scheduler

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/supernova0730/job/internal/models"
	"github.com/supernova0730/job/internal/task"
	"github.com/supernova0730/job/pkg/logger"
	"github.com/supernova0730/job/pkg/uuid"
	"go.uber.org/zap"
	"time"
)

const (
	ResultUnknown = "UNKNOWN"
	ResultSuccess = "SUCCESS"
	ResultError   = "ERROR"
	ResultPanic   = "PANIC"
)

type Scheduler struct {
	schd           *gocron.Scheduler
	jobRepo        JobRepository
	jobHistoryRepo JobHistoryRepository
}

func New(
	jobRepo JobRepository,
	jobHistoryRepo JobHistoryRepository,
) *Scheduler {
	return &Scheduler{
		schd:           gocron.NewScheduler(time.UTC),
		jobRepo:        jobRepo,
		jobHistoryRepo: jobHistoryRepo,
	}
}

func (s *Scheduler) registerTasks(ctx context.Context) (err error) {
	jobs, err := s.jobRepo.ListActive(ctx)
	if err != nil {
		return
	}

	for _, job := range jobs {
		registeredTask := task.FindByCode(job.Code)
		if registeredTask == nil {
			return fmt.Errorf("not registered task: %s", job.Code)
		}

		err = s.registerTask(ctx, registeredTask, job.Schedule)
		if err != nil {
			return
		}
	}

	return nil
}

func (s *Scheduler) registerTask(ctx context.Context, task task.Task, schedule string) error {
	_, err := s.schd.CronWithSeconds(schedule).Tag(task.Code()).Do(func() {
		id := uuid.Generate()
		log := logger.Log.
			Named(fmt.Sprintf("[%s]", task.Code())).
			With(zap.String("id", id))

		log.Info("started")
		started := time.Now()

		jobHistoryID, err := s.before(ctx, task.Code())
		if err != nil {
			log.Error("set job starting failed", zap.Error(err))
		}

		defer func() {
			if r := recover(); r != nil {
				log.Error("panic", zap.Any("recover", r))
				if err := s.after(ctx, task.Code(), jobHistoryID, ResultPanic, r.(error).Error()); err != nil {
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

		err = s.after(ctx, task.Code(), jobHistoryID, result, resultMessage)
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

func (s *Scheduler) refresh(ctx context.Context) (err error) {
	codes, err := s.jobRepo.ListCodes(ctx)
	if err != nil {
		return
	}

	for _, code := range codes {
		err = s.schd.RemoveByTag(code)
		if err != nil && !errors.Is(err, gocron.ErrJobNotFoundWithTag) {
			return
		}
	}

	return s.registerTasks(ctx)
}

func (s *Scheduler) before(ctx context.Context, code string) (jobHistoryID int64, err error) {
	err = s.jobRepo.SetRunning(ctx, code, true)
	if err != nil {
		return
	}

	return s.jobHistoryRepo.Insert(ctx, models.JobHistory{
		JobCode: code,
		Started: time.Now(),
		Result:  ResultUnknown,
	})
}

func (s *Scheduler) after(
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

func (s *Scheduler) Start(ctx context.Context, refreshRate time.Duration) {
	s.schd.StartAsync()

	ticker := time.NewTicker(refreshRate)
	for {
		err := s.refresh(ctx)
		if err != nil {
			logger.Log.Fatal("failed to reset scheduler", zap.Error(err))
		}

		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
		}
	}
}

func (s *Scheduler) Stop() {
	s.schd.Stop()
}
