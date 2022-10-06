package task

import "github.com/supernova0730/job/internal/task/echo"

type Task interface {
	Do() error
	Code() string
}

var registered []Task

func init() {
	registered = []Task{
		&echo.Success{},
		&echo.Error{},
		&echo.Panic{},
		// register new tasks here
	}
}

func FindByCode(code string) Task {
	for _, job := range registered {
		if code == job.Code() {
			return job
		}
	}
	return nil
}
