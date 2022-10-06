package models

import "time"

type JobHistory struct {
	ID            int64     `gorm:"column:id"`
	JobCode       string    `gorm:"column:job_code"`
	Started       time.Time `gorm:"column:started"`
	Finished      time.Time `gorm:"column:finished"`
	Result        string    `gorm:"column:result"`
	ResultMessage string    `gorm:"column:result_message"`
}

func (JobHistory) TableName() string {
	return "job_history"
}
