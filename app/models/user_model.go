package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username        string
	FirebaseID      string
	ProfilePhotoURL string
	Token           string `gorm:"index;unique"`

	CachedElo int

	UserMatches []UserMatch `gorm:"foreignKey:UserID"`
	UserInfo    UserInfo    `gorm:"foreignKey:UserID"`
}

type UserInfo struct {
	gorm.Model
	UserID uint

	SelectedChampion uint8
}
