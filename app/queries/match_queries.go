package queries

import (
	"sidewarslobby/app/models"

	"gorm.io/gorm"
)

type MatchQueries struct {
	DB *gorm.DB
}

func (q *MatchQueries) GetMatch(id int) *models.Match {
	var match models.Match
	res := q.DB.First(&match, id)

	if res.Error != nil {
		return nil
	}

	return &match
}

func (q *MatchQueries) GetMatchByMatchmakingID(mId string) *models.Match {
	var match models.Match
	res := q.DB.First(&match, "matchmaking_id = ?", mId)

	if res.Error != nil {
		return nil
	}

	return &match
}

func (q *MatchQueries) UpdateMatch(match *models.Match) error {
	return q.DB.Save(&match).Error
}

func (q *MatchQueries) FindOrCreateMatch(match *models.Match) (*models.Match, error) {
	if m := q.GetMatchByMatchmakingID(match.MatchmakingID); m != nil {
		return m, nil
	}

	res := q.DB.Create(match)
	return match, res.Error
}

func (q *MatchQueries) GetUserMatch(id int) (*models.UserMatch, error) {
	var userMatch models.UserMatch

	res := q.DB.Model(models.UserMatch{}).Preload("Match").Preload("User").Preload("Match.UserMatches").First(&userMatch, id)

	if res.Error != nil {
		return nil, res.Error
	}

	return &userMatch, nil
}

func (q *MatchQueries) UpdateUserMatch(userMatch *models.UserMatch) error {
	return q.DB.Save(&userMatch).Error
}

func (q *MatchQueries) CreateUserMatch(userMatch *models.UserMatch) error {
	return q.DB.Create(userMatch).Error
}

func (q *MatchQueries) GetMatchUsersByTeamID(match *models.Match, teamID uint8) ([]models.User, error) {
	var userMatches []models.UserMatch
	res := q.DB.Model(match).Preload("User").Association("UserMatches").Find(&userMatches)
	if res != nil {
		return nil, res
	}

	var ret []models.User
	for _, v := range userMatches {
		if v.TeamID == teamID {
			ret = append(ret, v.User)
		}
	}
	return ret, nil
}
