package models

import (
	"time"
)

type (
	User struct {
		UserID         int64     `json:"userID" gorm:"index:idx_user_id,unique;not null"`
		UserName       string    `json:"userName" gorm:"not null"`
		CreateDateTime time.Time `json:"createDateTime" gorm:"not null"`
		Active         byte      `json:"active" gorm:"not null"`
		UpdateDateTime time.Time `json:"updateDateTime" gorm:"not null"`
	}
)
