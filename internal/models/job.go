package models

type Job struct {
	ID        int64  `gorm:"column:id"`
	Code      string `gorm:"column:code"`
	Schedule  string `gorm:"column:schedule"`
	IsActive  bool   `gorm:"column:is_active"`
	IsRunning bool   `gorm:"column:is_running"`
}

func (Job) TableName() string {
	return "job"
}
