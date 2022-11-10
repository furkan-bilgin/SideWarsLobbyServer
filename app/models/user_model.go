package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username   string
	FirebaseID string
	Token      string `gorm:"index;unique"`

	UserMatches []UserMatch `gorm:"foreignKey:UserID"`
	UserInfo    UserInfo    `gorm:"foreignKey:UserID"`
}

type UserInfo struct {
	gorm.Model
	UserID uint

	SelectedChampion uint8
	CachedElo        int
}

func (u *UserInfo) Sanitize() UserInfo {
	return UserInfo{
		SelectedChampion: u.SelectedChampion,
		CachedElo:        u.CachedElo,
	}
}
