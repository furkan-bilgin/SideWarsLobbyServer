package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserMatch struct {
	gorm.Model

	ID uuid.UUID `gorm:"primaryKey"`

	UserID    uint
	MatchID   uint
	ScoreDiff int //Score diff after match is done

	Finished bool
	UserWon  bool

	UserChampion int

	Match Match `gorm:"foreignKey:MatchID;references:ID"`
	User  User  `gorm:"foreignKey:UserID;references:ID"`
}

type Match struct {
	gorm.Model

	UserMatches []UserMatch `gorm:"foreignKey:MatchID"`
}
