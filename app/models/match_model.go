package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserMatch struct {
	gorm.Model

	ID        uuid.UUID `gorm:"primaryKey"`
	CreatedAt time.Time

	UserID    uint
	MatchID   uint
	ScoreDiff int //Score diff after match is done

	Finished bool
	UserWon  bool

	Match Match `gorm:"foreignKey:MatchID;references:ID"`
}

type Match struct {
	ID uint

	CreatedAt   time.Time
	UserMatches []UserMatch `gorm:"foreignKey:MatchID"`
}
