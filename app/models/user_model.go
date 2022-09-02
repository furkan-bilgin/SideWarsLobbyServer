package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID              uint `gorm:"primaryKey"`
	CreatedAt       time.Time
	Username        string
	FirebaseID      string
	ProfilePhotoURL string

	UserMatches []UserMatch
}
