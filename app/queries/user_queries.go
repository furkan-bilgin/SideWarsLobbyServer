package queries

import "github.com/jmoiron/sqlx"

type UserQueries struct {
	*sqlx.DB
}

func (q *UserQueries) GetUserById() {

}
