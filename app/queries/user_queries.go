package queries

import "gorm.io/gorm"

type UserQueries struct {
	*gorm.DB
}

func (q *UserQueries) GetUserById() {

}
