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
	MatchmakingID uuid.UUID `gorm:"index;unique"`

	Finished    bool
	UserMatches []UserMatch
}

func (m *Match) GetUsersByTeamID(teamID uint8) []*User {
	var res []*User
	for _, v := range m.UserMatches {
		res = append(res, &v.User)
	}

	return res
}
