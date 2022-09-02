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
	Token           string

	CachedScore int

	UserMatches []UserMatch
}

func (u *User) CalculateScore() int {
	diff := 0
	for _, v := range u.UserMatches {
		diff += v.ScoreDiff
	}

	return diff
}
