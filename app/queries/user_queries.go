package queries

import (
	"sidewarslobby/app/models"

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

func (q *UserQueries) CacheUserScore(user models.User) {
	q.DB.Model(&user).Update("cached_score", user.CalculateScore())
}
