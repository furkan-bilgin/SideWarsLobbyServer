package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

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
		if v.Finished {
			diff += v.ScoreDiff
		}
	}

	return diff
}
