package models

import (
	"time"
)

type (
	User struct {
		UserId         int64     `json:"UserId" gorm:"index:idx_user_id,unique;not null"`
		UserName       string    `json:"UserName" gorm:"not null"`
		CreateDateTime time.Time `json:"createDateTime" gorm:"not null"`
		Active         byte      `json:"active" gorm:"not null"`
		UpdateDateTime time.Time `json:"updateDateTime" gorm:"not null"`
	}
)
