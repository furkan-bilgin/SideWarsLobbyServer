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

func (q *UserQueries) GetUserById(id uint) models.User {
	var user models.User
	q.DB.First(&user, id)

	return user
}

func (q *UserQueries) UpdateUserDetails(user models.User, updates models.User) {
	q.DB.Model(&user).Updates(updates)
}

func (q *UserQueries) CreateOrUpdateUser(firebaseUser *auth.UserRecord) models.User {
	var user models.User
	userUpdate := models.User{ProfilePhotoURL: firebaseUser.PhotoURL, Username: firebaseUser.DisplayName, Token: uuid.NewString()}

	res := q.DB.First(&user, "firebase_id = ?", firebaseUser.UID)

	if res.Error != nil {
		// User does not exist, create new one
		user := userUpdate
		user.FirebaseID = firebaseUser.UID

		q.DB.Create(&user)
		return user
	}

	// Else, update data
	q.UpdateUserDetails(user, userUpdate)
	return q.GetUserById(user.ID)
}

func (q *UserQueries) CacheUserScore(user *models.User) error {
	user.CachedElo = user.CalculateElo()
	return q.DB.Save(user).Error
}
