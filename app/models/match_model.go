package models

import (
	"gorm.io/gorm"
)

type UserMatch struct {
	gorm.Model

	MatchID uint
	UserID  uint
	TeamID  uint8

	UserWon   bool
	ScoreDiff int //Score diff after match is done

	UserChampion uint8

	Match Match
	User  User
}

type Match struct {
	gorm.Model
	MatchmakingID string `gorm:"index;unique"`
	Finished      bool
	UserMatches   []UserMatch
}
