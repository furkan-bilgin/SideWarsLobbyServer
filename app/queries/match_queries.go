package queries

import (
	"sidewarslobby/app/models"

	"gorm.io/gorm"
)

type MatchQueries struct {
	DB *gorm.DB
}

func (q *MatchQueries) GetUserMatch(token string) (*models.UserMatch, error) {
	var userMatch models.UserMatch

	res := q.DB.First(&userMatch, token)

	if res.Error != nil {
		return nil, res.Error
	}

	return &userMatch, nil
}

func (q *MatchQueries) UpdateUserMatch(userMatch *models.UserMatch) error {
	return q.DB.Save(&userMatch).Error
}

func (q *MatchQueries) UpdateMatch(match *models.Match) error {
	return q.DB.Save(&match).Error
}
