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

	UserMatches []UserMatch `gorm:"foreignKey:UserID"`
	UserInfo    UserInfo    `gorm:"foreignKey:UserID"`
}

type UserInfo struct {
	gorm.Model
	UserID uint

	SelectedChampion uint8
}

func (u *User) CalculateElo() int {
	diff := repository.BeginnerElo
	for _, v := range u.UserMatches {
		if v.Match.Finished {
			diff += v.ScoreDiff
		}
	}

	return diff
}
