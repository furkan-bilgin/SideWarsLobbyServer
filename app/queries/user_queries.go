package queries

import (
	"sidewarslobby/app/models"

	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserQueries struct {
	DB *gorm.DB
}

func (q *UserQueries) GetUserById(id uint) *models.User {
	var user models.User
	res := q.DB.First(&user, id)

	if res.Error != nil {
		return nil
	}

	return &user
}

func (q *UserQueries) GetUserByToken(token string) *models.User {
	var user models.User
	res := q.DB.First(&user, "token = ?", token)
	if res.Error != nil {
		return nil
	}

	return &user
}

func (q *UserQueries) UpdateUserDetails(user models.User, updates models.User) {
	q.DB.Model(&user).Updates(updates)
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
	user.CachedElo = user.CalculateElo()
	return q.DB.Save(user).Error
}
