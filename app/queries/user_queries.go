package queries

import "gorm.io/gorm"

type UserQueries struct {
	DB *gorm.DB
}

func (q *UserQueries) GetUserById() {

}
