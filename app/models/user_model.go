package models

import (
	"sidewarslobby/pkg/repository"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username        string
	FirebaseID      string
	ProfilePhotoURL string
	Token           string

	CachedElo int

	UserMatches []UserMatch
}

func (u *User) CalculateElo() int {
	diff := repository.BeginnerElo
	for _, v := range u.UserMatches {
		if v.Finished {
			diff += v.ScoreDiff
		}
	}

	return diff
}
