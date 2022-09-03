package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	TeamRed  uint8 = 0
	TeamBlue uint8 = 1
)

type UserMatch struct {
	gorm.Model

	ID uuid.UUID `gorm:"primaryKey"`

	UserID  uint
	MatchID uint
	TeamID  uint8

	UserWon   bool
	ScoreDiff int //Score diff after match is done

	UserChampion uint8

	Match Match `gorm:"foreignKey:MatchID;references:ID"`
	User  User  `gorm:"foreignKey:UserID;references:ID"`
}

type Match struct {
	gorm.Model

	Finished    bool
	UserMatches []UserMatch `gorm:"foreignKey:MatchID"`
}

func (m *Match) GetUsersByTeamID(teamID uint8) []*User {
	var res []*User
	for _, v := range m.UserMatches {
		res = append(res, &v.User)
	}

	return res
}
