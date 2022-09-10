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

	MatchmakingID uuid.UUID `gorm:"index"`
	UserID        uint
	TeamID        uint8
	Token         string

	UserWon   bool
	ScoreDiff int //Score diff after match is done

	UserChampion uint8

	Match Match `gorm:"foreignKey:MatchmakingID;references:MatchmakingID"`
	User  User
}

type Match struct {
	gorm.Model

	MatchmakingID uuid.UUID `gorm:"index;unique"`

	Finished    bool
	UserMatches []UserMatch `gorm:"foreignKey:MatchmakingID;references:MatchmakingID"`
}

func (m *Match) GetUsersByTeamID(teamID uint8) []*User {
	var res []*User
	for _, v := range m.UserMatches {
		res = append(res, &v.User)
	}

	return res
}
