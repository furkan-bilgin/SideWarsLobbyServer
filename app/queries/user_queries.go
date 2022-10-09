package queries

import (
	"sidewarslobby/app/models"
	"sidewarslobby/pkg/repository"

	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserQueries struct {
	DB *gorm.DB
}

func preloadUser(db *gorm.DB) *gorm.DB {
	return db.Model(models.User{}).Preload("UserMatches").Preload("UserInfo")
}

func (q *UserQueries) GetUserById(id uint) *models.User {
	var user models.User
	res := preloadUser(q.DB).First(&user, id)

	if res.Error != nil {
		return nil
	}

	return &user
}

func (q *UserQueries) GetUserByToken(token string) (*models.User, error) {
	var user models.User
	res := preloadUser(q.DB).First(&user, "token = ?", token)
	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func (q *UserQueries) UpdateUserDetails(user models.User, updates models.User) error {
	return q.DB.Model(&user).Updates(updates).Error
}

func (q *UserQueries) UpdateUserInfo(userInfo models.UserInfo, updates models.UserInfo) error {
	return q.DB.Model(&userInfo).Updates(updates).Error
}

// Creates or updates a user, also returns True if a new user record was created
func (q *UserQueries) CreateOrUpdateUser(firebaseUser *auth.UserRecord) (*models.User, bool) {
	var user models.User
	userUpdate := models.User{ProfilePhotoURL: firebaseUser.PhotoURL, Username: firebaseUser.DisplayName, Token: uuid.NewString()}

	res := q.DB.First(&user, "firebase_id = ?", firebaseUser.UID)

	if res.Error != nil {
		// User does not exist, create new one
		user := userUpdate
		user.FirebaseID = firebaseUser.UID

		q.DB.Create(&user)
		return &user, true
	}

	// Else, update data
	q.UpdateUserDetails(user, userUpdate)
	return q.GetUserById(user.ID), false
}

func (q *UserQueries) CacheUserElo(user *models.User) error {
	// Fetch UserMatches
	var userMatches []models.UserMatch
	res := q.DB.Model(user).Preload("Match").Association("UserMatches").Find(&userMatches)
	if res != nil {
		return res
	}

	// Add match diffs into beginner elo
	diff := repository.BeginnerElo
	for _, v := range userMatches {
		if v.Match.Finished {
			diff += v.ScoreDiff
		}
	}
	user.CachedElo = diff
	return q.DB.Save(user).Error
}
